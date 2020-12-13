package internal

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"sort"
	"time"
)

type Window struct {
	c           *Controller
	feeds       *tview.Table
	articles    *tview.Table
	status      *tview.Table
	help        *tview.Table
	preview     *tview.TextView
	app         *tview.Application
	showPreview bool
	showHelp    bool
}

func (w *Window) Init(c *Controller, inputFunc func(event *tcell.EventKey) *tcell.EventKey) {
	w.c = c

	w.showPreview = true
	w.showHelp = false

	// Feeds window
	w.feeds = tview.NewTable()
	w.feeds.SetBorder(true)
	w.feeds.SetBorderPadding(1, 1, 1, 1)
	w.feeds.SetBorderColor(tcell.GetColor(w.c.theme.FeedBorder))
	w.feeds.SetTitle(fmt.Sprintf("%s Feeds", w.c.theme.FeedIcon)).SetTitleColor(tcell.GetColor(w.c.theme.FeedBorderTitle))

	// Articles window
	w.articles = tview.NewTable()
	w.articles.SetTitleAlign(tview.AlignLeft)
	w.articles.SetBorder(true)
	w.articles.SetBorderPadding(1, 1, 1, 1)
	w.articles.SetBorderColor(tcell.GetColor(w.c.theme.ArticleBorder))
	w.articles.SetTitle(fmt.Sprintf("%s Articles", w.c.theme.ArticleIcon)).SetTitleColor(tcell.GetColor(w.c.theme.ArticleBorderTitle))

	w.InitHelpWindow()

	// Preview window
	w.preview = tview.NewTextView()
	w.preview.SetBorder(true)
	w.preview.SetBorderPadding(1, 1, 1, 1)
	w.preview.SetTitleAlign(tview.AlignLeft)
	w.preview.SetBorderColor(tcell.GetColor(w.c.theme.PreviewBorder))
	w.preview.SetScrollable(true)
	w.preview.SetWordWrap(true)
	w.preview.SetDynamicColors(true)
	w.preview.SetTitle(fmt.Sprintf("%s Preview", w.c.theme.PreviewIcon)).SetTitleColor(tcell.GetColor(w.c.theme.PreviewBorderTitle))

	w.InitStatusWindow()

	w.app = tview.NewApplication()
	w.app.SetInputCapture(inputFunc)
}

func (w *Window) InitHelpWindow() {
	// Help window
	w.help = tview.NewTable()
	w.help.SetTitleAlign(tview.AlignLeft)
	w.help.SetBorder(true)
	w.help.SetBorderPadding(1, 1, 1, 1)
	w.help.SetBorderColor(tcell.GetColor(w.c.theme.ArticleBorder))
	w.help.SetTitle("💡 Help").SetTitleColor(tcell.GetColor(w.c.theme.ArticleBorderTitle))

	ts := tview.NewTableCell("Key")
	ts.SetAlign(tview.AlignLeft)
	ts.Attributes |= tcell.AttrBold
	ts.SetSelectable(false)
	w.help.SetCell(0, 0, ts)

	ts = tview.NewTableCell("Action")
	ts.SetAlign(tview.AlignLeft)
	ts.Attributes |= tcell.AttrBold
	ts.SetSelectable(false)
	w.help.SetCell(0, 1, ts)

	i := 1
	configKeys := w.c.conf.GetConfigKeys()
	var keys []string
	for _, k := range configKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		i++
		ts = tview.NewTableCell(fmt.Sprintf("%s", configKeys[k]))
		ts.SetAlign(tview.AlignLeft)
		ts.Attributes |= tcell.AttrBold
		ts.SetSelectable(false)
		ts.SetTextColor(tcell.GetColor(w.c.theme.StatusKey))
		w.help.SetCell(i, 0, ts)

		ts = tview.NewTableCell(fmt.Sprintf("%s", k))
		ts.SetAlign(tview.AlignLeft)
		ts.SetTextColor(tcell.GetColor(w.c.theme.StatusText))
		ts.Attributes |= tcell.AttrBold
		ts.SetSelectable(false)
		w.help.SetCell(i, 1, ts)
	}
}

