package adapters

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"log"
	"snakers-bot/internal/service"
	"snakers-bot/internal/usecases"
)

type BotHandler struct {
	bot            *tgbotapi.BotAPI
	userService    *service.UserService
	productService *service.ProductService
}

func NewBotHandler(bot *tgbotapi.BotAPI, us *service.UserService, ps *service.ProductService) *BotHandler {
	return &BotHandler{
		bot:            bot,
		userService:    us,
		productService: ps,
	}
}

func (h *BotHandler) HandleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		ctx := context.Background()

		switch update.Message.Command() {
		case "start":
			h.handleStart(ctx, update.Message)
		case "catalog":
			h.handleCatalog(ctx, update.Message)
		default:
			h.handleDefault(update.Message)
		}
	}
}

func (h *BotHandler) handleStart(ctx context.Context, message *tgbotapi.Message) {
	telegramID := message.From.ID

	_, err := h.userService.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// ИСПРАВЛЕНО: Добавляем TelegramID при создании
			newUser := &usecases.User{
				Name:       message.From.UserName,
				TelegramID: telegramID,
			}
			if err := h.userService.CreateUser(ctx, newUser); err != nil {
				log.Printf("Failed to create user: %v", err)
				h.sendErrorMessage(message.Chat.ID, "Не удалось создать ваш профиль.")
				return
			}
			log.Printf("New user created: %s, with ID: %d", newUser.Name, newUser.TelegramID)
		} else {
			log.Printf("Failed to get user: %v", err)
			h.sendErrorMessage(message.Chat.ID, "Произошла ошибка при поиске вашего профиля.")
			return
		}
	}

	responseText := fmt.Sprintf("Привет, %s! Это бот для продажи самых топовых шузов. Используй /catalog, чтобы посмотреть товары.", message.From.FirstName)
	msg := tgbotapi.NewMessage(message.Chat.ID, responseText)
	h.bot.Send(msg)
}

func (h *BotHandler) handleCatalog(ctx context.Context, message *tgbotapi.Message) {
	products, err := h.productService.GetAllProducts(ctx)
	if err != nil {
		log.Printf("Failed to get products: %v", err)
		h.sendErrorMessage(message.Chat.ID, "Не удалось загрузить каталог товаров.")
		return
	}

	if len(products) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Каталог пока пуст. Загляните попозже!")
		h.bot.Send(msg)
		return
	}

	h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Вот что у нас есть:"))

	for _, p := range products {
		caption := fmt.Sprintf(
			"*%s*\n\n%s\n\nЦена: *%s руб.*",
			p.Name,
			p.Description,
			p.Price.StringFixed(2),
		)

		photoMsg := tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FileURL(p.ImageURL))
		photoMsg.Caption = caption
		photoMsg.ParseMode = tgbotapi.ModeMarkdown

		if _, err := h.bot.Send(photoMsg); err != nil {
			log.Printf("Failed to send product photo: %v", err)
		}
	}
}

func (h *BotHandler) handleDefault(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Неизвестная команда. Используйте /start или /catalog")
	h.bot.Send(msg)
}

func (h *BotHandler) sendErrorMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	h.bot.Send(msg)
}
