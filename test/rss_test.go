package test

import (
	"testing"

	"github.com/mmcdole/gofeed"

	"github.com/seasona/rssgo/internal"
)

func TestFetchRSS(t *testing.T) {
	url := "http://feeds.ign.com/ign/games-all"

	fp := gofeed.NewParser()

	rss := internal.RSS{}
	feed, err := rss.FetchURL(fp, url)

	if err != nil {
		t.Errorf("Can't fetch from RSS")
		panic(err)
	}

	for _, item := range feed.Items {
		if item == nil {
			continue
		}

		//fmt.Println("Title: ", item.Title)
		//fmt.Println("Date: ", item.Published)
	}
}

// there may be two error: err: Failed to detect feed type, which means the RSS url has problem;
// connectex: A connection attempt failed, which is network problem, may caused by GFW in china
func TestUpdate(t *testing.T) {
	url := "testdata/feedly.opml"

	rss := internal.RSS{}

	rss.GetTitleURLFromOPML(url)

	rss.Update()
}
