// models contains all models for the application
// meaning the configurable parts which will be used to render the webpage
package models

import (
	"github.com/bossm8/portfoli.go/models/config"
	"github.com/bossm8/portfoli.go/models/utils"
)

// TemplateData is the object passed to all of the html template renderings
type TemplateData struct {
	RenderContact bool
	Data          interface{}
	Profile       *config.ProfileConfig
	BasePath      string
}

// SetConfigDir sets the directory to search for yaml configurations to dir
func setConfigDir(dir string) {
	utils.SetYAMLDir(dir)
}

// LoadConfiguration loads the static configuration from configDir
func LoadConfiguration(configDir string) (cfg *config.Config, err error) {
	setConfigDir(configDir)
	cfg, err = config.GetConfig()
	return
}
