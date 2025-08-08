package interfaces

type ISecureTokenGenerator interface {
	GenerateSecureToken() string
}
