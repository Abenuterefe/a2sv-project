package interfaces

type PasswordService interface {
	HashPassword(password string) (string,error)
	VerifyPassword(hashedPassword, plainPassword string) error
}