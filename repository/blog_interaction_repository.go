package repository

import (
	"context"
	"time"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"github.com/Abenuterefe/a2sv-project/domain/interfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type blogInteractionRepository struct {
	collection *mongo.Collection
}

func NewBlogInteractionRepositoryMongo(collection *mongo.Collection) interfaces.BlogInteractionRepositoryInterface {
	return &blogInteractionRepository{collection: collection}
}

func (r *blogInteractionRepository) AddInteraction(ctx context.Context, interaction *entities.BlogInteraction) error {
	// For likes/dislikes, prevent duplicates with upsert
	if interaction.Type == "like" || interaction.Type == "dislike" {
		filter := bson.M{"blog_id": interaction.BlogID, "user_id": interaction.UserID, "type": interaction.Type}
		update := bson.M{"$setOnInsert": interaction}
		opts := options.Update().SetUpsert(true)
		_, err := r.collection.UpdateOne(ctx, filter, update, opts)
		return err
	}

	// For views, set expiration time (24 hours from now)
	if interaction.Type == "view" {
		expiresAt := time.Now().Add(24 * time.Hour)
		interaction.ExpiresAt = &expiresAt
	}

	// Insert the interaction
	_, err := r.collection.InsertOne(ctx, interaction)
	return err
}

func (r *blogInteractionRepository) RemoveInteraction(ctx context.Context, blogID string, userID string, interactionType string) error {
	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}
	filter := bson.M{"blog_id": blogObjID, "user_id": userID, "type": interactionType}
	_, err = r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *blogInteractionRepository) HasInteraction(ctx context.Context, blogID string, userID string, interactionType string) (bool, error) {
	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return false, err
	}
	filter := bson.M{"blog_id": blogObjID, "user_id": userID, "type": interactionType}
	count, err := r.collection.CountDocuments(ctx, filter)
	return count > 0, err
}

func (r *blogInteractionRepository) HasRecentView(ctx context.Context, blogID string, userID string, ipAddress string, userAgent string) (bool, error) {
	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return false, err
	}

	// Build filter for recent views (within 24 hours)
	filter := bson.M{
		"blog_id":    blogObjID,
		"type":       "view",
		"expires_at": bson.M{"$gt": time.Now()}, // Not expired yet
	}

	// For authenticated users, check by userID
	if userID != "" && userID != "anonymous" {
		filter["user_id"] = userID
	} else {
		// For anonymous users, check by IP + User-Agent combo
		filter["user_id"] = "anonymous"
		filter["ip_address"] = ipAddress
		filter["user_agent"] = userAgent
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	return count > 0, err
}
