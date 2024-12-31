package mail

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

// gomail docs:  https://pkg.go.dev/gopkg.in/gomail.v2#section-readme

func SendMail(to string, from string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("To", to)
	m.SetHeader("From", from)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	host := os.Getenv("EMAIL_HOST")
	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_APP_PASSWORD")
	port, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {
		log.Fatal("Error converting EMAIL_PORT from string to integer:", err)
	}

	d := gomail.NewDialer(host, port, username, password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendEmailVerification(to string, link string) error {
	err := SendMail(to, os.Getenv("EMAIL_FROM"), "Welcome!",
		fmt.Sprintf(EmailVerificationTemplate, link, link))
	if err != nil {
		return err
	}
	return nil
}
