package models

import (
	"errors"
	"log"
	"net/mail"
	"reflect"
	"strings"
)

const smtpConfigName = "mail.yml"

type SMTPConfig struct {
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type MailConfig struct {
	To   string      `yaml:"to"`
	SMTP *SMTPConfig `yaml:"smtp"`
}

func GetSMTPConfig() (cfg *MailConfig, err error) {
	mailCfg := &struct {
		Mail *MailConfig `yaml:"mail"`
	}{}
	if err = loadFromFile(smtpConfigName, mailCfg); nil != err {
		return
	}
	cfg = mailCfg.Mail
	if _, err = mail.ParseAddress(cfg.To); nil != err {
		log.Printf("[ERROR]: Invalid or missing email address %s", cfg.To)
	}

	val := reflect.ValueOf(*cfg.SMTP)
	for i := 0; i < val.NumField(); i++ {
		if v := val.Field(i); v.IsZero() {
			log.Printf(
				"[ERROR]: SMTP config lacking a correct value for '%s'\n",
				strings.ToLower(val.Type().Field(i).Name),
			)
			err = errors.New("missing key in SMTP config")
		}
	}
	return
}
