package models

const educationConfigName = "education.yml"

type Education struct {
	Base           `yaml:",inline"`
	Specialization string `yaml:"specialization"`
}

func GetEducations() (edu []*Education) {
	eduCfg := &struct {
		Educations []*Education `yaml:"educations"`
	}{}
	loadFromFile(educationConfigName, eduCfg)
	edu = eduCfg.Educations
	return
}
