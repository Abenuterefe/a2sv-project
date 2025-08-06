package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID           primitive.ObjectID    `bson:"_id,omitempty"`
	UserID       string    `bson:"user_id"`
	Title        string    `bson:"title"`
	Content      string    `bson:"content"`
	Tags         []string  `bson:"tags"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
	ViewCount    int       `bson:"view_count"`
	LikeCount    int       `bson:"like_count"`
	DislikeCount int       `bson:"dislike_count"`
}

