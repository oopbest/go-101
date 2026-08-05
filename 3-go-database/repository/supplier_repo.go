package repository

import (
	"database/sql"
	"fmt"

	"github.com/oopbest/go-test-db/models"
)

type SupplierRepository struct {
	DB *sql.DB
}

func NewSupplierRepository(db *sql.DB) *SupplierRepository {
	return &SupplierRepository{DB: db}
}

func (r *SupplierRepository) Create(s *models.Supplier) (int, error) {
	query := `INSERT INTO suppliers (name, phone, email) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := r.DB.QueryRow(query, s.Name, s.Phone, s.Email).Scan(&s.ID, &s.CreatedAt)
	if err != nil {
		return 0, err
	}
	return s.ID, nil
}

func (r *SupplierRepository) GetByID(id int) (*models.Supplier, error) {
	query := `SELECT id, name, phone, email, created_at FROM suppliers WHERE id = $1`
	var s models.Supplier
	err := r.DB.QueryRow(query, id).Scan(&s.ID, &s.Name, &s.Phone, &s.Email, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SupplierRepository) GetAll() ([]models.Supplier, error) {
	query := `SELECT id, name, phone, email, created_at FROM suppliers ORDER BY id`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []models.Supplier
	for rows.Next() {
		var s models.Supplier
		if err := rows.Scan(&s.ID, &s.Name, &s.Phone, &s.Email, &s.CreatedAt); err != nil {
			return nil, err
		}
		suppliers = append(suppliers, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return suppliers, nil
}

func (r *SupplierRepository) Update(s *models.Supplier) error {
	query := `UPDATE suppliers SET name = $1, phone = $2, email = $3 WHERE id = $4`
	res, err := r.DB.Exec(query, s.Name, s.Phone, s.Email, s.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("supplier with id %d not found", s.ID)
	}
	return nil
}

func (r *SupplierRepository) Delete(id int) error {
	query := `DELETE FROM suppliers WHERE id = $1`
	res, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("supplier with id %d not found", id)
	}
	return nil
}
