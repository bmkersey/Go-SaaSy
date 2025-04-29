package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"
)

type EmailSender struct {
	From     string
	Password string
	Host     string
	Port     string
}

func NewEmailSender() *EmailSender {
	return &EmailSender{
		From:     os.Getenv("EMAIL_FROM"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		Host:     "smtp.gmail.com",
		Port:     "587",
	}
}

func (e *EmailSender) SendTemplate(to, subject, templateName string, data any) error {
	tmplPath := filepath.Join("internal", "email", "templates", templateName)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	msg := []byte("Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		body.String())

	auth := smtp.PlainAuth("", e.From, e.Password, e.Host)
	err = smtp.SendMail(e.Host+":"+e.Port, auth, e.From, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}
