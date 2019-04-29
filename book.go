package main

import (
	"fmt"
	"time"
)

// Book 书信息
type Book struct {
	Platform string `json:"platform"`
	Title    string `json:"title"`
	BookID   int64  `json:"book_id"`
	URL      string `json:"url"`
	Author   string `json:"author"`
	id       int64
	stop     bool
}

func (b Book) String() string {
	return fmt.Sprintf("%s %s %s %s", b.Platform, b.Title, b.Author, b.URL)
}

// Chapter 章节信息
type Chapter struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func (c Chapter) String() string {
	return fmt.Sprintf("%s %s", c.Title, c.URL)
}

// Start 一直获取更新
func (b *Book) Start(c chan []Chapter) {
	for b.stop != true {
		select {
		case <-time.After(10 * time.Minute):
			b.refresh(c)
		}
	}
}

func (b *Book) refresh(c chan []Chapter) {
	c <- fetchBook(b.Platform, b.BookID)
}
