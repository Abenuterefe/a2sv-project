// file: usecase/profile_usecase_test.go
package usecase_test

import (
    "context"
    "errors"
    "testing"

    "github.com/Abenuterefe/a2sv-project/domain/entities"
    "github.com/Abenuterefe/a2sv-project/domain/interfaces"
    "github.com/stretchr/testify/assert"
)

// mockProfileRepo is a mock implementation of ProfileRepository
type mockProfileRepo struct {
    profile *entities.Profile
    err     error
}

func (m *mockProfileRepo) FindByUserID(ctx context.Context, userID string) (*entities.Profile, error) {
    return m.profile, m.err
}

// Unused methods - just dummy implementations
func (m *mockProfileRepo) UpdateProfile(ctx context.Context, profile *entities.Profile) error {
    return nil
}
func (m *mockProfileRepo) UpdateProfilePicture(ctx context.Context, userID string, picturePath string) error {
    return nil
}

// minimal ProfileUsecase implementation for test
type profileUsecase struct {
    repo interfaces.ProfileRepository
}

func (u *profileUsecase) GetProfile(ctx context.Context, userID string) (*entities.Profile, error) {
    return u.repo.FindByUserID(ctx, userID)
}
func (u *profileUsecase) UpdateProfile(ctx context.Context, userID, username, bio, profilePicture string) error {
    return nil
}
func (u *profileUsecase) UploadProfilePicture(ctx context.Context, userID string, file interface{}, fileHeader interface{}) (string, error) {
    return "", nil
}

func TestGetProfile(t *testing.T) {
    t.Run("Success", func(t *testing.T) {
        mockRepo := &mockProfileRepo{
            profile: &entities.Profile{
                UserID:         "123",
                Bio:            "Hello world",
                ProfilePicture: "pic.jpg",
            },
            err: nil,
        }
        usecase := &profileUsecase{repo: mockRepo}
        profile, err := usecase.GetProfile(context.Background(), "123")

        assert.NoError(t, err)
        assert.NotNil(t, profile)
        assert.Equal(t, "123", profile.UserID)
        assert.Equal(t, "Hello world", profile.Bio)
        assert.Equal(t, "pic.jpg", profile.ProfilePicture)
    })

    t.Run("Failure - repository error", func(t *testing.T) {
        mockRepo := &mockProfileRepo{
            profile: nil,
            err:     errors.New("database error"),
        }
        usecase := &profileUsecase{repo: mockRepo}
        profile, err := usecase.GetProfile(context.Background(), "123")

        assert.Error(t, err)
        assert.Nil(t, profile)
        assert.Equal(t, "database error", err.Error())
    })

    t.Run("Failure - profile not found", func(t *testing.T) {
        mockRepo := &mockProfileRepo{
            profile: nil,
            err:     nil, // no error but no profile found
        }
        usecase := &profileUsecase{repo: mockRepo}
        profile, err := usecase.GetProfile(context.Background(), "123")

        assert.NoError(t, err)
        assert.Nil(t, profile)
    })
}
