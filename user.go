package main

import "github.com/go-telegram-bot-api/telegram-bot-api"

//用户信息
type User struct {
	ID  int64
	bot *tgbotapi.BotAPI
}
