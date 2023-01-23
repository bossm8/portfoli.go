package models

import "html/template"

const projectConfigName = "projects.yml"

type ProjectConfig struct {
	Projects []*PortfolioCard `yaml:"projects"`
}

// Make sure the interface is implemented
var _ listConfig = &ProjectConfig{}

func (pc *ProjectConfig) GetElements() []*PortfolioCard {
	return pc.Projects
}

func (pc *ProjectConfig) GetConfigName() string {
	return projectConfigName
}

type Project struct {
	Base `yaml:",inline"`
}

// Make sure the interface is implemented
var _ PortfolioCard = &Project{}

func GetProjects() (proj []*Project) {
	return
}

func (p *Project) GetHTMLTemplate() template.HTML {
	return ""
}
