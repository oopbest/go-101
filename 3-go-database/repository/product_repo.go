package repository

import (
	"database/sql"
	"fmt"

	"github.com/oopbest/go-test-db/models"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) Create(p *models.Product) (int, error) {
	query := `INSERT INTO products (name, price, stock, supplier_id) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	err := r.DB.QueryRow(query, p.Name, p.Price, p.Stock, p.SupplierID).Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		return 0, err
	}
	return p.ID, nil
}

func (r *ProductRepository) GetByID(id int) (*models.ProductWithSupplier, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.supplier_id, p.created_at, s.name as supplier_name
		FROM products p
		LEFT JOIN suppliers s ON p.supplier_id = s.id
		WHERE p.id = $1`
	var p models.ProductWithSupplier
	err := r.DB.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.Price, &p.Stock, &p.SupplierID, &p.CreatedAt, &p.SupplierName,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) GetAll() ([]models.ProductWithSupplier, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.supplier_id, p.created_at, s.name as supplier_name
		FROM products p
		LEFT JOIN suppliers s ON p.supplier_id = s.id
		ORDER BY p.id`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.ProductWithSupplier
	for rows.Next() {
		var p models.ProductWithSupplier
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.SupplierID, &p.CreatedAt, &p.SupplierName); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) Update(p *models.Product) error {
	query := `UPDATE products SET name = $1, price = $2, stock = $3, supplier_id = $4 WHERE id = $5`
	res, err := r.DB.Exec(query, p.Name, p.Price, p.Stock, p.SupplierID, p.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", p.ID)
	}
	return nil
}

func (r *ProductRepository) Delete(id int) error {
	query := `DELETE FROM products WHERE id = $1`
	res, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", id)
	}
	return nil
}
