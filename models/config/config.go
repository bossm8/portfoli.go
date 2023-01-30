// Copyright (c) 2023, Boss Marco <bossm8@hotmail.com>
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
	"errors"
	"html/template"
	"log"
	"reflect"
	"strings"

	apputils "github.com/bossm8/portfoli.go/utils"

	"github.com/bossm8/portfoli.go/models/content"
	"github.com/bossm8/portfoli.go/models/utils"
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

var (
	cfg *Config = nil
)

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
	// Animations defines if animations should be added to the page or not
	Animations bool `yaml:"animations"`
}

// RenderHTML renders all HTML fields of the profile by passing them through the
// templates engine. This enables having e.g. the Assemble function the configs
func (p *ProfileConfig) RenderHTML() error {
	val := reflect.ValueOf(*p)
	for i := 0; i < val.NumField(); i++ {
		if res, ok := val.Field(i).Interface().(*template.HTML); ok {
			newHTML, err := apputils.ProcessHTMLContent(res)
			if err != nil {
				log.Printf("failed to process HTML template for %s", val.Type().Field(i).Name)
				return err
			}
			*res = *newHTML
		}
	}
	return nil
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

// Load loads and returns the configuration from <config.dir>/config.yaml
func Load() (*Config, error) {
	// Default values which well be used on first load when nothing is configured
	defaultBrandImage := template.HTML(
		"<img src='/static/img/portfoli.go-yellow.svg' style='width: 25px; margin-bottom: 4px;'/>",
	)
	cfg = &Config{
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

// Get returns the loaded config (Load must have been called at least once, else it will fail)
func Get() *Config {
	if cfg == nil {
		log.Fatalln("[ERROR] Cannot return config, please call LoadConfig first")
	}
	return cfg
}
