package test

import (
	"log"
	"os"
	"testing"

	"github.com/seasona/rssgo/internal"
)

func TestWindowInit(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	if flog, err := os.Create("/home/jijie/rssgo/test/testdata/test.log"); err != nil {
		log.Println("Can't create log file")
	} else {
		log.SetOutput(flog)
	}
	c := internal.Controller{}
	cfg := "../config.json"
	dbFile := "testdata/rssgo.db"
	theme := "../theme/default.json"
	c.Init(cfg, theme, dbFile)
}
