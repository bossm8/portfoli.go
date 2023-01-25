package models

import (
	"errors"
	"html/template"
	"log"
	"net/mail"
	"reflect"
	"strings"
)

// The directory where all (static and dynamic) configuration files are read from
var configDir string

// The name of the configuration file which contains the main application config
const configName = "config.yml"

// SocialMedia represents a generic social media type
type SocialMedia struct {
	// Make this one of the 'social' type icons of https://icons.getbootstrap.com/#icons
	Type string `yaml:"type"`
	Link string `yaml:"link"`
}

// ProfileConfig contains the configurations about the profile which will
// be highlighted in the portfolio
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

// Config contains the static configuration of the portfolio,
// meaning the mailing config and your profile settings
type Config struct {
	Profile *ProfileConfig `yaml:"profile"`
	SMTP    *SMTPConfig    `yaml:"smtp"`
}

// GetConfig loads and returns the configuration from <config.dir>/config.yaml
func GetConfig() (*Config, error) {
	// Default values which well be used on first load when nothing is configured
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

	// TODO: check for content types (enabled / disabled) if any and if they are valid

	// All configuration of smtp is required for the mailing service to be working
	// as yaml.v3 does not yet have a required tag, the check is made manually
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

// SetConfigDir sets the configuration directory where dynamic and static configurations
// should be read from
func SetConfigDir(dir string) {
	configDir = dir
}
