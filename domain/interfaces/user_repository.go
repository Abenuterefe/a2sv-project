package interfaces

import (
	"context"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
)

type UserRepository interface {
	// Functions to access user
	Create(ctx context.Context, user *entities.User) error
	FindByEmail(ctx context.Context, email string) (*entities.User, error)

	// Store, access, delete jwt token to the user
	StoreToken(ctx context.Context, token *entities.Token) error
	FindToken(ctx context.Context, refreshToken string) (*entities.Token, error)
	DeleteToken(ctx context.Context, refreshToken string) error
}
