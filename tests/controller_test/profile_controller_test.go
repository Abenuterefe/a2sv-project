package controller_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
    "mime/multipart"
	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"github.com/Abenuterefe/a2sv-project/delivery/controllers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProfileUsecase is a testify mock that implements interfaces.ProfileUsecase
type MockProfileUsecase struct {
	mock.Mock
}

func (m *MockProfileUsecase) UpdateProfile(ctx context.Context, userID, username, bio, profilePicture string) error {
	args := m.Called(ctx, userID, username, bio, profilePicture)
	return args.Error(0)
}

func (m *MockProfileUsecase) GetProfile(ctx context.Context, userID string) (*entities.Profile, error) {
	args := m.Called(ctx, userID)
	if profile, ok := args.Get(0).(*entities.Profile); ok {
		return profile, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProfileUsecase) UploadProfilePicture(ctx context.Context, userID string, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	args := m.Called(ctx, userID, file, fileHeader)
	return args.String(0), args.Error(1)
}

func TestGetProfile_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockUC := new(MockProfileUsecase)

	expectedProfile := &entities.Profile{
		ID:             "1",
		UserRole:       "user",
		UserID:         "12345",
		Bio:            "Hello World",
		ProfilePicture: "http://localhost:8080/pic.jpg",
		Username:       "johndoe",
		Email:          "john@example.com",
	}

	mockUC.On("GetProfile", mock.Anything, "12345").Return(expectedProfile, nil)

	controller := controllers.NewProfileController(mockUC)

	// Set up Gin with fake middleware to set userID
	router := gin.Default()
	router.GET("/user/profile/me", func(c *gin.Context) {
		c.Set("userID", "12345")
		controller.GetProfile(c)
	})

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/user/profile/me", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var got entities.Profile
	err := json.Unmarshal(rec.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, expectedProfile.Username, got.Username)
	assert.Equal(t, expectedProfile.Email, got.Email)

	mockUC.AssertExpectations(t)
}

func TestGetProfile_Error(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockUC := new(MockProfileUsecase)

	mockUC.On("GetProfile", mock.Anything, "12345").Return(nil, assert.AnError)

	controller := controllers.NewProfileController(mockUC)

	router := gin.Default()
	router.GET("/user/profile/me", func(c *gin.Context) {
		c.Set("userID", "12345")
		controller.GetProfile(c)
	})

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/user/profile/me", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockUC.AssertExpectations(t)
}
