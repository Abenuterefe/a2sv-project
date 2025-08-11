// moved from blog_usecase.go for clarity
package interfaces

import (
	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"context"
)

// BlogUseCaseInterface defines the contract for blog use case operations
// This interface should be used in the usecase implementation
type BlogUseCaseInterface interface {
	CreateBlog(ctx context.Context, blog *entities.Blog, userID string) error
	GetBlogsByUserID(ctx context.Context, userID string, page int64, limit int64) ([]*entities.Blog, error)
	// Get a single blog by its ID
	GetBlogByID(ctx context.Context, id string) (*entities.Blog, error)
	// Update an existing blog (fields must include ID)
	UpdateBlog(ctx context.Context, blog *entities.Blog) error
	// Delete a blog by its ID
	DeleteBlog(ctx context.Context, id string) error
	// Get popular blogs with popularity scores
	GetPopularBlogs(ctx context.Context, limit int64) ([]*entities.BlogWithPopularity, error)
	// Filter blogs based on criteria
	FilterBlogs(ctx context.Context, filter *entities.BlogFilter) (*entities.FilterResponse, error)
	// Search blogs based on title and/or author
	SearchBlogs(ctx context.Context, search *entities.BlogSearch) (*entities.SearchResponse, error)
}
