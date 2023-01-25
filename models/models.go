// models contains all models for the application
// meaning the configurable parts which will be used to render the webpage
package models

import (
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"strings"
)

const (
	contentKindExperience = iota
	contentKindEducation
	contenKindProject
	contentkindCertification
)

var (
	templatesDir    = filepath.Join("models", "templates", "html")
	contentKinds    = []string{"experience", "education", "projects", "certifications"}
	contentMappings = map[string]ContentConfig{
		contentKinds[contentKindExperience]:    &ExperienceConfig{},
		contentKinds[contentKindEducation]:     &EducationConfig{},
		contentKinds[contenKindProject]:        &ProjectConfig{},
		contentKinds[contentkindCertification]: &CertificationConfig{},
	}
)

// GetRenderedContent reads the content kind passed from its yaml configuration
// and returns all configured elements as html to be placed in the main
// template directly
func GetRenderedContent(kind string) ([]template.HTML, error) {
	// get the correct object to load
	obj := contentMappings[kind]

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
	return data, err
}

func GetRoutingRegexString() string {
	return fmt.Sprintf("(%s)", strings.Join(contentKinds, "|"))
}
