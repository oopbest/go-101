# การจัดการข้อผิดพลาดในภาษา Go (Error Handling in Go)

ในภาษา Go จะไม่มีระบบ **Exception (try-catch-finally)** เหมือนภาษาโปรแกรมอื่น แต่ Go ใช้วิธี **Explicit Error Handling** โดยฟังก์ชันที่อาจเกิดข้อผิดพลาดจะส่งคืนค่าผลลัพธ์คู่กับอินเทอร์เฟซ `error` เป็นค่าสุดท้ายเสมอ หากไม่มีข้อผิดพลาดเกิดขึ้นจะส่งคืนค่า `nil`

---

## 1. แนวคิดพื้นฐาน (Basic Concepts)

### 1.1 `error` Interface ใน Go
`error` ในภาษา Go เป็น Built-in Interface ที่มีข้อกำหนดง่ายๆ เพียงเมธอดเดียวดังนี้:
```go
type error interface {
    Error() string
}
```

### 1.2 การตรวจสอบ Error (`if err != nil`)
รูปแบบมาตรฐาน (Idiomatic Go) ในการจัดการ Error คือการตรวจสอบว่า `err != nil` หรือไม่ทันทีหลังเรียกใช้งานฟังก์ชัน:
```go
result, err := divide(10, 2)
if err != nil {
    // จัดการข้อผิดพลาดเมื่อเกิด error
    fmt.Println("Error:", err)
    return
}
// ทำงานต่อหากไม่มี error
fmt.Println("Result:", result)
```

---

## 2. การสร้าง Error ขึ้นใช้งาน (Creating Errors)

### 2.1 การใช้ `errors.New()`
เหมาะสำหรับสร้างข้อผิดพลาดที่มีข้อความคงที่ สั้นกระชับ:
```go
import "errors"

var ErrDivideByZero = errors.New("cannot divide by zero")
```

### 2.2 การใช้ `fmt.Errorf()`
เหมาะสำหรับสร้างข้อผิดพลาดที่ต้องการใส่ตัวแปรหรือฟอร์แมตข้อความแบบไดนามิก:
```go
import "fmt"

func validateAge(age int) error {
    if age < 0 {
        return fmt.Errorf("invalid age %d: age cannot be negative", age)
    }
    return nil
}
```

### 2.3 การสร้าง Custom Error ชนิด Struct
เมื่อต้องการเก็บบริบทหรือข้อมูลของ Error เพิ่มเติม (เช่น Username, Error Code, Timestamp) เราสามารถสร้าง Struct และเขียนเมธอด `Error() string` เพื่อทำตามสัญญาของ `error` interface ได้:

```go
type LoginError struct {
    Username string
    Message  string
}

// เมธอดนี้ทำให้ *LoginError ทำหน้าที่เป็น error interface
func (e *LoginError) Error() string {
    return fmt.Sprintf("Login error from username %s : %s", e.Username, e.Message)
}
```

---

## 3. ตัวอย่างโค้ดประกอบการใช้งาน (Code Examples)

### 3.1 ตัวอย่างที่ 1: การส่งคืน Error พื้นฐานด้วย `errors.New()`
```go
package main

import (
	"errors"
	"fmt"
)

// ฟังก์ชันหารตัวเลข ส่งคืนผลลัพธ์ (float64) และ error
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("cannot divide by zero") // ส่งคืน error เมื่อตัวหารเป็น 0
	}
	return a / b, nil // ส่งคืน nil เมื่อคำนวณสำเร็จไม่มีข้อผิดพลาด
}

func main() {
	res, err := divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err) // Output: Error: cannot divide by zero
	} else {
		fmt.Println("Result:", res)
	}
}
```

### 3.2 ตัวอย่างที่ 2: การใช้ Custom Error Struct (`LoginError`)
```go
package main

import "fmt"

// 1. ประกาศ Custom Struct สำหรับเก็บรายละเอียด Error
type LoginError struct {
	Username string
	Message  string
}

// 2. สร้างเมธอด Error() string เพื่อให้ LoginError ทำหน้าที่เป็น error interface
func (e *LoginError) Error() string {
	return fmt.Sprintf("Login error from username %s : %s", e.Username, e.Message)
}

// 3. ฟังก์ชันตรวจสอบการเข้าสู่ระบบ
func login(username, password string) error {
	if username != "admin" || password != "123456" {
		return &LoginError{
			Username: username,
			Message:  "Invalid credentials",
		}
	}
	return nil // คืนค่า nil แสดงว่าเข้าสู่ระบบสำเร็จ
}

func main() {
	// ทดลองเรียกใช้ด้วยข้อมูลที่ไม่ถูกต้อง
	err := login("testuser", "56789")
	if err != nil {
		fmt.Println(err) 
		// Output: Login error from username testuser : Invalid credentials
	} else {
		fmt.Println("Login successful")
	}
}
```

---

## 4. เทคนิคและแนวปฏิบัติที่ดี (Best Practices)

| หัวข้อ | คำอธิบาย |
| :--- | :--- |
| **Return Error Last** | กำหนดให้ค่า `error` เป็น Return Value ตัวสุดท้ายของฟังก์ชันเสมอ เช่น `(int, error)` |
| **Check Immediately** | ควรตรวจสอบ `if err != nil` ทันทีหลังเรียกฟังก์ชันที่อาจเกิดข้อผิดพลาด |
| **Early Return (Guard Clause)** | หากเกิด Error ให้จัดการแล้ว `return` ทันที เพื่อหลีกเลี่ยงการเขียน `if-else` ซ้อนกันหลายชั้น |
| **Custom Error Types** | สร้าง Struct ทำตาม `error` interface เพื่อเก็บข้อมูลบริบทของ Error เพิ่มเติม |
| **Panic & Recover** | ใช้เฉพาะกรณีวิกฤตที่โปรแกรมไม่สามารถรันต่อได้จริงๆ ไม่ควรใช้ `panic` แทนการคืนค่า `error` ทั่วไป |
