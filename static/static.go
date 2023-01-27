package static

import (
	"errors"
	"io/ioutil"
	"log"

	"bossm8.ch/portfolio/messages"
	"bossm8.ch/portfolio/models/config"
	"bossm8.ch/portfolio/utils"
)

var (
	cfg *config.Config
)

func Build(configDir string) {
	var err error
	cfg, err = utils.LoadConfiguration(configDir)
	if nil != err && !errors.Is(err, config.ErrInvalidSMTPConfig) {
		log.Fatalf("[ERROR] Loading configuration failed: %s\n", err)
	}
	messages.Compile(cfg.Profile.Email.Address)
	build()
}

func build() {
	templates, err := ioutil.ReadDir("templates/html")
	if nil != err {
		log.Fatalf("[ERROR] Could not read template directory: %s\n", err)
	}

	for _, tpl := range templates {
		if cfg.StaticIgnoreRegex().MatchString(tpl.Name()) {
			continue
		}
		log.Println(tpl.Name())
	}
}
