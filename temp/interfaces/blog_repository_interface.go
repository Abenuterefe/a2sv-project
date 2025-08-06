package interfaces

import (
	"a2sv-project/Domain/entities"
	"context"
)

// BlogRepositoryInterface defines the contract for blog repository operations
type BlogRepositoryInterface interface {
	CreateBlog(ctx context.Context, blog *entities.Blog) error
	GetBlogsByUserID(ctx context.Context, userID string, page int64, limit int64) ([]*entities.Blog, error)
	// Get a single blog by its ID
	GetBlogByID(ctx context.Context, id string) (*entities.Blog, error)
	// Update an existing blog
	UpdateBlog(ctx context.Context, blog *entities.Blog) error
	// Delete a blog by its ID
	DeleteBlog(ctx context.Context, id string) error
}
