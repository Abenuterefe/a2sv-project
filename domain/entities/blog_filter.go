package entities

import "time"

// BlogFilter represents the filtering criteria for blog posts
type BlogFilter struct {
	Tags           []string   `json:"tags,omitempty" form:"tags"`
	DateFrom       *time.Time `json:"date_from,omitempty" form:"date_from"`
	DateTo         *time.Time `json:"date_to,omitempty" form:"date_to"`
	PopularitySort string     `json:"popularity_sort,omitempty" form:"popularity_sort"` // "views", "likes", "engagement", "dislikes"
	SortOrder      string     `json:"sort_order,omitempty" form:"sort_order"`           // "asc", "desc"
	Limit          int        `json:"limit,omitempty" form:"limit"`
	Skip           int        `json:"skip,omitempty" form:"skip"`
}

// FilterResponse represents the response structure for filtered blogs
type FilterResponse struct {
	Blogs      []*Blog `json:"blogs"`
	Count      int     `json:"count"`
	TotalCount int64   `json:"total_count,omitempty"`
	Page       int     `json:"page,omitempty"`
	Limit      int     `json:"limit,omitempty"`
}
