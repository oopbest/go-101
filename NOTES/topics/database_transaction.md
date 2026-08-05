# 💳 การจัดการ Database Transaction ใน Go (Database Transaction in Go)

สรุปการใช้งาน **Database Transaction** ในภาษา Go ด้วยแพ็กเกจมาตรฐาน `database/sql` เพื่อรักษารวมความถูกต้องของข้อมูลตามหลัก ACID การจัดการ Deadlock และเทคนิคระบบ Retry

---

## 1. 🎯 แนวคิดหลักของ Database Transaction (ACID Properties)

Transaction คือการรวมชุดคำสั่ง SQL หลาย ๆ คำสั่งเข้าด้วยกันเป็น **หน่วยการทำงานเดียว (Single Unit of Work)** โดยปฏิบัติตามหลัก ACID:

* **Atomicity (All or Nothing)**: ต้องสำเร็จทั้งหมดทุกขั้นตอน หรือหากเกิดความผิดพลาดแม้แต่ขั้นตอนเดียว จะต้องยกเลิกคืนค่ากลับทั้งหมด (Rollback)
* **Consistency**: ข้อมูลต้องถูกต้องตามกฎเกณฑ์และ Constraints ของฐานข้อมูลทั้งก่อนและหลังทำรายการ
* **Isolation**: การทำงานของแต่ละ Transaction ที่เกิดขึ้นพร้อมกัน ต้องแยกออกจากกัน ไม่ส่งผลรบกวนกัน
* **Durability**: เมื่อทำการบันทึก (`Commit`) แล้ว ข้อมูลจะถูกบันทึกลงฐานข้อมูลอย่างถาวร

---

## 2. 🛠️ รูปแบบการใช้งานสไตล์ Go (Safe Transaction Pattern)

ใน Go จะเปิด Transaction ด้วย `db.Begin()` หรือ `db.BeginTx()` และใช้เทคนิค **`defer tx.Rollback()`** เพื่อความปลอดภัย:

```go
func TransferMoney(db *sql.DB, fromAccID, toAccID int, amount float64) error {
	// 1. เริ่มต้น Transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// 2. ป้องกัน Transaction ค้างด้วย defer tx.Rollback()
	// หากสั่ง tx.Commit() สำเร็จแล้ว สั่ง tx.Rollback() ใน defer จะไม่มีผลใดๆ (Safely ignored)
	defer tx.Rollback()

	// 3. ขั้นตอนที่ 1: ตัดเงินจากบัญชีผู้โอน (ใช้ tx.Exec ไม่ใช่ db.Exec)
	res1, err := tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND balance >= $1", amount, fromAccID)
	if err != nil {
		return fmt.Errorf("failed to deduct balance: %w", err)
	}
	rows1, _ := res1.RowsAffected()
	if rows1 == 0 {
		return fmt.Errorf("insufficient balance or account %d not found", fromAccID)
	}

	// 4. ขั้นตอนที่ 2: เพิ่มเงินเข้าบัญชีผู้รับ
	res2, err := tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, toAccID)
	if err != nil {
		return fmt.Errorf("failed to add balance: %w", err)
	}
	rows2, _ := res2.RowsAffected()
	if rows2 == 0 {
		return fmt.Errorf("recipient account %d not found", toAccID)
	}

	// 5. หากทุกขั้นตอนสำเร็จ สั่ง ยืนยันการทำรายการ (Commit)
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
```

---

## 3. 🛍️ ตัวอย่างระบบจริง: การสร้าง Order และหัก Stock สินค้า

ตัวอย่างการบันทึกรายการคำสั่งซื้อ (Create Order) พร้อมหักจำนวนสินค้าคงคลัง (Deduct Stock):

```go
type OrderReq struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

func CreateOrderTx(db *sql.DB, req OrderReq) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// 1. หัก Stock สินค้า (ตรวจสอบสต็อกเพียงพอก่อน)
	deductSQL := `UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1`
	res, err := tx.Exec(deductSQL, req.Quantity, req.ProductID)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil || affected == 0 {
		return 0, fmt.Errorf("product out of stock or product not found")
	}

	// 2. สร้างแถวรายการคำสั่งซื้อใหม่ (Order Record)
	var orderID int
	createOrderSQL := `INSERT INTO orders (product_id, quantity, status) VALUES ($1, $2, 'PAID') RETURNING id`
	err = tx.QueryRow(createOrderSQL, req.ProductID, req.Quantity).Scan(&orderID)
	if err != nil {
		return 0, err
	}

	// 3. ยืนยัน Transaction
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return orderID, nil
}
```

---

## 4. ⚙️ การตั้งค่า Isolation Level ด้วย `db.BeginTx()`

หากต้องการกำหนดระดับความเป็นอิสระของการทำรายการ (Isolation Level) หรือใส่ `context.Context` ให้ใช้ `db.BeginTx()`:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// กำหนด Isolation Level เป็น Read Committed (หรือ RepeatableRead / Serializable)
opts := &sql.TxOptions{
	Isolation: sql.LevelReadCommitted,
	ReadOnly:  false,
}

