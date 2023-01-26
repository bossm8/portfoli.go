// models contains all models for the application
// meaning the configurable parts which will be used to render the webpage
package models

import (
	"bossm8.ch/portfolio/models/config"
	"bossm8.ch/portfolio/models/utils"
)

// TemplateData is the object passed to all of the html template renderings
type TemplateData struct {
	RenderContact bool
	Data          interface{}
	Profile       *config.ProfileConfig
}

func SetConfigDir(dir string) {
	utils.SetYAMLDir(dir)
}
