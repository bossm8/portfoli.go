package models

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// var confPath = filepath.Join("..", "configs")
var confPath = "configs"

type Base struct {
	Img         string `yaml:"img"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

func loadFromFile(filename string, obj interface{}) {
	yamlFile, err := os.ReadFile(
		filepath.Join(confPath, filename),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := yaml.Unmarshal(yamlFile, obj); err != nil {
		log.Fatal(err)
	}
}
