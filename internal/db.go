package internal

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type DB struct {
	db       *sql.DB
	tableMap map[string]string
	c        *Controller
}

func (d *DB) Init(c *Controller, dbFile string) {
	d.c = c
	d.tableMap = make(map[string]string)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Can't open database %v, %v", dbFile, err)
	}
	d.db = db

	d.CreateTables(db, c.rss)
}

// the rss title may contain space or symbol can't be name of database table,
// so omit the character except alphabet and digit
func (d *DB) handleRSSTitle(title string) string {
	var tname string
	for _, c := range title {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			tname += string(c)
		}
	}
	d.tableMap[title] = tname

	return tname
}

func (d *DB) CreateTables(db *sql.DB, rss *RSS) {
	for _, tu := range rss.titleURLs {
		tname := d.handleRSSTitle(tu.Title)

		ct := fmt.Sprintf(`create table if not exists %v(
				id integer not null primary key,
				feed text,
				title text,
				content text,
				link text,
				read bool,
				display_name string,
				deleted bool,
				published DATETIME);`, tname)

		_, err := d.db.Exec(ct)

		if err != nil {
			log.Panicf("Can't create database %v because: %v", tname, err)
		}
	}
}
