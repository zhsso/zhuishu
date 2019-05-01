package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// TeleBot 机器人接口
type TeleBot struct {
	bot *tgbotapi.BotAPI
}

func teleSearchBook(keyword string) string {
	books := searchBook(keyword)
	str := "搜索结果"
	for _, book := range books {
		if book == nil {
			continue
		}
		str = fmt.Sprintf("%s\n%d: %s %s %s", str, book.id, book.Title, book.Author, book.Platform)
	}
	return str
}

func (t *TeleBot) run() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
			log.Panic(err)
	}
	t.bot = bot
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}
		user := userManager.getUser(update.Message.Chat.ID)
		// Create a new MessageConfig. We don't have text yet,
		// so we should leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "search":
			msg.Text = teleSearchBook(update.Message.CommandArguments())
		case "list":
			msg.Text = user.getBookListStr()
		case "add":
			if user.addBook(update.Message.CommandArguments()) {
				msg.Text = "添加成功"
			} else {
				msg.Text = "添加失败"
			}
		case "delete":
			if user.deleteBook(update.Message.CommandArguments()) {
				msg.Text = "删除成功"
			} else {
				msg.Text = "删除失败"
			}
		default:
			msg.Text = "I don't know that command"
		}
		t.send(msg)
	}
}

func (t *TeleBot) send(msg tgbotapi.MessageConfig) {
	if _, err := t.bot.Send(msg); err != nil {
		userID := msg.ChatID
		userManager.deleteUser(userID)
	}
}
