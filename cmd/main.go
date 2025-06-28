package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"snakers-bot/internal/config"
)

func main() {
	loadConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load loadConfig", err)
	}

	bot, err := tgbotapi.NewBotAPI(loadConfig.TelegramToken)
	if err != nil {
		log.Fatal("Cannot init telegram bot", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Васап это бот для продажи самых топовых шузов на рынке")
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}
