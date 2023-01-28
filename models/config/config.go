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

const (
	// The name of the configuration file which contains the main application config
	ConfigFile = "config.yml"
)

// SocialMedia represents a generic social media type
type SocialMedia struct {
	// Type of media, should be one of the 'social' type icons of https://icons.getbootstrap.com/#icons
	Type string `yaml:"type"`
	// Link to the social media profile
	Link string `yaml:"link"`
}

// ErrInvalidSMTPConfig signals that the app may continue, but without the contact form
var ErrInvalidSMTPConfig = errors.New("invalid SMTP configuration")

// ProfileConfig contains the configurations about the profile which will
// be highlighted in the portfolio - (optional) means if null, no page will be rendered
type ProfileConfig struct {
	// BrandName is the name displayed in the navigation bar
	BrandName string `yaml:"brandname"`
	// BrandImage is the image displayed in the navigation bar
	BrandImage *template.HTML `yaml:"brandimage"`
	// BannerImage is the image displayed on the index page
	BannerImage string `yaml:"bannerimage"`
	// Avatar displayed as profile image
	Avatar string `yaml:"avatar"`
	// FirstName displayed for the profile
	FirstName string `yaml:"firstname"`
	// LastName displayed for the profile
	LastName string `yaml:"lastname"`
	// Contact email address for the profile
	Email *EmailAddress `yaml:"email"`
	// Heading shown on the index page
	Heading *template.HTML `yaml:"heading"`
	// SubHeading shown on the index page
	SubHeading *template.HTML `yaml:"subheading"`
	// Slogan shown on the index page
	Slogan string `yaml:"slogan"`
	// Heading shown on the contact page
	ContactHeading string `yaml:"contactheading"`
	// All links to social media, displayed in the footer bar and on the index page
	SocialMedia []*SocialMedia `yaml:"social"`
	// ContentTypes enabled for the page (each element is optional)
	// - the element itself not and the list must be valid content types
	ContentTypes []string `yaml:"content"`
	// Me: Short introduction shown in the bio and skills page (optional)
	Me *template.HTML `yaml:"me"`
}

// Config contains the static configuration of the portfolio,
// meaning the mailing config and your profile settings
type Config struct {
	// Profile: the static configuration about the profile loaded on start
	Profile *ProfileConfig `yaml:"profile"`
	// SMTP configuration used to send emails via the contact form
	SMTP *SMTPConfig `yaml:"smtp"`
	// RenderContact signals if the contact form should be rendered or not
	RenderContact bool
}

// GetConfig loads and returns the configuration from <config.dir>/config.yaml
func GetConfig() (*Config, error) {
	// Default values which well be used on first load when nothing is configured
	defaultBrandImage := template.HTML(
		"<img src='/static/img/portfoli.go-yellow.svg' style='width: 25px; margin-bottom: 4px;'/>",
	)
	cfg := &Config{
		Profile: &ProfileConfig{
			BrandName:  "Portfoli.go",
			BrandImage: &defaultBrandImage,
		},
	}
	if err := utils.LoadFromYAMLFile(ConfigFile, cfg); nil != err {
		return nil, err
	}

	for _, contentType := range cfg.Profile.ContentTypes {
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
				"[ERROR] SMTP config lacking a correct value for '%s'\n",
				strings.ToLower(val.Type().Field(i).Name),
			)
			cfg.RenderContact = false
			return cfg, ErrInvalidSMTPConfig
		}
	}
	return cfg, nil
}
