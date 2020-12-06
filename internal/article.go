package internal

import "time"

// Article is for
type Article struct {
	c         *Controller
	id        int
	feed      string // feed name
	title     string // article title
	content   string
	link      string
	read      bool
	deleted   bool
	published time.Time
}
