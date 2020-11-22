package internal

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"strings"
)

type Config struct {
	// only exported fields will be encoded/decoded in JSON,
	// fields must start with capital letters to be exported.
	OPML                      string              `json:"opml"`
	Feeds                     []map[string]string `json:"feeds"`
	Theme                     string              `json:"theme"`
	KeyMoveDown               string              `json:"keyMoveDown"`
	KeyMoveUp                 string              `json:"keyMoveUp"`
	KeySwitchWindows          string              `json:"keySwitchWindows"`
	KeyQuit                   string              `json:"keyQuit"`
	KeyHelp                   string              `json:"keyHelp"`
	KeyMarkArticle            string              `json:"keyMarkArticle"`
	DaysKeepDeletedArticle    int                 `json:"daysKeepDeletedArticle"`
	DaysKeepReadArticle       int                 `json:"daysKeepReadArticle"`
	SkipArticlesOlderThanDays int                 `json:"skipArticlesOlderThanDays"`
	SecondsBetweenUpdates     int                 `json:"secondsBetweenUpdates"`
}

func (c *Config) LoadConfig(file string) {
	cf, err := os.Open(file)
	defer cf.Close()

	if err != nil {
		log.Fatal("Can't load config file:", err)
	}

	jsonParser := json.NewDecoder(cf)
	err = jsonParser.Decode(c)

	if err != nil {
		log.Fatal("Can't decode json file:", err)
	}

	keys := make(map[string]struct{})
	val := reflect.Indirect(reflect.ValueOf(c))

	// use reflect to detect if the key setting repeat
	for i := 0; i < val.NumField(); i++ {
		if strings.HasPrefix(val.Type().Field(i).Name, "key") {
			if _, ok := keys[val.Field(i).String()]; ok {
				log.Fatal("The key defined more than once: key:", val.Field(i).String())
			} else {
				keys[val.Field(i).String()] = struct{}{}
			}
		}
	}
}
