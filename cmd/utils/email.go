package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
)

type EmailData struct {
	URL       string
	Firstname string
	Subject   string
}

type Config struct {
	EmailFrom string
	SMTPPass  string
	SMTPUser  string
	SMTPHost  string
	SMTPPort  string
}

// ParseTemplateDir email template parser
func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmail(recipientEmail string, data *EmailData, templateName string) error {
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	emailFrom := os.Getenv("RESET_PWD_FROM")

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return err
	}

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		return err
	}

	template = template.Lookup(templateName)
	template.Execute(&body, &data)
	fmt.Println(template.Name())

	m := gomail.NewMessage()

	m.SetHeader("From", emailFrom)
	m.SetHeader("To", recipientEmail)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	dlr := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	dlr.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return nil
}
