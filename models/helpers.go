package models

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func loadFromYAMLFile(filename string, obj interface{}) (err error) {
	// log.Printf("[INFO] Loading yaml file '%s' from directory '%s'\n", filename, configDir)
	var yamlFile []byte
	if yamlFile, err = os.ReadFile(
		filepath.Join(configDir, filename),
	); err != nil {
		log.Printf("[ERROR]: Failed to load yaml file '%s': %s\n", filename, err)
		return
	}
	if err = yaml.Unmarshal(yamlFile, obj); err != nil {
		log.Printf("[ERROR]: Failed to parse yaml file '%s': %s\n", filename, err)
	}
	return
}

func unmarshal(obj listConfig) error {
	return loadFromYAMLFile(obj.GetConfigName(), obj)
}

func renderCard(card portfolioCard) (template.HTML, error) {
	htmlTpl := filepath.Join(templatesDir, card.GetTemplateName())
	tpl, err := template.ParseFiles(htmlTpl)
	if nil != err {
		log.Printf("[ERROR] Failed to parse template '%s': %s\n", htmlTpl, err)
		return "", err
	}
	rendered := bytes.Buffer{}
	if err := tpl.ExecuteTemplate(&rendered, "card", card); nil != err {
		log.Printf("[ERROR] Failed to process template %s with error %s\n", tpl.Name(), err)
		return "", err
	}

	return template.HTML(rendered.String()), nil
}

func castToCard[T portfolioCard](cards []T) []portfolioCard {
	casted := make([]portfolioCard, len(cards))
	for idx, crd := range cards {
		casted[idx] = crd
	}
	return casted
}
