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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

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

	var course Course
	modulesCount := 10

	for i := 1; i <= modulesCount; i++ {
		module, err := generateModule(reqBody.Language, i)
		if err != nil {
			http.Error(w, "Failed to generate module", http.StatusInternalServerError)
			return
		}
		course.Modules = append(course.Modules, module)

		time.Sleep(2 * time.Second)
	}

	courseJSON, _ := json.Marshal(course)

	w.Header().Set("Content-Type", "application/json")
	w.Write(courseJSON)
}

func generateModule(language string, moduleNumber int) (Module, error) {
	prompt := fmt.Sprintf(`
	Generate Module %d for a complete beginner to advanced course on %s.
	Provide a Module Heading and Module Content only.
	Format:
	- Heading: "Module Name"
	- Content: "Module Content"
	`, moduleNumber, language)

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

	apiKey := os.Getenv("GEMINI_API_KEY")
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", apiKey)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return Module{}, err
	}
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)

	var geminiResponse map[string]interface{}
	json.Unmarshal(responseBody, &geminiResponse)

	content := geminiResponse["candidates"].([]interface{})[0].(map[string]interface{})["content"].(map[string]interface{})["parts"].([]interface{})[0].(map[string]interface{})["text"].(string)

	return Module{
		Heading: fmt.Sprintf("Module %d", moduleNumber),
		Content: content,
	}, nil
}
