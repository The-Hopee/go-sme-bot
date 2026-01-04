package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не задан")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Авторизован как @%s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		text := update.Message.Text

		log.Printf("[%d] %s", chatID, text)

		var reply string
		switch text {
		case "/start":
			reply = "👋 Привет! Я помогу вам принимать платежи и автоматически отправлять чеки.\n\n" +
				"Сначала привяжите свой ИНН: /set_inn 123456789012"
		case "/help":
			reply = "Доступные команды:\n/start — начать\n/set_inn — указать ИНН\n/pay — создать платёж"
		default:
			reply = "Неизвестная команда. Напишите /help"
		}

		msg := tgbotapi.NewMessage(chatID, reply)

		msg.ReplyToMessageID = update.Message.MessageID
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		}
	}
}
