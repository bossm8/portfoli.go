package models

type EducationConfig struct {
	Educations []*Education `yaml:"educations"`
}

// Make sure the interface is implemented
var _ ContentConfig = &EducationConfig{}

func (ec *EducationConfig) GetElements() []Content {
	return castToContent(ec.Educations)
}

func (ec *EducationConfig) GetConfigName() string {
	return ec.GetContentKind() + ".yml"
}

func (ec *EducationConfig) GetContentKind() string {
	return contentKinds[contentKindEducation]
}

type Education struct {
	ContentBase      `yaml:",inline"`
	School           string `yaml:"school"`
	Specialization   string `yaml:"specialization"`
	ContentDateRange `yaml:",inline"`
}

// Make sure the interface is implemented
var _ Content = &Education{}

func (e *Education) GetTemplateName() string {
	return "education.html"
}
