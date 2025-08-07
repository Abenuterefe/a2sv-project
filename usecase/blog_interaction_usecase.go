package usecase

import (
	"context"
	"time"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"github.com/Abenuterefe/a2sv-project/domain/interfaces"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type blogInteractionUseCase struct {
	repo     interfaces.BlogInteractionRepositoryInterface
	blogRepo interfaces.BlogRepositoryInterface
}

func NewBlogInteractionUseCase(repo interfaces.BlogInteractionRepositoryInterface, blogRepo interfaces.BlogRepositoryInterface) interfaces.BlogInteractionUseCaseInterface {
	return &blogInteractionUseCase{
		repo:     repo,
		blogRepo: blogRepo,
	}
}

func (u *blogInteractionUseCase) LikeBlog(ctx context.Context, blogID string, userID string) error {
	interaction := &entities.BlogInteraction{
		ID:        primitive.NewObjectID(),
		BlogID:    toObjectID(blogID),
		UserID:    userID,
		Type:      "like",
		CreatedAt: time.Now(),
	}

	// First, add the interaction
	err := u.repo.AddInteraction(ctx, interaction)
	if err != nil {
		return err
	}

	// Then update the blog's like counter
	return u.blogRepo.UpdateBlogCounters(ctx, blogID, 1, 0, 0)
}

func (u *blogInteractionUseCase) DislikeBlog(ctx context.Context, blogID string, userID string) error {
	interaction := &entities.BlogInteraction{
		ID:        primitive.NewObjectID(),
		BlogID:    toObjectID(blogID),
		UserID:    userID,
		Type:      "dislike",
		CreatedAt: time.Now(),
	}

	// First, add the interaction
	err := u.repo.AddInteraction(ctx, interaction)
	if err != nil {
		return err
	}

	// Then update the blog's dislike counter
	return u.blogRepo.UpdateBlogCounters(ctx, blogID, 0, 1, 0)
}

func (u *blogInteractionUseCase) ViewBlog(ctx context.Context, blogID string, userID string) error {
	interaction := &entities.BlogInteraction{
		ID:        primitive.NewObjectID(),
		BlogID:    toObjectID(blogID),
		UserID:    userID,
		Type:      "view",
		CreatedAt: time.Now(),
	}

	// First, add the interaction
	err := u.repo.AddInteraction(ctx, interaction)
	if err != nil {
		return err
	}

	// Then update the blog's view counter
	return u.blogRepo.UpdateBlogCounters(ctx, blogID, 0, 0, 1)
}

func toObjectID(id string) primitive.ObjectID {
	objID, _ := primitive.ObjectIDFromHex(id)
	return objID
}
