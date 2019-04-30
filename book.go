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
	isStop   bool
	users    map[int64]*User
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

func newBook(BookID int64) *Book {
	return &Book{}
}

func (b *Book) addUser(u *User) {
	b.users[u.id] = u
}

func (b *Book) deleteUser(u *User) int {
	delete(b.users, u.id)
	return len(b.users)
}

func (b *Book) stop() {
	b.isStop = true
}

// Start 一直获取更新
func (b *Book) start() {
	for b.isStop != true {
		select {
		case <-time.After(10 * time.Minute):
			b.refresh()
		}
	}
}

func (b *Book) refresh() {
	chapters := fetchBook(b.Platform, b.BookID)
	if len(chapters) != 0 {
		go b.sync(chapters)
	}
}

func (b *Book) sync(chapters []Chapter) {
	for _, chapter := range chapters {
		for _, user := range b.users {
			user.sendChapter(b, chapter)
		}
	}
}
