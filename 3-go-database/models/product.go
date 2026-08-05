package models

import "time"

type Product struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Price      float64   `json:"price"`
	Stock      int       `json:"stock"`
	SupplierID *int      `json:"supplier_id"` // Pointer for handling NULL supplier_id
	CreatedAt  time.Time `json:"created_at"`
}

// ProductWithSupplier for JOIN queries
type ProductWithSupplier struct {
	Product
	SupplierName *string `json:"supplier_name"`
}
