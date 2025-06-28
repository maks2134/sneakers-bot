package service

import (
	"context"
	"snakers-bot/internal/repository"
	"snakers-bot/internal/usecases"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *usecases.User) error {
	return s.repo.Create(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*usecases.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, user *usecases.User) error {
	return s.repo.Update(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserService) GetUserByTelegramID(ctx context.Context, telegramID int64) (*usecases.User, error) {
	return s.repo.GetByTelegramID(ctx, telegramID)
}
