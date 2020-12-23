package internal

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"sort"
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
	c.db.CleanUp()

	c.win = &Window{}
	c.win.Init(c, c.InputFunc)
	//todo the first selected feed can't show articels
	c.win.FeedSelectedFunc(c.SelectFeed)
	c.win.ArticleSelectedFunc(c.SelectArticle)

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
	c.showFeeds()
}

func (c *Controller) GetAllArticlesFromDB() {
	c.articles = c.db.All()
}

func (c *Controller) Quit() {
	c.win.app.Stop()
	os.Exit(0)
}

func (c *Controller) GetSelectedArticle() *Article {
	if c.win.app.GetFocus() != c.win.articles {
		return nil
	}

	r, _ := c.win.articles.GetSelection()
	cell := c.win.articles.GetCell(r, 2)

	ref := cell.GetReference()
	if ref != nil {
		return ref.(*Article)
	}
	return nil
}

func (c *Controller) OpenLink(link string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", link).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", link).Start()
	}
	if err != nil {
		log.Println(err)
	}
}

func (c *Controller) showFeeds() {
	c.win.ClearFeeds()

	var feeds []string

	for feed := range c.articles {
		feeds = append(feeds, feed)
	}

	// keep ordinal of Feeds
	sort.Strings(feeds)

	for _, feedTitle := range feeds {
		unread := 0
		total := len(c.articles[feedTitle])

		for _, article := range c.articles[feedTitle] {
			if !article.read {
				unread++
			}
		}

		c.win.AddToFeeds(feedTitle, unread, total, feedTitle)
	}
}

func (c *Controller) showArticles(feedTitle string) {
	c.win.ClearArticles()
	// be carefull of the trap of range of golang, range will pass a value not a reference,
	// so the pointer to range loop is always same
	for i := 0; i < len(c.articles[feedTitle]); i++ {
		c.win.AddToArticles(&c.articles[feedTitle][i])
	}
	c.win.articles.ScrollToBeginning()
}

func (c *Controller) SelectFeed(row, col int) {
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

func (c *Controller) SelectArticle(row, col int) {
	if row < 0 {
		return
	}

	c.win.preview.Clear()

	ref := c.win.c.GetArticleForSelection()
	if ref != nil {
		c.win.AddToPreview(ref)
	}
}

func (c *Controller) MarkArticle() {
	a := c.GetSelectedArticle()
	if a == nil {
		return
	}
	r, _ := c.win.articles.GetSelection()
	cell := c.win.articles.GetCell(r, 0)
	if a.read {
		c.db.MarkUnread(a)
		cell.SetText("")
	} else {
		c.db.MarkRead(a)
		cell.SetText(c.theme.ReadMarker)
	}
	a.read = !a.read
	// update unread in Feeds
	c.showFeeds()
}

func (c *Controller) UpdateLoop() {
	c.GetAllArticlesFromDB()
	go c.UpdateFeeds()
	c.showFeeds()
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
	case c.conf.KeyMoveUp:
		c.win.MoveUp()
	case c.conf.KeyMoveDown:
		c.win.MoveDown()
	case c.conf.KeyOpenLink:
		a := c.GetSelectedArticle()
		if a != nil {
			c.OpenLink(a.link)
		}
	case c.conf.KeyMarkArticle:
		c.MarkArticle()
	default:
		return event
	}

	return nil
}
