package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"jaytaylor.com/html2text"
)

type Window struct {
	c              *Controller
	feeds          *tview.Table
	articles       *tview.Table
	status         *tview.Table
	help           *tview.Table
	preview        *tview.TextView
	app            *tview.Application
	flexFeed       *tview.Flex
	flexPreArticle *tview.Flex
	flexGlobal     *tview.Flex
	flexStatus     *tview.Flex
	layout         *tview.Flex
	showPreview    bool
	showHelp       bool
	nFeeds         int
	nArticles      int
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

	w.app = tview.NewApplication()
	w.app.SetInputCapture(inputFunc)

	w.InitStatusWindow()
	w.UpdateStatusTicker()
	w.InitFlex()
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

	for k, v := range configKeys {
		i++
		ts = tview.NewTableCell(fmt.Sprintf("%s", v))
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
}

func (w *Window) UpdateStatus() {
	// time now
	c := w.status.GetCell(0, 0)
	// tview support Color tags
	c.SetText(fmt.Sprintf("[%s][[%s]Time: [%s]%s[%s]]",
		w.c.theme.StatusBrackets,
		w.c.theme.StatusKey,
		w.c.theme.StatusText,
		time.Now().Format("15:04:05"),
		w.c.theme.StatusBrackets))

	// last update time
	c = w.status.GetCell(0, 1)
	c.SetText(fmt.Sprintf("[%s][[%s]Last Update: [%s]%s[%s]]",
		w.c.theme.StatusBrackets,
		w.c.theme.StatusKey,
		w.c.theme.StatusText,
		w.c.lastUpdate.Format("15:04"),
		w.c.theme.StatusBrackets))

	// todo need to update all articles and unread number
	// total articles
	c = w.status.GetCell(0, 2)
	c.SetText(
		fmt.Sprintf(
			"[%s][[%s]Total Articles: [%s]%d[%s]]",
			w.c.theme.StatusBrackets,
			w.c.theme.StatusKey,
			w.c.theme.StatusText,
			0,
			w.c.theme.StatusBrackets))

	// unread number
	c = w.status.GetCell(0, 3)
	c.SetText(
		fmt.Sprintf(
			"[%s][[%s]Total Unread: [%s]%d[%s]]",
			w.c.theme.StatusBrackets,
			w.c.theme.StatusKey,
			w.c.theme.StatusText,
			0,
			w.c.theme.StatusBrackets,
		),
	)

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

	// w.app.Draw()
}

func (w *Window) UpdateStatusTicker() {
	// the app in tview will draw first time during Run()
	// https://github.com/rivo/tview/wiki/Concurrency
	w.UpdateStatus()
	// w.app.Draw()
	//? can't use app.Draw() there
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				w.UpdateStatus()
				w.app.Draw()
			}
			// another way to write timer in tview
			// time.Sleep(1 * time.Second)
			// w.app.QueueUpdateDraw(w.UpdateStatus)
		}
	}()
}

func (w *Window) InitFlex() {
	w.flexFeed = tview.NewFlex().SetDirection(tview.FlexRow)
	w.flexFeed.AddItem(w.feeds, 0, 1, false)

	// article:preview = 2:1
	w.flexPreArticle = tview.NewFlex().SetDirection(tview.FlexRow)
	w.flexPreArticle.AddItem(w.articles, 0, 2, false)
	w.flexPreArticle.AddItem(w.preview, 0, 1, false)

	// feed:article = 2:5
	w.flexGlobal = tview.NewFlex().SetDirection(tview.FlexColumn)
	w.flexGlobal.AddItem(w.flexFeed, 0, 2, false)
	w.flexGlobal.AddItem(w.flexPreArticle, 0, 5, false)

	w.flexStatus = tview.NewFlex().SetDirection(tview.FlexRow)
	w.flexStatus.AddItem(w.flexGlobal, 0, 20, false)
	w.flexStatus.AddItem(w.status, 1, 1, false)

	w.layout = tview.NewFlex()
	w.layout.AddItem(w.flexStatus, 0, 1, false)
}

