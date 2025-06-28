package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"snakers-bot/internal/adapters"
	"snakers-bot/internal/config"
	"snakers-bot/internal/repository"
	"snakers-bot/internal/service"
	"snakers-bot/internal/usecases"
	"snakers-bot/internal/utils"
)

func main() {
	loadConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config", err)
	}

	db, err := config.InitDB(loadConfig)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	err = db.AutoMigrate(
		&usecases.User{},
		&usecases.Product{},
		&usecases.Order{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Вызываем сидер после миграций
	utils.SeedProducts(db)

	botAPI, err := tgbotapi.NewBotAPI(loadConfig.TelegramToken)
	if err != nil {
		log.Fatal("Cannot init telegram bot", err)
	}

	botAPI.Debug = false
	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	userRepo := repository.NewUserRepo(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(orderRepo)

	// Передаем новый сервис в конструктор
	botHandler := adapters.NewBotHandler(botAPI, userService, productService, orderService)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := botAPI.GetUpdatesChan(u)

	botHandler.HandleUpdates(updates)
}
