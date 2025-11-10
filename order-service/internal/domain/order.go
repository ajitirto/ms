package domain

import "time"

type Order struct {
    OrderID      int       `json:"order_id"`      // SERIAL PRIMARY KEY
    CustomerID   int       `json:"customer_id"`   // INTEGER NOT NULL
    OrderDate    time.Time `json:"order_date"`    // DATE NOT NULL
    TotalAmount  float64   `json:"total_amount"`  // NUMERIC(10, 2) NOT NULL
}

type OrderCreationRequest struct {
    CustomerID  int     `json:"customer_id"`
    TotalAmount float64 `json:"total_amount"`
}