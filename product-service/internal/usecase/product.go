package usecase

import (
	"github.com/ajitirto/ms/product-service/internal/domain"
)

type ProductUsecase interface {
	GetProduct(id int) (domain.Product, error)
	GetAllProducts() ([]domain.Product, error)
}