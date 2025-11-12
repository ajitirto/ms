package domain

import "time"

type Product struct {
	ProductID   int       `json:"product_id"`  // SERIAL PRIMARY KEY
	CustomerID  string    `json:"customer_id"` // VARCHAR(100) NOT NULL
	ProductDate time.Time `json:"product_date"`
	TotalAmount float64   `json:"total_amount"` // NUMERIC(10, 2) NOT NULL
}
