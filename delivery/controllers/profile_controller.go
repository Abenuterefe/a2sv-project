package controllers

import (
	"fmt"
	"net/http"

	"github.com/Abenuterefe/a2sv-project/domain/interfaces"
	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	profileUsecase interfaces.ProfileUsecase
}

func NewProfileController(profileUsecase interfaces.ProfileUsecase) *ProfileController {
	return &ProfileController{
		profileUsecase: profileUsecase,
	}
}

type UpdateProfileRequest struct {
	Username       string `json:"username"`
	Bio            string `json:"bio"`
	ProfilePicture string `json:"profilePicture"`
}

func (pc *ProfileController) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID := c.GetString("userID") // assume middleware already sets this

	err := pc.profileUsecase.UpdateProfile(
		c.Request.Context(),
		userID,
		req.Username,
		req.Bio,
		req.ProfilePicture,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}
func (pc *ProfileController) GetProfile(c *gin.Context) {
	userID := c.GetString("userID") // from auth middleware

	profileData, err := pc.profileUsecase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profileData)
}
func (pc *ProfileController) UploadProfilePicture(c *gin.Context) {
	fmt.Println("Reached UploadProfilePicture")
	userID := c.GetString("userID") // From auth middleware

	file, fileHeader, err := c.Request.FormFile("profilePicture")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read file"})
		return
	}
	defer file.Close()

	// Pass to usecase
	picturePath, err := pc.profileUsecase.UploadProfilePicture(
		c.Request.Context(),
		userID,
		file,
		fileHeader,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile picture uploaded successfully",
		"path":    picturePath,
	})
}
