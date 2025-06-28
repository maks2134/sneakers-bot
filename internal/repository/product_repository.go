package repository

import (
	"context"
	"gorm.io/gorm"
	models "snakers-bot/internal/usecases"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{db: db}
}

func (r productRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r productRepository) GetByID(ctx context.Context, id uint) (*models.Product, error) {
	var product models.Product
	return &product, r.db.WithContext(ctx).First(&product, id).Error
}

func (r productRepository) Update(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r productRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	return products, r.db.WithContext(ctx).Find(&products).Error
}

func (r productRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}
