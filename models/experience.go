package models

import "html/template"

const experienceConfigName = "experience.yml"

type ExperienceConfig struct {
	Experiences []*PortfolioCard `yaml:"experiences"`
}

// Make sure the interface is implemented
var _ listConfig = &ExperienceConfig{}

func (ec *ExperienceConfig) GetElements() []*PortfolioCard {
	return ec.Experiences
}

func (ed *ExperienceConfig) GetConfigName() string {
	return experienceConfigName
}

type Experience struct {
	Base    `yaml:",inline"`
	Company string `yaml:"company"`
}

// Make sure the interface is implemented
var _ PortfolioCard = &Experience{}

func GetExperiences() []*PortfolioCard {
	return unmarshal(&ExperienceConfig{})
}

func (e *Experience) GetHTMLTemplate() template.HTML {
	return ""
}
