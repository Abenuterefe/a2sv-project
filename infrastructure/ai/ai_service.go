package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Constants for OpenRouter
const (
	openRouterURL = "https://openrouter.ai/api/v1/chat/completions"
	apiKey        = "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"
	model         = "google/gemini-2.0-flash-exp:free"
)

type OpenAIService struct{}

func NewOpenAIService() *OpenAIService {
	return &OpenAIService{}
}

type openRouterRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string    `json:"role"`
	Content []content `json:"content"`
}

type content struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

type openRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (s *OpenAIService) GenerateBlog(prompt string) (string, error) {
	if prompt == "" {
		return "", errors.New("prompt cannot be empty")
	}

	// Build the request body
	reqBody := openRouterRequest{
		Model: model,
		Messages: []message{
			{
				Role: "user",
				Content: []content{
					{Type: "text", Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", openRouterURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("OpenRouter API error: " + string(body))
	}

	var aiResp openRouterResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		return "", err
	}

	if len(aiResp.Choices) == 0 {
		return "", errors.New("AI returned empty response")
	}

	return aiResp.Choices[0].Message.Content, nil
}