package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogWithPopularity struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title           string             `bson:"title" json:"title"`
	Content         string             `bson:"content" json:"content"`
	UserID          string             `bson:"user_id" json:"user_id"`
	LikeCount       int                `bson:"like_count" json:"like_count"`
	DislikeCount    int                `bson:"dislike_count" json:"dislike_count"`
	ViewCount       int                `bson:"view_count" json:"view_count"`
	CommentCount    int                `json:"comment_count"`
	PopularityScore float64            `json:"popularity_score"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}
