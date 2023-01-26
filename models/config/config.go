package config

import (
	"errors"
	"html/template"
	"log"
	"reflect"
	"strings"

	"bossm8.ch/portfolio/models/content"
	"bossm8.ch/portfolio/models/utils"
)

// The name of the configuration file which contains the main application config
const configName = "config.yml"

// SocialMedia represents a generic social media type
type SocialMedia struct {
	// Make this one of the 'social' type icons of https://icons.getbootstrap.com/#icons
	Type string `yaml:"type"`
	Link string `yaml:"link"`
}

var ErrInvalidSMTPConfig = errors.New("invalid SMTP configuration")

// ProfileConfig contains the configurations about the profile which will
// be highlighted in the portfolio
type ProfileConfig struct {
	BrandName      string         `yaml:"brandname"`
	BannerImage    string         `yaml:"bannerimage"`
	Avatar         string         `yaml:"avatar"`
	FirstName      string         `yaml:"firstname"`
	LastName       string         `yaml:"lastname"`
	Email          *EmailAddress  `yaml:"email"`
	Heading        *template.HTML `yaml:"heading"`
	SubHeading     *template.HTML `yaml:"subheading"`
	Slogan         string         `yaml:"slogan"`
	ContactHeading string         `yaml:"contactheading"`
	SocialMedia    []*SocialMedia `yaml:"social"`
	ContentKinds   []string       `yaml:"content"`
	Me             *template.HTML `yaml:"me"`
}

// Config contains the static configuration of the portfolio,
// meaning the mailing config and your profile settings
type Config struct {
	Profile       *ProfileConfig `yaml:"profile"`
	SMTP          *SMTPConfig    `yaml:"smtp"`
	RenderContact bool
}

// GetConfig loads and returns the configuration from <config.dir>/config.yaml
func GetConfig() (*Config, error) {
	// Default values which well be used on first load when nothing is configured
	cfg := &Config{
		Profile: &ProfileConfig{},
	}
	if err := utils.LoadFromYAMLFile(configName, cfg); nil != err {
		return nil, err
	}

	for _, contentType := range cfg.Profile.ContentKinds {
		if !content.IsValidContentType(contentType) {
			return nil, errors.New("invalid content kind " + contentType)
		}
	}

	// All configuration of smtp is required for the mailing service to be working
	// as yaml.v3 does not yet have a required tag, the check is made manually
	cfg.RenderContact = true
	val := reflect.ValueOf(*cfg.SMTP)
	for i := 0; i < val.NumField(); i++ {
		if v := val.Field(i); v.IsZero() {
			log.Printf(
				"[ERROR]: SMTP config lacking a correct value for '%s'\n",
				strings.ToLower(val.Type().Field(i).Name),
			)
			cfg.RenderContact = false
			return cfg, ErrInvalidSMTPConfig
		}
	}
	return cfg, nil
}
