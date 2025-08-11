package entities

type BlogResponse struct {
	Title          string   `json:"title,omitempty"`
	Paragraphs     []string `json:"paragraphs"`
	ParagraphCount int      `json:"paragraph_count"`
}
