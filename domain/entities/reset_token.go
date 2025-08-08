package entities

import "time"

type ResetToken struct {
	UserID    string `bson:"user_id"`
	Token     string `bson:"token"`
	ExpiresAt time.Time `bson:"expires_at"`
}