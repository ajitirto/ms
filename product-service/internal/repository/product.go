package repository

import "github.com/ajitirto/ms/product-service/internal/domain"


type ProductRepository interface {
	GetByID(id int) (domain.Product, error)
	GetAll() ([]domain.Product, error)
}