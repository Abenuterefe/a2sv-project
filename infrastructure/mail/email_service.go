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

//function to Password reset email
func (s *MailService) SendPasswordResetEmail(to string, resetLink string) error{
	subject := "Reset Your Password"
	body := fmt.Sprintf(`
		<html>
			<body>
				<p>Hello,</p>
				<p>You requested a password reset. Click the link below to set a new password:</p>
				<p><a href="%s">Reset Password</a></p>
				<p>If you did not request this, please ignore this email.</p>
			</body>
		</html>`, resetLink)
	
	//create message object
	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	//send message to usere 
	d:= gomail.NewDialer(s.SMTPHost, s.SMTPPort, s.From, s.Password)
	return d.DialAndSend(m)
}