package interfaces

type MailService interface{
	SendVerificationEmail(to, token string) error
}