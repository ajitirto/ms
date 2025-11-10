package server

import (
    "encoding/json"
    "net/http"
    "strconv"
    "github.com/ajitirto/ms/order-service/internal/domain"
    "github.com/ajitirto/ms/order-service/internal/usecase"
)

type OrderHandler struct {
    OrderUC usecase.OrderUsecase
}

func NewOrderHandler(uc usecase.OrderUsecase) *OrderHandler {
    return &OrderHandler{OrderUC: uc}
}

// StatusHandler untuk health check
func (h *OrderHandler) StatusHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Order Service is running and ready."))
}

// HandleOrders menangani GET /orders dan POST /orders
func (h *OrderHandler) HandleOrders(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        h.getAllOrders(w, r)
    case http.MethodPost:
        h.createNewOrder(w, r)
    default:
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
    }
}

func (h *OrderHandler) getAllOrders(w http.ResponseWriter, r *http.Request) {
    orders, err := h.OrderUC.GetAllOrders()
    if err != nil {
        http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) createNewOrder(w http.ResponseWriter, r *http.Request) {
    var req domain.OrderCreationRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    order, err := h.OrderUC.CreateNewOrder(req)
    if err != nil {
        // Cek jika errornya adalah business rule error
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(order)
}

// GetOrderByID handles GET /orders/{id}
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
    // Parsing ID dari path (misalnya /orders/123)
    idStr := r.URL.Path[len("/orders/"):]
    orderID, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid Order ID format", http.StatusBadRequest)
        return
    }

    order, err := h.OrderUC.GetOrder(orderID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound) // Asumsi error dari repo adalah not found
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(order)
}