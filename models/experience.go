package models

const experienceConfigName = "experience.yml"

type Experience struct {
	Base    `yaml:",inline"`
	Company string `yaml:"company"`
}

func GetExperiences() (exp []*Experience) {
	expCfg := &struct {
		Experiences []*Experience `yaml:"experiences"`
	}{}
	loadFromFile(experienceConfigName, expCfg)
	exp = expCfg.Experiences
	return
}
