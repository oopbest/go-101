# การใช้งาน GORM ร่วมกับ PostgreSQL (GORM in Go)

GORM เป็น Object-Relational Mapping (ORM) สำหรับภาษา Go ช่วยให้จัดการข้อมูลผ่าน Struct และ Method ของ Go โดยยังสามารถใช้ SQL, Transaction และ Connection Pool จาก `database/sql` ได้เมื่อจำเป็น

---

## 1. การติดตั้งและเชื่อมต่อ PostgreSQL

ติดตั้ง GORM และ PostgreSQL driver:

```bash
go get gorm.io/gorm
go get gorm.io/driver/postgres
```

ตัวอย่างการเชื่อมต่อ:

```go
package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := "host=localhost port=5432 user=postgres password=postgres123 dbname=godb sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	return db, nil
}
```

> **หมายเหตุ**
> `*gorm.DB` เป็นตัวใช้งาน GORM ส่วน Connection Pool จริงยังบริหารโดย `*sql.DB` ที่อยู่ด้านล่าง

ตั้งค่า Connection Pool ได้ดังนี้:

```go
sqlDB, err := db.DB()
if err != nil {
	return err
}

sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

---

## 2. การสร้าง Models และกำหนด Tags

ตัวอย่าง Models ที่สัมพันธ์กันแบบ Supplier หนึ่งรายมี Product ได้หลายรายการ:

```go
type Supplier struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Phone     string    `gorm:"size:20" json:"phone"`
	Email     string    `gorm:"size:100;uniqueIndex" json:"email"`
	Products  []Product `gorm:"foreignKey:SupplierID" json:"products,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Product struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Name         string     `gorm:"size:100;not null" json:"name"`
	Price        float64    `gorm:"type:numeric(10,2);not null;default:0" json:"price"`
	Stock        int        `gorm:"not null;default:0" json:"stock"`
	SupplierID   *uint      `json:"supplier_id"`
	Supplier     *Supplier  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"supplier,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
```

Tags ที่พบบ่อย:

| Tag | ความหมาย |
| :--- | :--- |
| `primaryKey` | กำหนด Primary Key |
| `not null` | ห้ามเก็บค่า `NULL` |
| `size:100` | กำหนดขนาดคอลัมน์ข้อความ |
| `uniqueIndex` | สร้าง Unique Index |
| `default:0` | กำหนดค่าเริ่มต้น |
| `column:product_name` | กำหนดชื่อคอลัมน์เอง |
| `foreignKey:SupplierID` | ระบุ Foreign Key ของความสัมพันธ์ |

---

## 3. Auto Migration

```go
func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&Supplier{},
		&Product{},
	)
}
```

`AutoMigrate()` เหมาะกับการทดลองและโปรเจกต์ขนาดเล็ก เพราะช่วยสร้างตาราง คอลัมน์ Index และ Constraints จาก Models

สำหรับ Production ควรใช้ migration files ที่มี version เช่น `golang-migrate` เพื่อให้ตรวจสอบลำดับการเปลี่ยน Schema และ Rollback ได้ชัดเจน

---

## 4. CRUD Operations

### 4.1 Create

```go
product := Product{
	Name:       "Laptop Pro 15",
	Price:      45000,
	Stock:      10,
	SupplierID: &supplierID,
}

result := db.Create(&product)
if result.Error != nil {
	return result.Error
}

fmt.Println("created product ID:", product.ID)
```

### 4.2 Read One และ Read All

```go
var product Product
result := db.First(&product, id)

if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	return fmt.Errorf("product not found")
}
if result.Error != nil {
	return result.Error
}
```

```go
var products []Product
if err := db.Order("id ASC").Find(&products).Error; err != nil {
	return err
}
```

### 4.3 Update

```go
result := db.Model(&Product{}).
	Where("id = ?", id).
	Updates(map[string]any{
		"price": 42900,
		"stock": 8,
	})

if result.Error != nil {
	return result.Error
}
if result.RowsAffected == 0 {
	return fmt.Errorf("product not found")
}
```

`Updates(struct)` จะข้าม Zero Value เช่น `0`, `false` และ `""` โดยค่าเริ่มต้น หากต้องการอัปเดตค่าเหล่านี้ให้ใช้ `map[string]any` หรือระบุ Fields ด้วย `Select()`

