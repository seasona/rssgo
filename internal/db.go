package internal

import (
	"database/sql"
	"log"

	"github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
	c  *Controller
}

func (d *DB) Init(c *Controller, dbFile string) {
	d.c = c

	db, err := sql.Open("sqlite3", dbFile)

	if err != nil {
		log.Fatalf("Can't open database %v, %v", dbFile, err)
	}

	d.db = db
}
