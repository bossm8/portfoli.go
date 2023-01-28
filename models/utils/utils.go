package utils

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// The directory where all (static and dynamic) configuration files are read from
var yamlDir string

// SetYAMLDir sets the configuration directory where dynamic and static configurations
// should be read from
func SetYAMLDir(dir string) {
	yamlDir = dir
}

// LoadFromYAMLFile loads the file with filename into obj
func LoadFromYAMLFile(filename string, obj interface{}) (err error) {
	// log.Printf("[INFO] Loading yaml file '%s' from directory '%s'\n", filename, configDir)
	var yamlFile []byte
	if yamlFile, err = os.ReadFile(
		filepath.Join(yamlDir, filename),
	); err != nil {
		log.Printf("[ERROR]: Failed to load yaml file '%s': %s\n", filename, err)
		return
	}
	if err = yaml.Unmarshal(yamlFile, obj); err != nil {
		log.Printf("[ERROR]: Failed to parse yaml file '%s': %s\n", filename, err)
	}
	return
}
