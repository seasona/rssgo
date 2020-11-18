package internal

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
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

func (d *DB) CreateTables(db *sql.DB, rss *RSS) {
	// the rss title may contain space 
	for _, tu := range rss.titleURLs {
		ct := fmt.Sprintf(`create table if not exists %v(
				id integer not null primary key,
				feed text,
				title text,
				content text,
				link text,
				read bool,
				display_name string,
				deleted bool,
				published DATETIME);`, tu.Title)

		_, err := d.db.Exec(ct)

		if err != nil {
			log.Panicf("Can't create database %v because: %v", tu.Title, err)
		}
	}
}