### 4.4 Delete

```go
result := db.Delete(&Product{}, id)
if result.Error != nil {
	return result.Error
}
if result.RowsAffected == 0 {
	return fmt.Errorf("product not found")
}
```

ควรมี `WHERE` หรือ Primary Key ทุกครั้ง เพราะ GORM ป้องกัน Global Update และ Global Delete โดยค่าเริ่มต้น

---

## 5. การโหลดข้อมูลที่สัมพันธ์กัน

ใช้ `Preload()` เมื่อต้องการโหลด Supplier พร้อม Product:

```go
var product Product
err := db.
	Preload("Supplier").
	First(&product, id).
	Error
```

โหลด Supplier ทั้งหมดพร้อมรายการ Products:

```go
var suppliers []Supplier
err := db.
	Preload("Products").
	Find(&suppliers).
	Error
```

`Preload()` มักรัน Query เพิ่มสำหรับ Association ส่วน `Joins()` เหมาะกับกรณีที่ต้องการ JOIN เพื่อกรองหรือเลือกข้อมูลใน Query เดียว ควรตรวจ SQL และจำนวน Query เมื่อข้อมูลมีปริมาณมาก

---

## 6. Transaction

ใช้ `Transaction()` เมื่อหลายคำสั่งต้องสำเร็จหรือยกเลิกพร้อมกัน:

```go
err := db.Transaction(func(tx *gorm.DB) error {
	if err := tx.Model(&Product{}).
		Where("id = ? AND stock >= ?", productID, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity)).
		Error; err != nil {
		return err
	}

	order := Order{
		ProductID: productID,
		Quantity:  quantity,
	}
	if err := tx.Create(&order).Error; err != nil {
		return err
	}

	return nil // Commit
})
```

ถ้าฟังก์ชันคืน `error` GORM จะ Rollback แต่ถ้าคืน `nil` จะ Commit

---

## 7. การใช้งานร่วมกับ context.Context

ใน REST API ควรส่ง Context จาก HTTP Request เข้า Query:

```go
func (r *ProductRepository) GetByID(
	ctx context.Context,
	id uint,
) (*Product, error) {
	var product Product
	err := r.DB.
		WithContext(ctx).
		Preload("Supplier").
		First(&product, id).
		Error

	return &product, err
}
```

เมื่อ Client ยกเลิก Request หรือหมดเวลา Database Query จะสามารถถูกยกเลิกตาม Context ได้

---

## 8. แนวปฏิบัติที่ดี

1. ตรวจ `result.Error` ทุกครั้งหลังเรียก GORM
2. ตรวจ `RowsAffected` สำหรับ Update และ Delete เพื่อแยกกรณีไม่พบข้อมูล
3. ใช้ `context.Context` ใน Repository ของ REST API
4. ใช้ Transaction เมื่อต้องแก้หลายตารางเป็นหนึ่งหน่วยงาน
5. หลีกเลี่ยงการส่ง Database Model รับ JSON โดยตรง ควรแยก Request/Response DTO
6. ใช้ `Preload()` เท่าที่จำเป็น เพื่อป้องกัน Query และข้อมูลตอบกลับมากเกินไป
7. เปิด SQL Logger ใน Development และลดระดับ Log ใน Production
8. ใช้ migration files ที่มี version สำหรับ Production แทนการพึ่ง `AutoMigrate()` เพียงอย่างเดียว
9. ใช้ Parameter Binding เช่น `Where("id = ?", id)` ไม่ต่อ SQL จาก input โดยตรง
10. ดึง `*sql.DB` จาก `db.DB()` เมื่อต้องตั้งค่า Connection Pool หรือปิดฐานข้อมูล

---

## เอกสารอ้างอิง

- [GORM Guides](https://gorm.io/docs/)
- [Connecting to PostgreSQL](https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL)
- [CRUD Interface](https://gorm.io/docs/create.html)
- [Associations](https://gorm.io/docs/associations.html)
- [Transactions](https://gorm.io/docs/transactions.html)
- [Context](https://gorm.io/docs/context.html)
