package entities

type PromptRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}