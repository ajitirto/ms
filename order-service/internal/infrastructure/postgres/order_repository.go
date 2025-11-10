package postgres

import (
    "database/sql"
    "fmt"
    "github.com/ajitirto/ms/order-service/internal/domain"
    "github.com/ajitirto/ms/order-service/internal/repository"
)

type OrderRepository struct {
    db *sql.DB
}

func NewOrderRepository(db *sql.DB) repository.OrderRepository {
    return &OrderRepository{db: db}
}

func (r *OrderRepository) GetByID(id int) (domain.Order, error) {
    var order domain.Order
    query := `SELECT order_id, customer_id, order_date, total_amount FROM orders WHERE order_id = $1`
    
    err := r.db.QueryRow(query, id).Scan(
        &order.OrderID,
        &order.CustomerID,
        &order.OrderDate,
        &order.TotalAmount,
    )
    if err == sql.ErrNoRows {
        return domain.Order{}, fmt.Errorf("order with ID %d not found", id)
    }
    return order, err
}

func (r *OrderRepository) Create(order domain.Order) (domain.Order, error) {
    query := `INSERT INTO orders (customer_id, order_date, total_amount) VALUES ($1, $2, $3) RETURNING order_id`
    
    // Catatan: OrderDate dari Domain sudah disetel di Usecase
    err := r.db.QueryRow(query, 
        order.CustomerID, 
        order.OrderDate, 
        order.TotalAmount,
    ).Scan(&order.OrderID) // Mengambil ID yang baru dibuat

    if err != nil {
        return domain.Order{}, fmt.Errorf("failed to create order: %w", err)
    }
    return order, nil
}

func (r *OrderRepository) GetAll() ([]domain.Order, error) {
    query := `SELECT order_id, customer_id, order_date, total_amount FROM orders`
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to query orders: %w", err)
    }
    defer rows.Close()

    var orders []domain.Order
    for rows.Next() {
        var order domain.Order
        if err := rows.Scan(
            &order.OrderID,
            &order.CustomerID,
            &order.OrderDate,
            &order.TotalAmount,
        ); err != nil {
            return nil, fmt.Errorf("failed to scan order row: %w", err)
        }
        orders = append(orders, order)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error during rows iteration: %w", err)
    }
    return orders, nil
}