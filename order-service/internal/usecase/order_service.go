package usecase

import (
    "time"
    "github.com/ajitirto/ms/order-service/internal/domain"
    "github.com/ajitirto/ms/order-service/internal/repository"
	"fmt"
)

type orderService struct {
    orderRepo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderUsecase {
    return &orderService{orderRepo: repo}
}

func (s *orderService) GetOrder(id int) (domain.Order, error) {
    // Contoh business logic: validasi ID
    if id <= 0 {
        return domain.Order{}, fmt.Errorf("invalid order ID")
    }
    return s.orderRepo.GetByID(id)
}

func (s *orderService) CreateNewOrder(req domain.OrderCreationRequest) (domain.Order, error) {
    // 1. Business Rule: Total amount tidak boleh nol atau negatif
    if req.TotalAmount <= 0 {
        return domain.Order{}, fmt.Errorf("total amount must be positive")
    }
    
    // 2. Mapping Request ke Domain Entity
    newOrder := domain.Order{
        CustomerID:  req.CustomerID,
        OrderDate:   time.Now(), // Menyetel tanggal saat ini
        TotalAmount: req.TotalAmount,
    }
    
    // 3. Panggil Repository
    return s.orderRepo.Create(newOrder)
}

func (s *orderService) GetAllOrders() ([]domain.Order, error) {
    return s.orderRepo.GetAll()
}