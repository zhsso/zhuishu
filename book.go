package main

import (
	"fmt"
	"strconv"
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
	users    map[int64]*User
	key      string
	needMark bool
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

func createBook(bookID int64, book *Book) *Book {
	book.key = fmt.Sprintf("book/%d", bookID)
	book.id = bookID
	//book.users =
	vals := map[string]interface{}{
		"platform": book.Platform,
		"title":    book.Title,
		"bookId":   book.BookID,
		"url":      book.URL,
		"author":   book.Author,
	}
	redisClient.HMSet(book.key, vals)
	redisClient.SAdd("books", bookID)
	book.needMark = true
	return book
}

func newBook(bookID int64) *Book {
	// 图书信息保存在 book/bookID
	b := &Book{
		id:    bookID,
		key:   fmt.Sprintf("book/%d", bookID),
		users: make(map[int64]*User),
	}
	//
	fields := redisClient.HGetAll(b.key)
	if fields == nil {
		return nil
	}
	//TODO:
	fieldsv := fields.Val()
	b.Platform = fieldsv["platform"]
	b.Title = fieldsv["title"]
	b.BookID, _ = strconv.ParseInt(fieldsv["bookId"], 10, 64)
	b.URL = fieldsv["url"]
	b.Author = fieldsv["author"]
	b.needMark = true
	return b
}

func (b *Book) addUser(u *User) {
	b.users[u.id] = u
	if len(b.users) == 1 {
		b.needMark = true
	}
}

func (b *Book) deleteUser(u *User) int {
	delete(b.users, u.id)
	return len(b.users)
}

// Start 一直获取更新
func (b *Book) start() {
	for {
		select {
		case <-time.After(time.Minute):
			b.refresh()
		}
	}
}

func (b *Book) refresh() {
	if len(b.users) > 0 {
		if b.needMark {
			markBook(b.Platform, b.BookID)
		} else {
			chapters := fetchBook(b.Platform, b.BookID)
			if len(chapters) != 0 {
				go b.sync(chapters)
			}
		}
	}
}

func (b *Book) sync(chapters []Chapter) {
	fmt.Printf("%s to users counts: %d\n", b.Title, len(b.users))
	for _, chapter := range chapters {
		for _, user := range b.users {
			fmt.Printf("%s to user: %d\n", b.Title, user.id)
			user.sendChapter(b, chapter)
		}
	}
}
