package content

type ProjectConfig struct {
	Projects []*Project `yaml:"projects"`
}

// Make sure the interface is implemented
var _ ContentConfig = &ProjectConfig{}

func (pc *ProjectConfig) GetElements() []Content {
	return castToContent(pc.Projects)
}

func (pc *ProjectConfig) GetConfigName() string {
	return pc.GetContentType() + ".yml"
}

func (pc *ProjectConfig) GetContentType() string {
	return ContentTypes[typeProject]
}

type Project struct {
	ContentBase `yaml:",inline"`
}

// Make sure the interface is implemented
var _ Content = &Project{}

func (p *Project) GetTemplateName() string {
	return "project.html"
}
