# การใช้งาน Interface ในภาษา Go (Interface in Go)

**Interface** คือ ชนิดข้อมูลชนิดหนึ่งใน Go ที่ใช้กำหนด **ข้อตกลงหรือสัญญา (Contract/Behavior)** ว่า Type ที่จะทำตาม Interface นี้จะต้องมี Method อะไรบ้าง (Signature ไหน) โดยใน Go จะใช้วิธี **Implicit Implementation** (ไม่ต้องใช้คีย์เวิร์ด `implements`) ชนิดข้อมูลใดที่มี Method ครบตามที่ Interface กำหนด จะถือว่าเป็น Interface นั้นทันที!

---

## 1. การประกาศและแนวคิดพื้นฐาน (Syntax & Concept)

### 1.1 การประกาศ Interface
```go
type Shape interface {
    Area() float64
}
```

### 1.2 Implicit Implementation (ทำตามสัญญาโดยอัตโนมัติ)
ใน Go หาก Struct หรือ Custom Type ใดมี Method ที่ชื่อ และ Signature (พารามิเตอร์และชนิดค่าที่ส่งกลับ) ตรงกับที่ Interface กำหนด จะถือว่าทำตาม Interface นั้นทันทีโดยไม่ต้องระบุคีย์เวิร์ดใดๆ เพิ่มเติม

```go
// Circle มี Method Area() float64 -> ถือว่าเป็น Shape โดยอัตโนมัติ
type Circle struct {
    R float64
}

func (c Circle) Area() float64 {
    return 3.14 * c.R * c.R
}

// Square มี Method Area() float64 -> ถือว่าเป็น Shape เช่นกัน
type Square struct {
    Side float64
}

func (s Square) Area() float64 {
    return s.Side * s.Side
}
```

---

## 2. พหุรูปร่าง (Polymorphism) ด้วย Interface

ประโยชน์หลักของ Interface คือทำให้เราสามารถเขียนฟังก์ชันที่รับ Parameter เป็น Interface เดียวกัน แต่ส่ง Struct ต่างชนิดกันที่มีพฤติกรรมตรงตาม Interface เข้ามาประมวลผลได้

```go
// ฟังก์ชัน printArea รับ Parameter เป็น Shape (รับได้ทั้ง Circle, Square หรือชนิดข้อมูลอื่นที่เป็น Shape)
func printArea(s Shape) {
    fmt.Println(s.Area())
}
```

---

## 3. ตัวอย่างโค้ดประกอบการใช้งาน (Code Examples)

### 3.1 ตัวอย่างที่ 1: Interface คำนวณพื้นที่รูปทรงเรขาคณิต (Shape & Area)
```go
package main

import "fmt"

// 1. ประกาศ Interface
type Shape interface {
	Area() float64
}

// 2. Struct ที่ 1: Circle (วงกลม)
type Circle struct {
	R float64
}

func (c Circle) Area() float64 {
	return 3.14 * c.R * c.R
}

// 3. Struct ที่ 2: Square (สี่เหลี่ยมจัตุรัส)
type Square struct {
	Side float64
}

func (s Square) Area() float64 {
	return s.Side * s.Side
}

// 4. ฟังก์ชันที่รับ Interface เป็น Parameter (Polymorphism)
func printArea(s Shape) {
	fmt.Println(s.Area())
}

func main() {
	circle := Circle{R: 5}
	square := Square{Side: 5}

	// ทั้ง circle และ square สามารถส่งเข้า printArea() ได้เพราะทั้งคู่มี Method Area() float64
	printArea(circle) // Output: 78.5 (3.14 * 5 * 5)
	printArea(square) // Output: 25   (5 * 5)
}
```

### 3.2 ตัวอย่างที่ 2: Interface สำหรับการส่งเสียงพูด (Speaker & Speak)
```go
package main

import "fmt"

// 1. ประกาศ Interface
type Speaker interface {
	Speak() string
}

// 2. Struct ที่ 1: Person
type Person struct {
	Name string
}

func (p Person) Speak() string {
	return "Hello, my name is " + p.Name
}

// 3. Struct ที่ 2: Dog
type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return "Woof! My name is " + d.Name
}

// 4. ฟังก์ชันที่รับ Interface เป็น Parameter
func makeSound(s Speaker) {
	fmt.Println(s.Speak())
}

func main() {
	person := Person{Name: "John"}
	dog := Dog{Name: "Buddy"}

	makeSound(person) // Output: Hello, my name is John
	makeSound(dog)    // Output: Woof! My name is Buddy
}
```

---

## 4. คุณสมบัติและแนวคิดสำคัญของ Interface ใน Go

| หัวข้อ | คำอธิบาย |
| :--- | :--- |
| **Implicit Implementation** | ไม่ต้องใช้คีย์เวิร์ด `implements` แค่สร้าง Method ให้ตรง Signature ก็ถือว่าทำตามสัญญาแล้ว |
| **Decoupling & Flexibility** | ช่วยลดความยึดติดกันของโค้ด (Loose Coupling) เพิ่มความยืดหยุ่น และเอื้อต่อการเขียน Unit Test / Mocking |
| **Empty Interface (`interface{}` / `any`)** | Interface ที่ไม่มี Method กำหนด สามารถรับค่าข้อมูลได้ทุกชนิดในภาษา Go (ใช้อักษรย่อ `any` ได้ตั้งแต่ Go 1.18+) |
| **Type Assertion** | เทคนิคการตรวจสอบและแปลงข้อมูลจาก Interface กลับเป็น Concrete Type เดิม เช่น `c, ok := s.(Circle)` |
