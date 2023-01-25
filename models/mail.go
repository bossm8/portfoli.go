package models

import (
	"fmt"
	"log"
	"net/mail"

	"gopkg.in/gomail.v2"
)

const (
	subject = "[Portfolio] New message from %s"
)

// SMTPConfig contains the configuration of the mailing service
type SMTPConfig struct {
	// User which is used to login to the smpt service
	// Should be an email, because emails will be sent with this address used
	// as the From header (will be checked on loading)
	User string `yaml:"user"`
	// Password which is used to login to the smpt service
	Pass string `yaml:"pass"`
	// The smtp host which will send the emails
	Host string `yaml:"host"`
	// The port on which the smtp host listens on
	Port int `yaml:"port"`
}

func (smtp *SMTPConfig) SendMail(to string, senderName string, replyTo *mail.Address, message string) error {
	mail := gomail.NewMessage()
	subject := fmt.Sprintf(subject, senderName)
	mail.SetHeaders(map[string][]string{
		"To":       {to},
		"From":     {smtp.User, "[Portfolio]: " + senderName},
		"Reply-To": {replyTo.Address},
		"Subject":  {subject},
	})
	mail.SetBody("text/plain", message)

	dialer := gomail.NewDialer(smtp.Host, smtp.Port, smtp.User, smtp.Pass)
	if err := dialer.DialAndSend(mail); err != nil {
		log.Printf("[ERROR] Could not send email: %s\n", err)
		return err
	}
	return nil
}
