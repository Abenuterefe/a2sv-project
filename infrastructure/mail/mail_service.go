package mail

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

type MailService struct {
	From     string
	SMTPHost string
	SMTPPort int
	Password string
}

func NewMailService() *MailService {
	return &MailService{
		From:     os.Getenv("EMAIL_SENDER"),
		SMTPHost: os.Getenv("SMTP_HOST"),
		SMTPPort: 587,
		Password: os.Getenv("EMAIL_PASSWORD"),
	}
}

// function to send verification email
func (s *MailService) SendVerificationEmail(to, token string) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", s.From)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", "Verify your email")

	verifyURL := fmt.Sprintf("http://localhost:8080/auth/verify?token=%s", token)
	body := fmt.Sprintf("Click here to verify your email: <a href=\"%s\">Verify</a>", verifyURL)
	mail.SetBody("text/html", body)

	dialer := gomail.NewDialer(s.SMTPHost, s.SMTPPort, s.From, s.Password)

	return dialer.DialAndSend(mail)
}
