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

package content

import (
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	apputils "github.com/bossm8/portfoli.go/utils"

	"github.com/bossm8/portfoli.go/config"
	"github.com/bossm8/portfoli.go/models/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	typeExperience = iota
	typeEducation
	typeProject
	typeCertification
	typeMe
)

var (
	// All possible content types
	ContentTypes = []string{"experience", "education", "projects", "certifications", "me"}
	// Mappings to easily get the correct content type based on the request path
	contentMappings = map[string]ContentConfig{
		ContentTypes[typeExperience]:    &ExperienceConfig{},
		ContentTypes[typeEducation]:     &EducationConfig{},
		ContentTypes[typeProject]:       &ProjectConfig{},
		ContentTypes[typeCertification]: &CertificationConfig{},
		ContentTypes[typeMe]:            &AboutMeConfig{},
	}
	// Regex which contains all possible content types
	rex = regexp.MustCompile(fmt.Sprintf("(%s)", strings.Join(ContentTypes, "|")))
)

// ContentTemplateData the data which must be passed to the content html templates
type ContentTemplateData struct {
	Type string
	HTML []template.HTML
}

// ContentConfig interface which represents a configuration of a content
type ContentConfig interface {
	// GetElements returns the elements loaded from the corresponding yaml file
	GetElements() []Content
	// GetConfigName returns the name of the corresponding yaml config file
	GetConfigName() string
	// GetContentType returns the type of content this object contains
	GetContentType() string
}

// Content which is contained in the yaml files
type Content interface {
	// GetTemplateName returns the name of the corresponding html template
	GetTemplateName() string
}

// ContentBase contains shared attributes for all content types
type ContentBase struct {
	// Image to render in the card
	Image string `yaml:"image"`
	// Name to display in the heading
	Name string `yaml:"name"`
	// Link to external content
	Link string `yaml:"link"`
	// Description displayed in the card body
	Description template.HTML `yaml:"description"`
}

// ContentDateRange specifies a range of two dates
type ContentDateRange struct {
	// From a date
	From time.Time `yaml:"from"`
	// To, may be string or date format
	To interface{} `yaml:"to"`
}

// GetFromDateAsStr returns the date formatted as string
func (d *ContentDateRange) GetFromDateAsStr() string {
	return d.From.Format("2006-01-02")
}

// GetToDateAsStr checks if to date is a date and formats if not a date
// the string content is returned, if not set this defaults to 'now'
func (d *ContentDateRange) GetToDateAsStr() string {
	if date, ok := d.To.(time.Time); ok {
		return date.Format("2006-01-02")
	} else if str, ok := d.To.(string); ok {
		return str
	}
	return "now"
}

// GetRenderedContent reads the content kind passed from its yaml configuration
// and returns all configured elements as html to be placed in the main
// template directly
func GetRenderedContent(contentType string) (*ContentTemplateData, error) {
	// Get the correct object to load
	// TODO validate so we do not have null values
	obj := contentMappings[contentType]

	err := loadContentConfig(obj)
	if nil != err {
		log.Printf("[ERROR] Generating content failed: %s\n", err)
		return nil, err
	}

	// render the content read from yaml into the html models
	cards := obj.GetElements()
	data := make([]template.HTML, 0)
	for _, crd := range cards {
		if tpl, err := renderContent(crd); nil != err {
			return nil, err
		} else {
			data = append(data, tpl)
		}

	}
	titledType := cases.Title(language.English).String(contentType)
	return &ContentTemplateData{Type: titledType, HTML: data}, err
}

// GetRoutingRegexString returns the regex which catches the endpoints for
// the content types as string
func GetRoutingRegexString() string {
	return rex.String()
}

// renderContent renders the passed content as html from its template
func renderContent(content Content) (template.HTML, error) {

	contentBaseTpl := filepath.Join(config.ContentTemplatesPath(), "base.html")
	htmlTpl := filepath.Join(config.ContentTemplatesPath(), content.GetTemplateName())

	rendered, err := apputils.RenderTemplate("content", contentBaseTpl, htmlTpl, content)
	if nil != err {
		log.Printf("[ERROR] Failed to parse template '%s': %s\n", htmlTpl, err)
		return "", err
	}

	return template.HTML(rendered), nil
}

// castToContent casts the array of a specific kind to an array of Content
// for further use
func castToContent[T Content](content []T) []Content {
	casted := make([]Content, len(content))
	for idx, crd := range content {
		casted[idx] = crd
	}
	return casted
}

// unmarshallContentConfig
func unmarshalContentConfig(content ContentConfig) error {
	return utils.LoadFromYAMLFile(content.GetConfigName(), content)
}

// loadContentConfig loads the specified content from it's yaml file
func loadContentConfig(content ContentConfig) error {
	err := unmarshalContentConfig(content)
	return err
}

// IsValidContentType returns if the content type passed is a valid one
func IsValidContentType(contentType string) bool {
	isValid := true
	// Check if all content kinds specified in the yaml config are valid
	if !rex.MatchString(contentType) {
		log.Printf("[ERROR] Invalid content kind '%s', allowed values are: %s\n", contentType, ContentTypes)
		isValid = false
	}
	return isValid
}
