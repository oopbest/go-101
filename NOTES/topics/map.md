# การใช้งาน Map ในภาษา Go (Map in Go)

**Map** คือ โครงสร้างข้อมูลที่เก็บข้อมูลในรูปแบบ **Key-Value Pairs** (เทียบเท่ากับ Dictionary ใน Python, Hash Table ใน C++ หรือ Object/Map ใน JavaScript) โดย Key แต่ละตัวจะต้องเป็นประเภทที่เปรียบเทียบกันได้ (`comparable`) และไม่ซ้ำกัน

---

## 1. ความสำคัญของ Map (Why Map is Important)

1. **ค้นหาข้อมูลรวดเร็วระดับ $O(1)$ (Fast Lookup Time Complexity)**:
   * เมื่อเทียบกับ Slice ที่ต้องวนลูปหาข้อมูลแบบ $O(n)$ การค้นหาด้วย Key ใน Map ใช้เวลาคงที่ $O(1)$ ไม่ว่าข้อมูลจะมีจำนวนหลักแสนหรือหลักล้านรายการ
2. **ความยืดหยุ่นในการใช้ Key (Flexible Key Types)**:
   * ต่างจาก Slice ที่ต้องอ้างอิงด้วย Index ตัวเลข $0, 1, 2...$ เท่านั้น แต่ Map สามารถใช้ `string`, `int`, `struct` หรือประเภทข้อมูลอื่นที่เปรียบเทียบกันได้มาเป็น Key
3. **ใช้สร้างโครงสร้างข้อมูลแบบ Custom (Emulating Sets & In-memory Cache)**:
   * ใช้ทำ Caching, Lookup Table, นับความถี่ข้อมูล (Frequency Counter), จัดกลุ่มข้อมูล (Group By), และจำลอง Set ป้องกันค่าซ้ำ

---

## 2. ลักษณะสำคัญและข้อควรระวัง (Characteristics & Pitfalls)

- **Unordered**: ลำดับข้อมูลใน Map ไม่แน่นอน การวนลูปด้วย `for key, val := range m` แต่ละครั้งอาจได้ลำดับไม่เหมือนกัน
- **Reference Type**: Map เป็น Reference Type เมื่อส่งเข้าฟังก์ชัน จะเป็นการส่ง Pointer อ้างอิงไปยัง memory ชุดเดิม (การแก้ไขภายในฟังก์ชันจะกระทบต้นทาง)
- **Zero Value คือ `nil`**: 
  * ประกาศ `var m map[string]int` จะได้ค่าเป็น `nil`
  * สามารถอ่านค่าจาก `nil` map ได้ (คืนค่า Zero Value ของชนิดนั้น) แต่ **ห้ามเขียน/เพิ่มข้อมูลลง `nil` map ทันที** เพราะจะเกิด `panic: assignment to entry in nil map`
  * ต้องสร้างผ่าน `make(map[KeyType]ValueType)` หรือ Map Literal (`map[string]int{}`) ก่อนเสมอ
- **Not Thread-Safe**: 
  * Map ใน Go ไม่รองรับการเขียนแบบ Concurrent หากหลาย Goroutine เขียน Map ตัวเดียวกันพร้อมกัน จะเกิด `fatal error: concurrent map writes`
  * ทางแก้ไข: ต้องใช้ `sync.RWMutex` คอยล็อกก่อนอ่าน/เขียน หรือใช้ `sync.Map`

---

## 3. ตัวอย่างการใช้งานพื้นฐาน (Basic Operations)

```go
package main

import "fmt"

func main() {
	// 1. การสร้าง Map ด้วย make()
	myMap := make(map[string]int)

	// การเพิ่มและอัปเดตข้อมูล (Insert & Update)
	myMap["apple"] = 3
	myMap["banana"] = 6
	myMap["banana"] = 12 // อัปเดตค่าเดิม

	// 2. การเข้าถึงข้อมูลด้วย Key (Access Value)
	fmt.Println("ค่าของ banana:", myMap["banana"]) // 12

	// 3. การลบข้อมูลด้วย delete() (Delete Key)
	delete(myMap, "apple")

	// 4. การตรวจสอบว่ามี Key อยู่หรือไม่ (Comma-ok Idiom)
	val, ok := myMap["apple"]
	if ok {
		fmt.Println("พบ apple:", val)
	} else {
		fmt.Println("ไม่พบ apple ใน Map") // ทำงานบรรทัดนี้
	}
}
```

