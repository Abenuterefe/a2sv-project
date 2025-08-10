package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"github.com/Abenuterefe/a2sv-project/domain/interfaces"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userUsecase struct {
	userRepo            interfaces.UserRepository
	passwordHasher      interfaces.PasswordService
	authService         interfaces.AuthService
	mailService         interfaces.MailService
	oauthService        interfaces.OAuthService
	secureToknGenerator interfaces.ISecureTokenGenerator
}

var AccessTokenTTL = time.Minute * 15

func NewUserUsecase(
	userRepo interfaces.UserRepository,
	hasher interfaces.PasswordService,
	auth interfaces.AuthService,
	mailService interfaces.MailService,
	oauthService interfaces.OAuthService,
	secureToknGenerator interfaces.ISecureTokenGenerator) interfaces.UserUsecase {
	return &userUsecase{
		userRepo:            userRepo,
		passwordHasher:      hasher,
		authService:         auth,
		mailService:         mailService,
		oauthService:        oauthService,
		secureToknGenerator: secureToknGenerator,
	}
}

func (u *userUsecase) Regiser(ctx context.Context, user *entities.User) error {
	// Clean input
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Username = strings.TrimSpace(user.Username)

	// Check if user exists
	existingUser, _ := u.userRepo.FindByEmail(ctx, user.Email)
	if existingUser != nil {
		return errors.New("email already registered")
	}

	// Hash password
	hashedPwd, err := u.passwordHasher.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPwd

	// Fill other fields of user
	user.ID = primitive.NewObjectID()
	user.Role = entities.RoleUser //by default role is user role
	user.Verified = false
	user.VerificationToken = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// save to database
	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return err
	}

	// Send verification email to confirm email
	return u.mailService.SendVerificationEmail(user.Email, user.VerificationToken)
}

func (u *userUsecase) Login(ctx context.Context, email, password string) (*entities.Token, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	// find user by his/her email
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Restrict unverified users from logging in
	if !user.Verified {
		return nil, errors.New("please verify your email befor login")
	}

	// verify pwd
	if err := u.passwordHasher.VerifyPassword(user.Password, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate jwt tokens (accesstoken and refresh token)
	accessToken, err := u.authService.CreateAccessToken(user.ID.Hex(), string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.authService.CreateRefreshToken(user.ID.Hex(), string(user.Role))
	if err != nil {
		return nil, err
	}

	// Populate token object
	token := &entities.Token{
		UserID:       user.ID.Hex(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(15 * time.Minute),
	}

	// Store refresh token in Database
	if err := u.userRepo.StoreToken(ctx, token); err != nil {
		return nil, err
	}

	// Return token to client device
	return token, nil
}

// Refresh token handler usecase
func (u *userUsecase) RefreshToken(ctx context.Context, refreshToken string) (*entities.Token, error) {
	// fins refresh token in db
	storedToken, err := u.userRepo.FindToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// verify referesh token validity
	claims, err := u.authService.VerifyToken(refreshToken, false)
	if err != nil {
		return nil, err
	}

	// Generate new access token
	accessToken, err := u.authService.CreateAccessToken(claims.UserID, claims.Role)
	if err != nil {
		return nil, err
	}

	return &entities.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(AccessTokenTTL),
		UserID:       storedToken.UserID,
	}, nil
}

// Implement verify email function
func (u *userUsecase) VerifyEmail(ctx context.Context, token string) error {
	user, err := u.userRepo.FindByVerificationToken(ctx, token)
	if err != nil || user == nil {
		return errors.New("invalid verification token")
	}

	// if we can find user registered with token verification token, change status and verToken
	user.Verified = true
	user.VerificationToken = "" //clear roken

	return u.userRepo.Update(ctx, user)
}

// Implement resending verification email funcion
func (u *userUsecase) ResendVerificationEmail(ctx context.Context, email string) error {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.Verified {
		return errors.New("user already verified")
	}

	// Generate new verification token
	user.VerificationToken = uuid.New().String()
	err = u.userRepo.Update(ctx, user)
	if err != nil {
		return errors.New("failed to update verification token")
	}

	// send email again
	return u.mailService.SendVerificationEmail(user.Email, user.VerificationToken)
}

// Promote user func
func (u *userUsecase) PromoteUser(ctx context.Context, userID string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	user, err := u.userRepo.FindByID(ctx, objID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	// promote user
	user.Role = entities.RoleAdmin
	user.UpdatedAt = time.Now()

	return u.userRepo.Update(ctx, user)
}

// demote user func usecase
func (u *userUsecase) DemoteUser(ctx context.Context, userID string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalide user ID")
	}

	user, err := u.userRepo.FindByID(ctx, objID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	// demote user
	user.Role = entities.RoleUser
	user.UpdatedAt = time.Now()

	return u.userRepo.Update(ctx, user)
}

// logout user
func (u *userUsecase) Logout(ctx context.Context, userID string) error {
	return u.userRepo.DeleteToken(ctx, userID)
}

// login using google
func (u *userUsecase) GoogleOAuthLogin(ctx context.Context, code string) (*entities.Token, error) {
	userInfo, err := u.oauthService.GetUserInfo(ctx, code)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.FindByEmail(ctx, userInfo.Email)
	if err != nil || user == nil {
		//if used doesnt exist, register and verify automatically
		newUser := &entities.User{
			Email:     userInfo.Email,
			Username:  userInfo.Name,
			Verified:  true,
			Role:      entities.RoleUser,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := u.userRepo.Create(ctx, newUser)
		if err != nil {
			return nil, err
		}

		//user = newUser
	}

	//Re access the user with new object id if it was not exist
	user, err = u.userRepo.FindByEmail(ctx, userInfo.Email)
	if err != nil || user == nil {
		return nil, errors.New("failed to retrieve user after creation")
	}

	accessToken, err := u.authService.CreateAccessToken(user.ID.Hex(), string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.authService.CreateRefreshToken(user.ID.Hex(), string(user.Role))
	if err != nil {
		return nil, err
	}

	// Populate token object
	token := &entities.Token{
		UserID:       user.ID.Hex(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(15 * time.Minute),
	}

	// Store refresh token in Database
	if err := u.userRepo.StoreToken(ctx, token); err != nil {
		return nil, err
	}

	return token, nil

}

// reset password request
func (u *userUsecase) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	//===generate secure cryptographic token=========
	//token := "hardcoded"
	token := u.secureToknGenerator.GenerateSecureToken()

	resetToken := &entities.ResetToken{
		UserID:    user.ID.Hex(),
		Token:     token,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	//save the token
	if err := u.userRepo.SaveResetToken(ctx, resetToken); err != nil {
		return err
	}

	//Send reset link(link to Reset password frontend page = http://localhost:3000/reset-password?token=abc123 via email
	resetURL := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", &token)
	return u.mailService.SendPasswordResetEmail(user.Email, resetURL)
}

// Reset password
func (u *userUsecase) ResetPassword(ctx context.Context, token, newPassword string) error {
	resetToken, err := u.userRepo.FindByResetToken(ctx, token)
	if err != nil || time.Now().After(resetToken.ExpiresAt) {
		return errors.New("invalid or expired token")
	}

	hashedNewPwd, err := u.passwordHasher.HashPassword(newPassword)
	if err != nil {
		return err
	}

	userID, _ := primitive.ObjectIDFromHex(resetToken.UserID)
	if err := u.userRepo.UpdatePassword(ctx, userID, hashedNewPwd); err != nil {
		return err
	}

	return u.userRepo.DeleteResetToken(ctx, token)
}
