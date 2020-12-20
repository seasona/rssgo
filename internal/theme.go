package internal

import (
	"encoding/json"
	"log"
	"os"
)

type Theme struct {
	FeedNames          []string `json:"feedNames"`
	Date               string   `json:"date"`
	Time               string   `json:"time"`
	ArticleBorder      string   `json:"articleBorder"`
	PreviewBorder      string   `json:"previewBorder"`
	FeedBorder         string   `json:"feedBorder"`
	ArticleBorderTitle string   `json:"articleBorderTitle"`
	FeedBorderTitle    string   `json:"feedBorderTitle"`
	PreviewBorderTitle string   `json:"previewBorderTitle"`
	Highlights         string   `json:"highlights"`
	TableHead          string   `json:"tableHead"`
	Title              string   `json:"title"`
	UnreadFeedName     string   `json:"unreadFeedName"`
	TotalColumn        string   `json:"totalColumn"`
	UnreadColumn       string   `json:"unreadColumn"`
	PreviewText        string   `json:"previewText"`
	PreviewLink        string   `json:"previewLink"`
	ReadMarker         string   `json:"unreadMarker"`
	FeedIcon           string   `json:"feedIcon"`
	ArticleIcon        string   `json:"articleIcon"`
	PreviewIcon        string   `json:"previewIcon"`
	StatusBackground   string   `json:"statusBackground"`
	StatusText         string   `json:"statusText"`
	StatusKey          string   `json:"statusKey"`
	StatusBrackets     string   `json:"statusBrackets"`
}

func (t *Theme) LoadTheme(file string) {
	tf, err := os.Open(file)
	defer tf.Close()

	if err != nil {
		log.Fatal("Can't open theme file:", err)
	}

	jsonParser := json.NewDecoder(tf)
	err = jsonParser.Decode(t)

	if err != nil {
		log.Fatal("Failed to parse theme file:", err)
	}
}
