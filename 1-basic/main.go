package main

import (
	"fmt"
	"net/http"

	"github.com/oopbest/go-basic/oopbest"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}

	fmt.Fprintf(w, "Hello, %s!", name)
}

func main() {
	// calling function from another package
	oopbest.SayHelloThailand()

	http.HandleFunc("/hello", helloHandler)
	fmt.Println("Server Starting on port 8080...")

	http.ListenAndServe(":8080", nil)
}
