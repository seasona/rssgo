package test

import (
	"testing"

	"github.com/seasona/rssgo/internal"
)

func TestUpdateFeeds(t *testing.T) {
	c := internal.Controller{}
	cfg := "../config.json"
	dbFile := "testdata/rssgo.db"
	c.Init(cfg, dbFile)

	c.UpdateFeeds()
}
