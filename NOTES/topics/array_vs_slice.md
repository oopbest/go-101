# เปรียบเทียบ Array vs Slice ในภาษา Go (Golang)

## 1. ตารางเปรียบเทียบข้อแตกต่าง (Comparison Table)

| หัวข้อ (Feature) | Array | Slice |
| :--- | :--- | :--- |
| **ขนาด (Size)** | คงที่ (Fixed size) ต้องระบุขนาดตั้งแต่แรก | ยืดหยุ่น ขยาย/ลดขนาดได้ (Dynamic size) |
| **การประกาศ (Syntax)** | `[5]int{1, 2, 3, 4, 5}` (ระบุขนาดใน `[]`) | `[]int{1, 2, 3, 4, 5}` (ไม่ระบุขนาดใน `[]`) |
| **ชนิดข้อมูล (Type)** | ขนาดเป็นส่วนหนึ่งของ Type<br>*(เช่น `[3]int` กับ `[5]int` ถือเป็นคนละชนิดข้อมูล)* | เป็นชนิดข้อมูลเดียวกันคือ `[]int`<br>*(ไม่ว่าจะมีความยาวเท่าใด)* |
| **การส่งค่าเข้าฟังก์ชัน (Passing Mechanism)** | **Pass by Value**<br>*(คัดลอกข้อมูลอาร์เรย์ทั้งหมดไปใหม่)* | **Pass by Reference / Header**<br>*(ส่งโครงสร้างที่ชี้ไปยัง Underlying Array เดิม)* |
| **โครงสร้างภายใน (Internal Structure)** | เก็บข้อมูลเรียงต่อกันโดยตรงในหน่วยความจำ | มี 3 ส่วน: Pointer, Length (`len`), Capacity (`cap`) |
| **การเพิ่มข้อมูล** | ไม่สามารถใช้ `append()` เพิ่มข้อมูลได้ | ใช้ `append()` เพิ่มข้อมูลได้อย่างยืดหยุ่น |
| **การแปลงข้อมูล (Conversion)** | แปลงเป็น Slice ได้ง่ายด้วย Slicing `arr[:]` | แปลงเป็น Array ได้ตั้งแต่ Go 1.20+ (`[N]T(slice)`) หรือ Pointer `(*[N]T)(slice)` |

---

## 2. โครงสร้างภายในของ Slice (Slice Header)

Slice ใน Go ทำหน้าที่เป็นตัวห่อหุ้ม (Wrapper/Abstraction) ซ้อนอยู่บน **Array เบื้องหลัง (Underlying Array)** ซึ่งประกอบด้วย 3 ฟิลด์หลัก:

1. **Pointer**: ตัวชี้ไปยังตำแหน่งแรกของ Underlying Array ที่ Slice นั้นอ้างอิงอยู่
2. **Length (`len`)**: จำนวนสมาชิกที่มีอยู่ใน Slice ณ ปัจจุบัน
3. **Capacity (`cap`)**: ความจุสูงสุดที่ Slice สามารถรับได้ โดยนับจากตำแหน่ง Pointer ไปจนถึงตัวสุดท้ายของ Underlying Array

---

## 3. ตัวอย่างโค้ดประกอบการใช้งาน (Code Examples)

### ตัวอย่าง 1: Array
```go
package main

import "fmt"

func main() {
    // ประกาศ Array ขนาด 3
    var arr [3]int = [3]int{10, 20, 30}
    fmt.Println(arr) // [10 20 30]

    // arr = append(arr, 40) ❌ คอมไพล์ไม่ผ่าน! ไม่สามารถใช้ append กับ Array ได้
}
```

### ตัวอย่าง 2: Slice
```go
package main

import "fmt"

func main() {
    // ประกาศ Slice (ไม่ระบุขนาดใน [])
    slice := []int{10, 20, 30}
    fmt.Println(slice) // [10 20 30]

    // เพิ่มสมาชิกด้วย append()
    slice = append(slice, 40)
    fmt.Println(slice)      // [10 20 30 40]
    fmt.Println(len(slice)) // 4
    fmt.Println(cap(slice)) // 6 (Go จะเพิ่ม Capacity ให้อัตโนมัติเมื่อข้อมูลเต็ม)
}
```

