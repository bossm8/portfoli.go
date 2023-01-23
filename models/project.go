package models

type ProjectConfig struct {
	Projects []*Project `yaml:"projects"`
}

// Make sure the interface is implemented
var _ listConfig = &ProjectConfig{}

func (pc *ProjectConfig) GetElements() []portfolioCard {
	return castToCard(pc.Projects)
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
