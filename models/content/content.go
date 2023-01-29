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
	"regexp"
	"strings"

	apputils "github.com/bossm8/portfoli.go/utils"

	"github.com/bossm8/portfoli.go/models/utils"
)

const (
	typeExperience = iota
	typeEducation
	typeProject
	typeCertification
	typeAbout
)

var (
	// All possible content types
	ContentTypes = []string{"experience", "education", "projects", "certifications", "bio"}
	// Mappings to easily get the correct content type based on the request path
	contentMappings = map[string]ContentConfig{
		ContentTypes[typeExperience]:    &ExperienceConfig{},
		ContentTypes[typeEducation]:     &EducationConfig{},
		ContentTypes[typeProject]:       &ProjectConfig{},
		ContentTypes[typeCertification]: &CertificationConfig{},
		ContentTypes[typeAbout]:         &AboutMeConfig{},
	}
	// Regex which contains all possible content types
	rex = regexp.MustCompile(fmt.Sprintf("(%s)", strings.Join(ContentTypes, "|")))
)

// ContentTemplateData the data which must be passed to the content html templates
type ContentTemplateData struct {
	Title string
	HTML  *template.HTML
}

// ContentConfig interface which represents a configuration of a content
type ContentConfig interface {
	// ConfigName returns the name of the corresponding yaml config file
	ConfigName() string
	// ContentType returns the type of content this object contains
	ContentType() string
	// Title returns the title of the content which can be used in the templates
	Title() string
	// Render returns the rendered html to be placed in the content template
	Render() (*template.HTML, error)
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
		log.Printf("[ERROR] Loading content failed: %s\n", err)
		return nil, err
	}

	data, err := obj.Render()
	if err != nil {
		log.Printf("[ERROR] Failed to render content for %s: %s\n", contentType, err)
		return nil, err
	}

	// make sure the html content is processed, so Assemble (for example) can
	// be used in the html configuration
	data, err = apputils.ProcessHTMLContent(data)
	if err != nil {
		log.Printf("[ERROR] HTML content prossecing of %s failed\n", contentType)
	}

	return &ContentTemplateData{Title: obj.Title(), HTML: data}, err
}

// GetRoutingRegexString returns the regex which catches the endpoints for
// the content types as string
func GetRoutingRegexString() string {
	return rex.String()
}

// unmarshallContentConfig
func unmarshalContentConfig(content ContentConfig) error {
	return utils.LoadFromYAMLFile(content.ConfigName(), content)
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
