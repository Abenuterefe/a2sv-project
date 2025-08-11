package mail

import (
	"fmt"
	"os"
	"time"

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
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Email Verification</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      background-color: #f4f4f4;
      margin: 0;
      padding: 0;
    }
    .container {
      max-width: 500px;
      margin: 30px auto;
      background-color: #ffffff;
      border-radius: 8px;
      box-shadow: 0px 4px 8px rgba(0,0,0,0.1);
      overflow: hidden;
    }
    .header {
      background-color: #4CAF50;
      color: white;
      padding: 20px;
      text-align: center;
      font-size: 24px;
    }
    .content {
      padding: 20px;
      font-size: 16px;
      color: #333333;
      line-height: 1.5;
    }
    .button {
      display: inline-block;
      background-color: #4CAF50;
      color: white;
      padding: 12px 20px;
      text-decoration: none;
      border-radius: 5px;
      font-weight: bold;
    }
    .footer {
      font-size: 12px;
      color: #777777;
      text-align: center;
      padding: 15px;
      border-top: 1px solid #eaeaea;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      Email Verification
    </div>
    <div class="content">
      <p>Hello,</p>
      <p>Thank you for registering! Please verify your email by clicking the button below:</p>
      <p style="text-align: center;">
        <a href="%s" class="button">Verify Email</a>
      </p>
      <p>If you didn’t create an account, you can safely ignore this email.</p>
    </div>
    <div class="footer">
      &copy; %d Backend team group 2. All rights reserved.
    </div>
  </div>
</body>
</html>
`, verifyURL, time.Now().Year())
	mail.SetBody("text/html", body)

	dialer := gomail.NewDialer(s.SMTPHost, s.SMTPPort, s.From, s.Password)

	return dialer.DialAndSend(mail)
}

// function to Password reset email
func (s *MailService) SendPasswordResetEmail(to string, resetLink string) error {
	subject := "Reset Your Password"
	body := fmt.Sprintf(`
<html>
  <body style="margin:0; padding:0; font-family: Arial, sans-serif; background-color:#f4f4f4;">
    <table width="100%%" border="0" cellspacing="0" cellpadding="0" style="padding: 20px;">
      <tr>
        <td align="center">
          <table width="600" style="background-color:#ffffff; border-radius:8px; overflow:hidden; box-shadow:0 2px 8px rgba(0,0,0,0.1);">
            <!-- Header -->
            <tr>
              <td style="background-color:#4CAF50; padding:20px; text-align:center; color:white; font-size:24px; font-weight:bold;">
                Password Reset Request
              </td>
            </tr>
            <!-- Body -->
            <tr>
              <td style="padding: 30px; color:#333333; font-size:16px; line-height:1.5;">
                <p>Hello,</p>
                <p>We received a request to reset your password. Please click the button below to set a new password:</p>
                <p style="text-align:center; margin: 30px 0;">
                  <a href="%s" style="background-color:#4CAF50; color:white; padding:12px 20px; text-decoration:none; font-weight:bold; border-radius:5px; display:inline-block;">
                    Reset Password
                  </a>
                </p>
                <p>If you didn’t request this, you can safely ignore this email.</p>
              </td>
            </tr>
            <!-- Footer -->
            <tr>
              <td style="background-color:#f4f4f4; padding:15px; text-align:center; color:#888888; font-size:12px;">
                &copy; %d Backend team Group 2. All rights reserved.
              </td>
            </tr>
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>
`, resetLink, time.Now().Year())

	//create message object
	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	//send message to usere
	d := gomail.NewDialer(s.SMTPHost, s.SMTPPort, s.From, s.Password)
	return d.DialAndSend(m)
}
