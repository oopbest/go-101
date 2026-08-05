# การใช้งาน Struct ในภาษา Go (Struct in Go)

**Struct** คือ โครงสร้างข้อมูลแบบกำหนดเอง (Custom Data Type) ที่เกิดจากการรวมกลุ่มของตัวแปรประเภทต่างๆ (Fields) เข้าด้วยกันเป็นชนิดข้อมูลใหม่ เหมาะสำหรับใช้แทนวัตถุที่มีหลายคุณสมบัติในโลกจริง (เช่น ข้อมูลนิสิตนักศึกษา ข้อมูลที่อยู่ ฯลฯ)

---

## 1. การประกาศและสร้าง Struct (Declaration & Instantiation)

### 1.1 การประกาศ Struct พื้นฐาน
```go
type Student struct {
    Name   string
    Age    int
    Weight float32
    Height int
    Grade  string
}
```

### 1.2 Nested Struct (Struct ซ้อน Struct)
ในภาษา Go เราสามารถนำ Struct หนึ่งไปเป็น Field ภายในอีก Struct หนึ่งได้ เพื่อจัดกลุ่มข้อมูลให้เป็นระเบียบมากยิ่งขึ้น

```go
type Address struct {
    City    string
    ZIPCode string
    Country string
}

type Student struct {
    Name    string
    Age     int
    Weight  float32
    Height  int
    Grade   string
    Address Address // Nested Struct
}
```

---

## 2. ตัวอย่างการใช้งาน (Code Examples)

### 2.1 การใช้งาน Struct ร่วมกับ Map และการเข้าถึง Nested Struct
เราสามารถสร้าง `Slice` หรือ `Map` ที่เก็บข้อมูลเป็น Struct เพื่อประยุกต์ใช้งานจริงได้

```go
package main

import "fmt"

type Address struct {
    City    string
    ZIPCode string
    Country string
}

type Student struct {
    Name    string
    Age     int
    Weight  float32
    Height  int
    Grade   string
    Address Address // Nested Struct
}

func main() {
    // 1. สร้าง Map ที่เก็บ Key เป็น string และ Value เป็น Student Struct
    students := make(map[string]Student)

    // 2. กำหนดค่าข้อมูลนิสิตพร้อมข้อมูล Address ที่ซ้อนอยู่ภายใน
    students["st001"] = Student{
        Name:   "John",
        Age:    20,
        Weight: 71.5,
        Height: 175,
        Grade:  "A",
        Address: Address{
            City:    "New York",
            ZIPCode: "10001",
            Country: "USA",
        },
    }

    students["st002"] = Student{
        Name:   "Adam",
        Age:    30,
        Weight: 75.5,
        Height: 180,
        Grade:  "B",
        Address: Address{
            City:    "Chicago",
            ZIPCode: "60601",
            Country: "USA",
        },
    }

    // 3. วนลูปอ่านข้อมูลใน Map และเข้าถึง Nested Struct Fields
    for id, student := range students {
        fmt.Println("ID:", id)
        fmt.Println("Name:", student.Name)
        fmt.Println("Age:", student.Age)
        fmt.Println("Grade:", student.Grade)
        // เข้าถึงข้อมูล Nested Struct ด้วย Dot Notation ซ้อนกัน
        fmt.Println("City:", student.Address.City)
        fmt.Println("ZIPCode:", student.Address.ZIPCode)
        fmt.Println("Country:", student.Address.Country)
        fmt.Println("--------------------------")
    }
}
```

---

## 3. สรุปจุดสำคัญของ Struct

| หัวข้อ | คำอธิบาย |
| :--- | :--- |
| **Field Access** | เข้าถึงหรือแก้ไขค่าใน Struct ใช้เครื่องหมายจุด (`.`) เช่น `student.Name` |
| **Nested Field Access** | เข้าถึง Struct ที่ซ้อนอยู่ด้านใน ใช้จุดซ้อนจุด เช่น `student.Address.City` |
| **Zero Value** | หากไม่ได้ระบุค่าให้ Field ใด ค่าเริ่มต้นจะเป็น Zero Value ของชนิดข้อมูลนั้น (เช่น `0`, `""`, `false`) |
| **Pass by Value** | โดยพื้นฐานการส่ง Struct เข้าฟังก์ชันจะเป็นการ **Pass by Value** (คัดลอกข้อมูลใหม่) หากต้องการให้แก้ไขข้อมูลเดิมได้ต้องใช้ Pointer (`*Student`) |
