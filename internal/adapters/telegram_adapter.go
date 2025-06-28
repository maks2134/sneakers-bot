package adapters

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"snakers-bot/internal/interfaces"
	"snakers-bot/internal/service"
	"snakers-bot/internal/usecases"
	"strconv"
	"strings"
	"time"
)

const (
	buttonCatalog = "Каталог 👟"
	buttonOrders  = "Мои заказы 📦"
	buttonProfile = "Профиль 👤"
	buttonBalance = "Баланс 💳"

	callbackBuy   = "buy_"
	callbackTopUp = "top_up"
)

var mainKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(buttonCatalog),
		tgbotapi.NewKeyboardButton(buttonOrders),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(buttonProfile),
		tgbotapi.NewKeyboardButton(buttonBalance),
	),
)

type BotHandler struct {
	bot             *tgbotapi.BotAPI
	userService     *service.UserService
	productService  *service.ProductService
	orderService    *service.OrderService
	paymentProvider interfaces.PaymentProvider
}

func NewBotHandler(bot *tgbotapi.BotAPI, us *service.UserService, ps *service.ProductService, os *service.OrderService, pp interfaces.PaymentProvider) *BotHandler {
	return &BotHandler{
		bot:             bot,
		userService:     us,
		productService:  ps,
		orderService:    os,
		paymentProvider: pp,
	}
}

func (h *BotHandler) HandleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		ctx := context.Background()

		if update.CallbackQuery != nil {
			h.handleCallbackQuery(ctx, update.CallbackQuery)
			continue
		}

		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				h.handleStart(ctx, update.Message)
			case buttonCatalog, "/catalog":
				h.handleCatalog(ctx, update.Message)
			case buttonOrders:
				h.handleMyOrders(ctx, update.Message)
			case buttonProfile:
				h.handleProfile(ctx, update.Message)
			case buttonBalance:
				h.handleBalance(ctx, update.Message)
			default:
				h.handleDefault(update.Message)
			}
		}
	}
}

func (h *BotHandler) handleCatalog(ctx context.Context, message *tgbotapi.Message) {
	products, err := h.productService.GetAllProducts(ctx)
	if err != nil {
		h.sendErrorMessage(message.Chat.ID, "Не удалось загрузить каталог товаров.")
		return
	}
	if len(products) == 0 {
		h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Каталог пока пуст."))
		return
	}

	h.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Вот что у нас есть:"))

	for _, p := range products {
		fileBytes, err := downloadFile(p.ImageURL)
		if err != nil {
			log.Printf("Failed to download image %s: %v", p.ImageURL, err)
			continue
		}

		photoFile := tgbotapi.FileBytes{Name: p.Name + ".jpg", Bytes: fileBytes}
		photoMsg := tgbotapi.NewPhoto(message.Chat.ID, photoFile)

		photoMsg.Caption = fmt.Sprintf(
			"*%s*\n\n%s\n\nЦена: *%s руб.*",
			p.Name,
			p.Description,
			p.Price.StringFixed(2),
		)
		photoMsg.ParseMode = tgbotapi.ModeMarkdown

		buyButton := tgbotapi.NewInlineKeyboardButtonData("Купить 🛒", fmt.Sprintf("%s%d", callbackBuy, p.ID))
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(buyButton),
		)
		photoMsg.ReplyMarkup = &inlineKeyboard

		if _, err := h.bot.Send(photoMsg); err != nil {
			log.Printf("Failed to send product photo for product ID %d: %v", p.ID, err)
		}

		time.Sleep(300 * time.Millisecond)
	}
}
func (h *BotHandler) handleStart(ctx context.Context, message *tgbotapi.Message) {
	if _, err := h.getOrCreateUser(ctx, message); err != nil {
		h.sendErrorMessage(message.Chat.ID, "Произошла ошибка при создании вашего профиля.")
		return
	}
	responseText := fmt.Sprintf("Привет, %s! Я бот для продажи самых топовых шузов. Выбери, что тебя интересует 👇", message.From.FirstName)
	msg := tgbotapi.NewMessage(message.Chat.ID, responseText)
	msg.ReplyMarkup = mainKeyboard
	h.bot.Send(msg)
}

func (h *BotHandler) handleBalance(ctx context.Context, message *tgbotapi.Message) {
	user, err := h.getOrCreateUser(ctx, message)
	if err != nil {
		h.sendErrorMessage(message.Chat.ID, "Не удалось найти ваш профиль.")
		return
	}
	responseText := fmt.Sprintf("На вашем балансе: *%s руб.*", user.Balance.StringFixed(2))
	msg := tgbotapi.NewMessage(message.Chat.ID, responseText)
	msg.ParseMode = tgbotapi.ModeMarkdown
	topUpButton := tgbotapi.NewInlineKeyboardButtonData("Пополнить", callbackTopUp)
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(topUpButton))
	msg.ReplyMarkup = &inlineKeyboard
	h.bot.Send(msg)
}

func (h *BotHandler) handleCallbackQuery(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	h.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
	if strings.HasPrefix(callback.Data, callbackBuy) {
		h.handleBuyCallback(ctx, callback)
	} else if callback.Data == callbackTopUp {
		h.handleTopUpCallback(ctx, callback)
	}
}

func (h *BotHandler) handleTopUpCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	dummyOrder := &usecases.Order{ID: 0}
	paymentDetails, err := h.paymentProvider.GeneratePaymentDetails(ctx, dummyOrder)
	if err != nil {
		log.Printf("Failed to generate payment details: %v", err)
		h.sendErrorMessage(callback.Message.Chat.ID, "Не удалось получить реквизиты для оплаты.")
		return
	}
	customText := "Для пополнения баланса, совершите перевод по одному из следующих реквизитов:\n\n" +
		strings.SplitN(paymentDetails, "Криптовалюты:", 2)[1]
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, customText)
	msg.ParseMode = tgbotapi.ModeMarkdown
	h.bot.Send(msg)
}

func (h *BotHandler) handleBuyCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	productIDStr := strings.TrimPrefix(callback.Data, callbackBuy)
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		h.sendErrorMessage(callback.Message.Chat.ID, "Ошибка: неверный ID товара.")
		return
	}
	user, err := h.getOrCreateUser(ctx, callback.Message)
	if err != nil {
		h.sendErrorMessage(callback.Message.Chat.ID, "Не удалось найти ваш профиль.")
		return
	}
	product, err := h.productService.GetProductByID(ctx, uint(productID))
	if err != nil {
		h.sendErrorMessage(callback.Message.Chat.ID, "Товар не найден.")
		return
	}
	newOrder := &usecases.Order{
		UserID:   user.ID,
		Status:   usecases.StatusAwaiting,
		Date:     time.Now(),
		Products: []usecases.Product{*product},
	}
	if err := h.orderService.CreateOrder(ctx, newOrder); err != nil {
		h.sendErrorMessage(callback.Message.Chat.ID, "Не удалось создать заказ.")
		return
	}
	responseText := fmt.Sprintf("✅ Отлично! Заказ *#%d* на кроссовки *%s* создан и ожидает оплаты.", newOrder.ID, product.Name)
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, responseText)
	msg.ParseMode = tgbotapi.ModeMarkdown
	h.bot.Send(msg)
}

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

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
