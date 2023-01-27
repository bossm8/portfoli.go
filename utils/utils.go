package utils

import (
	"bossm8.ch/portfolio/models"
	"bossm8.ch/portfolio/models/config"
)

func RenderTemplate() {

}

func LoadConfiguration(configDir string) (cfg *config.Config, err error) {
	models.SetConfigDir(configDir)
	cfg, err = config.GetConfig()
	return
}
