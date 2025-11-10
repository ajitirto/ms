package repository

import "github.com/ajitirto/ms/order-service/internal/domain"

// OrderRepository mendefinisikan kontrak untuk operasi penyimpanan data Order
type OrderRepository interface {
    GetByID(id int) (domain.Order, error)
    Create(order domain.Order) (domain.Order, error)
    GetAll() ([]domain.Order, error)
}