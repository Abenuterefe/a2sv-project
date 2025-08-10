package usecase

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"github.com/Abenuterefe/a2sv-project/domain/interfaces"
	"github.com/Abenuterefe/a2sv-project/infrastructure/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type profileUsecase struct {
	userRepo     interfaces.UserRepository
	profileRepo  interfaces.ProfileRepository
	fileStorage  storage.FileStorage
}

func NewProfileUsecase(
	userRepo interfaces.UserRepository,
	profileRepo interfaces.ProfileRepository,
	fileStorage storage.FileStorage,
) interfaces.ProfileUsecase {
	return &profileUsecase{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		fileStorage: fileStorage,
	}
}

// UpdateProfile updates username, bio, and optionally profile picture path
func (uc *profileUsecase) UpdateProfile(ctx context.Context, userID, username, bio, profilePicture string) error {
	if username != "" {
		objectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return errors.New("invalid user ID")
		}
		if err := uc.userRepo.UpdateUsername(ctx, objectID, username); err != nil {
			return err
		}
	}

	profile := &entities.Profile{
		UserID: userID,
	}

	shouldUpdateProfile := false

	if bio != "" {
		profile.Bio = bio
		shouldUpdateProfile = true
	}

	if profilePicture != "" {
		profile.ProfilePicture = profilePicture
		shouldUpdateProfile = true
	}

	if shouldUpdateProfile {
		if err := uc.profileRepo.UpdateProfile(ctx, profile); err != nil {
			return err
		}
	}

	return nil
}

// GetProfile returns merged user + profile info
func (uc *profileUsecase) GetProfile(ctx context.Context, userID string) (*entities.Profile, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := uc.userRepo.FindByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	profile, err := uc.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := &entities.Profile{
		UserID:         userID,
		UserRole:     	string(user.Role),
		Username:       user.Username,
		Email:          user.Email,
		Bio:            "",
		ProfilePicture: "",
	}

	if profile != nil {
		result.Bio = profile.Bio
		result.ProfilePicture = "http://localhost:8080/"+profile.ProfilePicture
	}

	return result, nil
}

// UploadProfilePicture delegates file saving to storage layer and updates DB with path
// UploadProfilePicture delegates file saving to storage layer and updates DB with path
func (uc *profileUsecase) UploadProfilePicture(
	ctx context.Context,
	userID string,
	file multipart.File,
	fileHeader *multipart.FileHeader,
) (string, error) {
	// Validate user ID format
	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", errors.New("invalid user ID")
	}

	// Store file via infrastructure
	filePath, err := uc.fileStorage.SaveProfilePicture(file, fileHeader)
	if err != nil {
		return "", err
	}

	// Update profile picture field in profile table/collection
	profile := &entities.Profile{
		UserID:         userID,
		ProfilePicture: filePath,
	}

	if err := uc.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return "", err
	}

	return filePath, nil
}
