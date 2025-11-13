package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ajitirto/ms/product-service/internal/usecase"
)

type ProductHandler struct {
	ProductUC usecase.ProductUsecase
}

func NewProductHandler(uc usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{ProductUC: uc}
}

// StatusHandler untuk health check
func (h *ProductHandler) StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product Service is running and ready."))
}

// HandleProducts menangani GET /products
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	h.getAllProducts(w, r)
}

func (h *ProductHandler) getAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.ProductUC.GetAllProducts(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}
	idStr := parts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	product, err := h.ProductUC.GetProduct(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to retrieve product", http.StatusInternalServerError)
		return
	}
	if product.ProductID == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}