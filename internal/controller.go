package internal

import "time"

type Controller struct {
	db         *DB
	rss        *RSS
	conf       Config
	theme      Theme
	articles   map[string][]Article
	lastUpdate time.Time
}

func (c *Controller) Init(cfg, theme, dbFile string) {
	c.articles = make(map[string][]Article)

	c.conf.LoadConfig(cfg)

	// rss depends on config's opml file
	c.rss = &RSS{}
	c.rss.Init(c)

	// db depends on rss feed, so must initialize behind it
	c.db = &DB{}
	c.db.Init(c, dbFile)
}

// UpdateFeeds will update all rss and save the result to database
func (c *Controller) UpdateFeeds() {
	c.rss.Update()
	for _, f := range c.rss.titleFeeds {
		if f.feed == nil {
			continue
		}

		for _, item := range f.feed.Items {
			if item == nil {
				continue
			}

			// the published should be last update time rather than publish time
			var published time.Time
			if item.UpdatedParsed != nil {
				published = *item.PublishedParsed
			} else if item.PublishedParsed != nil {
				published = *item.PublishedParsed
			} else {
				published = time.Now()
			}

			if int(time.Now().Sub(published).Hours()/24) > c.conf.SkipArticlesOlderThanDays {
				continue
			}

			content := item.Description
			if content == "" {
				content = item.Content
			}

			// the feed attribute which gofeed get may be different from it's title,
			// for example, The Verge's feed attribute is The Verge - All Feed,
			// so use f.title as article's feed
			a := Article{
				c:         c,
				feed:      f.title,
				title:     item.Title,
				content:   content,
				link:      item.Link,
				read:      false,
				deleted:   false,
				published: published,
			}

			exist := false
			for _, e := range c.articles[f.title] {
				if e.title == a.title {
					exist = true
					break
				}
			}

			if !exist {
				c.db.Save(a)
			}
		}
	}
}

func (c *Controller) GetAllArticleFromDB() {
	c.articles = c.db.All()
}