func (w *Window) Start() {
	if err := w.app.SetRoot(w.layout, true).SetFocus(w.feeds).Run(); err != nil {
		panic(err)
	}
}

// SwitchFocus feeds->article->preview->feeds
func (w *Window) SwitchFocus() {
	f := w.app.GetFocus()
	if f == w.feeds {
		w.app.SetFocus(w.articles)
	} else {
		w.app.SetFocus(w.feeds)
	}
}

func (w *Window) TriggerPreview() {
	if !w.showPreview {
		w.flexPreArticle.AddItem(w.preview, 0, 1, false)
		w.showPreview = true
	} else {
		w.flexPreArticle.RemoveItem(w.preview)
		w.showPreview = false
	}
}

func (w *Window) TriggerHelp() {
	if !w.showHelp {
		w.flexPreArticle.RemoveItem(w.articles)
		w.flexPreArticle.RemoveItem(w.preview)
		w.flexPreArticle.AddItem(w.help, 0, 1, false)
		w.showHelp = true
	} else {
		w.flexPreArticle.AddItem(w.articles, 0, 2, false)
		if w.showPreview {
			w.flexPreArticle.AddItem(w.preview, 0, 1, false)
		}
		w.flexPreArticle.RemoveItem(w.help)
		w.showHelp = false
	}
}

// ClearFeeds will clear the feed window and reset it
func (w *Window) ClearFeeds() {
	w.feeds.Clear()
	w.feeds.SetTitle(fmt.Sprintf("%s Feeds", w.c.theme.FeedIcon)).SetTitleColor(tcell.GetColor(w.c.theme.FeedBorderTitle))
	w.nFeeds = 0

	ts := tview.NewTableCell("Total")
	ts.SetAlign(tview.AlignLeft)
	ts.Attributes |= tcell.AttrBold
	ts.SetSelectable(false)
	ts.SetTextColor(tcell.GetColor(w.c.theme.TableHead))
	w.feeds.SetCell(0, 0, ts)

	ts = tview.NewTableCell("Unread")
	ts.SetAlign(tview.AlignLeft)
	ts.Attributes |= tcell.AttrBold
	ts.SetSelectable(false)
	ts.SetTextColor(tcell.GetColor(w.c.theme.TableHead))
	w.feeds.SetCell(0, 1, ts)

	ts = tview.NewTableCell("Feed")
	ts.SetAlign(tview.AlignLeft)
	ts.Attributes |= tcell.AttrBold
	ts.SetTextColor(tcell.GetColor(w.c.theme.TableHead))
	ts.SetSelectable(false)
	w.feeds.SetCell(0, 2, ts)

	w.feeds.SetSelectable(true, false)
}

func (w *Window) AddToFeeds(title string, unread, total int, ref string) {
	w.nFeeds++
	nc := tview.NewTableCell(fmt.Sprintf("%d", total))
	nc.SetAlign(tview.AlignLeft)
	w.feeds.SetCell(w.nFeeds, 0, nc)
	nc.SetSelectable(true)
	nc.SetTextColor(tcell.GetColor(w.c.theme.TotalColumn))

	// Display number of unread articles
	nc = tview.NewTableCell(fmt.Sprintf("%d", unread))
	nc.SetAlign(tview.AlignLeft)
	w.feeds.SetCell(w.nFeeds, 1, nc)
	nc.SetSelectable(true)
	nc.SetTextColor(tcell.GetColor(w.c.theme.UnreadColumn))

	nc = tview.NewTableCell(fmt.Sprintf("%s", title))
	nc.SetAlign(tview.AlignLeft)
	w.feeds.SetCell(w.nFeeds, 2, nc)
	nc.SetSelectable(true)
	nc.SetTextColor(tcell.GetColor("white"))
	nc.SetReference(ref)
}

