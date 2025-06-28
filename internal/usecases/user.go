package usecases

import "github.com/shopspring/decimal"

type User struct {
	ID            uint  `gorm:"primaryKey"`
	TelegramID    int64 `gorm:"uniqueIndex;column:telegram_id"`
	Name          string
	Balance       decimal.Decimal `gorm:"type:numeric"`
	LoyaltyPoints int64
	Orders        []Order `gorm:"foreignKey:UserID"`
}
