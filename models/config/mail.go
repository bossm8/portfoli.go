package config

import (
	"fmt"
	"log"
	"net/mail"

	"gopkg.in/gomail.v2"
	"gopkg.in/yaml.v3"
)

const (
	// Subject which will be used in the contact emails
	subject = "[Portfolio] New message from %s"
)

// SMTPConfig contains the configuration of the mailing service
type SMTPConfig struct {
	// User which is used to login to the smpt service
	// Should be an email, because emails will be sent with this address used
	// as the From header (will be checked on loading)
	User EmailAddress `yaml:"user"`
	// Password which is used to login to the smpt service
	Pass string `yaml:"pass"`
	// The smtp host which will send the emails
	Host string `yaml:"host"`
	// The port on which the smtp host listens on
	Port int `yaml:"port"`
}

// SendMail sends the email message to receiver via the configured smtp service
func (smtp *SMTPConfig) SendMail(
	receiver *mail.Address,
	replyTo *mail.Address,
	senderName string,
	message string,
) error {
	mail := gomail.NewMessage()
	subject := fmt.Sprintf(subject, senderName)
	mail.SetHeaders(map[string][]string{
		"To":       {receiver.Address},
		"From":     {smtp.User.Address.Address, "[Portfolio]: " + senderName},
		"Reply-To": {replyTo.Address},
		"Subject":  {subject},
	})
	mail.SetBody("text/plain", message)

	dialer := gomail.NewDialer(smtp.Host, smtp.Port, smtp.User.Address.Address, smtp.Pass)
	if err := dialer.DialAndSend(mail); err != nil {
		log.Printf("[ERROR] Could not send email: %s\n", err)
		return err
	}
	return nil
}

// EmailAddress is a wrapper around mail.Address, which does not implement
// the UnmarshalYAML function
type EmailAddress struct {
	*mail.Address
}

// UmarshalYAML unmarshals the string address from yaml into an EmailAddress
func (m *EmailAddress) UnmarshalYAML(value *yaml.Node) error {
	if addr, err := mail.ParseAddress(value.Value); nil != err {
		log.Fatalf("[ERROR]: Invalid profile email address %s", value.Value)
		return err
	} else {
		m.Address = addr
	}
	return nil
}

// Make sure the Unmarshaler interface is implemented
var _ yaml.Unmarshaler = &EmailAddress{}
