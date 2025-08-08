package entities

// BlogSearch represents the search criteria for blog posts
type BlogSearch struct {
	Title  string `json:"title,omitempty" form:"title"`
	Author string `json:"author,omitempty" form:"author"`
	Limit  int    `json:"limit,omitempty" form:"limit"`
	Skip   int    `json:"skip,omitempty" form:"skip"`
}

// SearchResponse represents the response structure for blog search results
type SearchResponse struct {
	Blogs      []*BlogWithAuthor `json:"blogs"`
	Count      int               `json:"count"`
	TotalCount int64             `json:"total_count,omitempty"`
	Query      *BlogSearch       `json:"query,omitempty"`
}

// BlogWithAuthor represents a blog with author information
type BlogWithAuthor struct {
	Blog       `bson:",inline"`
	AuthorName string `json:"author_name" bson:"author_name"`
}
