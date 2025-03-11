package main

import (
	"fmt"
	"net/http"

	"gemini/handler"
	"gemini/utils"
)

func main() {
	utils.LoadEnv()

	// 🔥 Change your API Endpoint here
	http.HandleFunc("/generate", handler.GenerateContent)

	// ✅ Fancy server logs
	apiURL := "http://localhost:8080"
	fmt.Println("===========================================")
	fmt.Printf("✅ Server running on: %s\n", apiURL)
	fmt.Printf("🔗 Access your API here: %s/generate\n", apiURL)
	fmt.Println("===========================================")

	// 🚀 Start the server
	http.ListenAndServe(":8080", nil)
}
