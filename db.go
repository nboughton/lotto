package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	qGen "github.com/nboughton/go-sqgenlite"
)

var (
	sqlPragmas = "PRAGMA journal_mode=WAL;	PRAGMA busy_timeout=5000"
	sqlSchema  = `
	CREATE TABLE IF NOT EXISTS results 
	(id INTEGER PRIMARY KEY AUTOINCREMENT, 
	date DATETIME,
	ball_set INT,
	ball_machine TEXT,
	num_1 INT, 
	num_2 INT, 
	num_3 INT, 
	num_4 INT, 
	num_5 INT, 
	num_6 INT, 
	bonus INT)`
	formatSqlite   = "2006-01-02 15:04:05-07:00"
	formatYYYYMMDD = "2006-01-02"
)

// AppDB is a wrapper for *sql.DB so I can extend it by adding my own methods
type AppDB struct {
	*sql.DB
}

type dbRow struct {
	Date    time.Time `json:"date"`
	Machine string    `json:"machine"`
	Set     int       `json:"set"`
	Num     []int     `json:"num"`
}

type queryParams struct {
	Type, Machine string
	Set, Query    int
	Start, End    time.Time
}

func connectDB(path string) *AppDB {
	// Connect to the database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Disable connection pooling
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(0)

	// Create DB schema if it doesn't exist
	if _, err := db.Exec(sqlSchema); err != nil {
		log.Fatal(err.Error())
	}

	// Set PRAGMAs
	if _, err := db.Exec(sqlPragmas); err != nil {
		log.Fatal(err.Error())
	}

	aDB := &AppDB{db}

	// Check if I have any data yet, otherwise populate the db
	log.Println("Checking database")
	if rows, _ := aDB.getRowCount(); rows == 0 {
		log.Println("No data, scraping site")
		if err := aDB.populateDB(); err != nil {
			log.Fatal(err.Error())
		}
	}

	// run update
	log.Println("Updating database")
	if err := aDB.updateDB(); err != nil {
		log.Fatal(err.Error())
	}

	return aDB
}

// Apply filters for queries, always run this before executing a query as it edits the
// query in place by pointer
func applyFilters(q *qGen.Query, p queryParams) {
	q.Where("date BETWEEN DATE(?) AND DATE(?, '+1 day')", p.Start, p.End) // I always constrain results by date

	if p.Machine != "all" && p.Set != 0 {
		q.Append("AND ball_machine = ? AND ball_set = ?", p.Machine, p.Set)
	} else if p.Machine != "all" && p.Set == 0 {
		q.Append("AND ball_machine = ?", p.Machine)
	} else if p.Machine == "all" && p.Set != 0 {
		q.Append("AND ball_set = ?", p.Set)
	}
}

func (db *AppDB) getResults(p queryParams) <-chan dbRow {
	c := make(chan dbRow)

	go func() {
		q := qGen.NewQuery().Select("results", "date", "ball_machine", "ball_set", "num_1", "num_2", "num_3", "num_4", "num_5", "num_6", "bonus")

		applyFilters(q, p)
		stmt, _ := db.Prepare(q.Order("DATE(date)").SQL)
		rows, err := stmt.Query(q.Args...)
		if err != nil {
			log.Println(err.Error())
		}
		defer rows.Close()

		for rows.Next() {
			var r dbRow
			r.Num = make([]int, balls)
			if err := rows.Scan(&r.Date, &r.Machine, &r.Set, &r.Num[0], &r.Num[1], &r.Num[2], &r.Num[3], &r.Num[4], &r.Num[5], &r.Num[6]); err != nil {
				log.Println(err.Error())
			}

			c <- r
		}

		close(c)
	}()

	return c
}

func (db *AppDB) getLastDraw() ([]int, error) {
	r, q := make([]int, balls), qGen.NewQuery().
		Select("results", "num_1", "num_2", "num_3", "num_4", "num_5", "num_6", "bonus").
		Order("DATE(date)").
		Append("DESC LIMIT 1")
	stmt, _ := db.Prepare(q.SQL)

	if err := stmt.QueryRow().Scan(&r[0], &r[1], &r[2], &r[3], &r[4], &r[5], &r[6]); err != nil {
		return r, err
	}

	return r, nil
}

