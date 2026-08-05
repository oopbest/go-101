# เปรียบเทียบ Function vs Method ในภาษา Go (Golang)

## 1. ตารางเปรียบเทียบข้อแตกต่าง (Comparison Table)

| หัวข้อ (Feature) | Function (ฟังก์ชัน) | Method (เมธอด) |
| :--- | :--- | :--- |
| **การผูกกับ Type** | เป็นอิสระ ไม่ผูกติดกับ Type หรือ Struct ใดๆ | ต้องผูกติดกับ Type (Receiver Type) เช่น Struct |
| **การเรียกใช้งาน (Invocation)** | เรียกผ่านชื่อฟังก์ชันโดยตรง เช่น `sum(10, 20)` | เรียกผ่านอินสแตนซ์ของ Type (Dot Notation) เช่น `student.FullName()` |
| **Syntax การประกาศ** | `func Name(params) ReturnType` | `func (r ReceiverType) Name(params) ReturnType` |
| **Receiver** | ไม่มี Receiver | มี Receiver (แบ่งเป็น Value Receiver และ Pointer Receiver) |
| **การเข้าถึงข้อมูลภายใน Struct** | ต้องส่ง Struct เข้ามาทาง Parameter | สามารถเข้าถึงฟิลด์ของ Receiver ได้โดยตรงผ่านตัวแปร Receiver |
| **วัตถุประสงค์ (Use Case)** | ใช้สำหรับประมวลผลทั่วไป (Utility/Generic logic) | ใช้สำหรับแสดงพฤติกรรม (Behavior) ของ Object/Type นั้นๆ |

---

## 2. โครงสร้างและการประกาศ (Syntax & Structure)

### 2.1 Function (ฟังก์ชันทั่วไป)
ฟังก์ชันคือบล็อกโค้ดอิสระ รับอินพุต (Parameters) ประมวลผล และส่งคืนค่า (Return value) โดยไม่ต้องมี Receiver

```go
func FullName(firstName string, lastName string) string {
    return firstName + " " + lastName
}
```

### 2.2 Method (เมธอด)
เมธอดคือฟังก์ชันที่มี **Receiver Argument** อยู่ระหว่างคีย์เวิร์ด `func` และชื่อเมธอด

```go
type Student struct {
    FirstName string
    LastName  string
}

// Receiver (s Student) ผูกกับ Student struct
func (s Student) FullName() string {
    return s.FirstName + " " + s.LastName
}
```

---

## 3. Value Receiver vs Pointer Receiver ใน Method

ใน Go การสร้าง Method บน Receiver สามารถแบ่งออกเป็น 2 ชนิดหลัก:

### 3.1 Value Receiver `func (s Student)`
- เป็นการส่งค่าแบบ **Pass by Value** (คัดลอกข้อมูลอินสแตนซ์เดิมมาใหม่)
- การแก้ไขฟิลด์ภายในเมธอด **จะไม่มีผล** ต่ออินสแตนซ์เดิมภายนอก
- เหมาะสำหรับเมธอดที่อ่านค่าอย่างเดียว (Read-only) หรือ Struct ที่มีขนาดเล็ก

```go
func (s Student) ChangeName(newName string) {
    s.FirstName = newName // ❌ เปลี่ยนเฉพาะสำเนา ตัวจริงไม่เปลี่ยน
}
```

### 3.2 Pointer Receiver `func (s *Student)`
- เป็นการส่งตัวชี้ตำแหน่งหน่วยความจำ **Pass by Reference/Pointer**
- การแก้ไขฟิลด์ภายในเมธอด **จะมีผล** ต่ออินสแตนซ์เดิมภายนอกทันที
- เหมาะสำหรับเมธอดที่ต้องการแก้ไขข้อมูล (Mutate State) หรือ Struct ที่มีขนาดใหญ่เพื่อหลีกเลี่ยงการคัดลอกข้อมูลในหน่วยความจำ

```go
func (s *Student) UpdateName(newName string) {
    s.FirstName = newName // ✅ เปลี่ยนแปลงข้อมูลในอินสแตนซ์เดิมจริง
}
```

---

## 4. ตัวอย่างโค้ดเปรียบเทียบการใช้งาน (Code Example)

```go
package main

import "fmt"

type Student struct {
	FirstName string
	LastName  string
}

// 1. Function: ส่งค่าเข้ามาทาง Parameter
func GetFullNameFunc(s Student) string {
	return s.FirstName + " " + s.LastName
}

// 2. Method (Value Receiver): อ่านค่าผ่าน Receiver
func (s Student) FullName() string {
	return s.FirstName + " " + s.LastName
}

// 3. Method (Pointer Receiver): แก้ไขค่าเดิมผ่าน Pointer Receiver
func (s *Student) SetFirstName(newName string) {
	s.FirstName = newName
}

func main() {
	student := Student{
		FirstName: "John",
		LastName:  "Doe",
	}

	// เรียกใช้งาน Function
	fmt.Println("Func Output:", GetFullNameFunc(student))

	// เรียกใช้งาน Method (Value Receiver)
	fmt.Println("Method Output:", student.FullName())

	// เรียกใช้งาน Method (Pointer Receiver) เพื่ออัปเดตชื่อ
	student.SetFirstName("Jane")
	fmt.Println("Updated Method Output:", student.FullName())
}
```

---

## 5. สรุปเมื่อไหร่ควรใช้ Function หรือ Method

1. **ใช้ Function เมื่อ**:
   - เป็นการประมวลผลทั่วไป เช่น คำนวณคณิตศาสตร์ (`math.Max`), จัดการ String (`strings.ToUpper`)
   - ไม่เกี่ยวข้องกับ State หรือคุณสมบัติเฉพาะของ Struct ชนิดใดชนิดหนึ่ง

2. **ใช้ Method เมื่อ**:
   - ต้องการกำหนดพฤติกรรม (Behavior) ให้กับ Struct หรือ Custom Type นั้นๆ
   - ต้องการทำให้โค้ดอ่านง่ายในรูปแบบ Object-Oriented Style (เช่น `student.FullName()`, `user.Save()`)
   - ต้องการเตรียมไว้รองรับการทำ **Interface** ใน Go