### ตัวอย่าง 3: การตัดแบ่ง (Slicing) และค่า `len` / `cap`
```go
package main

import "fmt"

func main() {
    mySlice := []int{1, 2, 3, 4, 5}
    fmt.Println(mySlice)      // [1 2 3 4 5]
    fmt.Println(len(mySlice)) // 5
    fmt.Println(cap(mySlice)) // 5

    // ตัดเอาเฉพาะ Index 1 ถึงก่อน Index 3
    subSlice := mySlice[1:3]
    fmt.Println(subSlice)      // [2 3]
    fmt.Println(len(subSlice)) // 2 (มีสมาชิก 2 ตัว คือ 2 และ 3)
    fmt.Println(cap(subSlice)) // 4 (นับจาก Index 1 คือค่า 2 ไปจนถึงตัวสุดท้าย 5)
}
```

---

## 4. การแปลงข้อมูลระหว่าง Array และ Slice (Array & Slice Conversion)

### 4.1 การแปลง Array เป็น Slice (Array to Slice)
สามารถใช้ **Slicing Expression** (`arr[:]` หรือ `arr[low:high]`) เพื่อแปลง Array ให้กลายเป็น Slice ได้ทันที

> [!NOTE]
> Slice ที่ได้จะชี้ไปยัง **Underlying Array** ตัวเดิม ดังนั้น หากมีการแก้ไขสมาชิกใน Slice ข้อมูลใน Array ต้นทางจะเปลี่ยนแปลงตามไปด้วย

```go
package main

import "fmt"

func main() {
    arr := [5]int{10, 20, 30, 40, 50}
    
    // แปลง Array เป็น Slice
    slice := arr[:]

    // แก้ไขข้อมูลผ่าน Slice
    slice[0] = 99

    fmt.Println("Array:", arr)   // Array: [99 20 30 40 50] (ค่าเปลี่ยนตาม)
    fmt.Println("Slice:", slice) // Slice: [99 20 30 40 50]
}
```

### 4.2 การแปลง Slice เป็น Array (Slice to Array)

#### 1) การแปลงเป็น Array โดยตรง (`[N]T(slice)`) - Go 1.20+
ตั้งแต่ **Go 1.20** เป็นต้นไป เราสามารถแปลง Slice เป็น Array ได้โดยตรงด้วยไวยากรณ์ `[N]T(slice)` โดยการแปลงวิธีนี้จะเป็นการ **คัดลอกค่า (Copy / Pass by Value)** ไปสร้าง Array ชุดใหม่

```go
package main

import "fmt"

func main() {
    slice := []int{1, 2, 3, 4, 5}

    // แปลง Slice เป็น Array ขนาด 3
    arr := [3]int(slice[:3])

    arr[0] = 999
    fmt.Println("Array:", arr)   // Array: [999 2 3]
    fmt.Println("Slice:", slice) // Slice: [1 2 3 4 5] ( Slice เดิมไม่เปลี่ยน)
}
```

> [!WARNING]
> หากขนาดของ Array ที่ต้องการแปลง มีขนาดมากกว่าความยาว (`len`) ของ Slice จะทำให้เกิด **Runtime Panic** (เช่น `[10]int(slice)` ทั้งที่ slice มี len เป็น 5)

#### 2) การแปลงเป็น Array Pointer (`(*[N]T)(slice)`) - Go 1.17+
ตั้งแต่ **Go 1.17** เป็นต้นไป สามารถแปลง Slice ให้เป็น Pointer ของ Array ได้ ไวยากรณ์นี้ไม่ได้ทำการคัดลอกข้อมูลใหม่ แต่ชี้ไปยังหน่วยความจำเดิม

```go
package main

import "fmt"

func main() {
    slice := []int{1, 2, 3, 4, 5}

    // แปลง Slice เป็น Pointer ของ Array ขนาด 5
    arrPtr := (*[5]int)(slice)

    fmt.Println("Array Pointer:", *arrPtr) // [1 2 3 4 5]
}
```
