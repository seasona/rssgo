package internal

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"strings"
)

type Config struct {
	highLights                []string          `json:"highLights"`
	feeds                     map[string]string `json:"feeds"`
	keyMoveDown               string            `json:"keyMoveDown"`
	keyMoveUp                 string            `json:"keyMoveUp"`
	keySwitchWindows          string            `json:"keySwitchWindows"`
	keyQuit                   string            `json:"keyQuit"`
	keyHelp                   string            `json:"keyHelp"`
	keyMarkArticle            string            `json:"keyMarkArticle"`
	daysKeepDeletedArticle    int               `json:"daysKeepDeletedArticle"`
	daysKeepReadArticle       int               `json:"daysKeepReadArticle"`
	skipArticlesOlderThanDays int               `json:"skipArticlesOlderThanDays"`
	secondsBetweenUpdates     int               `json:"secondsBetweenUpdates"`
}

func (c *Config) LoadConfig(file string) Config {
	var conf Config
	cf, err := os.Open(file)
	defer cf.Close()

	if err != nil {
		log.Fatal("Can't load config file:", err)
	}

	jsonParser := json.NewDecoder(cf)
	err = jsonParser.Decode(&conf)

	if err != nil {
		log.Fatal("Can't decode json file:", err)
	}

	keys := make(map[string]struct{})
	val := reflect.Indirect(reflect.ValueOf(conf))

	// use reflect to detect if the key setting repeat
	for i := 0; i < val.NumField(); i++ {
		if strings.HasPrefix(val.Type().Name(), "key") {
			if _, ok := keys[val.Field(i).String()]; ok {
				log.Fatal("The key defined more than once: key:", val.Field(i).String())
			} else {
				keys[val.Field(i).String()] = struct{}{}
			}
		}
	}

	return conf
}
