// Package models interfaces with the SQLite DB
package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

// DB is our connection so we dont' have to keep opening / closing
var DB *sql.DB

// Open opens the connection and stashes it in DB
func Open(file string) error {
	var err error
	DB, err = sql.Open("sqlite3", file)
	return err
}

// Close just wraps the regular database Close function
func Close() {
	DB.Close()
}
