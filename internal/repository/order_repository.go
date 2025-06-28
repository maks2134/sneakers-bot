package repository

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	models "snakers-bot/internal/usecases"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) orderRepository {
	return orderRepository{db: db}
}

func (r orderRepository) Create(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r orderRepository) GetByID(ctx context.Context, id uint) (*models.Order, error) {
	var order models.Order
	return &order, r.db.WithContext(ctx).First(&order, id).Error
}

func (r orderRepository) Update(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

func (r orderRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Order{}, id).Error
}

func (r *orderRepository) GetByUserID(ctx context.Context, userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.WithContext(ctx).Preload(clause.Associations).Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}
