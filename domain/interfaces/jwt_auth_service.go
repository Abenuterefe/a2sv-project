package interfaces

import "github.com/Abenuterefe/a2sv-project/domain/entities"

type AuthService interface {
	CreateAccessToken(userID, role string) (string, error)
	CreateRefreshToken(userID, role string) (string, error)
	VerifyToken(tokenStr string, isAccess bool) (*entities.JWTClaims, error)
}
