package static

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"bossm8.ch/portfolio/messages"
	"bossm8.ch/portfolio/models"
	"bossm8.ch/portfolio/models/config"
	"bossm8.ch/portfolio/models/content"
	"bossm8.ch/portfolio/utils"
	"github.com/yosssi/gohtml"
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
	buildGeneric()
	buildContent()
	buildErrors()
}

func buildGeneric() {
	templates, err := os.ReadDir("templates/html")
	if nil != err {
		log.Fatalf("[ERROR] Could not read template directory: %s\n", err)
	}

	for _, tpl := range templates {
		if cfg.StaticIgnoreRegex().MatchString(tpl.Name()) ||
			cfg.Profile.Me == nil && tpl.Name() == cfg.AboutMeTemplateName+".html" {
			continue
		}

		build(
			tpl.Name(),
			tpl.Name(),
			nil,
		)
	}
}

func buildContent() {
	for _, contentType := range cfg.Profile.ContentKinds {
		content, err := content.GetRenderedContent(contentType)
		if nil != err {
			log.Fatalf("[ERROR] Rendering content %s: %s\n", contentType, err)
		}
		build(
			cfg.ContentTemplateName+".html",
			contentType+".html",
			content,
		)
	}
}

func buildErrors() {
	msg := messages.Get(string(messages.EndpointFail), string(messages.MsgNotFound))
	build(
		cfg.StatusTemplateName+".html",
		fmt.Sprintf("%d.html", msg.HttpStatus),
		msg,
	)
}

func build(tplFileName string, outputFileName string, data interface{}) {
	log.Printf(
		"[INFO] Rendering template %s to %s in %s\n",
		tplFileName,
		outputFileName,
		cfg.DistDir(),
	)

	htmlTpl := filepath.Join(cfg.HTMLTeplatesDir, tplFileName)

	tplData := &models.TemplateData{
		Profile:       cfg.Profile,
		RenderContact: false,
		Data:          data,
	}

	resp, err := utils.RenderTemplate(
		cfg.BaseTemplateName,
		cfg.BaseTemplatePath,
		htmlTpl,
		tplData,
	)
	if nil != err {
		log.Fatalf("[Error] Failed to render template: %s\n", err)
	}

	outputFile := filepath.Join(cfg.DistDir(), outputFileName)
	prettyHTML := gohtml.FormatBytes(resp.Bytes())
	if err := os.WriteFile(outputFile, prettyHTML, 0664); nil != err {
		log.Fatalf("[ERROR] Failed to write template: %s\n", err)
	}

}
