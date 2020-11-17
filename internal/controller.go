package internal

type Controller struct {
	db    *DB
	rss   *RSS
	conf  Config
	theme Theme
}

func (c *Controller) Init() {

}
