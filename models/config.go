package models

import (
	"errors"
	"log"
	"net/mail"
	"reflect"
	"strings"
)

const smtpConfigName = "config.yml"

// The configuration ot the smtp service
type SMTPConfig struct {
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// The configurations about yourself
type ProfileConfig struct {
	BrandName   string `yaml:"brandname"`
	BannerImage string `yaml:"bannerimage"`
	FirstName   string `yaml:"firstname"`
	LastName    string `yaml:"lastname"`
	Email       string `yaml:"email"`
}

// The configuration of the mailing service
type Config struct {
	Profile *ProfileConfig `yaml:"profile"`
	SMTP    *SMTPConfig    `yaml:"smtp"`
}

// Loads and returns the configuration from configs/mail.yaml
func GetConfig() (cfg *Config, err error) {
	cfg = &Config{
		Profile: &ProfileConfig{
			BrandName: "Queen",
			FirstName: "Freddy",
			LastName:  "Mercury",
			Email:     "freddy@mercury.me",
		},
	}
	if err = loadFromFile(smtpConfigName, cfg); nil != err {
		return
	}
	if _, err = mail.ParseAddress(cfg.Profile.Email); nil != err {
		log.Fatalf("[ERROR]: Invalid or missing profile email address %s", cfg.Profile.Email)
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
