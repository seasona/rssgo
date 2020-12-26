package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/OpenPeeDeeP/xdg"

	"github.com/seasona/rssgo/internal"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	defaultConf := "config.json"
	defaultTheme := "theme/default.json"
	defaultDB := "rssgo.db"
	defaultLog := "rssgo.log"

	config := flag.String("config", defaultConf, "Configuration file")
	theme := flag.String("theme", defaultTheme, "Theme file")
	db := flag.String("db", defaultDB, "Database file")
	logFile := flag.String("log", defaultLog, "Log file")
	version := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *version {
		fmt.Printf("rssgo version: %s", "v0.1.0")
		os.Exit(0)
	}

	// mkdir rssgo in xdg directory
	xdgDirs := xdg.New("", "rssgo")
	configHome := xdgDirs.ConfigHome()
	dataHome := xdgDirs.DataHome()

	for _, path := range []string{configHome, dataHome} {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.MkdirAll(path, 0700); err != nil {
				log.Println("Failed to create dir:", path)
			}
		}
	}

	if *config == defaultConf {
		path := xdgDirs.QueryConfig(*config)
		if path != "" {
			*config = path
		}
	}

	if *theme == defaultTheme {
		path := xdgDirs.QueryConfig(*theme)
		if path != "" {
			*theme = path
		}
	}

	if *db == defaultDB {
		path := xdgDirs.QueryData(*db)
		if path != "" {
			*db = path
		} else {
			*db = xdgDirs.DataHome() + string(os.PathSeparator) + defaultDB
		}
	}

	if *logFile == defaultLog {
		path := xdgDirs.QueryData(*logFile)
		if path != "" {
			*logFile = path
		} else {
			*logFile = xdgDirs.DataHome() + string(os.PathSeparator) + defaultLog
		}
	}

	log.Println("Using config:", *config)
	log.Println("Using theme:", *theme)
	log.Println("Using DB:", *db)
	log.Println("Using log file:", *logFile)

	if flog, err := os.Create(*logFile); err != nil {
		log.Println("Can't create log file:", *logFile)
	} else {
		log.SetOutput(flog)
	}

	co := internal.Controller{}
	co.Init(*config, *theme, *db)
}
