# การจัดการ Environment Variables ในภาษา Go (Environment Variables in Go)

**Environment Variables (ตัวแปรสภาพแวดล้อม)** เป็นกลไกสำคัญในการแยกการตั้งค่าของแอปพลิเคชัน (Configuration) ออกจากซอร์สโค้ดตามหลักการ **12-Factor App** เพื่อความปลอดภัยในการเก็บข้อมูลลับ (Secrets, API Keys, Database Password) และรองรับการทำงานหลากสภาพแวดล้อม (Development, Staging, Production)

---

## 1. การอ่านค่าด้วย Standard Library (`os` Package)

### 1.1 การใช้ `os.Getenv()`
ดึงค่าตัวแปรสภาพแวดล้อมตาม Key ที่กำหนด หากไม่มีการตั้งค่าไว้จะคืนค่าเป็น String ว่าง (`""`):

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // ค่า Default กรณีไม่ได้กำหนดไว้
	}
	fmt.Println("Server running on port:", port)
}
```

### 1.2 การใช้ `os.LookupEnv()`
ใช้เมื่อต้องการแยกแยะระหว่าง **"ตัวแปรที่ไม่ได้ถูกตั้งค่า"** กับ **"ตัวแปรที่ตั้งค่าเป็น String ว่าง"**:

```go
val, exists := os.LookupEnv("SECRET_KEY")
if !exists {
	fmt.Println("ไม่มีการตั้งค่า SECRET_KEY ในระบบ")
}
```

### 1.3 การเขียน Helper Function ยุบรวมรองรับทั้ง บังคับมีค่า (Require) และ ค่าเริ่มต้น (Default)
สามารถใช้ **Variadic Parameters (`...string`)** เพื่อสร้างฟังก์ชันอ่าน Env แบบอเนกประสงค์:

```go
func getEnv(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	log.Fatalf("Environment variable %s is required but not set", key)
	return ""
}

// วิธีใช้งาน:
adminPass := getEnv("ADMIN_PASSWORD")            // บังคับต้องมี ถ้าไม่มีจะ log.Fatal ทันที
jwtSecret := getEnv("JWT_SECRET", "default_secret") // ถ้าไม่มีจะใช้ "default_secret" แทน
```


---

## 2. การใช้งานไฟล์ `.env` ด้วยแพ็กเกจ `godotenv`

ในขั้นตอนการพัฒนา (Development) นิยมเก็บค่าตัวแปรไว้ในไฟล์ `.env` และใช้แพ็กเกจ [`github.com/joho/godotenv`](https://github.com/joho/godotenv) ในการโหลดเข้าสู่โปรเซสโดยอัตโนมัติ

### 2.1 การติดตั้ง
```bash
go get github.com/joho/godotenv
```

### 2.2 ตัวอย่างไฟล์ `.env`
```env
PORT=8080
SECRET=my_super_secret_key_123
DB_URL=postgres://user:pass@localhost:5432/mydb
```

### 2.3 การโหลดเข้าใช้งานใน `main.go`
```go
package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// โหลดไฟล์ .env เข้าสู่ os.Getenv
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	secret := os.Getenv("SECRET")
	log.Printf("Loaded SECRET: %s", secret)
}
```

> [!CAUTION]
> **ห้าม commit ไฟล์ `.env` เข้า Git Repository โดยเด็ดขาด!** ควรเพิ่ม `.env` ไว้ในไฟล์ `.gitignore` เสมอ และสร้างไฟล์ตัวอย่างชื่อ `.env.example` ไว้เป็นเทมเพลตแทน

---

## 3. การกำหนด Environment Variables ผ่าน CLI สั่งรัน

สามารถส่งค่าตัวแปรสภาพแวดล้อมได้ผ่านคำสั่ง Terminal ในแต่ละระบบปฏิบัติการ:

### Git Bash / Linux / macOS
```bash
SECRET=mysecretkey go run main.go
# หรือใช้กับ air (Live Reload)
SECRET=mysecretkey air
```

### PowerShell (Windows)
```powershell
$env:SECRET="mysecretkey"; go run main.go
```

### Command Prompt (CMD Windows)
```cmd
set SECRET=mysecretkey && go run main.go
```

---

## 4. สรุปแนวปฏิบัติที่ดี (Best Practices)

1. **ไม่ Hardcode ข้อมูลลับ**: ห้ามเก็บ Secret, Passwords หรือ API Keys ไว้ในโค้ด Go
2. **ใช้ `.env` เฉพาะการพัฒนาในเครื่อง**: ในระบบ Production ควรส่งค่าผ่าน Docker, Kubernetes Secrets หรือ Cloud Provider (AWS Secrets Manager)
3. **ใช้ `.env.example`**: สร้างไฟล์ `.env.example` ระบุเฉพาะชื่อ Key โดยไม่ใส่ค่าจริง เพื่อให้คนในทีมรู้ว่าต้องตั้งค่าตัวแปรใดบ้าง
