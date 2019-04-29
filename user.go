package main

import "github.com/go-telegram-bot-api/telegram-bot-api"

// User 用户信息
type User struct {
	ID  int64
	bot *tgbotapi.BotAPI
}

func getUser(id int64) *User {
	return &User{}
}

func (u User) getBookListStr() string {
	return ""
}

func (u User) addBook(bookID string) string {
	return ""
}

func (u User) deleteBook(bookID string) string {
	return ""
}
