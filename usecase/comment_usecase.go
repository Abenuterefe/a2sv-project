package usecase

import (
	"context"
	"time"

	"github.com/Abenuterefe/a2sv-project/domain/interfaces"

	"github.com/Abenuterefe/a2sv-project/domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// commentUseCase implements the CommentUseCaseInterface
type commentUseCase struct {
	repo interfaces.CommentRepositoryInterface
}

func NewCommentUseCase(repo interfaces.CommentRepositoryInterface) interfaces.CommentUseCaseInterface {
	return &commentUseCase{repo: repo}
}

func (u *commentUseCase) CreateComment(ctx context.Context, comment *entities.Comment, userID string, blogID string) error {
	// Generate a new ObjectID for the comment
	comment.ID = primitive.NewObjectID()

	// Convert blogID string to ObjectID
	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}

	// Set the user ID and blog ID
	comment.UserID = userID
	comment.BlogID = blogObjID

	// Set timestamps
	now := time.Now()
	comment.CreatedAt = now
	comment.UpdatedAt = now

	return u.repo.CreateComment(ctx, comment)
}

// GetCommentsByBlogID returns all comments for a blog
func (u *commentUseCase) GetCommentsByBlogID(ctx context.Context, blogID string) ([]*entities.Comment, error) {
	return u.repo.GetCommentsByBlogID(ctx, blogID)
}

// GetCommentByID returns a single comment by ID
func (u *commentUseCase) GetCommentByID(ctx context.Context, id string) (*entities.Comment, error) {
	return u.repo.GetCommentByID(ctx, id)
}

// UpdateComment updates an existing comment
func (u *commentUseCase) UpdateComment(ctx context.Context, comment *entities.Comment) error {
	// Update timestamp
	comment.UpdatedAt = time.Now()
	return u.repo.UpdateComment(ctx, comment)
}

// DeleteComment removes a comment by ID
func (u *commentUseCase) DeleteComment(ctx context.Context, id string) error {
	return u.repo.DeleteComment(ctx, id)
}
