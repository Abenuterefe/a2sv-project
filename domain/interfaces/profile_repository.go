package interfaces

import (
	"context"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
)

type ProfileRepository interface {
	UpdateProfile(ctx context.Context, profile *entities.Profile) error
	FindByUserID(ctx context.Context, userID string) (*entities.Profile, error)
	UpdateProfilePicture(ctx context.Context, userID string, picturePath string) error // NEW
}
