package entities

import "github.com/golang-jwt/jwt/v5"

// Custom claims structure
type JWTClaims struct {
	UserID string `json:"sub"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
 