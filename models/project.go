package models

import "html/template"

type ProjectConfig struct {
	Projects []*Project `yaml:"projects"`
}

// Make sure the interface is implemented
var _ listConfig = &ProjectConfig{}

func (pc *ProjectConfig) GetRenderedElements() ([]template.HTML, error) {
	data := make([]template.HTML, len(pc.Projects))
	for idx, proj := range pc.Projects {
		rendered, err := renderCard(proj)
		if nil != err {
			return nil, err
		}
		data[idx] = rendered
	}
	return data, nil
}

func (pc *ProjectConfig) GetConfigName() string {
	return pc.GetContentKind() + ".yml"
}

func (pc *ProjectConfig) Load() error {
	return unmarshal(pc)
}

func (pc *ProjectConfig) GetContentKind() string {
	return kinds[kindProj]
}

type Project struct {
	Base `yaml:",inline"`
}

// Make sure the interface is implemented
var _ portfolioCard = &Project{}

func (p *Project) GetTemplateName() string {
	return "project.html"
}

func LoadProjects() (listConfig, error) {
	proj := &ProjectConfig{}
	err := proj.Load()
	return proj, err
}
