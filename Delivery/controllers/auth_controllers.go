package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"github.com/Abenuterefe/a2sv-project/domain/interfaces"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	UserUsecase interfaces.UserUsecase
}

func NewAuthController(userUC interfaces.UserUsecase) *AuthController {
	return &AuthController{UserUsecase: userUC}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (ac *AuthController) Regiser(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user := &entities.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := ac.UserUsecase.Regiser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// mock mail (optional)
	go func() {
		println("[MOCK EMAIL] Sent verification mail to", user.Email)
	}()

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful. Please verify your email."})

}

// Login handler
func (ac *AuthController) Login(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create context for login handler
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Generate tokens
	token, err := ac.UserUsecase.Login(ctx, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "login seccussful",
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"expires_at":    token.ExpiresAt,
	})
}

// Refresh Handler
func (ac *AuthController) Refresh(c *gin.Context) {
	var req refreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, err := ac.UserUsecase.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Bind token to response body
	c.JSON(http.StatusOK, gin.H{
		"access_token": token.AccessToken,
		"expires_at":   token.ExpiresAt,
	})
}

// User profile handler
func (a *AuthController) Profile(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"message": "Access granted to user profile",
		"userID":  userID,
		"role":    role,
	})
}

// Admin dashboard(profile) handler
func (a *AuthController) AdminDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the Admin Dashboard",
	})
}