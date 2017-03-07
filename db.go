package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	ql "github.com/nboughton/go-stupidqlite"
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
	formatSqlite = "2006-01-02 15:04:05-07:00"
	formatYYYYMMDD = "2006-01-02"
)

// AppDB is a wrapper for *sql.DB so we can extend it by adding our own methods
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

	// Check if we have any data yet, otherwise populate the db
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

func (db *AppDB) getResults(p queryParams) <-chan dbRow {
	var c = make(chan dbRow)

	go func() {
		var (
			rows *sql.Rows
			err  error
			q    = ql.NewQuery().
				Select("date", "ball_machine", "ball_set", "num_1", "num_2", "num_3", "num_4", "num_5", "num_6", "bonus").
				From("results")
		)

		f := ql.NewFilterSet().Add(ql.Between, "date:DATE")
		if p.Machine != "all" && p.Set != 0 {
			qu := q.Where(f.Add(ql.Eq, "ball_machine").Add(ql.Eq, "ball_set")).Order("DATE(date)").SQL

			rows, err = db.Query(qu, p.Start, p.End, p.Machine, p.Set)
		} else if p.Machine != "all" && p.Set == 0 {
			qu := q.Where(f.Add(ql.Eq, "ball_machine")).Order("DATE(date)").SQL

			rows, err = db.Query(qu, p.Start, p.End, p.Machine)
		} else if p.Machine == "all" && p.Set != 0 {
			qu := q.Where(f.Add(ql.Eq, "ball_set")).Order("DATE(date)").SQL

			rows, err = db.Query(qu, p.Start, p.End, p.Set)
		} else {
			qu := q.Where(f).Order("DATE(date)").SQL

			rows, err = db.Query(qu, p.Start, p.End)
		}

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

func (db *AppDB) getResultsAverage(p queryParams) ([]int, error) {
	var (
		r = make([]int, 7)
		q = ql.NewQuery().
			Select("SUM(num_1)/COUNT(num_1)",
				"SUM(num_2)/COUNT(num_2)",
				"SUM(num_3)/COUNT(num_3)",
				"SUM(num_4)/COUNT(num_4)",
				"SUM(num_5)/COUNT(num_5)",
				"SUM(num_6)/COUNT(num_6)",
				"SUM(bonus)/COUNT(bonus)").
			From("results")
	)

	c := ql.NewFilterSet().Add(ql.Between, "date:DATE")
	if p.Machine != "all" && p.Set != 0 {
		c.Add(ql.Eq, "ball_machine").Add(ql.Eq, "ball_set")
		qu := q.Where(c).SQL

		if err := db.QueryRow(qu, p.Start, p.End, p.Machine, p.Set).Scan(&r[0], &r[1], &r[2], &r[3], &r[4], &r[5], &r[6]); err != nil {
			return r, err
		}
	} else if p.Machine != "all" && p.Set == 0 {
		c.Add(ql.Eq, "ball_machine")
		qu := q.Where(c).SQL

		if err := db.QueryRow(qu, p.Start, p.End, p.Machine).Scan(&r[0], &r[1], &r[2], &r[3], &r[4], &r[5], &r[6]); err != nil {
			return r, err
		}
	} else if p.Machine == "all" && p.Set != 0 {
		c.Add(ql.Eq, "ball_set")
		qu := q.Where(c).SQL

		if err := db.QueryRow(qu, p.Start, p.End, p.Set).Scan(&r[0], &r[1], &r[2], &r[3], &r[4], &r[5], &r[6]); err != nil {
			return r, err
		}
	} else {
		qu := q.Where(c).SQL

		if err := db.QueryRow(qu, p.Start, p.End).Scan(&r[0], &r[1], &r[2], &r[3], &r[4], &r[5], &r[6]); err != nil {
			return r, err
		}
	}

	return r, nil
}

func (db *AppDB) getRowCount() (int, error) {
	var rows int
	if err := db.QueryRow(ql.NewQuery().Select("COUNT(*)").From("results").SQL).Scan(&rows); err != nil {
		return rows, err
	}

	return rows, nil
}

func (db *AppDB) getDataRange() (time.Time, time.Time, error) {
	var (
		first string
		last  string
		q     = ql.NewQuery().
			Select("MIN(date)", "MAX(date)").
			From("results")
	)

	if err := db.QueryRow(q.SQL).Scan(&first, &last); err != nil {
		return time.Now(), time.Now(), err
	}

	f, _ := time.Parse(formatSqlite, first)
	l, _ := time.Parse(formatSqlite, last)
	return f, l, nil
}

func (db *AppDB) getMachineList(bySet int) ([]string, error) {
	var (
		result []string
		rows   *sql.Rows
		err    error
		q      = ql.NewQuery().
			Select("DISTINCT(ball_machine)").
			From("results")
	)

	if bySet != 0 {
		f := ql.NewFilterSet().Add(ql.Eq, "ball_set")
		q.Where(f).Order("ball_machine")
		rows, err = db.Query(q.SQL, bySet)
	} else {
		rows, err = db.Query(q.Order("ball_machine").SQL)
	}
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var m string
		rows.Scan(&m)
		result = append(result, m)
	}

	return result, nil
}

func (db *AppDB) getSetList(byMachine string) ([]int, error) {
	var (
		result []int
		rows   *sql.Rows
		err    error
		q      = ql.NewQuery().
			Select("DISTINCT(ball_set)").
			From("results")
	)

	if byMachine != "all" {
		f := ql.NewFilterSet().Add(ql.Eq, "ball_machine")
		q.Where(f).Order("ball_set")
		rows, err = db.Query(q.SQL, byMachine)
	} else {
		rows, err = db.Query(q.Order("ball_set").SQL)
	}
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var s int
		rows.Scan(&s)
		result = append(result, s)
	}

	return result, nil
}

func (db *AppDB) updateDB() error {
	// Prepare insert
	qInsert, err := db.Prepare(ql.NewQuery().
		Insert("results", "date", "ball_set", "ball_machine", "num_1", "num_2", "num_3", "num_4", "num_5", "num_6", "bonus").SQL)
	if err != nil {
		return err
	}
	// Generate Select
	qSelect := ql.NewQuery().
		Select("COUNT(*)").
		From("results").
		Where(ql.NewFilterSet().Add(ql.Eq, "date")).SQL

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Iterate scrape data
	for d := range updateScraper() {
		// Check we don't already have this record
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
	q, err := db.Prepare(ql.NewQuery().
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
