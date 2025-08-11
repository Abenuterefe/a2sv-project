package interfaces

type MailService interface{
	SendVerificationEmail(to, token string) error
	SendPasswordResetEmail(to string, resetLink string) error
}