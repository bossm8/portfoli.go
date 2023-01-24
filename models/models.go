package models

import (
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"strings"
)

const (
	kindExp = iota
	kindEdu
	kindProj
	kindCert
)

var (
	templatesDir = filepath.Join("models", "templates", "html")
	kinds        = []string{"experience", "education", "projects", "certifications"}
	mapping      = map[string]listConfig{
		kinds[kindExp]:  &ExperienceConfig{},
		kinds[kindEdu]:  &EducationConfig{},
		kinds[kindProj]: &ProjectConfig{},
		kinds[kindCert]: &CertificationConfig{},
	}
)

func GetContent(kind string) ([]template.HTML, error) {
	obj := mapping[kind]
	err := loadListConfig(obj)
	if nil != err {
		log.Printf("[ERROR] Generating content failed: %s\n", err)
		return nil, err
	}
	cards := obj.GetElements()
	data := make([]template.HTML, 0)
	for _, crd := range cards {
		if tpl, err := renderCard(crd); nil != err {
			return nil, err
		} else {
			data = append(data, tpl)
		}

	}
	return data, err
}

func GetRoutingRegex() string {
	return fmt.Sprintf("/(%s)", strings.Join(kinds, "|"))
}