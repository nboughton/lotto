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
	bset INT,
	bmachine TEXT,
	num_1 INT, 
	num_2 INT, 
	num_3 INT, 
	num_4 INT, 
	num_5 INT, 
	num_6 INT, 
	bonus INT)`
	sqliteDateFormat = "2006-01-02 15:04:05-07:00"
)

// AppDB is a wrapper for *sql.DB so we can extend it by adding our own methods
type AppDB struct {
	*sql.DB
}

type dbRow struct {
	date    time.Time
	machine string
	set     int
	num     []int
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
	if rows, _ := aDB.getRowCount(); rows == 0 {
		log.Println("No data, scraping site")
		if err := aDB.populateDB(); err != nil {
			log.Fatal(err.Error())
		}
	}

	return aDB
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

	f, _ := time.Parse(sqliteDateFormat, first)
	l, _ := time.Parse(sqliteDateFormat, last)
	return f, l, nil
}

func (db *AppDB) getMachineList() ([]string, error) {
	var (
		result []string
		q      = ql.NewQuery().
			Select("DISTINCT(bmachine)").
			From("results").
			Order("bmachine")
	)

	rows, err := db.Query(q.SQL)
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

func (db *AppDB) getSetList() ([]int, error) {
	var (
		result []int
		q      = ql.NewQuery().
			Select("DISTINCT(bset)").
			From("results").
			Order("bset")
	)

	rows, err := db.Query(q.SQL)
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

func (db *AppDB) populateDB() error {
	// Prepare statement
	q, err := db.Prepare(ql.NewQuery().
		Insert("results", "date", "bset", "bmachine", "num_1", "num_2", "num_3", "num_4", "num_5", "num_6", "bonus").SQL)
	if err != nil {
		return err
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Iterate scrape data
	for d := range scraper() {
		if _, err := tx.Stmt(q).Exec(d.date, d.set, d.machine, d.num[0], d.num[1], d.num[2], d.num[3], d.num[4], d.num[5], d.num[6]); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	return nil
}
