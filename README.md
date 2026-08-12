# Go Learning Notes (สรุปเนื้อหาบทเรียน Golang)

ยินดีต้อนรับสู่คลังสรุปการเรียนรู้ภาษา Go (Golang) ครับ สามารถเลือกดูหัวข้อเนื้อหาแยกตามหมวดหมู่และเข้าดูโปรเจกต์ตัวอย่างได้ด้านล่างนี้:

---

## 📁 โครงสร้างโปรเจกต์และแบบฝึกหัด (Course Modules)

- 📂 [1-basic](file:///c:/Users/oopbest-pc-rog/Sites/1_learning/1_golang/mikelopster/1-basic) : ปูพื้นฐานไวยากรณ์ภาษา Go, Data Types, Slice, Map, Struct, Interface และ Pointer
- 📂 [2-fiber](file:///c:/Users/oopbest-pc-rog/Sites/1_learning/1_golang/mikelopster/2-fiber) : การพัฒนา Web RESTful API ด้วย Go Fiber Framework, JWT Auth, Role-Based Access Control, Middleware และ Swagger Documentation
- 📂 [3-go-database](file:///c:/Users/oopbest-pc-rog/Sites/1_learning/1_golang/mikelopster/3-go-database) : การทำ Docker Compose (PostgreSQL, Adminer, Redis), การใช้ `database/sql` driver `lib/pq`, Repository Pattern, Transactions และ Stored Procedures
- 📂 [4-go-gorm](file:///c:/Users/oopbest-pc-rog/Sites/1_learning/1_golang/mikelopster/4-go-gorm) : การใช้งาน GORM (ORM) ร่วมกับ PostgreSQL, Migration, Associations, Preload, Soft Delete และ GORM Logger

---

## 🛠️ คำสั่ง Go CLI พื้นฐานที่จำเป็นต้องรู้ (Essential Go Commands)

รวมคำสั่งภาษา Go (Command-Line) ที่ใช้งานบ่อยในการจัดการโปรเจกต์ Module และ Dependencies:

### 1. การจัดการ Module และ Dependencies
- **`go mod init <module-name>`** : สร้างไฟล์ `go.mod` เพื่อเริ่มต้นจัดการ Go Module ในโปรเจกต์ใหม่ (เช่น `go mod init myapp`)
- **`go get <package-path>`** : ดาวน์โหลดและติดตั้งแพ็กเกจ/ไลบรารีภายนอกเข้ามาใช้งานในโปรเจกต์ (เช่น `go get github.com/gofiber/fiber/v2`)
- **`go mod tidy`** : ลบ dependencies ที่ไม่ได้ใช้ออกจาก `go.mod` และเติม dependencies ที่ยังขาดอยู่อัตโนมัติ (ควรใช้เสมอเมื่อเพิ่ม/ลบไลบรารี)
- **`go mod download`** : ดาวน์โหลดแพ็กเกจทั้งหมดที่ระบุไว้ใน `go.mod` ลงเครื่อง

### 2. การรันและการคอมไพล์โปรเจกต์ (Build & Run)
- **`go run <file.go>`** หรือ **`go run .`** : คอมไพล์และสั่งรันโปรเจกต์ทันทีโดยไม่สร้างไฟล์ Executable ถาวร (เหมาะสำหรับทดสอบ/พัฒนา)
- **`go build`** : คอมไพล์โปรเจกต์ให้ออกมาเป็นไฟล์ Binary Executable (เช่น `.exe` บน Windows) สำหรับนำไป Deploy
- **`go install`** : คอมไพล์และติดตั้งไฟล์ Binary ไว้ในโฟลเดอร์ `$GOPATH/bin` เพื่อให้เรียกใช้งานคำสั่งนั้นได้จากทุกที่บนระบบ

### 3. การทดสอบ จัดฟอร์แมต และตรวจสอบระบบ (Testing & Tools)
- **`go test ./...`** : รัน Unit Test ทั้งหมดในโปรเจกต์และ sub-packages
- **`go fmt ./...`** : จัดรูปแบบจัดหน้าโค้ด (Formatting) ทุกไฟล์ในโปรเจกต์ให้เป็นมาตรฐานของ Go โดยอัตโนมัติ
- **`go env`** : แสดงค่าการตั้งค่า Environment Variables ทั้งหมดของ Go (เช่น `GOPATH`, `GOROOT`)

### 4. 🔄 ตารางเปรียบเทียบคำสั่ง Go CLI vs npm CLI (Node.js)

สำหรับผู้ที่คุ้นเคยกับ Node.js / npm สามารถเปรียบเทียบคำสั่งที่ทำหน้าที่ใกล้เคียงกันได้ดังนี้:

| การทำงาน (Action / Task) | Go CLI (Golang) | npm CLI (Node.js) | คำอธิบายเพิ่มเติม |
| :--- | :--- | :--- | :--- |
| **สร้างโปรเจกต์ใหม่** | `go mod init <name>` | `npm init` / `npm init -y` | สร้างไฟล์ Manifest (`go.mod` vs `package.json`) |
| **ติดตั้ง Library เพิ่มเติม** | `go get <package>` | `npm install <package>` | ดาวน์โหลดและบันทึกแพ็กเกจลง Manifest File |
| **ติดตั้ง Dependencies ทั้งหมด** | `go mod download` | `npm install` | ดาวน์โหลดแพ็กเกจตามรายการในไฟล์ Manifest |
| **จัดการ/เคลียร์ Dependencies** | `go mod tidy` | `npm prune` / `npm uninstall` | `go mod tidy` จะช่วยลบตัวไม่ใช้ และเติมตัวที่ขาดให้อัตโนมัติ |
| **สั่งรันโปรเจกต์** | `go run .` / `go run main.go` | `npm start` / `node index.js` | คอมไพล์และสั่งรันโปรเจกต์ในโหมดพัฒนา |
| **คอมไพล์สำหรับ Production** | `go build` | `npm run build` | สร้างไฟล์คอมไพล์ (`Binary Executable` vs `JS Bundle`) |
| **ติดตั้ง CLI Tool แบบ Global** | `go install <path>` | `npm install -g <package>` | ติดตั้งคำสั่งไว้ในระบบเพื่อเรียกใช้ได้ทุกที่ |
| **รัน Unit Test** | `go test ./...` | `npm test` | รันชุดทดสอบความถูกต้องของโค้ด |
| **จัดฟอร์แมตโค้ดอัตโนมัติ** | `go fmt ./...` | `npx prettier --write .` | จัดหน้าโค้ดให้เรียบร้อยตามมาตรฐานภาษา |
| **ไฟล์ข้อมูล Dependencies** | `go.mod` | `package.json` | ไฟล์หลักระบุเวอร์ชันและรายชื่อแพ็กเกจ |
| **ไฟล์ Lock Version & Hash** | `go.sum` | `package-lock.json` | ไฟล์เก็บ Checksum ยืนยันความถูกต้องของเวอร์ชัน |

---

## 📚 สารบัญเนื้อหา (Table of Contents)

### 1. [เปรียบเทียบ Array vs Slice (Array & Slice)](NOTES/topics/array_vs_slice.md)
สรุปโครงสร้างข้อมูลแบบลำดับใน Go:
- **ข้อแตกต่างระหว่าง Array และ Slice**: ตารางเปรียบเทียบขนาด, ชนิดข้อมูล, การส่งค่าเข้าฟังก์ชัน และการเพิ่มข้อมูล
- **โครงสร้างภายใน Slice Header**: ทำความเข้าใจ Pointer, Length (`len`), และ Capacity (`cap`)
- **ตัวอย่างโค้ดการใช้งาน**: การประกาศ, `append()`, และการตัดแบ่ง Slicing (`[low:high]`)
- **การแปลงข้อมูล (Conversion)**: แปลง Array เป็น Slice (`arr[:]`) และแปลง Slice เป็น Array (`[N]T(slice)` / Pointer)

### 2. [การใช้งาน Map (Map in Go)](NOTES/topics/map.md)
สรุปโครงสร้างข้อมูลแบบ Key-Value Pairs พร้อมความสำคัญและตัวอย่างการประยุกต์ใช้งานจริง:
- **ความสำคัญและการค้นหาข้อมูล $O(1)$**: ประสิทธิภาพการค้นหาข้อมูลระดับ $O(1)$ ด้วย Hash Table และความยืดหยุ่นในการใช้ Key ประเภทต่างๆ
- **คุณสมบัติและข้อควรระวังสำคัญ**: Unordered, Reference Type, Zero Value (`nil` map panic) และเรื่อง **Thread-Safety** (`concurrent map writes`)
- **การจัดการข้อมูลและ Comma-ok Idiom**: การใช้ `make()`, `delete()`, `len()` และการเช็คการมีอยู่ด้วย `val, ok := m[key]`
- **ตัวอย่างการใช้งานในระบบจริง (Real-World Use Cases)**: การใช้ Map สร้าง Set ป้องกันค่าซ้ำ (`map[string]bool`), การทำ Frequency Counter (`counts[word]++`), In-Memory Index/Caching และการจัดกลุ่มข้อมูล (Group By)

### 3. [การใช้งาน Struct (Struct in Go)](NOTES/topics/struct.md)
สรุปการใช้งานโครงสร้างข้อมูลแบบกำหนดเอง (Custom Data Types):
- **การประกาศและสร้าง Struct**: การกำหนดชนิดข้อมูลเฉพาะสำหรับอธิบายวัตถุ
- **Nested Struct**: การสร้าง Struct ซ้อน Struct เพื่อจัดกลุ่มข้อมูลให้เป็นระบบ
- **การใช้งานร่วมกับ Map/Slice**: ตัวอย่างการสร้าง `map[string]Student` และการอ่านค่าด้วย Dot Notation (`student.Address.City`)

### 4. [เปรียบเทียบ Function vs Method (Function vs Method)](NOTES/topics/function_vs_method.md)
สรุปความแตกต่างระหว่างฟังก์ชันทั่วไปและเมธอดในภาษา Go:
- **ตารางเปรียบเทียบข้อแตกต่าง**: การผูกติดกับ Type, การเรียกใช้งาน, การประกาศ Syntax และ Receiver
- **Value Receiver vs Pointer Receiver**: ข้อแตกต่างของการส่งค่าแบบ Pass by Value (`func (s Student)`) และ Pass by Reference (`func (s *Student)`)
- **ตัวอย่างโค้ดและเกณฑ์การใช้งาน**: ตัวอย่างเปรียบเทียบการเรียกใช้งาน Function และ Method พร้อมคำแนะนำเมื่อไหร่ควรเลือกใช้แบบใด

### 5. [การใช้งาน Interface (Interface in Go)](NOTES/topics/interface.md)
สรุปการใช้งาน Interface และพฤติกรรม Polymorphism ใน Go:
- **Implicit Implementation**: เข้าใจการทำตามสัญญา Interface โดยไม่ต้องใช้คีย์เวิร์ด `implements`
- **Polymorphism**: ตัวอย่างการส่งอินสแตนซ์ต่าง Struct เช่น `Circle` และ `Square` เข้าฟังก์ชันคำนวณพื้นที่ `printArea(s Shape)` รวมถึงตัวอย่าง `Speaker`
- **แนวคิดสำคัญ**: Empty Interface (`interface{}` / `any`), Type Assertion และการลดความยึดติดของโค้ด (Decoupling)

### 6. [การใช้งาน Pointer (Pointer in Go)](NOTES/topics/pointer.md)
สรุปการจัดการตำแหน่งหน่วยความจำและการใช้งาน Pointer ใน Go:
- **สัญลักษณ์สำคัญ**: สัญลักษณ์ `&` (Address-of) และ `*` (Dereference / Pointer Type)
- **หลักการทำงานและ Struct Pointer**: การส่ง Pointer เข้าฟังก์ชันเพื่อแก้ไขข้อมูลต้นทาง (`giveRaise(&emp, 15)`) พร้อมความสามารถ **Automatic Dereferencing** (`e.Salary`)
- **ประโยชน์และการเลือกใช้งาน**: ประโยชน์ของการประหยัดหน่วยความจำ และเมื่อไหร่ที่ควรเลือกใช้ Pointer

### 7. [การจัดการข้อผิดพลาด (Error Handling in Go)](NOTES/topics/error_handling.md)
สรุปแนวคิดและการจัดการข้อผิดพลาดในภาษา Go:
- **`error` Interface & `if err != nil`**: ทำความเข้าใจการจัดการข้อผิดพลาดแบบ Explicit Error Handling
- **การสร้าง Error**: วิธีการใช้ `errors.New()`, `fmt.Errorf()` และการสร้าง **Custom Error Struct** (`LoginError`)
- **แนวปฏิบัติที่ดี (Best Practices)**: การวางตำแหน่ง Return Error ตัวสุดท้าย, Early Return (Guard Clause) และข้อควรระวังเรื่อง `panic`

### 8. [การใช้งาน HTTP Server & Client (HTTP in Go)](NOTES/topics/http.md)
สรุปการสร้าง HTTP Web Server และการส่ง Request ด้วย Package `net/http`:
- **การสร้าง Web Server**: การใช้ `http.HandleFunc()` และ `http.ListenAndServe()` เพื่อให้บริการ API
- **การจัดการ Request & Response**: การอ่าน Query Parameters, Headers, Request Body และการส่ง JSON Response
- **การใช้งาน HTTP Client**: การส่ง GET/POST Request ด้วย `http.Get()`, `http.Post()` และ `http.Client`

### 9. [การใช้งาน Go Fiber Framework (Fiber in Go)](NOTES/topics/fiber.md)
สรุปการพัฒนา Web API ด้วย Fiber Framework (Express-style Web Framework บน Fasthttp):
- **การจัดเส้นทาง (Routing & Params)**: การรับ Path Parameter (`c.Params`) และ Query Parameter (`c.Query`)
- **การรับส่งข้อมูล (JSON & BodyParser)**: การใช้ `c.JSON()` และการแปลง Request Body ด้วย `c.BodyParser()`
- **Middleware & Grouping**: การตั้งค่า Logger, CORS, Recover และการจัดกลุ่ม API Endpoint ด้วย `app.Group()`

### 10. [การสร้าง RESTful API ด้วย Fiber (Fiber CRUD API)](NOTES/topics/fiber_crud.md)
สรุปการพัฒนา RESTful API ครบวงจรด้วย Fiber ร่วมกับแนวคิด Modular Package Structure:
- **โครงสร้างโปรเจกต์แบบแยก Package (Modular Structure)**: การแบ่งสถาปัตยกรรมเป็นสัดส่วน (`model`, `request`, `store`, `handler`) ให้โค้ดสะอาดและขยายง่าย
- **การจัดการ Concurrency (Thread-Safety)**: การใช้ `sync.RWMutex` (`RLock`/`Lock`) ป้องกัน Data Race เมื่อมี Goroutines เข้ามาทำงานพร้อมกัน
- **การจัดการ Slice & Memory**: เทคนิคการคัดลอก (Copy) Slice และเทคนิคการลบข้อมูลออกด้วย `append(s[:i], s[i+1:]...)`
- **การสร้าง RESTful Endpoints & Handlers**: การทำ Dependency Injection สู่ `main.go` และลงทะเบียน CRUD Routes (`GET`, `POST`, `PUT`, `DELETE`)
- **Unit Testing & Validation**: การเขียนชุดทดสอบ Unit Test แยกตามแพ็กเกจ (`store_test`, `handler_test`) และการจัดการ Error Response มาตรฐาน

### 11. [การจัดการ Environment Variables (Env in Go)](NOTES/topics/env.md)
สรุปการจัดการตัวแปรสภาพแวดล้อมเพื่อความปลอดภัยและการตั้งค่าระบบ:
- **การอ่านค่าด้วย Standard Library**: การใช้ `os.Getenv()`, `os.LookupEnv()` และการเขียน Helper Function ตรวจสอบค่าบังคับ (`log.Fatal`)
- **การใช้งานไฟล์ `.env`**: การโหลดค่าด้วยแพ็กเกจ `github.com/joho/godotenv` ในโหมด Development
- **การส่งค่าผ่าน CLI**: การตั้งค่าตัวแปรในระบบปฏิบัติต่างๆ (Bash: `KEY=val`, PowerShell: `$env:KEY="val"`, CMD: `set KEY=val`)

### 12. [การใช้งาน Middleware (Middleware in Go & Fiber)](NOTES/topics/middleware.md)
สรุปการใช้งานและเขียน Middleware ตัวคัดกรองใน Request-Response Pipeline:
- **Built-in Middlewares ใน Fiber**: การใช้งาน `logger.New()`, `recover.New()` ป้องกัน Server ล่ม และ `cors.New()`
- **การเขียน Custom Middleware**: ตัวอย่างการทำ Secret Key Authentication และการคำนวณระยะเวลาประมวลผล (Timing)
- **Middleware & Route Grouping**: การนำ Middleware ไปประยุกต์ใช้กับเฉพาะบางกลุ่ม Endpoint ด้วย `app.Group()`

### 13. [การสร้างเอกสาร API ด้วย Swagger (Swagger Documentation)](NOTES/topics/swagger.md)
สรุปการสร้างเอกสาร API แบบโต้ตอบได้ (Interactive UI) ด้วย Swagger และ OpenAPI Specs:
- **แพ็กเกจและการติดตั้ง**: การใช้งาน `gofiber/swagger` และ `swaggo/swag`
- **การเขียน Swag Annotations**: การระบุ General API Info (`@title`, `@securityDefinitions.apikey`) และ Annotations กำกับแต่ละ Handler Function
- **การสร้าง Docs & การทดสอบ**: การใช้คำสั่ง `swag init` สร้างไฟล์ Specs และการเข้าใช้งานผ่าน `http://localhost:8080/swagger/index.html`

### 14. [การใช้งาน Go database/sql & PostgreSQL (Go Database/SQL)](NOTES/topics/database_sql.md)
สรุปการใช้งานแพ็กเกจมาตรฐาน `database/sql` ร่วมกับ PostgreSQL และการจัดโครงสร้างโปรเจกต์ด้วย Repository Pattern:
- **การเชื่อมต่อ PostgreSQL**: การใช้งาน `database/sql` ร่วมกับ Driver `github.com/lib/pq` และการอ่านค่าจาก `.env`
- **การสร้าง Schema และ Constraints**: การใช้ `db.Exec()` สร้างตาราง `suppliers` และ `products` พร้อม Foreign Key Constraints
- **การทำ CRUD & Repository Pattern**: การจัด Folder Structure แบบ Clean Architecture (`config`, `models`, `repository`) และการสร้าง Constructor (`NewProductRepository`)
- **การเปรียบเทียบคำสั่ง SQL**: ข้อแตกต่างระหว่าง `db.Exec()`, `db.Query()`, และ `db.QueryRow()`
- **การจัดการ Pointer & NULL**: การดึงค่าและแสดงผลฟิลด์ Pointer (`*int`, `*string`) ด้วย `fmt.Printf` อย่างปลอดภัย
- **การแก้ปัญหาพบบ่อย (Troubleshooting)**: การจัดการ `rows.Err()`, แก้ไข Foreign Key Violation (Code `23503`), และเรื่อง Variable Scope / Type Safety ใน Go

### 15. [การจัดการ Database Transaction (Database Transaction in Go)](NOTES/topics/database_transaction.md)
สรุปการใช้งาน Database Transaction เพื่อรักษารวมความถูกต้องของข้อมูลตามหลัก ACID:
- **หลักการ ACID**: ทำความเข้าใจ Atomicity, Consistency, Isolation, และ Durability
- **Safe Transaction Pattern**: เทคนิคการใช้ `db.Begin()`, `defer tx.Rollback()`, และ `tx.Commit()` เพื่อป้องกันข้อมูลค้าง
- **ตัวอย่างระบบจริง**: โค้ดตัวอย่างระบบโอนเงินข้ามบัญชี และระบบสร้าง Order พร้อมหัก Stock สินค้า
- **การจัดการ Deadlock & Retry**: สาเหตุการเกิด Deadlock (Code `40P01`), เทคนิค Consistent Lock Ordering และการทำ Retry Mechanism ฝั่ง Application
- **การกำหนด Isolation Level**: การใช้ `db.BeginTx()` กำหนด Isolation Level และการใช้ร่วมกับ `context.Context`

### 16. [การใช้งาน Stored Procedure & Function (Stored Procedure in Go)](NOTES/topics/stored_procedure.md)
สรุปการสร้าง Stored Procedure และ Function ใน PostgreSQL ร่วมกับการเรียกใช้งานจากภาษา Go:
- **Procedure vs Function**: ความแตกต่างของ `CREATE FUNCTION` (คืนค่า) และ `CREATE PROCEDURE` (ควบคุม Transaction)
- **การเรียกใช้งานใน Go**: เทคนิคการใช้ `SELECT function_name()` ดึงค่า Scalar/Table และการใช้ `CALL procedure_name()` สั่งรันด้วย `db.Exec()`
- **ตัวอย่างโค้ดระบบจริง**: Function คำนวณ VAT, Function ดึงรายการสินค้าแบบ Table และ Procedure ย้ายสินค้าข้าม Supplier
- **ข้อดีและข้อเสีย**: เปรียบเทียบ Performance, Security, Vendor Lock-in, และความยากง่ายในการ Maintenance/Testing

### 17. [แนวคิด ORM, ODM และการใช้งาน GORM ร่วมกับ PostgreSQL](NOTES/topics/gorm.md)
สรุปแนวคิดการเชื่อมโลกของ Object ในโปรแกรมเข้ากับรูปแบบข้อมูลในฐานข้อมูล พร้อมการใช้ GORM ในภาษา Go:
- **ORM (Object-Relational Mapping)**: แปลง Struct/Object ให้สัมพันธ์กับ Table, Row, Column, Primary Key และ Foreign Key ในฐานข้อมูลเชิงสัมพันธ์ เช่น PostgreSQL และ MySQL ช่วยลด SQL ซ้ำ ๆ แต่ผู้พัฒนายังควรเข้าใจ SQL, Index และ Query Plan
- **ODM (Object-Document Mapping)**: แปลง Struct/Object ให้สัมพันธ์กับ Document และ Collection ในฐานข้อมูลแบบ Document เช่น MongoDB ซึ่งเหมาะกับข้อมูล JSON/BSON ที่มีโครงสร้างยืดหยุ่นและใช้ Embedded Document หรือ Reference แทนการ JOIN แบบ Relational
- **ORM vs ODM**: ORM ทำงานกับ Schema และ Relationships ของฐานข้อมูล SQL ส่วน ODM ทำงานกับ Document ที่สามารถซ้อนข้อมูลและมี Field แตกต่างกันได้ ทั้งสองแนวคิดช่วยเรื่อง Mapping, Validation และ CRUD แต่ไม่สามารถใช้แทนกันโดยตรง
- **ตำแหน่งของ GORM**: GORM เป็น **ORM ไม่ใช่ ODM** ออกแบบมาสำหรับฐานข้อมูลเชิงสัมพันธ์ จึงเหมาะกับ PostgreSQL ในโปรเจกต์นี้ หากใช้ MongoDB ควรเลือก MongoDB Driver หรือไลบรารี ODM ที่รองรับ Document Database โดยเฉพาะ
- **ข้อดีและข้อแลกเปลี่ยน**: ORM/ODM ช่วยลด Boilerplate และเพิ่มความเร็วในการพัฒนา แต่ abstraction อาจซ่อน Query ที่ไม่มีประสิทธิภาพ เช่น N+1 Queries จึงควรเปิด SQL Log และตรวจสอบ Query ที่ระบบสร้างขึ้น
- **การติดตั้งและเชื่อมต่อฐานข้อมูล**: การใช้ `gorm.io/gorm` ร่วมกับ `gorm.io/driver/postgres` และการตั้งค่า Connection Pool
- **Models, Tags และ Auto Migration**: การกำหนด Primary Key, Column Constraints, Relationships และการสร้างหรือปรับ Schema ด้วย `AutoMigrate()`
- **CRUD, Associations และ Transactions**: การใช้ `Create()`, `First()`, `Find()`, `Updates()`, `Delete()`, `Preload()` และ `Transaction()` พร้อมตรวจสอบ `Error` และ `RowsAffected`
- **Soft Delete**: การใช้ `gorm.Model` หรือ `gorm.DeletedAt` เพื่อกำหนดค่า `deleted_at` แทนการลบแถวจริง โดย Query ปกติจะไม่คืนข้อมูลที่ถูกลบ สามารถใช้ `Unscoped()` เพื่อค้นหาหรือลบถาวร และควรพิจารณาผลต่อ Unique Index กับ Foreign Key เพราะ Soft Delete เป็นคำสั่ง `UPDATE` ไม่ใช่ `DELETE`
- **GORM Logger และ SQL Logs**: การตั้งค่า `logger.Default.LogMode()` ด้วยระดับ `Silent`, `Error`, `Warn` และ `Info` เพื่อดู SQL, ระยะเวลาทำงาน, จำนวนแถว และ Error ระหว่างพัฒนา พร้อมข้อควรระวังไม่แสดง SQL หรือค่าพารามิเตอร์ที่อาจมีข้อมูลสำคัญใน Production
- **แนวปฏิบัติที่ดี**: การส่ง `context.Context`, ป้องกันการ Update/Delete โดยไม่มีเงื่อนไข แยก Database Model ออกจาก API DTO และเลือกใช้ Raw SQL เมื่อ Query ซับซ้อนหรือจำเป็นต้องปรับประสิทธิภาพเฉพาะจุด

---

> [!TIP]
> คลิกที่หัวข้อด้านบนเพื่อเปิดอ่านบทเรียนและตัวอย่างโค้ดของแต่ละเรื่องได้เลยครับ!
