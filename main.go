package main

import (
	"fmt"
	"net/http"

	"course_sdk/handler"
	"course_sdk/utils"
)

func main() {
	utils.LoadEnv()

	http.HandleFunc("/generate", handler.GenerateContent)

	apiURL := "http://localhost:8080"
	fmt.Println("===========================================")
	fmt.Printf("✅ Server running on: %s\n", apiURL)
	fmt.Printf("🔗 Access your API here: %s/generate\n", apiURL)
	fmt.Println("===========================================")

	http.ListenAndServe(":8080", nil)
}
