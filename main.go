package main

import (
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type user_settings struct {
	User_ID int64
}

func validInn(inputString string) bool {
	for i := 0; i < len(inputString); i++ {
		if inputString[i] < '0' || inputString[i] > '9' {
			return false
		}
	}

	return true
}

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

	var reply string
	storage := make(map[int64]string)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		text := update.Message.Text

		log.Printf("[%d] %s", chatID, text)

		inputData := strings.Fields(text)

		switch inputData[0] {
		case "/start":
			reply = "👋 Привет! Я помогу вам принимать платежи и автоматически отправлять чеки.\n\n" +
				"Сначала привяжите свой ИНН: /set_inn 123456789012"
		case "/help":
			reply = "Доступные команды:\n/start — начать\n/set_inn — указать ИНН\n/pay — создать платёж\n/my_inn — узнать свой ИНН. Если его нет в базе, то об этом сообщит система"
		case "/set_inn":
			if len(inputData) == 2 && len(inputData[1]) == 12 && validInn(inputData[1]) {
				reply = "ИНН добавлен в бд!"
				storage[chatID] = inputData[1]
			} else {
				reply = "ИНН некорректен!"
			}
		case "/my_inn":
			inn, ok := storage[chatID]
			if ok {
				reply = "Ваш ИНН: " + inn
			} else {
				reply = "ИНН отсуствует в бд"
			}
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
