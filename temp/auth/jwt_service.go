package auth

import (
	"errors"
	"os"
	"time"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"github.com/golang-jwt/jwt/v5"
)

// Create jwt service object
type JWTService struct{}

// Creating new object
func NewJWTService() *JWTService {
	return &JWTService{}
}

var (
	AccessTokenTTL  = time.Minute * 15
	RefreshTokenTTL = time.Hour * 24 * 7
)

// CreateAccessToken generates a new JWT access token
func (s *JWTService) CreateAccessToken(userID, role string) (string, error) {
	claims := entities.JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
}

// CreateRefreshToken generates a new JWT refresh token
func (s *JWTService) CreateRefreshToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
}

// VerifyToken parses and validates a JWT token
func (s *JWTService) VerifyToken(tokenStr string, isAccess bool) (*entities.JWTClaims, error) {
	secret := os.Getenv("ACCESS_SECRET")
	if !isAccess {
		secret = os.Getenv("REFRESH_SECRET")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &entities.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(*entities.JWTClaims)
	if !ok {
		return nil, errors.New("could not parse claims")
	}

	return claims, nil
}
