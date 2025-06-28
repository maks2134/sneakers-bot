package usecases

import "time"

type StatusOrder string

const (
	StatusArrive   StatusOrder = "В доставке"
	StatusAwaiting StatusOrder = "Ожидает оплаты"
	StatusComplete StatusOrder = "Приехал"
)

type Order struct {
	OrderId int64              `json:"order_id"`
	Status  string             `json:"status"`
	Date    time.Time          `json:"date"`
	Product map[string]Product `json:"product"`
}
