package models

import (
	"html/template"
	"time"
)

type ContentConfig interface {
	GetElements() []Content
	GetConfigName() string
	GetContentKind() string
}

type Content interface {
	GetTemplateName() string
}

type ContentBase struct {
	Image       string        `yaml:"image"`
	Name        string        `yaml:"name"`
	Link        string        `yaml:"link"`
	Description template.HTML `yaml:"description"`
}

type ContentDateRange struct {
	From time.Time   `yaml:"from"`
	To   interface{} `yaml:"to"`
}

func (d *ContentDateRange) GetFromDateAsStr() string {
	return d.From.Format("2006-01-02")
}

func (d *ContentDateRange) GetToDateAsStr() string {
	if date, ok := d.To.(time.Time); ok {
		return date.Format("2006-01-02")
	} else if str, ok := d.To.(string); ok {
		return str
	}
	return "now"
}
