package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"github.com/Abenuterefe/a2sv-project/domain/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/google/uuid"
)

type userUsecase struct {
	userRepo       interfaces.UserRepository
	passwordHasher interfaces.PasswordService
	authService    interfaces.AuthService
	mailService interfaces.MailService
}

var AccessTokenTTL = time.Minute * 15

func NewUserUsecase(
	userRepo interfaces.UserRepository,
	hasher interfaces.PasswordService,
	auth interfaces.AuthService,
	mailService interfaces.MailService) interfaces.UserUsecase {
	return &userUsecase{
		userRepo:       userRepo,
		passwordHasher: hasher,
		authService:    auth,
		mailService: mailService,
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
	if err != nil{
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
	if !user.Verified{
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

	refreshToken, err := u.authService.CreateRefreshToken(user.ID.Hex())
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
	user, err := u.userRepo.FindByVerificationToken(ctx,token)
	if err != nil || user == nil{
		return errors.New("invalid verification token")
	}

	// if we can find user registered with token verification token, change status and verToken
	user.Verified = true
	user.VerificationToken="" //clear roken

	return u.userRepo.Update(ctx,user)
}

// Implement resending verification email funcion
func (u *userUsecase) ResendVerificationEmail(ctx context.Context,email string) error {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.Verified{
		return errors.New("user already verified")
	}

	// Generate new verification token
	user.VerificationToken = uuid.New().String()
	err = u.userRepo.Update(ctx,user)
	if err != nil {
		return errors.New("failed to update verification token")
	}

	// send email again
	return u.mailService.SendVerificationEmail(user.Email, user.VerificationToken)
}