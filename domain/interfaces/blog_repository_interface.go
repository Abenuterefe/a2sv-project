package interfaces

import (
	"github.com/Abenuterefe/a2sv-project/domain/entities"
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
	// Update blog interaction counters (likes, dislikes, views)
	UpdateBlogCounters(ctx context.Context, blogID string, likeChange int, dislikeChange int, viewChange int) error
	// Get all blogs for popularity calculation
	GetAllBlogs(ctx context.Context) ([]*entities.Blog, error)
	// Filter blogs based on criteria
	FilterBlogs(ctx context.Context, filter *entities.BlogFilter) ([]*entities.Blog, int64, error)
	// Search blogs based on title and/or author
	SearchBlogs(ctx context.Context, search *entities.BlogSearch) ([]*entities.BlogWithAuthor, int64, error)
}
