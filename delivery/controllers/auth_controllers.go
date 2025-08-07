package controllers

import (
	"context"
	"net/http"
	"strings"

	"time"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"github.com/Abenuterefe/a2sv-project/domain/interfaces"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	UserUsecase  interfaces.UserUsecase
	oauthService interfaces.OAuthService
}

func NewAuthController(userUC interfaces.UserUsecase,
	oauthServ interfaces.OAuthService) *AuthController {
	return &AuthController{
		UserUsecase:  userUC,
		oauthService: oauthServ,
	}
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

// type refreshRequest struct {
// 	RefreshToken string `json:"refresh_token" binding:"required"`
// }

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

		c.JSON(http.StatusBadRequest, gin.H{"error1 malli maali": err.Error()})

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
		"message":       "login with Cridential (pwd, email)  seccussful",
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"expires_at":    token.ExpiresAt,
	})
}

// Refresh Handler
func (ac *AuthController) Refresh(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header missing"})
		return
	}

	// Expect header format: "Bearer <token>"
	fields := strings.Fields(authHeader)
	if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error1": "Invalid Authorization header format"})
		return
	}

	refreshToken := fields[1]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, err := ac.UserUsecase.RefreshToken(ctx, refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error2": err.Error()})

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

// Verify email handler
func (a *AuthController) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing verification token"})
		return
	}

	err := a.UserUsecase.VerifyEmail(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// Resend verification handler
func (a *AuthController) ResendVerification(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	err := a.UserUsecase.ResendVerificationEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email resent"})

}

// Promote user handler
func (a *AuthController) PromoteUser(c *gin.Context) {
	UserID := c.Param("id")
	err := a.UserUsecase.PromoteUser(c.Request.Context(), UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin"})
}

// Demote user handler
func (a *AuthController) DemoteUser(c *gin.Context) {
	userID := c.Param("id")

	err := a.UserUsecase.DemoteUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User demoted to regular user"})
}

// Logout handler
func (a *AuthController) Logout(c *gin.Context) {
	// extract user id from context (set by middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in token"})
		return
	}

	if err := a.UserUsecase.Logout(c.Request.Context(), userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// Google authentication handler
func (a *AuthController) GoogleLogin(c *gin.Context) {
	state := "my_state" // Ideally random & stored in cookie/session
	authURL := a.oauthService.GetAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// Google call back handler (request redirected from google)
func (a *AuthController) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found in query"})
		return
	}

	//generate authentication token for later login like we did for Login handler
	token, err := a.UserUsecase.GoogleOAuthLogin(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":       "login with GOOGLE seccussful!",
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"expires_at":    token.ExpiresAt,
	})

}
