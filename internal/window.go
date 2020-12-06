package internal

import "github.com/rivo/tview"

type Window struct {
	c        *Controller
	feeds    *tview.Table
	articles *tview.Table
	status   *tview.Table
	help     *tview.Table
	preview  *tview.TextView
}
