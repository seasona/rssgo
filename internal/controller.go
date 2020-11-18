package internal

type Controller struct {
	db    *DB
	rss   *RSS
	conf  Config
	theme Theme
}

func (c *Controller) Init(cfg, theme, dbFile string) {
	c.conf.LoadConfig(cfg)

	// rss depends on config's opml file
	c.rss = &RSS{}
	c.rss.Init(c)

	// db depends on rss feed, so must initialize behind it
	c.db = &DB{}
	c.db.Init(c, dbFile)
}
