package usecase

import (
	"github.com/ajitirto/ms/product-service/internal/repository"
	"github.com/ajitirto/ms/product-service/internal/domain"
)

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductUsecase {
	return &productService{productRepo: repo}
}

func (s *productService) GetProduct(id int) (domain.Product, error) {
	return s.productRepo.GetByID(id)
}
func (s *productService) GetAllProducts() ([]domain.Product, error) {
	return s.productRepo.GetAll()
}