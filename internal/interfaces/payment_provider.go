package interfaces

import (
	"context"
	entity "snakers-bot/internal/usecases"
)

type PaymentProvider interface {
	GeneratePaymentDetails(ctx context.Context, order *entity.Order) (string, error)
}
