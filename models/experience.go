package models

type ExperienceConfig struct {
	Experiences []*Experience `yaml:"experiences"`
}

// Make sure the interface is implemented
var _ listConfig = &ExperienceConfig{}

func (ed *ExperienceConfig) GetElements() []portfolioCard {
	return castToCard(ed.Experiences)
}

func (ed *ExperienceConfig) GetConfigName() string {
	return ed.GetContentKind() + ".yml"
}

func (ed *ExperienceConfig) Load() error {
	return unmarshal(ed)
}

func (ed *ExperienceConfig) GetContentKind() string {
	return kinds[kindExp]
}

type Experience struct {
	Base    `yaml:",inline"`
	Company string `yaml:"company"`
}

// Make sure the interface is implemented
var _ portfolioCard = &Experience{}

func (e *Experience) GetTemplateName() string {
	return "experience.html"
}

func LoadExperiences() (listConfig, error) {
	exp := &ExperienceConfig{}
	err := exp.Load()
	return exp, err
}
