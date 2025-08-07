package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogInteraction struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	BlogID    primitive.ObjectID `bson:"blog_id" json:"blog_id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Type      string             `bson:"type" json:"type"` // "like", "dislike", "view"
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
