package main

import (
	"fmt"
	"strconv"
)

// BookManager 书管理
type BookManager struct {
	books map[int64]*Book
}

func newBookManager() *BookManager {
	return &BookManager{
		books: make(map[int64]*Book),
	}
}

//加载所有
func (b *BookManager) loadAll() {
	//所有书籍id存在 books
	booksIDStrs := redisClient.SMembers("books")
	if booksIDStrs != nil {
		for _, bookIDStr := range booksIDStrs.Val() {
			bookID, err := strconv.ParseInt(bookIDStr, 10, 64)
			if err == nil {
				book := newBook(bookID)
				b.books[bookID] = book
			}
		}
	}
}

//增加书籍
func (b *BookManager) addBook(u *User, bookID int64) *Book {
	if book, ok := b.books[bookID]; ok {
		book.addUser(u)
		return book
	}
	book := newBook(bookID)
	b.books[bookID] = book
	book.addUser(u)
	go book.start()
	return book
}

//删除书籍
func (b *BookManager) deleteBook(u *User, bookID int64) {
	if book, ok := b.books[bookID]; ok {
		book.deleteUser(u)
	}
}

//获取书籍信息
func (b *BookManager) run() {
	for _, book := range b.books {
		go book.start()
	}
}

func (b *BookManager) getBook(bookID int64) *Book {
	if book, ok := b.books[bookID]; ok {
		return book
	}
	return nil
}

func (b *BookManager) checkBook(book *Book) *Book {
	key := fmt.Sprintf("bookm/%d", book.BookID)
	if redisClient.Exists(key).Val() == 1 {
		//存在
		bookIDStr := redisClient.Get(key).Val()
		bookID, _ := strconv.ParseInt(bookIDStr, 10, 64)
		book1 := bookManager.getBook(bookID)
		if book1 == nil {
			book.id = bookID
			book1 = book
		}
		return book1
	}
	//不存在,创建
	bookID := redisClient.Incr("bookID").Val()
	redisClient.Set(key, bookID, 0)
	//设置信息
	return createBook(bookID, book)
}
