package models

import "html/template"

const educationConfigName = "education.yml"

type EducationConfig struct {
	Educations []*PortfolioCard `yaml:"educations"`
}

// Make sure the interface is implemented
var _ listConfig = &EducationConfig{}

func (cc *EducationConfig) GetElements() []*PortfolioCard {
	return cc.Educations
}

func (ec *EducationConfig) GetConfigName() string {
	return educationConfigName
}

type Education struct {
	Base           `yaml:",inline"`
	Specialization string `yaml:"specialization"`
}

// Make sure the interface is implemented
var _ PortfolioCard = &Education{}

func GetEducations() []*PortfolioCard {
	return unmarshal(&EducationConfig{})
}

func (e *Education) GetHTMLTemplate() template.HTML {
	return ""
}
