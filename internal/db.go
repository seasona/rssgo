package internal

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type DB struct {
	db       *sql.DB
	tableMap map[string]string
	c        *Controller
}

func (d *DB) Init(c *Controller, dbFile string) {
	d.c = c
	d.tableMap = make(map[string]string)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Can't open database %v, %v", dbFile, err)
	}
	d.db = db

	d.CreateTables(c.rss)
}

// the rss title may contain space or symbol can't be name of database table,
// so omit the character except alphabet and digit
func (d *DB) handleRSSTitle(title string) string {
	var tname string
	for _, c := range title {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			tname += string(c)
		}
	}
	d.tableMap[title] = tname

	return tname
}

func (d *DB) CreateTables(rss *RSS) {
	for _, tu := range rss.titleURLs {
		tname := d.handleRSSTitle(tu.Title)

		ct := fmt.Sprintf(`create table if not exists %v(
				id integer not null primary key,
				feed text,
				title text,
				content text,
				link text,
				read bool,
				display_name string,
				deleted bool,
				published DATETIME);`, tname)

		_, err := d.db.Exec(ct)

		if err != nil {
			log.Panicf("Can't create database %v because: %v", tname, err)
		}
	}
}

func (d *DB) CleanUp() {
	for _, tname := range d.tableMap {
		// sqlite3 is not support bool, but go-sqlite3 can transform bool to 1 or 0
		st1, err := d.db.Prepare(fmt.Sprintf("delete from %v where published < date('now', '-%d day') and deleted = true",
			tname, d.c.conf.DaysKeepDeletedArticle))

		if err != nil {
			log.Println(err)
		}

		if _, err = st1.Exec(); err != nil {
			log.Println(err)
		}

		st1.Close()

		st2, err := d.db.Prepare(fmt.Sprintf("delete from %v where published < date('now', '-%d day') and read = true",
			tname, d.c.conf.DaysKeepReadArticle))

		if err != nil {
			log.Println(err)
		}

		if _, err = st2.Exec(); err != nil {
			log.Println(err)
		}

		st2.Close()
	}
}

func (d *DB) Save(a Article) {
	tname := d.tableMap[a.feed]
	st, err := d.db.Prepare(fmt.Sprintf("select title from %v where feed = ? and title = ? order by id", tname))
	if err != nil {
		log.Println(fmt.Sprintf("select title from %v where feed = ? and title = ? order by id", tname), a.feed)
		log.Println(err)
	}
	defer st.Close()

	res, err := st.Query(a.feed, a.title)
	if err != nil {
		log.Println(err)
	}
	defer res.Close()
	for res.Next() {
		return
	}

	st2, err := d.db.Prepare(fmt.Sprintf("insert into %v(feed, title, content, link, read, display_name, published,deleted) values(?, ?, ?, ?, ?, ?, ?,?)", tname))
	if err != nil {
		log.Println(err)
	}
	defer st2.Close()

	if _, err = st2.Exec(a.feed, a.title, a.content, a.link, false, a.feedDisplay, a.published, false); err != nil {
		log.Println(err)
	}
}

func (d *DB) Delete(a Article) {
	tname := d.tableMap[a.feed]
	st, err := d.db.Prepare(fmt.Sprintf("update %v set deleted = true where id = ?", tname))
	if err != nil {
		log.Println(err)
	}
	defer st.Close()

	if _, err = st.Exec(a.id); err != nil {
		log.Println(err)
	}
}

func (d *DB) MarkRead(a Article) {
	tname := d.tableMap[a.feed]
	st, err := d.db.Prepare(fmt.Sprintf("update %v set read = true where id = ?", tname))
	if err != nil {
		log.Println(err)
	}
	defer st.Close()

	if _, err = st.Exec(a.id); err != nil {
		log.Println(err)
	}
}

func (d *DB) MarkUnread(a Article) {
	tname := d.tableMap[a.feed]
	st, err := d.db.Prepare(fmt.Sprintf("update %v set read = false where id = ?", tname))
	if err != nil {
		log.Println(err)
	}
	defer st.Close()

	if _, err = st.Exec(a.id); err != nil {
		log.Println(err)
	}
}

func (d *DB) MarkAllRead(a Article) {
	tname := d.tableMap[a.feed]
	st, err := d.db.Prepare(fmt.Sprintf("update %v set read = ture", tname))
	if err != nil {
		log.Println(err)
	}
	defer st.Close()

	if _, err := st.Exec(); err != nil {
		log.Println(err)
	}
}

func (d *DB) MarkAllUnread(a Article) {
	tname := d.tableMap[a.feed]
	st, err := d.db.Prepare(fmt.Sprintf("update %v set read = false", tname))
	if err != nil {
		log.Println(err)
	}
	defer st.Close()

	if _, err := st.Exec(); err != nil {
		log.Println(err)
	}
}
