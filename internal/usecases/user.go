package usecases

import "github.com/shopspring/decimal"

type User struct {
	ID            uint `gorm:"primaryKey"`
	Name          string
	Balance       decimal.Decimal `gorm:"type:numeric"`
	LoyaltyPoints int64
	Orders        []Order `gorm:"foreignKey:UserID"`
}
