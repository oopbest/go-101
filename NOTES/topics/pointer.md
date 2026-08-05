# การใช้งาน Pointer ในภาษา Go (Pointer in Go)

**Pointer (พอยน์เตอร์)** คือตัวแปรประเภทหนึ่งในภาษา Go ที่ใช้เก็บ **ตำแหน่งที่อยู่บนหน่วยความจำ (Memory Address)** ของตัวแปรอื่น แทนที่จะเก็บค่าข้อมูลโดยตรง การใช้งาน Pointer ช่วยให้เราสามารถอ้างอิงและแก้ไขข้อมูลต้นทางได้โดยตรงโดยไม่ต้องคัดลอกข้อมูล (Pass by Reference)

---

## 1. สัญลักษณ์สำคัญของ Pointer (Key Operators)

| เครื่องหมาย | ชื่อเรียก | คำอธิบายการใช้งาน | ตัวอย่าง |
| :--- | :--- | :--- | :--- |
| **`&`** | **Address-of Operator** | ใช้ดึง **ตำแหน่งที่อยู่หน่วยความจำ (Memory Address)** ของตัวแปร | `p := &x` |
| **`*`** | **Dereference Operator** /<br>**Pointer Type** | 1. **ใช้ระบุชนิดตัวแปร Pointer**: เช่น `v *int`<br>2. **ใช้ดึง/แก้ไข ค่าที่อยู่ในที่อยู่นั้น**: เช่น `*p = 77` | `var p *int`<br>`*p = 77` |

> [!NOTE]
> **Automatic Dereferencing ใน Go**: สำหรับ Struct Pointer เช่น `e *Employee` เราสามารถเข้าถึงฟิลด์ได้โดยตรงด้วย `e.Salary` โดยไม่จำเป็นต้องเขียน `(*e).Salary` เพราะ Go Compiler จะจัดการดึงค่าจาก Pointer ให้อัตโนมัติ

---

## 2. การทำงานของ Pointer ในหน่วยความจำ (Memory Concept)

เมื่อเรากำหนด `y := &x` ตัวแปร `y` จะเก็บตำแหน่งหน่วยความจำของ `x` (เช่น `0xc0000a2008`)

```text
   [ตัวแปร x]  --->  ตำแหน่งในหน่วยความจำ: 0xc0000a2008  --->  ค่าข้างใน: 50
   [ตัวแปร y]  --->  เก็บตำแหน่ง: 0xc0000a2008 (ชี้ไปยัง x)
```

เมื่อสั่งคำสั่ง `*y = 77` คอมพิวเตอร์จะเดินตามที่อยู่ใน `y` ไปยังตำแหน่ง `0xc0000a2008` แล้วเปลี่ยนค่าข้างในกล่องนั้นเป็น `77` ส่งผลให้เมื่ออ่านค่าของ `x` และ `*y` จะได้ค่า `77` เช่นเดียวกัน

---

## 3. ตัวอย่างโค้ดประกอบการใช้งาน (Code Examples)

### 3.1 ตัวอย่างที่ 1: การใช้ Pointer กับชนิดข้อมูลพื้นฐาน (`int`)
```go
package main

import "fmt"

// 1. ฟังก์ชันรับ Parameter เป็น Pointer (Pass by Reference)
func changeValue(v *int) {
	*v = 50 // เปลี่ยนค่าในตำแหน่งหน่วยความจำที่ v ชี้อยู่ให้เป็น 50
}

func main() {
	x := 3
	fmt.Println("x เริ่มต้น:", x) // Output: 3

	// 2. ส่ง Memory Address ของ x (&x) เข้าไปในฟังก์ชัน
	changeValue(&x)
	fmt.Println("x หลังเรียก changeValue:", x) // Output: 50

	// 3. ประกาศตัวแปร Pointer y ชี้ไปยังตำแหน่งของ x
	y := &x
	fmt.Println("ตำแหน่งที่อยู่หน่วยความจำของ x (y):", y) // Output: 0xc000... (Memory Address)
	fmt.Println("ค่าที่ y ชี้อยู่ (*y):", *y)                 // Output: 50

	// 4. แก้ไขค่าผ่าน Pointer y (Dereferencing)
	*y = 77
	fmt.Println("x หลังถูกเปลี่ยนผ่าน *y:", x)            // Output: 77
	fmt.Println("ค่าที่ y ชี้อยู่ (*y):", *y)               // Output: 77
}
```

### 3.2 ตัวอย่างที่ 2: การใช้ Pointer ร่วมกับ Struct (`Employee`)
```go
package main

import "fmt"

type Employee struct {
	Name   string
	Age    int
	Salary float64
}

// ฟังก์ชันรับ Struct Pointer e *Employee เพื่อแก้ไขเงินเดือนของ Employee ต้นทาง
func giveRaise(e *Employee, percent float64) {
	// ใช้ e.Salary ได้เลยโดยไม่ต้องใช้ (*e).Salary (Automatic Dereferencing)
	e.Salary += e.Salary * (percent / 100)
}

func main() {
	emp := Employee{
		Name:   "John",
		Age:    32,
		Salary: 67500,
	}

	fmt.Println("Before:", emp) // Output: Before: {John 32 67500}

	// ส่ง Memory Address ของ emp (&emp) เพื่อให้ฟังก์ชันปรับแก้ไขเงินเดือนต้นทาง
	giveRaise(&emp, 15)

	fmt.Println("After:", emp)  // Output: After: {John 32 77625}
}
```

---

## 4. ประโยชน์และช่วงเวลาที่ควรเลือกใช้ Pointer

1. **ต้องการให้ฟังก์ชันแก้ไขข้อมูลต้นทางได้ (Mutate State)**: เพื่อให้การเปลี่ยนแปลงข้อมูลภายในฟังก์ชันส่งผลต่อตัวแปรหรือ Struct ภายนอก
2. **ประสิทธิภาพและการประหยัดหน่วยความจำ (Performance Optimization)**: หลีกเลี่ยงการคัดลอกข้อมูลขนาดใหญ่ (เช่น Struct ที่มีหลายฟิลด์) เมื่อต้องส่งเข้าฟังก์ชัน
3. **การใช้งานร่วมกัน (Shared Reference)**: เพื่อให้อยู่บนข้อมูลชุดเดียวกันจากหลายๆ จุดของโค้ด
