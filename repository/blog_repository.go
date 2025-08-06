package repository

import (
	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"github.com/Abenuterefe/a2sv-project/domain/interfaces"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type blogRepository struct {
	collection *mongo.Collection
}

func NewBlogRepositoryMongo(collection *mongo.Collection) interfaces.BlogRepositoryInterface {
	return &blogRepository{collection: collection}
}

func (r *blogRepository) CreateBlog(ctx context.Context, blog *entities.Blog) error {
	_, err := r.collection.InsertOne(ctx, blog)
	return err
}

// GetBlogsByUserID retrieves paginated blogs for a user
func (r *blogRepository) GetBlogsByUserID(ctx context.Context, userID string, page int64, limit int64) ([]*entities.Blog, error) {
	filter := bson.M{"user_id": userID}
	if page < 1 {
		page = 1
	}
	skip := (page - 1) * limit
	opts := options.Find().SetSkip(skip).SetLimit(limit)
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var blogs []*entities.Blog
	for cursor.Next(ctx) {
		var blog entities.Blog
		if err := cursor.Decode(&blog); err != nil {
			return nil, err
		}
		blogs = append(blogs, &blog)
	}
	return blogs, cursor.Err()
}

// GetBlogByID retrieves a single blog by its ID
func (r *blogRepository) GetBlogByID(ctx context.Context, id string) (*entities.Blog, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": oid}
	var blog entities.Blog
	if err := r.collection.FindOne(ctx, filter).Decode(&blog); err != nil {
		return nil, err
	}
	return &blog, nil
}

// UpdateBlog replaces an existing blog (matched by ID)
func (r *blogRepository) UpdateBlog(ctx context.Context, blog *entities.Blog) error {
	filter := bson.M{"_id": blog.ID}
	_, err := r.collection.ReplaceOne(ctx, filter, blog)
	return err
}

// DeleteBlog deletes a blog by its ID
func (r *blogRepository) DeleteBlog(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	_, err = r.collection.DeleteOne(ctx, filter)
	return err
}
