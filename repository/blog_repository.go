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

// UpdateBlogCounters increments/decrements the interaction counters for a blog
func (r *blogRepository) UpdateBlogCounters(ctx context.Context, blogID string, likeChange int, dislikeChange int, viewChange int) error {
	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	update := bson.M{
		"$inc": bson.M{
			"like_count":    likeChange,
			"dislike_count": dislikeChange,
			"view_count":    viewChange,
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// GetAllBlogs retrieves all blogs for popularity calculation
func (r *blogRepository) GetAllBlogs(ctx context.Context) ([]*entities.Blog, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
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

// FilterBlogs filters blogs based on provided criteria
func (r *blogRepository) FilterBlogs(ctx context.Context, filter *entities.BlogFilter) ([]*entities.Blog, int64, error) {
	// Build MongoDB filter
	mongoFilter := bson.M{}

	// Filter by tags
	if len(filter.Tags) > 0 {
		mongoFilter["tags"] = bson.M{"$in": filter.Tags}
	}

	// Filter by date range
	if filter.DateFrom != nil || filter.DateTo != nil {
		dateFilter := bson.M{}
		if filter.DateFrom != nil {
			dateFilter["$gte"] = *filter.DateFrom
		}
		if filter.DateTo != nil {
			dateFilter["$lte"] = *filter.DateTo
		}
		mongoFilter["created_at"] = dateFilter
	}

	// Build sort options
	sortOptions := bson.D{}
	if filter.PopularitySort != "" {
		sortField := "view_count"
		switch filter.PopularitySort {
		case "likes":
			sortField = "like_count"
		case "dislikes":
			sortField = "dislike_count"
		case "engagement":
			// For engagement, we'll use a combination (likes - dislikes + views)
			// MongoDB doesn't directly support computed sort, so we'll sort by likes first
			sortField = "like_count"
		case "views":
			sortField = "view_count"
		}

		sortOrder := -1 // desc by default
		if filter.SortOrder == "asc" {
			sortOrder = 1
		}
		sortOptions = append(sortOptions, bson.E{Key: sortField, Value: sortOrder})
	}

	// Default sort by created_at if no popularity sort
	if len(sortOptions) == 0 {
		sortOptions = append(sortOptions, bson.E{Key: "created_at", Value: -1})
	}

	// Get total count for pagination info
	totalCount, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, err
	}

	// Build find options
	findOptions := options.Find().SetSort(sortOptions)
	if filter.Limit > 0 {
		findOptions.SetLimit(int64(filter.Limit))
	}
	if filter.Skip > 0 {
		findOptions.SetSkip(int64(filter.Skip))
	}

	cursor, err := r.collection.Find(ctx, mongoFilter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var blogs []*entities.Blog
	for cursor.Next(ctx) {
		var blog entities.Blog
		if err := cursor.Decode(&blog); err != nil {
			return nil, 0, err
		}
		blogs = append(blogs, &blog)
	}

	return blogs, totalCount, cursor.Err()
}

// SearchBlogs searches for blogs based on title and/or author name
func (r *blogRepository) SearchBlogs(ctx context.Context, search *entities.BlogSearch) ([]*entities.BlogWithAuthor, int64, error) {
	// Build aggregation pipeline for searching with author lookup
	pipeline := []bson.M{}
	
	// Match stage - build search criteria
	matchStage := bson.M{}
	searchConditions := []bson.M{}
	
	// Search by title (case-insensitive partial match)
	if search.Title != "" {
		searchConditions = append(searchConditions, bson.M{
			"title": bson.M{
				"$regex":   search.Title,
				"$options": "i", // case-insensitive
			},
		})
	}
	
	// Search by author requires user lookup, so we'll add author filter after lookup
	if len(searchConditions) > 0 {
		if len(searchConditions) == 1 {
			matchStage = searchConditions[0]
		} else {
			matchStage["$and"] = searchConditions
		}
		pipeline = append(pipeline, bson.M{"$match": matchStage})
	}
	
	// Convert string user_id to ObjectID for lookup compatibility
	// Handle empty/invalid user_id gracefully
	pipeline = append(pipeline, bson.M{
		"$addFields": bson.M{
			"user_id_obj": bson.M{
				"$cond": bson.M{
					"if": bson.M{
						"$and": []bson.M{
							{"$ne": []interface{}{"$user_id", ""}},     // user_id is not empty
							{"$ne": []interface{}{"$user_id", nil}},    // user_id is not null
							{"$eq": []interface{}{bson.M{"$strLenCP": "$user_id"}, 24}}, // user_id length is 24 (valid ObjectID length)
						},
					},
					"then": bson.M{"$toObjectId": "$user_id"}, // Convert to ObjectID
					"else": nil,                                // Set to null if invalid
				},
			},
		},
	})
	
	// Lookup stage - join with user collection to get author name
	pipeline = append(pipeline, bson.M{
		"$lookup": bson.M{
			"from":         "user", // Fixed: collection name is "user" not "users"
			"localField":   "user_id_obj", // Use converted ObjectID field
			"foreignField": "_id",
			"as":           "author",
		},
	})
	
	// Unwind the author array (since lookup returns array)
	pipeline = append(pipeline, bson.M{
		"$unwind": bson.M{
			"path":                       "$author",
			"preserveNullAndEmptyArrays": true, // Keep blogs even if author not found
		},
	})
	
	// Add author_name field
	pipeline = append(pipeline, bson.M{
		"$addFields": bson.M{
			"author_name": bson.M{
				"$ifNull": []interface{}{
					"$author.username", // Use username as author name
					"Unknown Author",
				},
			},
		},
	})
	
	// Filter by author name if specified (after lookup)
	if search.Author != "" {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"author_name": bson.M{
					"$regex":   search.Author,
					"$options": "i", // case-insensitive
				},
			},
		})
	}
	
	// Count total documents (before skip/limit)
	countPipeline := append(pipeline, bson.M{"$count": "total"})
	countCursor, err := r.collection.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, err
	}
	defer countCursor.Close(ctx)
	
	var countResult []bson.M
	if err = countCursor.All(ctx, &countResult); err != nil {
		return nil, 0, err
	}
	
	var totalCount int64
	if len(countResult) > 0 {
		if count, ok := countResult[0]["total"].(int32); ok {
			totalCount = int64(count)
		}
	}
	
	// Add pagination
	if search.Skip > 0 {
		pipeline = append(pipeline, bson.M{"$skip": search.Skip})
	}
	if search.Limit > 0 {
		pipeline = append(pipeline, bson.M{"$limit": search.Limit})
	}
	
	// Remove the author object and temporary fields from final result
	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"author":     0, // Remove author object (we only need author_name)
			"user_id_obj": 0, // Remove temporary ObjectID field
		},
	})
	
	// Execute the aggregation
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	var results []*entities.BlogWithAuthor
	for cursor.Next(ctx) {
		var result entities.BlogWithAuthor
		if err := cursor.Decode(&result); err != nil {
			return nil, 0, err
		}
		results = append(results, &result)
	}
	
	return results, totalCount, cursor.Err()
}
