package usecase

import (
	"github.com/Abenuterefe/a2sv-project/domain/interfaces"
	"context"
	"time"

	"github.com/Abenuterefe/a2sv-project/domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// blogUseCase implements the BlogUseCaseInterface
type blogUseCase struct {
	repo interfaces.BlogRepositoryInterface
}

func NewBlogUseCase(repo interfaces.BlogRepositoryInterface) interfaces.BlogUseCaseInterface {
	return &blogUseCase{repo: repo}
}

func (u *blogUseCase) CreateBlog(ctx context.Context, blog *entities.Blog, userID string) error {
	// Generate a new ObjectID for the blog
	blog.ID = primitive.NewObjectID()

	// Set the user ID
	blog.UserID = userID

	// Set timestamps
	now := time.Now()
	blog.CreatedAt = now
	blog.UpdatedAt = now

	return u.repo.CreateBlog(ctx, blog)
}

// GetBlogsByUserID returns paginated blogs for a user
func (u *blogUseCase) GetBlogsByUserID(ctx context.Context, userID string, page int64, limit int64) ([]*entities.Blog, error) {
	return u.repo.GetBlogsByUserID(ctx, userID, page, limit)
}

// GetBlogByID returns a single blog by ID
func (u *blogUseCase) GetBlogByID(ctx context.Context, id string) (*entities.Blog, error) {
	return u.repo.GetBlogByID(ctx, id)
}

// UpdateBlog updates an existing blog
func (u *blogUseCase) UpdateBlog(ctx context.Context, blog *entities.Blog) error {
	// Update timestamp
	blog.UpdatedAt = time.Now()
	return u.repo.UpdateBlog(ctx, blog)
}

// DeleteBlog removes a blog by ID
func (u *blogUseCase) DeleteBlog(ctx context.Context, id string) error {
	return u.repo.DeleteBlog(ctx, id)
}
