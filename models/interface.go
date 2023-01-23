package models

import (
	"html/template"
	"time"
)

type listConfig interface {
	GetElements() []portfolioCard
	GetConfigName() string
	GetContentKind() string
}

type portfolioCard interface {
	GetTemplateName() string
}

type Base struct {
	Image       string        `yaml:"image"`
	Name        string        `yaml:"name"`
	Link        string        `yaml:"link"`
	Description template.HTML `yaml:"description"`
}

type DateRange struct {
	From time.Time   `yaml:"from"`
	To   interface{} `yaml:"to"`
}

func (d *DateRange) GetFromDateAsStr() string {
	return d.From.Format("2006-01-02")
}

func (d *DateRange) GetToDateAsStr() string {
	if date, ok := d.To.(time.Time); ok {
		return date.Format("2006-01-02")
	} else if str, ok := d.To.(string); ok {
		return str
	}
	return "now"
}
