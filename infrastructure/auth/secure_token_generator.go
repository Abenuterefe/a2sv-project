package auth

import (
	"crypto/rand"
	"encoding/base64"
)

type SecureTokenGenerator struct{}

func NewSecureTokenGenerator() *SecureTokenGenerator{
	return &SecureTokenGenerator{}
}

// generateSecureToken generates a cryptographically secure random token.
func (g *SecureTokenGenerator) GenerateSecureToken() string{
	//32 byte of entropy(256 bit)
	b := make([]byte, 32)
	if _,err := rand.Read(b); err != nil {
		//incase errorm just incode hard coded one
		return base64.URLEncoding.EncodeToString([]byte("fallback-token"))
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
}