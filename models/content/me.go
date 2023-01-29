package content

import "html/template"

type AboutMeConfig struct {
	Content template.HTML `yaml:"me"`
}

func (a *AboutMeConfig) GetConfigName() string {
	return a.GetContentType() + ".yml"
}

func (a *AboutMeConfig) GetElements() []Content {
	return nil
}

func (a *AboutMeConfig) GetContentType() string {
	return ContentTypes[typeMe]
}
