package adapters

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"log"
	"snakers-bot/internal/service"
	"snakers-bot/internal/usecases"
	"strings"
)

// Добавляем константы для текста кнопок. Это хорошая практика.
const (
	buttonCatalog = "Каталог 👟"
	buttonOrders  = "Мои заказы 📦"
	buttonProfile = "Профиль 👤"
)

// Определяем нашу клавиатуру один раз, чтобы использовать ее везде.
var mainKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(buttonCatalog),
		tgbotapi.NewKeyboardButton(buttonOrders),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(buttonProfile),
	),
)

type BotHandler struct {
	bot            *tgbotapi.BotAPI
	userService    *service.UserService
	productService *service.ProductService
	orderService   *service.OrderService // Добавляем сервис заказов
}

// Обновляем конструктор
func NewBotHandler(bot *tgbotapi.BotAPI, us *service.UserService, ps *service.ProductService, os *service.OrderService) *BotHandler {
	return &BotHandler{
		bot:            bot,
		userService:    us,
		productService: ps,
		orderService:   os,
	}
}

// Обновляем главный обработчик, чтобы он реагировал на текст
func (h *BotHandler) HandleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		ctx := context.Background()

		// Теперь мы обрабатываем и команды, и текст кнопок в одном месте
		switch update.Message.Text {
		case "/start":
			h.handleStart(ctx, update.Message)
		case buttonCatalog, "/catalog": // Реагируем и на кнопку, и на команду
			h.handleCatalog(ctx, update.Message)
		case buttonOrders:
			h.handleMyOrders(ctx, update.Message)
		case buttonProfile:
			h.handleProfile(ctx, update.Message)
		default:
			h.handleDefault(update.Message)
		}
	}
}

func (h *BotHandler) handleStart(ctx context.Context, message *tgbotapi.Message) {

	if _, err := h.getOrCreateUser(ctx, message); err != nil {
		h.sendErrorMessage(message.Chat.ID, "Произошла ошибка при создании вашего профиля.")
		return
	}

	responseText := fmt.Sprintf("Привет, %s! Я бот для продажи самых топовых шузов. Выбери, что тебя интересует 👇", message.From.FirstName)
	msg := tgbotapi.NewMessage(message.Chat.ID, responseText)

	// Прикрепляем нашу клавиатуру к сообщению
	msg.ReplyMarkup = mainKeyboard

	h.bot.Send(msg)
}

func (h *BotHandler) handleCatalog(ctx context.Context, message *tgbotapi.Message) {
	// ... код этой функции остается без изменений ...
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

// Новая функция для обработки кнопки "Мои заказы"
func (h *BotHandler) handleMyOrders(ctx context.Context, message *tgbotapi.Message) {
	user, err := h.getOrCreateUser(ctx, message)
	if err != nil {
		h.sendErrorMessage(message.Chat.ID, "Не удалось найти ваш профиль.")
		return
	}

	orders, err := h.orderService.GetUserOrders(ctx, user.ID)
	if err != nil {
		log.Printf("Failed to get user orders: %v", err)
		h.sendErrorMessage(message.Chat.ID, "Не удалось загрузить ваши заказы.")
		return
	}

	if len(orders) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "У вас пока нет заказов. Пора это исправить! 😉\nНажмите 'Каталог 👟', чтобы выбрать кроссовки.")
		h.bot.Send(msg)
		return
	}

	var responseText strings.Builder
	responseText.WriteString("🧾 *Ваши заказы:*\n\n")
	for _, order := range orders {
		responseText.WriteString(fmt.Sprintf(
			"📦 *Заказ #%d* от %s\n*Статус:* %s\n\n",
			order.ID,
			order.Date.Format("02.01.2006"),
			order.Status,
		))
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, responseText.String())
	msg.ParseMode = tgbotapi.ModeMarkdown
	h.bot.Send(msg)
}

// Новая функция для обработки кнопки "Профиль"
func (h *BotHandler) handleProfile(ctx context.Context, message *tgbotapi.Message) {
	user, err := h.getOrCreateUser(ctx, message)
	if err != nil {
		h.sendErrorMessage(message.Chat.ID, "Не удалось найти ваш профиль.")
		return
	}

	responseText := fmt.Sprintf(
		"👤 *Ваш профиль*\n\n"+
			"Имя: *%s*\n"+
			"Баланс: *%s руб.*\n"+
			"Бонусные баллы: *%d*",
		user.Name,
		user.Balance.StringFixed(2),
		user.LoyaltyPoints,
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, responseText)
	msg.ParseMode = tgbotapi.ModeMarkdown
	h.bot.Send(msg)
}

func (h *BotHandler) handleDefault(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "🤔 Не понимаю вас. Пожалуйста, используйте кнопки меню.")
	h.bot.Send(msg)
}

func (h *BotHandler) sendErrorMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	h.bot.Send(msg)
}

// Вспомогательная функция, чтобы не дублировать код получения/создания юзера
func (h *BotHandler) getOrCreateUser(ctx context.Context, message *tgbotapi.Message) (*usecases.User, error) {
	user, err := h.userService.GetUserByTelegramID(ctx, message.From.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			newUser := &usecases.User{
				Name:       message.From.UserName,
				TelegramID: message.From.ID,
			}
			if err := h.userService.CreateUser(ctx, newUser); err != nil {
				log.Printf("Failed to create user: %v", err)
				return nil, err
			}
			log.Printf("New user created: %s, with ID: %d", newUser.Name, newUser.TelegramID)
			return newUser, nil
		}
		log.Printf("Failed to get user: %v", err)
		return nil, err
	}
	return user, nil
}
