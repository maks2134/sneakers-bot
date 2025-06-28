package usecases

import "time"

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
