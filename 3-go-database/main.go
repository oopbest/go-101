package main

import (
	"fmt"
	"log"

	"github.com/oopbest/go-test-db/config"
	"github.com/oopbest/go-test-db/repository"
)

func main() {
	// 1. Connect to PostgreSQL Database
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close()
	fmt.Println("Successfully connected to PostgreSQL database!")

	// 2. Initialize Tables
	if err := config.InitTables(db); err != nil {
		log.Fatalf("Failed to initialize tables: %v", err)
	}
	fmt.Println("Tables initialized successfully.")

	// 3. Initialize Repositories
	supplierRepo := repository.NewSupplierRepository(db)
	productRepo := repository.NewProductRepository(db)

	// ==========================================
	// 1. SUPPLIER CRUD
	// ==========================================
	// DELETE Supplier
	// idToDelete := 8
	// if err := supplierRepo.Delete(idToDelete); err != nil {
	// 	log.Fatalf("Error deleting supplier: %v \n", err)
	// }
	// fmt.Printf("[DELETE] Supplier with ID %d deleted successfully.\n", idToDelete)

	// CREATE Supplier
	// supplier := &models.Supplier{
	// 	Name:  "Apple Thailand",
	// 	Phone: "088-444-5555",
	// 	Email: "apple@test.com",
	// }
	// id, err := supplierRepo.Create(supplier)
	// if err != nil {
	// 	log.Fatalf("Error creating supplier: %v", err)
	// }
	// fmt.Printf("[CREATE] Supplier created with ID: %d\n", id)

	// UPDATE Supplier
	// supplier := &models.Supplier{
	// 	ID:    9,
	// 	Name:  "Asus Technology",
	// 	Phone: "081-222-3333",
	// 	Email: "asus@gmail.com",
	// }
	// if err := supplierRepo.Update(supplier); err != nil {
	// 	log.Fatalf("Error updating supplier: %v", err)
	// }
	// fmt.Printf("[UPDATE] Supplier with ID %d updated successfully.\n", supplier.ID)

	// READ All Suppliers
	suppliers, err := supplierRepo.GetAll()
	if err != nil {
		log.Fatalf("Error getting all suppliers: %v", err)
	}
	fmt.Printf("[READ ALL] Total Suppliers: %d\n", len(suppliers))
	for _, s := range suppliers {
		fmt.Printf("  - ID: %d, Name: %s, Phone: %s, Email: %s\n", s.ID, s.Name, s.Phone, s.Email)
	}

	// ==========================================
	// 2. PRODUCT CRUD
	// ==========================================

	// DELETE Product
	// idToDelete := 4
	// if err := productRepo.Delete(idToDelete); err != nil {
	// 	log.Fatalf("Error deleting product: %v \n", err)
	// }
	// fmt.Printf("[DELETE] Product with ID %d deleted successfully.\n", idToDelete)

	// CREATE Product
	// supplierID1 := 1
	// product := &models.Product{
	// 	Name:       "Macbook Pro 14 M3",
	// 	Price:      52900.00,
	// 	Stock:      44,
	// 	SupplierID: &supplierID1,
	// }
	// if _, err := productRepo.Create(product); err != nil {
	// 	log.Fatalf("Error creating product: %v\n", err)
	// }
	// fmt.Printf("[CREATE] Product created with ID: %d\n", product.ID)

	// UPDATE Product and supplier
	// supplierID := 10
	// updateProd := &models.Product{
	// 	ID:         17,
	// 	Name:       "Macbook Pro M5",
	// 	Price:      99000.00,
	// 	Stock:      69,
	// 	SupplierID: &supplierID,
	// }
	// if err := productRepo.Update(updateProd); err != nil {
	// 	log.Fatalf("Error updating product: %v", err)
	// }
	// fmt.Printf("[UPDATE] Product ID %d updated: Name=%s, Price=%.2f, Stock=%d, SupplierID=%d\n",
	// 	updateProd.ID, updateProd.Name, updateProd.Price, updateProd.Stock, *updateProd.SupplierID)

	// READ All Products (With Supplier Name via JOIN)
	products, err := productRepo.GetAll()
	if err != nil {
		log.Fatalf("Error getting all products: %v", err)
	}
	fmt.Printf("[READ ALL With Supplier Name] Total Products: %d\n", len(products))
	for _, p := range products {
		supName := "N/A"
		if p.SupplierName != nil {
			supName = *p.SupplierName
		}
		fmt.Printf("  - ID: %d, Name: %s, Price: %.2f THB, Stock: %d, Supplier: %s\n",
			p.ID, p.Name, p.Price, p.Stock, supName)
	}

	// GET PRODUCT BY ID
	// singleProduct, err := productRepo.GetByID(1)
	// if err != nil {
	// 	log.Fatalf("Error getting product by ID: %v", err)
	// }
	// fmt.Printf("[READ ONE] Found Product ID %d: %s, Price: %.2f THB\n", singleProduct.ID, singleProduct.Name, singleProduct.Price)
}
