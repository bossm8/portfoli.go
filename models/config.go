package models

import (
	"errors"
	"html/template"
	"log"
	"net/mail"
	"reflect"
	"strings"
)

var configDir string

const configName = "config.yml"

// The configuration ot the smtp service
type SMTPConfig struct {
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// Generic social media type
type SocialMedia struct {
	// Make this one of the 'social' type icons of https://icons.getbootstrap.com/#icons
	Type string `yaml:"type"`
	Link string `yaml:"link"`
}

// The configurations about yourself
type ProfileConfig struct {
	BrandName      string         `yaml:"brandname"`
	BannerImage    string         `yaml:"bannerimage"`
	Avatar         string         `yaml:"avatar"`
	FirstName      string         `yaml:"firstname"`
	LastName       string         `yaml:"lastname"`
	Email          string         `yaml:"email"`
	Heading        template.HTML  `yaml:"heading"`
	SubHeading     template.HTML  `yaml:"subheading"`
	Slogan         string         `yaml:"slogan"`
	ContactHeading string         `yaml:"contactheading"`
	SocialMedia    []*SocialMedia `yaml:"social"`
	ContentKinds   []string       `yaml:"content"`
	Me             template.HTML  `yaml:"me"`
}

// The configuration of the mailing service
type Config struct {
	Profile *ProfileConfig `yaml:"profile"`
	SMTP    *SMTPConfig    `yaml:"smtp"`
}

// Loads and returns the configuration from configs/mail.yaml
func GetConfig() (*Config, error) {
	cfg := &Config{
		Profile: &ProfileConfig{
			BrandName: "Queen",
			FirstName: "Freddy",
			LastName:  "Mercury",
			Email:     "freddy@mercury.me",
		},
	}
	if err := loadFromYAMLFile(configName, cfg); nil != err {
		return nil, err
	}
	if _, err := mail.ParseAddress(cfg.Profile.Email); nil != err {
		log.Fatalf("[ERROR]: Invalid or missing profile email address %s", cfg.Profile.Email)
		return nil, err
	}

	var err error
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
	return cfg, err
}

func SetConfigDir(dir string) {
	configDir = dir
}
