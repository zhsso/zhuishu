package main

import (
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// TeleBot 机器人接口
type TeleBot struct {
}

func (t *TeleBot) run() {

}

func teleSearchBook(keyword string) string {
	return keyword
}

func runBotServer() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

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
		user := getUser(update.Message.Chat.ID)
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
			msg.Text = user.addBook(update.Message.CommandArguments())
		case "delete":
			msg.Text = user.deleteBook(update.Message.CommandArguments())
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
