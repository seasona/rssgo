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
	KeyMoveDown               string              `json:"keyMoveDown"`
	KeyMoveUp                 string              `json:"keyMoveUp"`
	KeySwitchWindows          string              `json:"keySwitchWindows"`
	KeyQuit                   string              `json:"keyQuit"`
	KeyHelp                   string              `json:"keyHelp"`
	KeyPreview                string              `json:"keyPreview"`
	KeyMarkArticle            string              `json:"keyMarkArticle"`
	KeyOpenLink               string              `json:"keyOpenLink"`
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

func (c *Config) GetConfigKeys() map[string]string {
	keys := make(map[string]string)
	
	keys["Open Link"] = c.KeyOpenLink
	keys["Mark Article"] = c.KeyMarkArticle
	//keys["Delete"] = c.KeyDeleteArticle
	keys["Up"] = c.KeyMoveUp
	keys["Down"] = c.KeyMoveDown
	keys["Toggle Preview"] = c.KeyPreview
	//keys["Mark All Read"] = c.KeyMarkAllRead
	//keys["Mark All UnRead"] = c.KeyMarkAllUnread
	keys["Toggle Help"] = c.KeyHelp
	//keys["Update Feeds"] = c.KeyUpdateFeeds
	keys["Switch Windows"] = c.KeySwitchWindows
	keys["Quit"] = c.KeyQuit

	return keys
}
