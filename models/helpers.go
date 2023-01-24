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

func unmarshalContentConfig(content ContentConfig) error {
	return loadFromYAMLFile(content.GetConfigName(), content)
}

func loadContentConfig(content ContentConfig) error {
	err := unmarshalContentConfig(content)
	return err
}

func renderContent(content Content) (template.HTML, error) {
	contentBaseTpl := filepath.Join(templatesDir, "base.html")
	htmlTpl := filepath.Join(templatesDir, content.GetTemplateName())
	tpl, err := template.ParseFiles(contentBaseTpl, htmlTpl)
	if nil != err {
		log.Printf("[ERROR] Failed to parse template '%s': %s\n", htmlTpl, err)
		return "", err
	}
	rendered := bytes.Buffer{}
	if err := tpl.ExecuteTemplate(&rendered, "content", content); nil != err {
		log.Printf("[ERROR] Failed to process template %s with error %s\n", tpl.Name(), err)
		return "", err
	}

	return template.HTML(rendered.String()), nil
}

func castToContent[T Content](content []T) []Content {
	casted := make([]Content, len(content))
	for idx, crd := range content {
		casted[idx] = crd
	}
	return casted
}
