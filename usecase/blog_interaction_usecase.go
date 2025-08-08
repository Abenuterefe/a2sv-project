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
	// Check if user already has any interaction with this blog
	hasLike, _ := u.repo.HasInteraction(ctx, blogID, userID, "like")
	hasDislike, _ := u.repo.HasInteraction(ctx, blogID, userID, "dislike")
	
	if hasLike {
		// User already liked - remove like (toggle off)
		err := u.repo.RemoveInteraction(ctx, blogID, userID, "like")
		if err != nil {
			return err
		}
		// Decrease like count
		return u.blogRepo.UpdateBlogCounters(ctx, blogID, -1, 0, 0)
	}
	
	if hasDislike {
		// User had disliked - remove dislike, add like (switch)
		err := u.repo.RemoveInteraction(ctx, blogID, userID, "dislike")
		if err != nil {
			return err
		}
		// Add new like
		interaction := &entities.BlogInteraction{
			ID:        primitive.NewObjectID(),
			BlogID:    toObjectID(blogID),
			UserID:    userID,
			Type:      "like",
			CreatedAt: time.Now(),
		}
		err = u.repo.AddInteraction(ctx, interaction)
		if err != nil {
			return err
		}
		// Update counters: -1 dislike, +1 like
		return u.blogRepo.UpdateBlogCounters(ctx, blogID, 1, -1, 0)
	}
	
	// User hasn't interacted before - add like
	interaction := &entities.BlogInteraction{
		ID:        primitive.NewObjectID(),
		BlogID:    toObjectID(blogID),
		UserID:    userID,
		Type:      "like",
		CreatedAt: time.Now(),
	}
	err := u.repo.AddInteraction(ctx, interaction)
	if err != nil {
		return err
	}
	// Increase like count
	return u.blogRepo.UpdateBlogCounters(ctx, blogID, 1, 0, 0)
}

func (u *blogInteractionUseCase) DislikeBlog(ctx context.Context, blogID string, userID string) error {
	// Check if user already has any interaction with this blog
	hasLike, _ := u.repo.HasInteraction(ctx, blogID, userID, "like")
	hasDislike, _ := u.repo.HasInteraction(ctx, blogID, userID, "dislike")
	
	if hasDislike {
		// User already disliked - remove dislike (toggle off)
		err := u.repo.RemoveInteraction(ctx, blogID, userID, "dislike")
		if err != nil {
			return err
		}
		// Decrease dislike count
		return u.blogRepo.UpdateBlogCounters(ctx, blogID, 0, -1, 0)
	}
	
	if hasLike {
		// User had liked - remove like, add dislike (switch)
		err := u.repo.RemoveInteraction(ctx, blogID, userID, "like")
		if err != nil {
			return err
		}
		// Add new dislike
		interaction := &entities.BlogInteraction{
			ID:        primitive.NewObjectID(),
			BlogID:    toObjectID(blogID),
			UserID:    userID,
			Type:      "dislike",
			CreatedAt: time.Now(),
		}
		err = u.repo.AddInteraction(ctx, interaction)
		if err != nil {
			return err
		}
		// Update counters: -1 like, +1 dislike
		return u.blogRepo.UpdateBlogCounters(ctx, blogID, -1, 1, 0)
	}
	
	// User hasn't interacted before - add dislike
	interaction := &entities.BlogInteraction{
		ID:        primitive.NewObjectID(),
		BlogID:    toObjectID(blogID),
		UserID:    userID,
		Type:      "dislike",
		CreatedAt: time.Now(),
	}
	err := u.repo.AddInteraction(ctx, interaction)
	if err != nil {
		return err
	}
	// Increase dislike count
	return u.blogRepo.UpdateBlogCounters(ctx, blogID, 0, 1, 0)
}

func (u *blogInteractionUseCase) ViewBlog(ctx context.Context, blogID string, userID string, ipAddress string, userAgent string) error {
	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}
	
	// If not authenticated, use anonymous tracking
	if userID == "" {
		userID = "anonymous"
	}
	
	// Check if this is a recent duplicate view before adding
	hasRecent, err := u.repo.HasRecentView(ctx, blogID, userID, ipAddress, userAgent)
	if err != nil {
		return err
	}
	
	if hasRecent {
		// Already viewed recently, don't record or increment
		return nil
	}
	
	// Create view interaction with IP tracking
	interaction := &entities.BlogInteraction{
		ID:        primitive.NewObjectID(),
		BlogID:    blogObjID,
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Type:      "view",
		CreatedAt: time.Now(),
	}
	
	// Add the view interaction
	err = u.repo.AddInteraction(ctx, interaction)
	if err != nil {
		return err
	}
	
	// Increment the view counter
	return u.blogRepo.UpdateBlogCounters(ctx, blogID, 0, 0, 1)
}

func toObjectID(id string) primitive.ObjectID {
	objID, _ := primitive.ObjectIDFromHex(id)
	return objID
}
