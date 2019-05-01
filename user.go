package main

import (
	"fmt"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

//用户信息

// User 用户信息
type User struct {
	id      int64
	books   map[int64]*Book
	bookKey string
}

func newUser(id int64) *User {
	return &User{
		id:      id,
		books:   make(map[int64]*Book),
		bookKey: fmt.Sprintf("user/book/%d", id),
	}
}

func (u User) getBookListStr() string {
	str := "关注列表"
	for _, book := range u.books {
		str = fmt.Sprintf("%s\n%d: %s(%s@%s)", str, book.id, book.Title, book.Author, book.Platform)
	}
	return str
}

func (u *User) loadBooks() {
	bookIds := redisClient.SMembers(u.bookKey).Val()
	if bookIds != nil {
		for _, bookIDstr := range bookIds {
			bookID, err := strconv.ParseInt(bookIDstr, 10, 64)
			if err == nil {
				book := bookManager.getBook(bookID)
				if book != nil {
					u.books[bookID] = book
					book.addUser(u)
				}
			}
		}
	}
}

func (u *User) addBook(bookIDStr string) bool {
	bookID, err := strconv.ParseInt(bookIDStr, 10, 64)
	if err != nil {
		println(bookIDStr)
		return false
	}
	//加进列表
	if redisClient.SAdd(u.bookKey, bookID).Val() == 0 {
		println("sadd fail!")
		println(u.bookKey)
		return false
	}

	book := bookManager.addBook(u, bookID)
	u.books[bookID] = book
	return true
}

func (u *User) deleteBook(bookIDStr string) bool {
	bookID, err := strconv.ParseInt(bookIDStr, 10, 64)
	if err != nil {
		return false
	}
	book := u.books[bookID]
	if book == nil {
		return false
	}
	redisClient.SRem(u.bookKey, bookID)
	delete(u.books, bookID)

	bookManager.deleteBook(u, bookID)
	return true
}

func (u *User) delete() {
	for bookID := range u.books {
		bookManager.deleteBook(u, bookID)
	}
}

func (u *User) sendChapter(b *Book, c Chapter) {
	msgTXT := fmt.Sprintf("%s: %s\n%s", b.Title, c.Title, c.URL)
	msg := tgbotapi.NewMessage(u.id, msgTXT)
	teleBot.send(msg)
}
