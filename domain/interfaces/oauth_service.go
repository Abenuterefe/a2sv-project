package interfaces

import (
	"context"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
)

type OAuthService interface {
	GetAuthURL(state string) string
	GetUserInfo(ctx context.Context, code string) (*entities.GoogleUser, error)
}
