package interfaces

import (
	"context"
)

type BlogInteractionUseCaseInterface interface {
	LikeBlog(ctx context.Context, blogID string, userID string) error
	DislikeBlog(ctx context.Context, blogID string, userID string) error
	ViewBlog(ctx context.Context, blogID string, userID string) error
}
