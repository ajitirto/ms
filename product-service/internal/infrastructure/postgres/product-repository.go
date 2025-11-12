package postgres


import (	
	"database/sql"
	"github.com/ajitirto/ms/product-service/internal/domain"
	"github.com/ajitirto/ms/product-service/internal/repository"
)

type ProductRepository struct {
	db *sql.DB
}	

func NewProductRepository(db *sql.DB) repository.ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetByID(id int) (domain.Product, error) {
	var product domain.Product
	query := `SELECT product_id, customer_id, product_date, total_amount FROM products WHERE product_id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&product.ProductID,
		&product.CustomerID,
		&product.ProductDate,
		&product.TotalAmount,
	)
	if err == sql.ErrNoRows {
		return domain.Product{}, nil
	}
	return product, err
}
func (r *ProductRepository) GetAll() ([]domain.Product, error) {
	query := `SELECT product_id, customer_id, product_date, total_amount FROM products`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		err := rows.Scan(
			&product.ProductID,
			&product.CustomerID,
			&product.ProductDate,
			&product.TotalAmount,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}