package models

type TemplateData struct {
	RenderContact bool
	Data          interface{}
	Profile       *ProfileConfig
}
