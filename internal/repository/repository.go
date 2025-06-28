package repository

import (
	"context"
	models "snakers-bot/internal/usecases"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
}

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id uint) (*models.Product, error)
	GetAll(ctx context.Context) ([]models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id uint) error
}

type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, id uint) (*models.Order, error)
	GetByUserID(ctx context.Context, userID uint) ([]models.Order, error)
	Update(ctx context.Context, order *models.Order) error
	Delete(ctx context.Context, id uint) error
}
