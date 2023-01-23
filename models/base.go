package models

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var confPath = "configs"

type Base struct {
	Img         string `yaml:"img"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

func loadFromFile(filename string, obj interface{}) (err error) {
	log.Printf("[INFO] Loading yaml file '%s' from directory '%s'\n", filename, confPath)
	var yamlFile []byte
	if yamlFile, err = os.ReadFile(
		filepath.Join(confPath, filename),
	); err != nil {
		log.Printf("[ERROR]: Failed to load yaml file '%s': %s\n", filename, err)
		return
	}
	if err = yaml.Unmarshal(yamlFile, obj); err != nil {
		log.Printf("[ERROR]: Failed to parse yaml file '%s': %s\n", filename, err)
	}
	return
}
