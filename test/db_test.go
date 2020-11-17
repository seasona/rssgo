package test

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/seasona/rssgo/internal"
)

func TestDBInit(t *testing.T) {
	url := "testdata/feedly.opml"
	rss := internal.RSS{}
	rss.GetTitleURLFromOPML(url)

	c := internal.Controller{}
	dbFile := "testdata/rssgo.db"

	sqldb, err := sql.Open("sqlite3", dbFile)

	if err != nil {
		log.Fatalf("Can't open database %v, %v", dbFile, err)
	}

	db := internal.DB{}
	db.Init(&c, dbFile)
	db.CreateTables(sqldb, &rss)
}
