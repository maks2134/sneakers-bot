package usecases

import (
	"github.com/shopspring/decimal"
	"time"
)

type StatusOrder string

const (
	StatusArrive   StatusOrder = "В доставке"
	StatusAwaiting StatusOrder = "Ожидает оплаты"
	StatusComplete StatusOrder = "Приехал"
)

type Order struct {
	ID       uint        `gorm:"primaryKey"`
	Status   StatusOrder `gorm:"type:varchar(20)"`
	Date     time.Time
	UserID   uint
	Products []Product `gorm:"many2many:order_products;"`
}

func (o *Order) TotalPrice() decimal.Decimal {
	total := decimal.NewFromInt(0)
	for _, p := range o.Products {
		total = total.Add(p.Price)
	}
	return total
}
