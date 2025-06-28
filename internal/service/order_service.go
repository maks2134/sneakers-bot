package service

import (
	"context"
	"snakers-bot/internal/repository"
	"snakers-bot/internal/usecases"
)

type OrderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *usecases.Order) error {
	return s.repo.Create(ctx, order)
}

func (s *OrderService) GetOrderByID(ctx context.Context, id uint) (*usecases.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrderService) GetUserOrders(ctx context.Context, userID uint) ([]usecases.Order, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *OrderService) UpdateOrder(ctx context.Context, order *usecases.Order) error {
	return s.repo.Update(ctx, order)
}

func (s *OrderService) DeleteOrder(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
