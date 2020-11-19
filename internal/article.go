package internal

import "time"

// Article is for
type Article struct {
	c           *Controller
	id          int
	feed        string
	feedDisplay string
	title       string
	content     string
	link        string
	read        bool
	deleted     bool
	published   time.Time
}
