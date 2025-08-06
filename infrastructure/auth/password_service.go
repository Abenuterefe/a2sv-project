package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordService struct{}

func NewBcryptPasswordService() *BcryptPasswordService {
	return &BcryptPasswordService{}
}

// HashPassword hashes a plaintext password using bcrypt.
func (s *BcryptPasswordService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// VerifyPassword compares a hashed password with a plaintext one.
func (s *BcryptPasswordService) VerifyPassword(hashedPassword, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return errors.New("invalid password")
	}
	return nil
}
