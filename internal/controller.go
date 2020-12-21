package internal

import (
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Controller struct {
	db         *DB
	rss        *RSS
	win        *Window
	conf       Config
	theme      Theme
	articles   map[string][]Article
	lastUpdate time.Time
}

func (c *Controller) Init(cfg, theme, dbFile string) {
	c.articles = make(map[string][]Article)

	c.conf.LoadConfig(cfg)
	c.theme.LoadTheme(theme)

	// rss depends on config's opml file
	c.rss = &RSS{}
	c.rss.Init(c)

	// db depends on rss feed, so must initialize behind it
	c.db = &DB{}
	c.db.Init(c, dbFile)

	c.win = &Window{}
	c.win.Init(c, c.InputFunc)
	c.win.RegisterSelectedFeedFunc(c.FeedSelectionChanged)

	c.UpdateLoop()

	c.win.Start()
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

			c.db.Save(a)
		}
	}
	c.lastUpdate = time.Now()
	// update articles in controller every updateFeeds
	c.GetAllArticlesFromDB()
}

func (c *Controller) GetAllArticlesFromDB() {
	c.articles = c.db.All()
}

func (c *Controller) Quit() {
	c.win.app.Stop()
	os.Exit(0)
}

func (c *Controller) showFeeds() {
	c.win.ClearFeeds()

	for feedTitle, articles := range c.articles {
		unread := 0
		total := len(articles)

		for _, article := range articles {
			if !article.read {
				unread++
			}
		}

		c.win.AddToFeeds(feedTitle, unread, total, feedTitle)
	}
}

func (c *Controller) showArticles(feedTitle string) {
	c.win.ClearArticles()
	for _, article := range c.articles[feedTitle] {
		c.win.AddToArticles(&article)
	}
	c.win.articles.ScrollToBeginning()
}

func (c *Controller) FeedSelectionChanged(row, col int) {
	if row <= 0 {
		return
	}

	r, _ := c.win.feeds.GetSelection()
	cell := c.win.feeds.GetCell(r, 2)
	ref := cell.GetReference()
	if ref != nil {
		c.showArticles(ref.(string))
	}
}

func (c *Controller) UpdateLoop() {
	c.GetAllArticlesFromDB()
	go c.UpdateFeeds()
	c.showFeeds()
	go func() {
		updateWin := time.NewTicker(10 * time.Second)
		select {
		case <-updateWin.C:
			c.showFeeds()
		}
	}()
}

func (c *Controller) InputFunc(event *tcell.EventKey) *tcell.EventKey {
	keyName := event.Name()
	// the keyName may be Rune[p], so need remove rune
	if strings.Contains(keyName, "Rune") {
		keyName = string(event.Rune())
	}

	switch keyName {
	case c.conf.KeySwitchWindows:
		c.win.SwitchFocus()
	case c.conf.KeyHelp:
		c.win.TriggerHelp()
	case c.conf.KeyPreview:
		c.win.TriggerPreview()
	case c.conf.KeyQuit:
		c.Quit()
	default:
		return event
	}

	return nil
}
