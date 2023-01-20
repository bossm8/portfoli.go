package models

import "fmt"

const experienceConfig = "experience.yml"

type Experience struct {
	Base    `yaml:",inline"`
	Company string
}

type Experiences struct {
	Experiences []*Experience `yaml:"experiences"`
}

func GetExperiences() (exp Experiences) {
	loadFromFile(experienceConfig, &exp)
	fmt.Printf("%s\n", exp.Experiences[0].Company)
	fmt.Printf("%s\n", exp.Experiences[0].Img)
	fmt.Printf("%s\n", exp.Experiences[0].Name)
	return
}
