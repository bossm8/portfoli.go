// Package static is used to build a static version of this template
package static

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	appconfig "github.com/bossm8/portfoli.go/config"

	"github.com/bossm8/portfoli.go/messages"
	"github.com/bossm8/portfoli.go/models"
	"github.com/bossm8/portfoli.go/models/config"
	"github.com/bossm8/portfoli.go/models/content"
	"github.com/bossm8/portfoli.go/utils"

	"github.com/yosssi/gohtml"
)

var (
	cfg *config.Config
)

// Build builds the static website by using the configs found in configDir
func Build(srvBasePath string, configDir string) {

	var err error
	cfg, err = models.LoadConfiguration(configDir)
	if nil != err && !errors.Is(err, config.ErrInvalidSMTPConfig) {
		log.Fatalf("[ERROR] Loading configuration failed: %s\n", err)
	}

	utils.Init(srvBasePath)
	messages.Compile(cfg.Profile.Email.Address)

	buildGeneric()
	buildContent()
	buildErrors()
}

// buildGeneric builds every page except contents and error
func buildGeneric() {
	templates, err := os.ReadDir("templates/html")
	if nil != err {
		log.Fatalf("[ERROR] Could not read template directory: %s\n", err)
	}

	for _, tpl := range templates {
		if appconfig.StaticIgnoreRegex().MatchString(tpl.Name()) ||
			cfg.Profile.Me == nil && tpl.Name() == appconfig.AboutMeTemplateName+".html" {
			continue
		}

		build(
			tpl.Name(),
			tpl.Name(),
			nil,
		)
	}
}

// buildContent builds the content pages
func buildContent() {
	for _, contentType := range cfg.Profile.ContentTypes {
		content, err := content.GetRenderedContent(contentType)
		if nil != err {
			log.Fatalf("[ERROR] Rendering content %s: %s\n", contentType, err)
		}
		build(
			appconfig.ContentTemplateName+".html",
			contentType+".html",
			content,
		)
	}
}

// buildError builds the error pages (which in case of static is 404 only)
func buildErrors() {
	msg := messages.Get(string(messages.EndpointFail), string(messages.MsgNotFound))
	build(
		appconfig.StatusTemplateName+".html",
		fmt.Sprintf("%d.html", msg.HttpStatus),
		msg,
	)
}

// build - generic method to build the template tplFileName to outputFileName
// with data
func build(tplFileName string, outputFileName string, data interface{}) {
	log.Printf(
		"[INFO] Rendering template %s to %s in %s\n",
		tplFileName,
		outputFileName,
		appconfig.DistDir(),
	)

	htmlTpl := filepath.Join(appconfig.HTMLTemplatesPath(), tplFileName)

	tplData := &models.TemplateData{
		Data:          data,
		Profile:       cfg.Profile,
		RenderContact: false,
	}

	resp, err := utils.RenderTemplate(
		appconfig.BaseTemplateName,
		appconfig.BaseTemplatePath(),
		htmlTpl,
		tplData,
	)
	if nil != err {
		log.Fatalf("[Error] Failed to render template: %s\n", err)
	}

	outputFile := filepath.Join(appconfig.DistDir(), outputFileName)
	prettyHTML := gohtml.FormatBytes(resp.Bytes())
	if err := os.WriteFile(outputFile, prettyHTML, 0664); nil != err {
		log.Fatalf("[ERROR] Failed to write template: %s\n", err)
	}

}
