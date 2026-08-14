// Copyright (c) 2023, Boss Marco <bossm8+portfoligo@hotmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package config

import (
	"fmt"
	"log"
	"net/mail"

	appconfig "github.com/bossm8/portfoli.go/config"
	apputils "github.com/bossm8/portfoli.go/utils"

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
	if html, err := renderMailHTML(senderName, replyTo.Address, message); nil != err {
		log.Printf("[WARNING] Failed to render HTML mail body, sending plain text only: %s\n", err)
	} else {
		mail.AddAlternative("text/html", html)
	}

	dialer := gomail.NewDialer(smtp.Host, smtp.Port, smtp.User.Address.Address, smtp.Pass)
	if err := dialer.DialAndSend(mail); err != nil {
		log.Printf("[ERROR] Could not send email: %s\n", err)
		return err
	}
	return nil
}

// mailData is passed to the mail html template
type mailData struct {
	Name    string
	Email   string
	Message string
}

// renderMailHTML renders the html alternative body sent alongside the plain
// text one, so contact form notifications also look decent in mail clients
// which prefer HTML.
func renderMailHTML(name, email, message string) (string, error) {
	data := mailData{Name: name, Email: email, Message: message}
	rendered, err := apputils.RenderTemplate("mail", data, appconfig.MailTemplatePath())
	if nil != err {
		return "", err
	}
	return string(rendered), nil
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
