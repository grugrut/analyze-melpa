package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	db, err := sql.Open("sqlite3", "./asset/melpa.db")
	if err != nil {
		log.Fatal(err)
	}

	q := `
CREATE TABLE packages (
name TEXT PRIMARY KEY,
desc TEXT,
type TEXT,
url TEXT,
maintainer TEXT,
repo TEXT,
fetcher TEXT,
dl INTEGER
	)
`
	execQuery(db, q)

	q = `
CREATE TABLE dependencies (
name TEXT,
depends TEXT,
PRIMARY KEY(name, depends)
)
`
	execQuery(db, q)

	q = `
CREATE TABLE keywords (
name TEXT,
keyword TEXT,
PRIMARY KEY(name, keyword)
)
`
	execQuery(db, q)

	q = `
CREATE TABLE authors (
name TEXT,
author TEXT,
PRIMARY KEY(name, author)
)
`
	execQuery(db, q)
}

func execQuery(db *sql.DB, q string) {
	if _, err := db.Exec(q); err != nil {
		log.Fatal(err)
	}
}