---

## 4. ตัวอย่างการใช้งานในระบบจริง (Real-world Code Examples)

### ตัวอย่างที่ 1: การจำลอง Set เพื่อเช็คและป้องกันค่าซ้ำ (Set Emulation)
ภาษา Go ไม่มีชนิดข้อมูล `Set` ในตัว เรานิยมใช้ `map[KeyType]bool` หรือ `map[KeyType]struct{}` ในการเก็บค่าไม่ให้ซ้ำ:

```go
package main

import "fmt"

func removeDuplicates(items []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range items {
		if !seen[item] { // ถ้ายังไม่เคยเจอใน Map
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	emails := []string{"a@test.com", "b@test.com", "a@test.com", "c@test.com"}
	uniqueEmails := removeDuplicates(emails)
	fmt.Println("Emails แบบไม่ซ้ำ:", uniqueEmails) // [a@test.com b@test.com c@test.com]
}
```

### ตัวอย่างที่ 2: การนับความถี่ของข้อมูล (Frequency Counter)
ใช้นับจำนวนสถิติ เช่น นับจำนวนคำ หรือนับสถิติการใช้งาน:

```go
package main

import "fmt"

func countWords(text []string) map[string]int {
	counts := make(map[string]int)
	for _, word := range text {
		counts[word]++ // หากไม่มี Key ค่าจะเริ่มจาก Zero Value (0) แล้ว +1
	}
	return counts
}

func main() {
	words := []string{"go", "fiber", "go", "gin", "go", "fiber"}
	fmt.Println("จำนวนคำ:", countWords(words)) // map[fiber:2 gin:1 go:3]
}
```

### ตัวอย่างที่ 3: การทำ In-Memory Lookup Index / Cache
ค้นหาข้อมูลวัตถุด้วย Unique ID ระดับความเร็ว $O(1)$:

```go
type User struct {
	ID   int
	Name string
}

// แปลง Slice เป็น Map เพื่อให้ค้นหาด้วย ID ได้รวดเร็ว O(1)
func buildUserIndex(users []User) map[int]User {
	userMap := make(map[int]User)
	for _, u := range users {
		userMap[u.ID] = u
	}
	return userMap
}
```

### ตัวอย่างที่ 4: การจัดกลุ่มข้อมูล (Group By)
```go
type Product struct {
	Name     string
	Category string
}

func groupByCategory(products []Product) map[string][]Product {
	grouped := make(map[string][]Product)
	for _, p := range products {
		grouped[p.Category] = append(grouped[p.Category], p)
	}
	return grouped
}
```

---

## 5. ตารางสรุปคำสั่งคำนวณที่สำคัญ

| การทำงาน (Operation) | ไวยากรณ์ (Syntax) | คำอธิบาย |
| :--- | :--- | :--- |
| **สร้าง Map** | `make(map[K]V)` | สร้าง Map พร้อมใช้งาน (ป้องกัน nil panic) |
| **เพิ่ม/อัปเดต** | `m[key] = value` | กำหนดค่าให้ผูกกับ Key |
| **อ่านค่า** | `val := m[key]` | คืนค่า Value หรือ Zero Value หากไม่มี Key |
| **เช็คการมีอยู่** | `val, ok := m[key]` | Comma-ok idiom (`ok` เป็น `true` เมื่อมี Key) |
| **ลบ Key** | `delete(m, key)` | ลบ Key ออกจาก Map |
| **นับจำนวน Key** | `len(m)` | คืนค่าจำนวน Key ทั้งหมดใน Map |
