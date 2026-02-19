package entity

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/Higor-ViniciusDev/auth/internal/internal_error"
)

type Email struct {
	To      []string
	Subject string
	Body    string
}

func (e *Email) SendEmail() *internal_error.InternalError {
	from := os.Getenv("SMTP_FROM")
	username := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	message := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n"+
			"%s",
		from,
		e.Subject,
		e.Body,
	))

	auth := smtp.PlainAuth("", username, password, smtpHost)

	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		e.To,
		message,
	)

	if err != nil {
		return internal_error.NewInternalServerError("error sending email via mailtrap")
	}

	return nil
}

type EmailSenderInterface interface {
	SendEmail() *internal_error.InternalError
}