func (w *Window) ClearArticles() {
	w.nArticles = 0
	w.articles.Clear()

	ts := tview.NewTableCell("")
	ts.SetAlign(tview.AlignLeft)
	ts.Attributes |= tcell.AttrBold
	//	ts.SetTextColor(tcell.ColorGreen)
	ts.SetSelectable(false)
	w.articles.SetCell(0, 0, ts)

	ts = tview.NewTableCell("Feed")
	ts.SetAlign(tview.AlignLeft)
	ts.Attributes |= tcell.AttrBold
	ts.SetTextColor(tcell.GetColor(w.c.theme.TableHead))
	ts.SetSelectable(false)
	w.articles.SetCell(0, 1, ts)

	ts = tview.NewTableCell("Title")
	ts.Attributes |= tcell.AttrBold
	ts.SetTextColor(tcell.GetColor(w.c.theme.TableHead))
	ts.SetSelectable(false)
	w.articles.SetCell(0, 2, ts)

	ts = tview.NewTableCell("Published")
	ts.Attributes |= tcell.AttrBold
	ts.SetTextColor(tcell.GetColor(w.c.theme.TableHead))
	ts.SetSelectable(false)
	w.articles.SetCell(0, 3, ts)

	w.articles.SetSelectable(true, false)
}

func (w *Window) AddToArticles(a *Article) {
	if a == nil {
		return
	}
	w.nArticles++

	nc := tview.NewTableCell("")
	nc.SetAlign(tview.AlignLeft)
	w.articles.SetCell(w.nArticles, 0, nc)
	nc.SetSelectable(false)
	if a.read {
		nc.SetText(w.c.theme.ReadMarker)
	}

	color := "white"
	fc := tview.NewTableCell(fmt.Sprintf("[%s]%s", color, a.feed))
	fc.SetTextColor(tcell.GetColor(color))
	fc.SetAlign(tview.AlignLeft)
	fc.SetMaxWidth(20)
	w.articles.SetCell(w.nArticles, 1, fc)

	tc := tview.NewTableCell(fmt.Sprintf("[%s]%s", w.c.theme.Title, a.title))
	tc.SetTextColor(tcell.GetColor(w.c.theme.Title))
	tc.SetSelectable(true)
	tc.SetMaxWidth(80)
	tc.SetAlign(tview.AlignLeft)
	tc.SetReference(a)
	w.articles.SetCell(w.nArticles, 2, tc)

	dc := tview.NewTableCell(
		fmt.Sprintf(
			"[%s]%s", w.c.theme.Time, a.published.Format("2006-01-02 15:04:05"),
		),
	)
	dc.SetTextColor(tcell.GetColor(w.c.theme.Date))
	dc.SetAlign(tview.AlignLeft)
	w.articles.SetCell(w.nArticles, 3, dc)
}

func (w *Window) ArticlesHasFocus() bool {
	if w.app.GetFocus() == w.articles {
		return true
	}
	return false
}

func (c *Controller) GetArticleForSelection() *Article {
	if !c.win.ArticlesHasFocus() {
		return nil
	}

	var cell *tview.TableCell

	r, _ := c.win.articles.GetSelection()
	cell = c.win.articles.GetCell(r, 2)
	ref := cell.GetReference()
	if ref != nil {
		return ref.(*Article)
	}
	return nil
}

func (w *Window) AddToPreview(a *Article) {
	parsed, err := html2text.FromString(a.content, html2text.Options{PrettyTables: true})
	if err != nil {
		log.Printf("Failed to parse html to text, rendering original.")
		parsed = a.content
	}

	w.preview.Clear()

	text := fmt.Sprintf(
		"[%s][%s][%s] %s [white]([%s]%s[white])\n\n[%s]%s\n\nLink: [%s]%s",
		"white",
		a.feed,
		w.c.theme.Title,
		a.title,
		w.c.theme.Date,
		a.published,
		w.c.theme.PreviewText,
		parsed,
		w.c.theme.PreviewLink,
		a.link,
	)
	w.preview.SetText(text)
	w.preview.ScrollToBeginning()
}

func (w *Window) FeedSelectedFunc(f func(r, c int)) {
	w.feeds.SetSelectedFunc(f)
	w.feeds.SetSelectionChangedFunc(f)
}

func (w *Window) ArticleSelectedFunc(f func(r, c int)) {
	w.articles.SetSelectedFunc(f)
	w.articles.SetSelectionChangedFunc(f)
}
