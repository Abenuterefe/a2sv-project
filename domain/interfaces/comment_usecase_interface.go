package interfaces

import (
	"context"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
)

// CommentUseCaseInterface defines the contract for comment use case operations
type CommentUseCaseInterface interface {
	CreateComment(ctx context.Context, comment *entities.Comment, userID string, blogID string) error
	GetCommentsByBlogID(ctx context.Context, blogID string) ([]*entities.Comment, error)
	GetCommentByID(ctx context.Context, id string) (*entities.Comment, error)
	UpdateComment(ctx context.Context, comment *entities.Comment) error
	DeleteComment(ctx context.Context, id string) error
}
