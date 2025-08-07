package interfaces

import (
	"context"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
)

type UserUsecase interface {
	Regiser(ctx context.Context, user *entities.User) error
	Login(ctx context.Context, email, password string) (*entities.Token, error)
	RefreshToken(ctx context.Context, refreshToken string) (*entities.Token, error)
	VerifyEmail(ctx context.Context, token string) error
	ResendVerificationEmail(ctx context.Context, email string) error
	PromoteUser(ctx context.Context, userID string) error
	DemoteUser(ctx context.Context, userID string) error
	Logout(ctx context.Context, userID string) error
}
