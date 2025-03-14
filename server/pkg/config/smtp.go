package config

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

const (
	CONFIG_SMTP_HOST        = "MAIL_HOST"
	CONFIG_SMTP_PORT        = "MAIL_PORT"
	CONFIG_SMTP_USERNAME    = "MAIL_USERNAME"
	CONFIG_SMTP_PASSWORD    = "MAIL_PASSWORD"
	CONFIG_SMTP_SENDER_NAME = "MAIL_FROM_NAME"
	CONFIG_SMTP_SENDER_EMAIL = "MAIL_FROM_ADDRESS"
)

type SMTPConfig struct {
	Host        string
	Port        string
	Username    string
	Password    string
	SenderName  string
	SenderEmail string
}

func LoadSMTPConfig() (*SMTPConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	config := &SMTPConfig{
		Host:        os.Getenv(CONFIG_SMTP_HOST),
		Port:        os.Getenv(CONFIG_SMTP_PORT),
		Username:    os.Getenv(CONFIG_SMTP_USERNAME),
		Password:    os.Getenv(CONFIG_SMTP_PASSWORD),
		SenderName:  os.Getenv(CONFIG_SMTP_SENDER_NAME),
		SenderEmail: os.Getenv(CONFIG_SMTP_SENDER_EMAIL),
	}

	return config, nil
}

func (config *SMTPConfig) SendEmail(from string, to []string, cc []string, subject string, body string) error {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	header := make(map[string]string)
	header["From"] = fmt.Sprintf("%s <%s>", config.SenderName, config.SenderEmail)
	header["To"] = joinAddresses(to)
	if len(cc) > 0 {
		header["Cc"] = joinAddresses(cc)
	}
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	addresses := append(to, cc...)
	err := smtp.SendMail(config.Host+":"+config.Port, auth, from, addresses, []byte(message))
	if err != nil {
		return err
	}
	return nil
}

func joinAddresses(addresses []string) string {
	return strings.Join(addresses, ", ")
}
