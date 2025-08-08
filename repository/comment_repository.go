package repository

import (
	"context"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"github.com/Abenuterefe/a2sv-project/domain/interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type commentRepository struct {
	collection *mongo.Collection
}

func NewCommentRepositoryMongo(collection *mongo.Collection) interfaces.CommentRepositoryInterface {
	return &commentRepository{collection: collection}
}

func (r *commentRepository) CreateComment(ctx context.Context, comment *entities.Comment) error {
	_, err := r.collection.InsertOne(ctx, comment)
	return err
}

// GetCommentsByBlogID retrieves all comments for a specific blog
func (r *commentRepository) GetCommentsByBlogID(ctx context.Context, blogID string) ([]*entities.Comment, error) {
	// Convert string blogID to ObjectID
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"blog_id": objID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []*entities.Comment
	for cursor.Next(ctx) {
		var comment entities.Comment
		if err := cursor.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, cursor.Err()
}

// GetCommentByID retrieves a single comment by its ID
func (r *commentRepository) GetCommentByID(ctx context.Context, id string) (*entities.Comment, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": oid}
	var comment entities.Comment
	if err := r.collection.FindOne(ctx, filter).Decode(&comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

// UpdateComment replaces an existing comment (matched by ID)
func (r *commentRepository) UpdateComment(ctx context.Context, comment *entities.Comment) error {
	filter := bson.M{"_id": comment.ID}
	_, err := r.collection.ReplaceOne(ctx, filter, comment)
	return err
}

// DeleteComment removes a comment by its ID
func (r *commentRepository) DeleteComment(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	_, err = r.collection.DeleteOne(ctx, filter)
	return err
}

// GetCommentCountByBlogID counts comments for a specific blog
func (r *commentRepository) GetCommentCountByBlogID(ctx context.Context, blogID string) (int64, error) {
	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return 0, err
	}
	filter := bson.M{"blog_id": blogObjID}
	return r.collection.CountDocuments(ctx, filter)
}
