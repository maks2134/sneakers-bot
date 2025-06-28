package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shopspring/decimal"
	"log"
	"snakers-bot/internal/adapters"
	"snakers-bot/internal/config"
	"snakers-bot/internal/repository"
	"snakers-bot/internal/service"
	"snakers-bot/internal/usecases"
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

	productRepo := repository.NewProductRepository(db)
	products, _ := productRepo.GetAll(context.Background())
	if len(products) == 0 {
		productRepo.Create(context.Background(), &usecases.Product{
			Name:  "Крутые кроссовки",
			Price: decimal.NewFromFloat(9990.99),
		})
		fmt.Println("Test product created.")
	}

	botAPI, err := tgbotapi.NewBotAPI(loadConfig.TelegramToken)
	if err != nil {
		log.Fatal("Cannot init telegram bot", err)
	}

	botAPI.Debug = true
	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	userRepo := repository.NewUserRepo(db)

	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)

	botHandler := adapters.NewBotHandler(botAPI, userService, productService)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := botAPI.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Failed to get updates channel: %v", err)
	}

	botHandler.HandleUpdates(updates)
}
