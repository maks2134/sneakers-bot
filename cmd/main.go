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

	utils.SeedProducts(db)

	botAPI, err := tgbotapi.NewBotAPI(loadConfig.TelegramToken)
	if err != nil {
		log.Fatal("Cannot init telegram bot", err)
	}

	botAPI.Debug = false
	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	// Создаем репозитории
	userRepo := repository.NewUserRepo(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Создаем сервисы
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(orderRepo)

	// НОВЫЙ БЛОК: Создаем и настраиваем адаптер для платежей
	paymentProvider := adapters.NewSimplePaymentAdapter(
		"SNEAKERS SHOP",
		"СБЕР 2202 2002 2002 2002",
		map[string]string{
			"USDT (TRC-20)": "T...ВАШ_АДРЕС...XYZ",
			"BTC":           "bc1...ВАШ_АДРЕС...xyz",
		},
	)

	// Передаем все зависимости в обработчик
	botHandler := adapters.NewBotHandler(botAPI, userService, productService, orderService, paymentProvider)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := botAPI.GetUpdatesChan(u)

	botHandler.HandleUpdates(updates)
}
