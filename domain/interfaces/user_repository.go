package interfaces

import (
	"context"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	// Functions to access user
	Create(ctx context.Context, user *entities.User) error
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByID(ctx context.Context, id primitive.ObjectID)(*entities.User, error)
	UpdateUsername(ctx context.Context, userID primitive.ObjectID, username string) error
	// Store, access, delete jwt token to the user
	StoreToken(ctx context.Context, token *entities.Token) error
	FindToken(ctx context.Context, refreshToken string) (*entities.Token, error)
	DeleteToken(ctx context.Context, refreshToken string) error

	// emial verification funcs
	FindByVerificationToken(ctx context.Context, token string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	UpdatePassword(ctx context.Context, userID primitive.ObjectID, hashedPassword string) error

	DeleteTokenByUserID(ctx context.Context, userID string) error

	//Reset token repository
	SaveResetToken(ctx context.Context, token *entities.ResetToken) error
	FindByResetToken(ctx context.Context, token string)(*entities.ResetToken, error)
	DeleteResetToken(ctx context.Context,token string) error
}
