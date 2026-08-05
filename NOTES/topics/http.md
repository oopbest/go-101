# การใช้งาน HTTP Server และ Client ในภาษา Go (HTTP in Go)

ภาษา Go มี Standard Library ที่ชื่อว่า `net/http` ซึ่งมีความสามารถสูงและประสิทธิภาพดีเยี่ยม สามารถใช้สร้างทั้ง HTTP Web Server และ HTTP Client ได้โดยไม่ต้องพึ่งพา Third-party Library เพิ่มเติม

---

## 1. การสร้าง HTTP Server (HTTP Web Server)

### 1.1 การสร้าง Server เบื้องต้นด้วย `http.HandleFunc`
เราสามารถลงทะเบียน Route และ Handler Function ได้ง่ายๆ ผ่าน `http.HandleFunc` และสั่งให้ Server เริ่มทำงานด้วย `http.ListenAndServe`:

```go
package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	// กำหนด Handler สำหรับ Route "/"
	http.HandleFunc("/", helloHandler)

	fmt.Println("Server running on http://localhost:8080")
	// เริ่ม Server ที่ Port 8080
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
```

---

## 2. การจัดการ Request และ Response (Handling Requests & Responses)

### 2.1 การรับ-ส่งข้อมูลแบบ JSON
การสร้าง REST API ใน Go นิยมใช้ `json.NewDecoder` สำหรับอ่าน Request Body และ `json.NewEncoder` สำหรับส่ง JSON Response

```go
package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	// กำหนด Header ให้เป็น application/json
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		var u User
		// แกะข้อมูล JSON จาก Request Body
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		// ส่ง Response กลับ
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
		return
	}

	// กรณี GET Method
	u := User{ID: 1, Name: "Mikelopster"}
	json.NewEncoder(w).Encode(u)
}
```

### 2.2 การอ่าน Query Parameters
```go
func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q") // ดึงค่าจาก ?q=...
	fmt.Fprintf(w, "Search query: %s", query)
}
```

---

## 3. การใช้งาน HTTP Client (Sending HTTP Requests)

### 3.1 การส่ง HTTP GET Request
```go
package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	resp, err := http.Get("https://api.github.com/users/octocat")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close() // ต้องปิด Body หลังใช้งานเสร็จเสมอ

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}
```

### 3.2 การสร้าง HTTP Client พร้อมกำหนด Timeout
เพื่อป้องกันไม่ให้ Client รอนานเกินไป ควรสร้าง `http.Client` และกำหนด Timeout:

```go
import (
	"net/http"
	"time"
)

client := &http.Client{
	Timeout: 10 * time.Second,
}

resp, err := client.Get("https://example.com")
```

---

## 4. สรุปแนวปฏิบัติที่ดี (Best Practices)
1. **ปิด Response Body เสมอ**: เมื่อใช้งาน HTTP Client ต้องใช้ `defer resp.Body.Close()` เพื่อป้องกัน Resource Leak
2. **ตั้งค่า Timeout**: ควรตั้งค่า Timeout สำหรับทั้ง HTTP Server และ Client เสมอเพื่อความปลอดภัย
3. **ตรวจสอบ HTTP Status Code**: ควรตรวจสอบ `resp.StatusCode` หรือกำหนด `w.WriteHeader()` ให้ถูกต้องตามมาตรฐาน HTTP REST API
