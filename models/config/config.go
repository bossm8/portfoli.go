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
	"encoding/json"
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

// ImagesConfig contains configuration for image caching behavior
type ImagesConfig struct {
	// Cache controls if remote images should be cached locally
	Cache bool `yaml:"cache"`
	// Force forces a re-download of cached images on startup
	Force bool `yaml:"force"`
}

// SEOConfig contains configuration for page title/metadata and social
// share previews
type SEOConfig struct {
	// SiteName is appended to every page's <title> (e.g. "Name - Page - SiteName")
	SiteName string `yaml:"sitename"`
	// Description is the default meta/og/twitter description for pages
	// which do not define their own (e.g. the about page does)
	Description string `yaml:"description"`
	// Image used for social share previews (og:image/twitter:image),
	// falls back to profile.avatar when unset
	Image string `yaml:"image"`
	// SiteURL is the absolute base URL of the deployed site (e.g.
	// https://example.com), optional - if set, it is used to build an
	// absolute image URL for maximum share-card compatibility
	SiteURL string `yaml:"siteurl"`
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

// ApplyImageCache updates image fields to use a cached local path when enabled.
func (p *ProfileConfig) ApplyImageCache() {
	if p == nil {
		return
	}
	if p.BannerImage != "" {
		p.BannerImage = apputils.MaybeCacheImage(p.BannerImage)
	}
	if p.Avatar != "" {
		p.Avatar = apputils.MaybeCacheImage(p.Avatar)
	}
}

// personSchema is the JSON-LD payload describing the profile as a
// schema.org Person, embedded on every page for richer search engine
// results (e.g. a knowledge-panel-style entry).
type personSchema struct {
	Context  string   `json:"@context"`
	Type     string   `json:"@type"`
	Name     string   `json:"name"`
	JobTitle string   `json:"jobTitle,omitempty"`
	Email    string   `json:"email,omitempty"`
	Image    string   `json:"image,omitempty"`
	URL      string   `json:"url,omitempty"`
	SameAs   []string `json:"sameAs,omitempty"`
}

// PersonJSONLD returns a JSON-LD <script> payload describing the profile as
// a schema.org Person. Built via encoding/json rather than assembled as
// text in the template, so field values (which may contain quotes,
// ampersands, etc.) are guaranteed to be correctly JSON-escaped rather than
// HTML-escaped, which would produce invalid JSON. image should already be
// fully resolved (e.g. via the Assemble template function) since this
// method has no access to the request's base path.
func (p *ProfileConfig) PersonJSONLD(seo *SEOConfig, image string) template.JS {
	if p == nil {
		return ""
	}
	name := strings.TrimSpace(p.FirstName + " " + p.LastName)
	if name == "" {
		name = p.BrandName
	}
	var email string
	if p.Email != nil && p.Email.Address != nil {
		email = p.Email.Address.Address
	}
	var sameAs []string
	for _, s := range p.SocialMedia {
		if s.Link != "" {
			sameAs = append(sameAs, s.Link)
		}
	}
	var siteURL string
	if seo != nil {
		siteURL = seo.SiteURL
	}
	schema := personSchema{
		Context:  "https://schema.org",
		Type:     "Person",
		Name:     name,
		JobTitle: p.Slogan,
		Email:    email,
		Image:    image,
		URL:      siteURL,
		SameAs:   sameAs,
	}
	b, err := json.Marshal(schema)
	if err != nil {
		log.Printf("[WARNING] Failed to marshal person JSON-LD: %s\n", err)
		return ""
	}
	return template.JS(b)
}

// Config contains the static configuration of the portfolio,
// meaning the mailing config and your profile settings
type Config struct {
	// Profile: the static configuration about the profile loaded on start
	Profile *ProfileConfig `yaml:"profile"`
	// SEO configuration for page title/metadata and social share previews
	SEO *SEOConfig `yaml:"seo"`
	// SMTP configuration used to send emails via the contact form
	SMTP *SMTPConfig `yaml:"smtp"`
	// Images configuration for caching remote images
	Images *ImagesConfig `yaml:"images"`
	// RenderContact signals if the contact form should be rendered or not
	RenderContact bool
}

// Load loads and returns the configuration from <config.dir>/config.yaml
func Load() (*Config, error) {
	// Default values which well be used on first load when nothing is configured
	defaultBrandImage := template.HTML(
		"<img src='/static/img/portfoli.go-yellow.svg' style='height: 26px;'/>",
	)
	cfg = &Config{
		Profile: &ProfileConfig{
			BrandName:  "Portfoli.go",
			BrandImage: &defaultBrandImage,
		},
		Images: &ImagesConfig{},
	}
	if err := utils.LoadFromYAMLFile(ConfigFile, cfg); nil != err {
		return nil, err
	}
	if cfg.Images == nil {
		cfg.Images = &ImagesConfig{}
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
