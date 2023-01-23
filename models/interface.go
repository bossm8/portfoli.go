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
	mapping      = map[string]func() (listConfig, error){
		kinds[kindExp]:  LoadExperiences,
		kinds[kindEdu]:  LoadEducation,
		kinds[kindProj]: LoadProjects,
		kinds[kindCert]: LoadCertifications,
	}
)

type Base struct {
	Img         string `yaml:"img"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type listConfig interface {
	GetElements() []portfolioCard
	GetConfigName() string
	GetContentKind() string
	Load() error
}

type portfolioCard interface {
	GetTemplateName() string
}

func GetContent(kind string) ([]template.HTML, error) {
	content, err := mapping[kind]()
	if nil != err {
		log.Printf("[ERROR] Generating content failed: %s\n", err)
		return nil, err
	}
	cards := content.GetElements()
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
