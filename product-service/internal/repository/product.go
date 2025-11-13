package repository

import (
	"context"
	"errors"

	"github.com/ajitirto/ms/product-service/internal/domain"
)

type ProductRepository interface {
	GetByID(ctx context.Context,id int) (*domain.Product, error)
	GetAll(ctx context.Context) ([]domain.Product, error)
}

type ProductCache interface {
	Get(ctx context.Context, id int) (*domain.Product, error)
	Set(ctx context.Context,  product *domain.Product) error

	GetAll(ctx context.Context, key string) ([]domain.Product, error)
    SetAll(ctx context.Context, key string, products []domain.Product) error
}

var ErrProductNotFound = errors.New("product not found")

