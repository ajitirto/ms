package usecase

import "github.com/ajitirto/ms/order-service/internal/domain"

type OrderUsecase interface {
    GetOrder(id int) (domain.Order, error)
    CreateNewOrder(req domain.OrderCreationRequest) (domain.Order, error)
    GetAllOrders() ([]domain.Order, error)
}