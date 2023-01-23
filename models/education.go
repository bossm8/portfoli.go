package models

import "html/template"

type EducationConfig struct {
	Educations []*Education `yaml:"educations"`
}

// Make sure the interface is implemented
var _ listConfig = &EducationConfig{}

func (ec *EducationConfig) GetRenderedElements() ([]template.HTML, error) {
	data := make([]template.HTML, len(ec.Educations))
	for idx, edu := range ec.Educations {
		rendered, err := renderCard(edu)
		if nil != err {
			return nil, err
		}
		data[idx] = rendered
	}
	return data, nil
}

func (ec *EducationConfig) GetConfigName() string {
	return ec.GetContentKind() + ".yml"
}

func (ec *EducationConfig) Load() error {
	return unmarshal(ec)
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

func LoadEducation() (listConfig, error) {
	edu := &EducationConfig{}
	err := edu.Load()
	return edu, err
}
