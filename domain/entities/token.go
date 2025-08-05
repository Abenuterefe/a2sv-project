package entities

import (
	"time"
)

// Token holds access & refresh tokens and their metadata
type Token struct {
	UserID       string    `bson:"user_id" json:"user_id"`
	AccessToken  string    `bson:"access_token" json:"access_token"`
	RefreshToken string    `bson:"refresh_token" json:"refresh_token"`
	ExpiresAt    time.Time `bson:"expires_at" json:"expires_at"`
}
