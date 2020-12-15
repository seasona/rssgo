package test

import (
	"github.com/seasona/rssgo/internal"
	"testing"
)

func TestWindowInit(t *testing.T) {
	c := internal.Controller{}
	cfg := "../config.json"
	dbFile := "testdata/rssgo.db"
	theme := "../theme/default.json"
	c.Init(cfg, theme, dbFile)

	win := internal.Window{}
	win.Init(&c)
	win.Start()
}
