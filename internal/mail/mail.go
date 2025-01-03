package mail

import (
	"fmt"
	"go-fiber-api-starter/internal/config"
	"log"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendEmailVerification(to string, link string) error {
	err := SendMail(to, config.EMAIL_FROM, "Welcome!",
		fmt.Sprintf(EmailVerificationTemplate, link, link))
	if err != nil {
		return err
	}
	return nil
}

func SendPasswordReset(to string, link string) error {
	err := SendMail(to, config.EMAIL_FROM, "Reset Password",
		fmt.Sprintf(ResetPasswordTemplate, link, link))
	if err != nil {
		return err
	}
	return nil
}

// gomail docs:  https://pkg.go.dev/gopkg.in/gomail.v2#section-readme
func SendMail(to string, from string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("To", to)
	m.SetHeader("From", from)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	host := config.EMAIL_HOST
	username := config.EMAIL_USERNAME
	password := config.EMAIL_APP_PASSWORD
	port, err := strconv.Atoi(config.EMAIL_PORT)
	if err != nil {
		log.Fatal("Error converting EMAIL_PORT from string to integer:", err)
	}

	d := gomail.NewDialer(host, port, username, password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
