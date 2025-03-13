package util

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func GetResponseFromLLM(prompt string) string {
	requestBody := []byte(`
	{
		"model": "gemma3:1b", 
		"prompt": "` + prompt + `", 
		"stream": false, 
		"temperature": 0.8,
		"max_new_tokens": 75 
	}`)

	resp, err := http.Post("http://ollama:11434/api/generate", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		log.Printf("Error calling AI service: %v", err)
		return "Error generating response"
	}

	defer resp.Body.Close()
	var response struct {
		Answer string `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("Error decoding AI response: %v", err)
		return "Error processing response"
	}
	log.Printf("RAW Response from LLM: %+v", response)

	return response.Answer
}
