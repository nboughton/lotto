package main

import (
	"database/sql"
	//"fmt"
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
)

// AppDB is a wrapper for *sql.DB so we can extend it by adding our own methods
type AppDB struct {
	*sql.DB
}

type dbRow struct {
	date        time.Time
	ballMachine string
	ballSet     int
	num         []int
}

func connectDB(path string) *AppDB {
	// Connect to the database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}

	// Disable connection pooling
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(0)

	// Create DB schema if it doesn't exist
	if _, err := db.Exec(sqlSchema); err != nil {
		log.Fatal(err)
	}

	// Set PRAGMAs
	if _, err := db.Exec(sqlPragmas); err != nil {
		log.Fatal(err)
	}

	// Check if we have any data yet
	var rows int
	row := db.QueryRow(ql.NewQuery().Select("COUNT(*)").From("results").SQL)
	row.Scan(&rows)

	if rows == 0 {
		log.Println("No data, scraping site")
		if err := populateDB(db); err != nil {
			log.Fatal(err.Error())
		}
	}

	return &AppDB{db}
}

func populateDB(db *sql.DB) error {
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare statement
	q, err := tx.Prepare(ql.NewQuery().Insert("results", "date", "ball_set", "ball_machine", "num_1", "num_2", "num_3", "num_4", "num_5", "num_6", "bonus").SQL)
	if err != nil {
		return err
	}

	// Iterate scrape data
	for d := range scraper() {
		if _, err := q.Exec(d.date, d.ballSet, d.ballMachine, d.num[0], d.num[1], d.num[2], d.num[3], d.num[4], d.num[5], d.num[6]); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	// Finalise commit
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
