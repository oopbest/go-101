# 🐘 การใช้งาน Go database/sql & PostgreSQL (Go Database/SQL)

สรุปการใช้งานแพ็กเกจมาตรฐาน `database/sql` ร่วมกับ PostgreSQL Driver (`github.com/lib/pq`) การจัดโครงสร้างโค้ดด้วย Repository Pattern และเทคนิคการแก้ปัญหาข้อผิดพลาดยอดฮิตในการพัฒนาโปรเจกต์ Go + Database

---

## 1. ⚙️ การตั้งค่าและการเชื่อมต่อฐานข้อมูล (Database Connection)

การเชื่อมต่อ PostgreSQL ใน Go จะใช้แพ็กเกจ `database/sql` ร่วมกับ Driver `github.com/lib/pq` โดยทำการโหลดค่าคอนฟิกจากไฟล์ `.env` ผ่าน `godotenv`

```go
package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL Driver (Anonymous Import)
)

func ConnectDB() (*sql.DB, error) {
	_ = godotenv.Load()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// ทดสอบการเชื่อมต่อกับฐานข้อมูล
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
```

> [!NOTE]
> `_ "github.com/lib/pq"` เป็นการนำเข้าแพ็กเกจแบบ Blank Identifier เพื่อให้ฟังก์ชัน `init()` ของ Driver ทำการลงทะเบียน (Register) เข้ากับ `database/sql`

---

## 2. 🏗️ การสร้าง Schema และตารางข้อมูล (Table Initialization)

การสร้างตารางและกำหนดความสัมพันธ์ (Foreign Key Constraint) ด้วย `db.Exec()`:

```go
func InitTables(db *sql.DB) error {
	supplierTableSQL := `
	CREATE TABLE IF NOT EXISTS suppliers (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		phone VARCHAR(20),
		email VARCHAR(100),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	productTableSQL := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
		stock INT NOT NULL DEFAULT 0,
		supplier_id INT REFERENCES suppliers(id) ON DELETE SET NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(supplierTableSQL); err != nil {
		return fmt.Errorf("failed to create suppliers table: %w", err)
	}
	if _, err := db.Exec(productTableSQL); err != nil {
		return fmt.Errorf("failed to create products table: %w", err)
	}
	return nil
}
```

---

## 3. 📂 การออกแบบโครงสร้างโปรเจกต์ (Repository Pattern & Constructor)

การแยกส่วนการทำงานออกเป็น Package ตามหลัก Clean Architecture เพื่อความง่ายในการดูแลรักษาและทดสอบ:

```text
3-go-database/
├── config/
│   └── db.go             # จัดการ Connection และ Schema Initialization
├── models/
│   ├── supplier.go       # Struct ของ Supplier
│   └── product.go        # Struct ของ Product และ ProductWithSupplier
├── repository/
│   ├── supplier_repo.go  # CRUD Operations สำหรับ Supplier
│   └── product_repo.go   # CRUD Operations สำหรับ Product
├── .env                  # Environment Variables
└── main.go               # Entry Point สั่งรันโปรเจกต์
```

### การสร้าง Constructor และ Dependency Injection
การใช้ดีไซน์แบบ **Dependency Injection** ช่วยให้ทุก Repository สามารถใช้วัตถุ `*sql.DB` จาก `main.go` ร่วมกันได้:

```go
type ProductRepository struct {
	DB *sql.DB
}

