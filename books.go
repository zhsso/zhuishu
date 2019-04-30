package main

// BookManager 书管理
type BookManager struct {
	books map[int64]*Book
}

//加载所有
func (b *BookManager) loadAll() {

}

//增加书籍
func (b *BookManager) addBook(u *User, bookID int64) *Book {
	if book, ok := b.books[bookID]; ok {
		book.addUser(u)
		return book
	} else {
		book := newBook(bookID)
		b.books[bookID] = book
		go book.start()
		return book
	}
}

//删除书籍
func (b *BookManager) deleteBook(u *User, bookID int64) {
	if book, ok := b.books[bookID]; ok {
		if book.deleteUser(u) == 0 {
			//没有追更的用户
			book.stop()
			delete(b.books, bookID)
		}
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
