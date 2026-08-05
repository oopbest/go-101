# ⚙️ การใช้งาน Stored Procedure & Function ใน Go (Stored Procedure & Function in Go)

สรุปการสร้าง **Stored Procedure / Function** ใน PostgreSQL (PL/pgSQL) และวิธีการเรียกใช้งานจากภาษา Go ด้วยแพ็กเกจมาตรฐาน `database/sql`

---

## 1. 🎯 ความแตกต่างระหว่าง Stored Procedure และ Function ใน PostgreSQL

| คุณลักษณะ (Feature) | Function (`CREATE FUNCTION`) | Stored Procedure (`CREATE PROCEDURE`) |
| :--- | :--- | :--- |
| **การเรียกใช้งานใน SQL** | เรียกผ่าน `SELECT function_name(...)` | เรียกผ่าน `CALL procedure_name(...)` |
| **การคืนค่า (Return Value)** | **ต้อง** คืนค่าเสมอ (`RETURNS scalar / TABLE / VOID`) | ไม่บังคับให้คืนค่า (สามารถใช้ `INOUT` / `OUT` parameters ได้) |
| **การควบคุม Transaction** | ไม่สามารถใช้ `COMMIT` / `ROLLBACK` ภายในได้ | **สามารถใช้ `COMMIT` / `ROLLBACK` ภายในขั้นตอน** ได้ (PostgreSQL 11+) |
| **การใช้ใน Query** | สามารถนำไปต่อใน `WHERE`, `JOIN` หรือ `SELECT` ได้ | ไม่สามารถนำไปซ้อนในคำสั่ง `SELECT` อื่นได้ |

---

## 2. 📝 การสร้าง Function & Stored Procedure ใน PostgreSQL

### 1) การสร้าง PL/pgSQL Function (คำนวณและคืนค่า)

```sql
-- Function คำนวณราคารวมสินค้าพร้อมภาษี (VAT 7%)
CREATE OR REPLACE FUNCTION calculate_total_with_vat(p_price NUMERIC, p_qty INT)
RETURNS NUMERIC AS $$
BEGIN
    RETURN (p_price * p_qty) * 1.07;
END;
$$ LANGUAGE plpgsql;

-- Function ดึงรายการสินค้าของ Supplier (คืนค่าเป็น TABLE)
CREATE OR REPLACE FUNCTION get_products_by_supplier(p_supplier_id INT)
RETURNS TABLE (product_id INT, product_name VARCHAR, price NUMERIC, stock INT) AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, p.name, p.price, p.stock
    FROM products p
    WHERE p.supplier_id = p_supplier_id
    ORDER BY p.id;
END;
$$ LANGUAGE plpgsql;
```

### 2) การสร้าง Stored Procedure (ย้าย Stock ระหว่าง Supplier)

```sql
CREATE OR REPLACE PROCEDURE transfer_product_supplier(
    p_product_id INT,
    p_new_supplier_id INT
)
AS $$
BEGIN
    UPDATE products
    SET supplier_id = p_new_supplier_id
    WHERE id = p_product_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Product ID % not found', p_product_id;
    END IF;
END;
$$ LANGUAGE plpgsql;
```

---

## 3. 💻 การเรียกใช้งานจากภาษา Go (`database/sql`)

### 1) การเรียก Function ที่คืนค่าเดียว (Scalar Value)

ใช้คำสั่ง `SELECT function_name($1, $2)` ร่วมกับ `db.QueryRow().Scan()`:

```go
func GetTotalWithVat(db *sql.DB, price float64, qty int) (float64, error) {
	var total float64
	query := `SELECT calculate_total_with_vat($1, $2)`
	err := db.QueryRow(query, price, qty).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to call function: %w", err)
	}
	return total, nil
}
```

### 2) การเรียก Function ที่คืนค่าเป็นตาราง (Table Result)

ใช้คำสั่ง `SELECT * FROM function_name($1)` ร่วมกับ `db.Query()` และ `for rows.Next()`:

```go
type ProductSimple struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

func GetProductsBySupplierFunc(db *sql.DB, supplierID int) ([]ProductSimple, error) {
	query := `SELECT * FROM get_products_by_supplier($1)`
	rows, err := db.Query(query, supplierID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []ProductSimple
	for rows.Next() {
		var p ProductSimple
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
```

### 3) การเรียก Stored Procedure (`CALL`)

ใช้คำสั่ง `CALL procedure_name($1, $2)` ร่วมกับ `db.Exec()`:

```go
func TransferProductSupplier(db *sql.DB, productID, newSupplierID int) error {
	query := `CALL transfer_product_supplier($1, $2)`
	_, err := db.Exec(query, productID, newSupplierID)
	if err != nil {
		return fmt.Errorf("failed to call procedure: %w", err)
	}
	return nil
}
```

---

## 4. ⚖️ ข้อดีและข้อเสียของการใช้ Stored Procedure

| ข้อดี (Pros) | ข้อเสีย (Cons) |
| :--- | :--- |
| **Performance สูง**: ประมวลผลบน Database Server โดยตรง ไม่ต้องส่งข้อมูลไปกลับหลายรอบข้าม Network | **Vendor Lock-in**: ย้ายไปใช้ Database Engine อื่นยาก (เช่น ย้ายจาก Postgres ไป MySQL ต้องเขียนใหม่) |
| **Security**: เพิ่มความปลอดภัยโดยการจำกัดสิทธิ์ให้ User เรียกใช้เฉพาะ Procedure แทนการเข้าถึง Table ตรงๆ | **Version Control & Maintenance ยาก**: จัดการและ Debug โค้ดใน DB ยากกว่าโค้ดฝั่ง Go |
| **Reduce Business Logic Duplication**: แชร์ logic ให้แอปหลายภาษา (Go, Node.js, Python) ใช้ร่วมกันได้ | **Testing ยาก**: ไม่สามารถทำ Unit Test แบบ Mock ได้ง่ายเท่า Business Logic ฝั่ง Application |

---

## 5. ⚠️ ข้อควรจำสำคัญ (Best Practices)

> [!IMPORTANT]
> 1. **การเรียกใช้ใน Go**:
>    - **Function**: เรียกด้วย `SELECT function_name(...)` หรือ `SELECT * FROM function_name(...)`
>    - **Procedure**: เรียกด้วย `CALL procedure_name(...)` ร่วมกับ `db.Exec()`
> 2. **Parameterized Query**: ส่งพารามิเตอร์แบบ `$1`, `$2` เสมอ เพื่อป้องกัน SQL Injection
> 3. **Error Handling**: หาก Procedure มีการโยน Error ออกมาด้วย `RAISE EXCEPTION` ใน PostgreSQL ฝั่ง Go จะได้รับ Error นั้นในตัวแปร `err` ทันที
