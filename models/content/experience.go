package content

type ExperienceConfig struct {
	Experiences []*Experience `yaml:"experiences"`
}

// Make sure the interface is implemented
var _ ContentConfig = &ExperienceConfig{}

func (ed *ExperienceConfig) GetElements() []Content {
	return castToContent(ed.Experiences)
}

func (ed *ExperienceConfig) GetConfigName() string {
	return ed.GetContentKind() + ".yml"
}

func (ed *ExperienceConfig) GetContentKind() string {
	return ContentTypes[typeExperience]
}

type Experience struct {
	ContentBase      `yaml:",inline"`
	Company          string `yaml:"company"`
	ContentDateRange `yaml:",inline"`
}

// Make sure the interface is implemented
var _ Content = &Experience{}

func (e *Experience) GetTemplateName() string {
	return "experience.html"
}
