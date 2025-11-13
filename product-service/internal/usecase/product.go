package usecase

import (
	"context"

	"github.com/ajitirto/ms/product-service/internal/domain"
)

type ProductUsecase interface {
	GetProduct(ctx context.Context, id int) (*domain.Product, error)
	GetAllProducts(ctx context.Context) ([]domain.Product, error)
}