package usecases

import "github.com/shopspring/decimal"

type Product struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Price    decimal.Decimal `gorm:"type:numeric"`
	ImageURL string
}
