package test

import (
	"fmt"
	"testing"

	"github.com/mmcdole/gofeed"

	"github.com/seasona/rssgo/internal"
)

func TestFetchRSS(t *testing.T) {
	url := "http://feeds.twit.tv/twit.xml"

	fp := gofeed.NewParser()

	rss := internal.RSS{}
	feed, err := rss.FetchURL(fp, url)

	if err != nil {
		t.Errorf("Can't fetch from RSS")
	}

	for _, item := range feed.Items {
		if item == nil {
			continue
		}

		fmt.Println("Title: ", item.Title)
		fmt.Println("Date: ", item.Published)
	}
}

