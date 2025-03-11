package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type RequestBody struct {
	Language string `json:"language"`
}

type Module struct {
	Heading string `json:"Heading"`
	Content string `json:"Content"`
}

type Course struct {
	Modules []Module `json:"modules"`
}

func GenerateContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read user input
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Parse the user's input (Language Name)
	var reqBody RequestBody
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if reqBody.Language == "" {
		http.Error(w, "Language cannot be empty", http.StatusBadRequest)
		return
	}

	// Generate 9-11 modules by calling Gemini API in a loop
	var course Course
	modulesCount := 10

	for i := 1; i <= modulesCount; i++ {
		module, err := generateModule(reqBody.Language, i)
		if err != nil {
			http.Error(w, "Failed to generate module", http.StatusInternalServerError)
			return
		}
		course.Modules = append(course.Modules, module)

		// Sleep for 2 seconds to avoid rate limiting from Google
		time.Sleep(2 * time.Second)
	}

	// Convert the entire course to JSON
	courseJSON, _ := json.Marshal(course)

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(courseJSON)
}

func generateModule(language string, moduleNumber int) (Module, error) {
	// Prepare the prompt dynamically for each module
	prompt := fmt.Sprintf(`
	Generate Module %d for a complete beginner to advanced course on %s.
	Provide a Module Heading and Module Content only.
	Format:
	- Heading: "Module Name"
	- Content: "Module Content"
	`, moduleNumber, language)

	// Construct the request body for Gemini API
	geminiBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{
						"text": prompt,
					},
				},
			},
		},
	}

	jsonBody, _ := json.Marshal(geminiBody)

	// Fetch API Key
	apiKey := os.Getenv("GEMINI_API_KEY")
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", apiKey)

	// Call the Gemini API
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return Module{}, err
	}
	defer resp.Body.Close()

	// Read the response
	responseBody, _ := ioutil.ReadAll(resp.Body)

	// Extract Module Content from Gemini API Response
	var geminiResponse map[string]interface{}
	json.Unmarshal(responseBody, &geminiResponse)

	content := geminiResponse["candidates"].([]interface{})[0].(map[string]interface{})["content"].(map[string]interface{})["parts"].([]interface{})[0].(map[string]interface{})["text"].(string)

	// Split Heading and Content
	return Module{
		Heading: fmt.Sprintf("Module %d", moduleNumber),
		Content: content,
	}, nil
}
