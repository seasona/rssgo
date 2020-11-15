package internal

import (
	"log"
	"net/http"
	"sync"

	"github.com/gilliek/go-opml/opml"
	"github.com/mmcdole/gofeed"
)

type RSS struct {
	feeds []struct {
		title string
		feed  *gofeed.Feed
	}
	c *Controller
}

type Feed struct {
	Title string
	URL   string
}

func (r *RSS) Init(c *Controller) {
	r.c = c
	r.feeds = []struct {
		title string
		feed  *gofeed.Feed
	}{}

	if c.conf.OpmlFile != "" {
		op, err := opml.NewOPMLFromFile(r.c.conf.OpmlFile)
		if err != nil {
			log.Fatal("Can't load opml file: ", r.c.conf.OpmlFile)
		}

		for _, b := range op.Body.Outlines {
			if b.Outlines != nil {
				for _, ib := range b.Outlines {
					url := r.getURLFromOPML(ib)
					if url != "" {
						r.c.feeds = append(r.c.feeds, Feed{ib.Title, url})
					}
				}
			} else {
				url := r.getURLFromOPML(b)
				if url != "" {
					r.c.feeds = append(r.c.feeds, Feed{b.Title, url})
				}
			}
		}
	}
}

func (r *RSS) getURLFromOPML(b opml.Outline) string {
	str := ""
	if b.XMLURL != "" {
		str = b.XMLURL
	} else if b.HTMLURL != "" {
		str = b.HTMLURL
	} else if b.URL != "" {
		str = b.URL
	}
	return str
}

// FetchURL will send a request to url and use gofeed parse response's body
func (r *RSS) FetchURL(fp *gofeed.Parser, url string) (*gofeed.Feed, error) {
	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.183 Safari/537.36")
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp != nil {
		defer func() {
			// there should be a error handle after defer
			resp.Body.Close()
		}()
	}

	return fp.Parse(resp.Body)
}

// Update will
func (r *RSS) Update() {
	var m sync.Mutex
	var wg sync.WaitGroup

	for _, f := range r.c.feeds {
		wg.Add(1)
		go func(f Feed) {
			fp := gofeed.NewParser()
			feed, err := r.FetchURL(fp, f.URL)
			if err != nil {
				log.Printf("Error occur when fetch url: %s, err: %v", f.URL, err)
			} else {
				m.Lock()
				r.feeds = append(r.feeds, struct {
					title string
					feed  *gofeed.Feed
				}{f.Title, feed})
				m.Unlock()
			}

			wg.Done()
		}(f)
	}
	wg.Wait()
}
