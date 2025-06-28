package service

import (
	"context"
	"snakers-bot/internal/repository"
	"snakers-bot/internal/usecases"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, product *usecases.Product) error {
	return s.repo.Create(ctx, product)
}

func (s *ProductService) GetProductByID(ctx context.Context, id uint) (*usecases.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]usecases.Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *usecases.Product) error {
	return s.repo.Update(ctx, product)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
