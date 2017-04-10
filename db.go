package main

import (
	"database/sql"
	"log"
	"time"

	"fmt"
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
	Set           int
	Start, End    string
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
func applyFilters(q *qGen.Query, p queryParams) []interface{} {
	qp := []interface{}{p.Start, p.End} // I always constrain results by date

	f := qGen.NewFilterSet().Add("date:DATE", qGen.Between)
	if p.Machine != "all" && p.Set != 0 {
		q.Where(f.Add("ball_machine", qGen.Eq).Add("ball_set", qGen.Eq))
		qp = append(qp, p.Machine, p.Set)

	} else if p.Machine != "all" && p.Set == 0 {
		q.Where(f.Add("ball_machine", qGen.Eq))
		qp = append(qp, p.Machine)

	} else if p.Machine == "all" && p.Set != 0 {
		q.Where(f.Add("ball_set", qGen.Eq))
		qp = append(qp, p.Set)

	} else {
		q.Where(f)

	}

	return qp
}

func (db *AppDB) getResults(p queryParams) <-chan dbRow {
	c := make(chan dbRow)

	go func() {
		q := qGen.NewQuery().
			Select("date", "ball_machine", "ball_set", "num_1", "num_2", "num_3", "num_4", "num_5", "num_6", "bonus").
			From("results")

		qp := applyFilters(q, p)
		stmt, _ := db.Prepare(q.Order("DATE(date)").SQL)
		rows, err := stmt.Query(qp...)
		if err != nil {
			log.Println(err.Error())
		}

		for rows.Next() {
			var r dbRow
			r.Num = make([]int, 7)
			if err := rows.Scan(&r.Date, &r.Machine, &r.Set, &r.Num[0], &r.Num[1], &r.Num[2], &r.Num[3], &r.Num[4], &r.Num[5], &r.Num[6]); err != nil {
				log.Println(err.Error())
			}

			c <- r
		}

		close(c)
	}()

	return c
}

func (db *AppDB) getLastDraw() []int {
	var r []int
	q := qGen.NewQuery().
		Select("num_1", "num_2", "num_3", "num_4", "num_5", "num_6", "bonus").
		From("results").
		Order("DATE(date)").
		Append("DESC LIMIT 1")
	stmt, err := db.Prepare(q.SQL)
	if err != nil {
		log.Println(err)
	}

	err = stmt.QueryRow().Scan(&r[0], &r[1], &r[2], &r[3], &r[4], &r[5], &r[6])
	if err != nil {
		log.Println(err)
	}

	return r
}

func (db *AppDB) getMachineSetCombinations(p queryParams) map[string]int {
	r := make(map[string]int)

	for row := range db.getResults(p) {
		r[fmt.Sprintf("%s:%d", row.Machine, row.Set)]++
	}

	return r
}

func (db *AppDB) getResultsAverage(p queryParams) ([]int, error) {
	r := make([]int, 7)
	q := qGen.NewQuery()

	qSelect := []string{}
	for i := 1; i <= 6; i++ {
		qSelect = append(qSelect, fmt.Sprintf("SUM(num_%d)/COUNT(num_%d)", i, i))
	}
	qSelect = append(qSelect, "SUM(bonus)/COUNT(bonus)")

	q.Select(qSelect...).From("results")

	qp := applyFilters(q, p)
	stmt, _ := db.Prepare(q.SQL)
	if err := stmt.QueryRow(qp...).Scan(&r[0], &r[1], &r[2], &r[3], &r[4], &r[5], &r[6]); err != nil {
		return r, err
	}

	return r, nil
}

func (db *AppDB) getResultsAverageRanges(p queryParams) ([]string, error) {
	r := []string{}
	q := qGen.NewQuery()

	qSelect := []string{}
	for i := 1; i <= 6; i++ {
		qSelect = append(qSelect, fmt.Sprintf("MIN(num_%d)", i), fmt.Sprintf("MAX(num_%d)", i))
	}
	qSelect = append(qSelect, "MIN(bonus)", "MAX(bonus)")

	q.Select(qSelect...).From("results")
	qp := applyFilters(q, p)
	stmt, _ := db.Prepare(q.SQL)

	rI := make([]int, 14)
	if err := stmt.QueryRow(qp...).Scan(&rI[0], &rI[1], &rI[2], &rI[3], &rI[4], &rI[5], &rI[6], &rI[7], &rI[8], &rI[9], &rI[10], &rI[11], &rI[12], &rI[13]); err != nil {
		return r, err
	}

	for i := 1; i < len(rI); i++ {
		if i%2 != 0 {
			r = append(r, fmt.Sprintf("%d-%d", rI[i-1], rI[i]))
		}
	}

	return r, nil
}

func (db *AppDB) getMachineList(p queryParams) ([]string, error) {
	r := []string{}
	q := qGen.NewQuery().
		Select("DISTINCT(ball_machine)").
		From("results")

	p.Machine = "all" // Ensure queryParams are right for this
	qp := applyFilters(q, p)
	stmt, _ := db.Prepare(q.Order("ball_machine").SQL)
	rows, err := stmt.Query(qp...)
	if err != nil {
		return r, err
	}

	for rows.Next() {
		var m string
		rows.Scan(&m)
		r = append(r, m)
	}

	return r, nil
}

func (db *AppDB) getSetList(p queryParams) ([]int, error) {
	r := []int{}
	q := qGen.NewQuery().
		Select("DISTINCT(ball_set)").
		From("results")

	p.Set = 0 // Ensure queryParams are right for this
	qp := applyFilters(q, p)
	stmt, _ := db.Prepare(q.Order("ball_set").SQL)
	rows, err := stmt.Query(qp...)
	if err != nil {
		return r, err
	}

	for rows.Next() {
		var s int
		rows.Scan(&s)
		r = append(r, s)
	}

	return r, nil
}

func (db *AppDB) getRowCount() (int, error) {
	var rows int
	if err := db.QueryRow(qGen.NewQuery().Select("COUNT(*)").From("results").SQL).Scan(&rows); err != nil {
		return rows, err
	}

	return rows, nil
}

func (db *AppDB) getDataRange() (time.Time, time.Time, error) {
	var first, last string
	q := qGen.NewQuery().
		Select("MIN(date)", "MAX(date)").
		From("results")
	stmt, _ := db.Prepare(q.SQL)

	if err := stmt.QueryRow().Scan(&first, &last); err != nil {
		return time.Now(), time.Now(), err
	}

	f, _ := time.Parse(formatSqlite, first)
	l, _ := time.Parse(formatSqlite, last)
	return f, l, nil
}

func (db *AppDB) updateDB() error {
	// Prepare insert
	qInsert, err := db.Prepare(qGen.NewQuery().
		Insert("results", "date", "ball_set", "ball_machine", "num_1", "num_2", "num_3", "num_4", "num_5", "num_6", "bonus").SQL)
	if err != nil {
		return err
	}
	// Generate Select
	qSelect := qGen.NewQuery().
		Select("COUNT(*)").
		From("results").
		Where(qGen.NewFilterSet().Add("date", qGen.Eq)).SQL

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Iterate scrape data
	for d := range updateScraper() {
		// Check I don't already have this record
		var i int
		if err := tx.QueryRow(qSelect, d.Date.Format(formatSqlite)).Scan(&i); err != nil {
			tx.Rollback()
			return err
		}
		if i != 0 {
			log.Println(d.Date.String(), "already in DB, skipping")
			continue
		}
		// Insert new record
		if _, err := tx.Stmt(qInsert).Exec(d.Date, d.Set, d.Machine, d.Num[0], d.Num[1], d.Num[2], d.Num[3], d.Num[4], d.Num[5], d.Num[6]); err != nil {
			tx.Rollback()
			return err
		}
		log.Println(d, " inserted")
	}
	tx.Commit()

	return nil
}

func (db *AppDB) populateDB() error {
	// Prepare statement
	q, err := db.Prepare(qGen.NewQuery().
		Insert("results", "date", "ball_set", "ball_machine", "num_1", "num_2", "num_3", "num_4", "num_5", "num_6", "bonus").SQL)
	if err != nil {
		return err
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Iterate scrape data
	for d := range archiveScraper() {
		if _, err := tx.Stmt(q).Exec(d.Date, d.Set, d.Machine, d.Num[0], d.Num[1], d.Num[2], d.Num[3], d.Num[4], d.Num[5], d.Num[6]); err != nil {
			tx.Rollback()
			return err
		}
		log.Println(d)
	}
	tx.Commit()

	return nil
}
