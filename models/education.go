package models

type EducationConfig struct {
	Educations []*Education `yaml:"educations"`
}

// Make sure the interface is implemented
var _ listConfig = &EducationConfig{}

func (ec *EducationConfig) GetElements() []portfolioCard {
	return castToCard(ec.Educations)
}

func (ec *EducationConfig) GetConfigName() string {
	return ec.GetContentKind() + ".yml"
}

func (ec *EducationConfig) GetContentKind() string {
	return kinds[kindEdu]
}

type Education struct {
	Base           `yaml:",inline"`
	School         string `yaml:"school"`
	Specialization string `yaml:"specialization"`
}

// Make sure the interface is implemented
var _ portfolioCard = &Education{}

func (e *Education) GetTemplateName() string {
	return "education.html"
}
