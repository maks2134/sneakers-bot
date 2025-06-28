package repository

import (
	"context"
	"gorm.io/gorm"
	models "snakers-bot/internal/usecases"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepo) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	return &user, r.db.WithContext(ctx).First(&user, id).Error
}

func (r *UserRepo) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *UserRepo) GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	var user models.User
	return &user, r.db.WithContext(ctx).Where("telegram_id = ?", telegramID).First(&user).Error
}