tx, err := db.BeginTx(ctx, opts)
if err != nil {
	return err
}
defer tx.Rollback()
```

---

## 5. 🔒 ปัญหา Deadlock ใน Database Transaction และแนวทางแก้ไข

### Deadlock คืออะไร?
Deadlock เกิดขึ้นเมื่อ Transaction ตั้งแต่ 2 ตัวขึ้นไปมีการแย่งชิง Lock ข้อมูลในลักษณะเป็นวงกลม (Circular Wait) โดยต่างฝ่ายต่างรอให้อีกฝ่ายคืน Lock ข้อมูลที่ตนเองต้องการ ทำให้ระบบติดค้าง (Stuck) ไม่สามารถทำงานต่อได้

#### ตัวอย่างสถานการณ์ Deadlock (Classic Scenario):
- **Transaction A**: ล็อคบัญชี 1 ➡️ กำลังรอเพื่อล็อคบัญชี 2
- **Transaction B**: ล็อคบัญชี 2 ➡️ กำลังรอเพื่อล็อคบัญชี 1
- **ผลลัพธ์ใน PostgreSQL**: จะเกิดข้อผิดพลาด Code `40P01` (`deadlock detected`) และระบบจะสั่ง Abort สุ่มยกเลิก Transaction ตัวใดตัวหนึ่งทันที

---

### แนวทางป้องกันและการแก้ไข Deadlock (Deadlock Prevention Strategies)

#### 1) การจัดเรียงลำดับการ Lock ให้เหมือนกันเสมอ (Consistent Lock Ordering)
บังคับให้ทุก Transaction อัปเดตหรือเข้าถึงข้อมูลตามลำดับ ID จากน้อยไปมากเสมอ ไม่ว่าจะโอนเงินจากใครไปหาใคร:

```go
func TransferMoneySafe(db *sql.DB, fromID, toID int, amount float64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// จัดเรียง ID จากน้อยไปมากก่อนอัปเดตเสมอ เพื่อป้องกันการติด Deadlock แบบ Circular Wait
	firstID, secondID := fromID, toID
	if firstID > secondID {
		firstID, secondID = toID, fromID
	}

	// ล็อคและอัปเดตเรียงตาม ID
	if _, err := tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", getAmount(firstID, fromID, amount), firstID); err != nil {
		return err
	}
	if _, err := tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", getAmount(secondID, fromID, amount), secondID); err != nil {
		return err
	}

	return tx.Commit()
}
```

#### 2) การทำระบบ Retry Mechanism (เมื่อเจอ Deadlock Error `40P01`)
เมื่อเจอ Deadlock ฝั่ง Application ควรทำการสั่งรันลองใหม่ (Retry) โดยอัตโนมัติ:

```go
func ExecTxWithRetry(db *sql.DB, maxRetries int, fn func(tx *sql.Tx) error) error {
	for i := 0; i < maxRetries; i++ {
		tx, err := db.Begin()
		if err != nil {
			return err
		}

		err = fn(tx)
		if err == nil {
			return tx.Commit()
		}

		tx.Rollback()
		// หากเป็น Deadlock Error (40P01) ให้รอเวลาแล้ววนลูปทำงานใหม่ (Retry)
		if strings.Contains(err.Error(), "40P01") || strings.Contains(err.Error(), "deadlock") {
			time.Sleep(time.Duration(i+1) * 50 * time.Millisecond) // Exponential Backoff
			continue
		}

		return err // Error อื่นๆ ที่ไม่ใช่ Deadlock คืนค่าออกไปทันที
	}
	return fmt.Errorf("transaction failed after max retries")
}
```

---

## 6. ⚠️ ข้อควรระวังและแนวปฏิบัติที่ดี (Best Practices)

> [!IMPORTANT]
> 1. **เรียกใช้เมธอดผ่านวัตถุ `tx` เสมอ**: ต้องใช้ `tx.Exec()`, `tx.Query()`, หรือ `tx.QueryRow()` แทน `db.Exec()` / `db.Query()` มิฉะนั้นคำสั่งนั้นจะอยู่นอก Transaction
> 2. **ใช้ `defer tx.Rollback()` เสมอ**: วางไว้บรรทัดถัดจาก `db.Begin()` ทันที เพื่อการันตีว่าหากเกิด error หรือ `panic` กลางคราว ระบบจะ Rollback คืนค่าข้อมูลให้อัตโนมัติ
> 3. **อย่าเปิด Transaction ค้างไว้นานเกินไป**: ไม่ควรใส่คำสั่งที่ใช้เวลานาน (เช่น การส่ง Email, การยิง HTTP Request ภายนอก) ไว้ในระหว่าง Transaction เพราะจะทำให้เกิดการ Lock แถวในตาราง DB นานเกินไป
> 4. **จัดลำดับการเรียง ID**: เมื่อทำ Transaction ที่มีการอัปเดตหลายแถว ควรเรียงลำดับ ID จากน้อยไปมากเสมอเพื่อป้องกัน Deadlock
