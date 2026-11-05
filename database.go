package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "lizibeth.db")
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password TEXT
	);
	CREATE TABLE IF NOT EXISTS collections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		description TEXT,
		image TEXT
	);
	CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		collection_id INTEGER,
		name TEXT,
		price TEXT,
		image TEXT,
		FOREIGN KEY(collection_id) REFERENCES collections(id)
	);
	CREATE TABLE IF NOT EXISTS bookings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		customer_name TEXT,
		service TEXT,
		date TEXT,
		status TEXT
	);
	`
	_, err = DB.Exec(sqlStmt)
	if err != nil {
		log.Fatal(fmt.Sprintf("%q: %s\n", err, sqlStmt))
	}
}
