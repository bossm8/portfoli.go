package models

import (
	"errors"
	"html/template"
	"log"
	"reflect"
	"regexp"
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

var InvalidSMTPConfigError = errors.New("Invalid SMTP configuration")

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
	Profile *ProfileConfig `yaml:"profile"`
	SMTP    *SMTPConfig    `yaml:"smtp"`
}

// GetConfig loads and returns the configuration from <config.dir>/config.yaml
func GetConfig() (*Config, error) {
	// Default values which well be used on first load when nothing is configured
	cfg := &Config{
		Profile: &ProfileConfig{},
	}
	if err := loadFromYAMLFile(configName, cfg); nil != err {
		return nil, err
	}

	// Check if all content kinds specified in the yaml config are valid
	rex := regexp.MustCompile("(" + strings.Join(contentKinds, "|") + ")")
	for _, contentKind := range cfg.Profile.ContentKinds {
		if !rex.MatchString(contentKind) {
			log.Printf("[ERROR] Invalid content kind '%s', allowed values are: %s\n", contentKind, contentKinds)
			return nil, errors.New("invalid content kind " + contentKind)
		}
	}

	// All configuration of smtp is required for the mailing service to be working
	// as yaml.v3 does not yet have a required tag, the check is made manually
	val := reflect.ValueOf(*cfg.SMTP)
	for i := 0; i < val.NumField(); i++ {
		if v := val.Field(i); v.IsZero() {
			log.Printf(
				"[ERROR]: SMTP config lacking a correct value for '%s'\n",
				strings.ToLower(val.Type().Field(i).Name),
			)
			return cfg, InvalidSMTPConfigError
		}
	}
	return cfg, nil
}

// SetConfigDir sets the configuration directory where dynamic and static configurations
// should be read from
func SetConfigDir(dir string) {
	configDir = dir
}
