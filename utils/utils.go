package utils

import (
	"bytes"
	"html/template"
	"log"

	"bossm8.ch/portfolio/models"
	"bossm8.ch/portfolio/models/config"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func RenderTemplate(name string, baseTemplate string, childTemplate string, data *models.TemplateData) (*bytes.Buffer, error) {

	var tpl *template.Template
	var err error

	// Title is used in templates to title case content kind names
	funcMap := template.FuncMap{
		"Title": cases.Title(language.English).String,
	}
	if tpl, err = template.New(childTemplate).Funcs(funcMap).ParseFiles(baseTemplate, childTemplate); nil != err {
		log.Printf("[ERROR] Failed to parse template: %s with error %s\n", childTemplate, err)
		return nil, err
	}

	resp := &bytes.Buffer{}
	if err = tpl.ExecuteTemplate(resp, name, data); nil != err {
		log.Printf("[ERROR] Failed to process template %s with error %s\n", tpl.Name(), err)
		return nil, err
	}

	return resp, nil
}

func LoadConfiguration(configDir string) (cfg *config.Config, err error) {
	models.SetConfigDir(configDir)
	cfg, err = config.GetConfig()
	return
}
