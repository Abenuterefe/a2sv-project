package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogInteraction struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	BlogID    primitive.ObjectID `bson:"blog_id" json:"blog_id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	IPAddress string             `bson:"ip_address,omitempty" json:"ip_address,omitempty"` // For anonymous users
	UserAgent string             `bson:"user_agent,omitempty" json:"user_agent,omitempty"` // For anonymous users
	Type      string             `bson:"type" json:"type"`                                 // "like", "dislike", "view"
	ExpiresAt *time.Time         `bson:"expires_at,omitempty" json:"expires_at,omitempty"` // For view expiration (24h)
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
