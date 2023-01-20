package models

const educationConfig = "education.yml"

type Education struct {
	Base `yaml:",inline"`
}

type Educations struct {
	Educations []*Education `yaml:"educations"`
}
