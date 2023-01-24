package models

// TemplateData is the object passed to all of the html template renderings
type TemplateData struct {
	RenderContact bool
	Data          interface{}
	Profile       *ProfileConfig
}