func (db *AppDB) getMachineList(p queryParams) ([]string, error) {
	r, q := []string{}, qGen.NewQuery().Select("results", "DISTINCT(ball_machine)")

	p.Machine = "all" // Ensure queryParams are right for this
	applyFilters(q, p)
	stmt, _ := db.Prepare(q.Order("ball_machine").SQL)
	rows, err := stmt.Query(q.Args...)
	if err != nil {
		return r, err
	}
	defer rows.Close()

	for rows.Next() {
		var m string
		rows.Scan(&m)
		r = append(r, m)
	}

	return r, nil
}

func (db *AppDB) getSetList(p queryParams) ([]int, error) {
	r, q := []int{}, qGen.NewQuery().Select("results", "DISTINCT(ball_set)")

	p.Set = 0 // Ensure queryParams are right for this
	applyFilters(q, p)
	stmt, _ := db.Prepare(q.Order("ball_set").SQL)
	rows, err := stmt.Query(q.Args...)
	if err != nil {
		return r, err
	}
	defer rows.Close()

	for rows.Next() {
		var s int
		rows.Scan(&s)
		r = append(r, s)
	}

	return r, nil
}

func (db *AppDB) getRowCount() (int, error) {
	var c int
	if err := db.QueryRow(qGen.NewQuery().Select("results", "COUNT(*)").SQL).Scan(&c); err != nil {
		return c, err
	}

	return c, nil
}

func (db *AppDB) getDataRange() (time.Time, time.Time, error) {
	var first, last string
	q := qGen.NewQuery().Select("results", "MIN(date)", "MAX(date)")
	stmt, _ := db.Prepare(q.SQL)

	if err := stmt.QueryRow().Scan(&first, &last); err != nil {
		return time.Now(), time.Now(), err
	}

	f, _ := time.Parse(formatSqlite, first)
	l, _ := time.Parse(formatSqlite, last)
	return f, l, nil
}

func (db *AppDB) updateDB() error {
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare queries
	qSel := qGen.NewQuery().Select("results", "COUNT(*)").Where("date = ?")
	sel, err := tx.Prepare(qSel.SQL)
	if err != nil {
		return err
	}

	qIns := qGen.NewQuery().Insert("results", []string{"date", "ball_set", "ball_machine", "num_1", "num_2", "num_3", "num_4", "num_5", "num_6", "bonus"})
	ins, err := tx.Prepare(qIns.SQL)
	if err != nil {
		return err
	}

	// Iterate scrape data
	for d := range updateScraper() {
		// Check I don't already have this record
		var i int

		// Set args for queries
		qSel.Args = []interface{}{d.Date.Format(formatSqlite)}
		qIns.Args = []interface{}{d.Date, d.Set, d.Machine, d.Num[0], d.Num[1], d.Num[2], d.Num[3], d.Num[4], d.Num[5], d.Num[6]}

		if err := sel.QueryRow(qSel.Args...).Scan(&i); err != nil {
			tx.Rollback()
			return err
		}
		if i != 0 {
			log.Println(d.Date.String(), "already in DB, skipping")
			continue
		}

		// Insert new record
		if _, err := ins.Exec(qIns.Args...); err != nil {
			tx.Rollback()
			return err
		}
		log.Println(d, " inserted")
	}
	tx.Commit()

	return nil
}

func (db *AppDB) populateDB() error {
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare statement
	qIns := qGen.NewQuery().Insert("results", []string{"date", "ball_set", "ball_machine", "num_1", "num_2", "num_3", "num_4", "num_5", "num_6", "bonus"})
	insert, err := tx.Prepare(qIns.SQL)
	if err != nil {
		return err
	}

	// Iterate scrape data
	for d := range archiveScraper() {
		// Set args
		qIns.Args = []interface{}{d.Date, d.Set, d.Machine, d.Num[0], d.Num[1], d.Num[2], d.Num[3], d.Num[4], d.Num[5], d.Num[6]}

		// Exec
		if _, err := insert.Exec(qIns.Args...); err != nil {
			tx.Rollback()
			return err
		}
		log.Println(d)
	}
	tx.Commit()

	return nil
}
