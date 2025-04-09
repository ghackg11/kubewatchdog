package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GenerateLlmResponse(prompt string) (string, error) {
	var llmUrl = "http://phi2:11435/api/generate"
	var model = "phi"
	requestBody := map[string]any{
		"model":  model,
		"prompt": prompt,
		"stream": false,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post(llmUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	var response struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	return response.Response, nil
}
