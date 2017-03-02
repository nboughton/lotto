package main

import (
	"database/sql"
	//"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	//ql "github.com/nboughton/go-stupidqlite"
)

var (
	baseURL    = "https://www.lottery.co.uk/lotto/results/archive-1994"
	sqlPragmas = "PRAGMA journal_mode=WAL;	PRAGMA busy_timeout=5000"
	sqlSchema  = `
	CREATE TABLE IF NOT EXISTS results 
	(id INTEGER PRIMARY KEY AUTOINCREMENT, 
	date DATETIME, 
	num1 INT, 
	num2 INT, 
	num3 INT, 
	num4 INT, 
	num5 INT, 
	num6 INT, 
	num7 INT)`
)

// AppDB is a wrapper for *sql.DB so we can extend it by adding our own methods
type AppDB struct {
	*sql.DB
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

	return &AppDB{db}
}
