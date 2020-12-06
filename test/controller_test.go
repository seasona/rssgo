package test

import (
	"testing"

	"github.com/seasona/rssgo/internal"
)

func TestUpdateFeeds(t *testing.T) {
	c := internal.Controller{}
	cfg := "../config.json"
	dbFile := "testdata/rssgo.db"
	theme := "../theme/default.json"
	c.Init(cfg, theme, dbFile)
	c.UpdateFeeds()
}

func TestGetAll(t *testing.T) {
	c := internal.Controller{}
	cfg := "../config.json"
	dbFile := "./testdata/rssgo.db"
	theme := "../theme/default.json"
	c.Init(cfg, theme, dbFile)
	c.GetAllArticleFromDB()
}