// Constructor Function สไตล์ Go (ขึ้นต้นด้วย New...)
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}
```

---

## 4. 📊 เปรียบเทียบ 3 คำสั่งหลักในการรัน SQL ใน Go

| คำสั่ง (Method) | การใช้งาน (Use Case) | คืนค่า (Return) | ตัวอย่างการนำไปใช้ |
| :--- | :--- | :--- | :--- |
| **`db.Exec()`** | คำสั่งที่ไม่ต้องการผลลัพธ์ข้อมูลกลับมา | `sql.Result`, `error` | `INSERT` (ปกติ), `UPDATE`, `DELETE`, `CREATE TABLE` |
| **`db.Query()`** | คำสั่งที่ดึงข้อมูล **หลายแถว (Multiple Rows)** | `*sql.Rows`, `error` | `SELECT` รายการทั้งหมด (`for rows.Next()`) |
| **`db.QueryRow()`** | คำสั่งที่ดึงข้อมูล **แถวเดียว (Single Row)** | `*sql.Row` | `SELECT WHERE id = $1`, `INSERT ... RETURNING id` |

---

## 5. 🔄 การทำ CRUD Operations (โค้ดตัวอย่าง)

### 1) CREATE & RETURNING ID (`QueryRow`)

เมื่อต้องการเพิ่มข้อมูลและรับค่า `id` ที่ถูกสร้างแบบ Auto-increment (SERIAL) กลับมา:

```go
func (r *SupplierRepository) Create(s *models.Supplier) (int, error) {
	query := `INSERT INTO suppliers (name, phone, email) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := r.DB.QueryRow(query, s.Name, s.Phone, s.Email).Scan(&s.ID, &s.CreatedAt)
	if err != nil {
		return 0, err
	}
	return s.ID, nil
}
```

### 2) READ ALL & JOIN (`Query`)

การอ่านข้อมูลหลายแถวร่วมกับการเชื่อมตาราง (`LEFT JOIN`) และการตรวจสอบ Data Stream Error:

```go
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

	// ⚠️ ตรวจสอบ error ที่อาจเกิดขึ้นระหว่างวนลูปอ่าน rows Stream
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
```

### 3) UPDATE & DELETE (`Exec`)

คำสั่งแก้ไขและลบข้อมูลใช้ `db.Exec()` และตรวจสอบผลลัพธ์การอัปเดตผ่าน `RowsAffected()`:

```go
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
```

---

## 6. 💡 เทคนิคการจัดการ Pointer และการพิมพ์ผลลัพธ์ (`fmt.Printf`)

ฟิลด์ที่ยอมให้เป็น `NULL` ในตาราง (เช่น Foreign Key `supplier_id`) ใน Go จะนิยมใช้ชนิดข้อมูลเป็น **Pointer** เช่น `*int` หรือ `*string`

### การแสดงผลลัพธ์ที่ปลอดภัย (Safely Printing Pointers)
การนำ `*string` หรือ `*int` ไปใส่ใน `fmt.Printf` โดยตรงจะเกิด Warning หรือพิมพ์ Memory Address ออกมา จึงต้องตรวจสอบ `nil` และ Dereference ค่าเสมอ:

```go
for _, p := range products {
	supID := 0
	if p.SupplierID != nil {
		supID = *p.SupplierID // Dereference ค่า int จาก Pointer
	}
	supName := "N/A"
	if p.SupplierName != nil {
		supName = *p.SupplierName // Dereference ค่า string จาก Pointer
	}

	fmt.Printf("  - ID: %d, Name: %s, Price: %.2f THB, Stock: %d, Supplier ID: %d, Supplier Name: %s\n",
		p.ID, p.Name, p.Price, p.Stock, supID, supName)
}
```

---

## 7. 🛠️ สรุปข้อผิดพลาดพบบ่อย และแนวทางแก้ไข (Troubleshooting)

### 1. `sql.Rows "rows" is used in Next loop without final check of rows.Err()`
* **สาเหตุ**: `rows.Next()` จะคืนค่า `false` ทั้งเมื่ออ่านจบปกติ และเมื่อเกิด Network/Stream Error ระหว่างทาง
* **แก้ไข**: ต้องเขียน `if err := rows.Err(); err != nil` ปิดท้ายลูป `for rows.Next()` เสมอ

### 2. `pq: insert or update on table violates foreign key constraint (23503)`
* **สาเหตุ**: สั่งเพิ่มหรืออัปเดตข้อมูลโดยอ้างอิง Foreign Key ID (เช่น `supplier_id`) ที่ไม่มีอยู่จริงในตารางอ้างอิง
* **แก้ไข**: ตรวจสอบว่า ID ตารางแม่นั้นมีอยู่จริงใน Database ก่อนทำรายการ

### 3. `undefined: productId` หรือ `declared and not used: productId`
* **สาเหตุ**: ประกาศตัวแปรในหัวเงื่อนไข `if productId, err := ...; err != nil` ทำให้ตัวแปร `productId` มีตัวตนแค่ภายในบล็อก `if` เท่านั้น
* **แก้ไข**: ประกาศตัวแปรรับค่านอกบล็อก `if` ก่อนเรียกเช็ค `err`

### 4. `cannot use ... (value of type *models.ProductWithSupplier) as *models.Product value in assignment`
* **สาเหตุ**: Go เป็นภาษา Strictly Typed ไม่ยอมให้นำ Struct ต่างชนิดกันมา Re-assign ใส่ตัวแปรเดิม
* **แก้ไข**: ประกาศตัวแปรรับค่าใหม่ด้วยชื่อแยกต่างหาก เช่น `singleProduct`

---

## 8. ⚠️ ข้อควรจำสำคัญ (Best Practices)

> [!IMPORTANT]
> 1. **การปิด Connection (`defer rows.Close()`)**: ทุกครั้งที่ใช้ `db.Query()` ต้องสั่ง `defer rows.Close()` เพื่อคืน Resource ให้กับ Connection Pool
> 2. **การเช็ค `rows.Err()`**: ต้องตรวจสอบ `rows.Err()` หลังจบลูป `for rows.Next()` ทุกครั้ง
> 3. **การจัดการค่า `NULL`**: ใช้ Pointer (`*int`, `*string`) ใน Struct เพื่อรับค่า `NULL` ป้องกัน `Scan error`
> 4. **Parameterized Queries ($1, $2)**: ใช้การส่งค่าแบบ `$1`, `$2` เสมอ เพื่อป้องกันปัญหา **SQL Injection**
