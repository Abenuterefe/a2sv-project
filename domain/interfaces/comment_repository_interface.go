package interfaces

import (
	"context"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
)

// CommentRepositoryInterface defines the contract for comment repository operations
type CommentRepositoryInterface interface {
	CreateComment(ctx context.Context, comment *entities.Comment) error
	GetCommentsByBlogID(ctx context.Context, blogID string) ([]*entities.Comment, error)
	GetCommentByID(ctx context.Context, id string) (*entities.Comment, error)
	UpdateComment(ctx context.Context, comment *entities.Comment) error
	DeleteComment(ctx context.Context, id string) error
	GetCommentCountByBlogID(ctx context.Context, blogID string) (int64, error)
}
