package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
)

const (
	openRouterURL = "https://openrouter.ai/api/v1/chat/completions"
	apiKey        = "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"
	model         = "deepseek/deepseek-chat-v3-0324:free"
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
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (s *OpenAIService) GenerateBlog(prompt string) (*entities.BlogResponse, error) {
	if prompt == "" {
		return nil, errors.New("prompt cannot be empty")
	}

	if !strings.Contains(strings.ToLower(prompt), "blog") {
		prompt = "Generate a blog regarding: " + prompt
	}

	reqBody := openRouterRequest{
		Model: model,
		Messages: []message{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", openRouterURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("OpenRouter API error: " + string(body))
	}

	var aiResp openRouterResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		return nil, err
	}

	if len(aiResp.Choices) == 0 {
		return nil, errors.New("AI returned empty response")
	}

	content := strings.TrimSpace(aiResp.Choices[0].Message.Content)
	paragraphs := splitIntoParagraphs(content)

	return &entities.BlogResponse{
		Title:          extractTitle(content),
		Paragraphs:     paragraphs,
		ParagraphCount: len(paragraphs),
	}, nil
}

func splitIntoParagraphs(text string) []string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	parts := strings.Split(text, "\n\n")

	var cleaned []string
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			cleaned = append(cleaned, p)
		}
	}
	return cleaned
}

func extractTitle(text string) string {
	re := regexp.MustCompile(`^(.*?)(\.|\n|$)`)
	match := re.FindStringSubmatch(text)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}