func (w *Window) InitStatusWindow() {
	w.status = tview.NewTable()
	w.status.SetBackgroundColor(tcell.GetColor(w.c.theme.StatusBackground))
	w.status.SetFixed(1, 6)

	for i := 0; i < 7; i++ {
		ts := tview.NewTableCell("")
		ts.SetAlign(tview.AlignLeft)
		ts.Attributes |= tcell.AttrBold
		ts.SetSelectable(false)
		w.status.SetCell(0, i, ts)
	}

	w.UpdateStatusTicker()
}

func (w *Window) UpdateStatus() {
	// time now
	c := w.status.GetCell(0, 0)
	// tview support Color tags
	c.SetText(fmt.Sprintf("[%s][[%s]Time: [%s]%s[%s]]",
		w.c.theme.StatusBrackets,
		w.c.theme.StatusKey,
		w.c.theme.StatusText,
		time.Now().Format("15:04"),
		w.c.theme.StatusBrackets))

	// last update time
	c = w.status.GetCell(0, 1)
	c.SetText(fmt.Sprintf("[%s][[%s]Last Update: [%s]%s[%s]]",
		w.c.theme.StatusBrackets,
		w.c.theme.StatusKey,
		w.c.theme.StatusText,
		w.c.lastUpdate.Format("15:04")))

	// todo need to update all articles and unread number
	// total articles
	//c = w.status.GetCell(0, 2)
	//c.SetText(
	//	fmt.Sprintf(
	//		"[%s][[%s]Total Articles: [%s]%d[%s]]",
	//		w.c.theme.StatusBrackets,
	//		w.c.theme.StatusKey,
	//		w.c.theme.StatusText,
	//		len(w.c.articles),
	//		w.c.theme.StatusBrackets))

	// unread number
	//c = w.status.GetCell(0, 3)
	//unread := 0
	//feeds := make(map[string]struct{})
	//for _, a := range w.c.articles {
	//	if _, ok := feeds[a.feed]; !ok {
	//		feeds[a.feed] = struct{}{}
	//	}
	//	if !a.read {
	//		unread++
	//	}
	//}
	//c.SetText(
	//	fmt.Sprintf(
	//		"[%s][[%s]Total Unread: [%s]%d[%s]]",
	//		w.c.theme.StatusBrackets,
	//		w.c.theme.StatusKey,
	//		w.c.theme.StatusText,
	//		unread,
	//		w.c.theme.StatusBrackets,
	//	),
	//)

	c = w.status.GetCell(0, 4)
	c.SetText(
		fmt.Sprintf(
			"[%s][[%s]Feeds: [%s]%d[%s]]",
			w.c.theme.StatusBrackets,
			w.c.theme.StatusKey,
			w.c.theme.StatusText,
			len(w.c.rss.titleFeeds),
			w.c.theme.StatusBrackets,
		),
	)

	c = w.status.GetCell(0, 5)
	c.SetText(
		fmt.Sprintf(
			"[%s][[%s]Help: [%s]%s[%s]]",
			w.c.theme.StatusBrackets,
			w.c.theme.StatusKey,
			w.c.theme.StatusText,
			w.c.conf.KeyHelp,
			w.c.theme.StatusBrackets,
		),
	)

	c = w.status.GetCell(0, 6)
	c.SetText(
		fmt.Sprintf(
			"[%s][[%s]Version: [%s]%s[%s]]",
			w.c.theme.StatusBrackets,
			w.c.theme.StatusKey,
			w.c.theme.StatusText,
			"0.1.0", //todo need to change version automatically
			w.c.theme.StatusBrackets,
		),
	)

	w.app.Draw()
}

func (w *Window) UpdateStatusTicker() {
	w.UpdateStatus()
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				w.UpdateStatus()
			}
		}
	}()
}
