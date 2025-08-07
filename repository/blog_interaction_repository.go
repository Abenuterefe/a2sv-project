package repository

import (
	"context"

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
	// Prevent duplicate like/dislike/view for the same user/blog/type
	filter := bson.M{"blog_id": interaction.BlogID, "user_id": interaction.UserID, "type": interaction.Type}
	update := bson.M{"$setOnInsert": interaction}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
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
