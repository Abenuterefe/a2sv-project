package interfaces

import (
	"context"
	"mime/multipart"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
)

type ProfileUsecase interface {
	UpdateProfile(ctx context.Context, userID, username, bio, profilePicture string) error
	GetProfile(ctx context.Context, userID string) (*entities.Profile, error)
	UploadProfilePicture(ctx context.Context, userID string, file multipart.File, fileHeader *multipart.FileHeader) (string, error) // NEW
}
