package interfaces

import (
	"context"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
)

type UserUsecase interface {
	Regiser(ctx context.Context, user *entities.User) error
	Login(ctx context.Context, email, password string) (*entities.Token, error)
	RefreshToken(ctx context.Context, refreshToken string) (*entities.Token, error)
}
