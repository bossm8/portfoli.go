// Copyright (c) 2023, Boss Marco <bossm8@hotmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package content

import (
	"html/template"
	"log"
	"path/filepath"
	"time"

	"github.com/bossm8/portfoli.go/config"
	apputils "github.com/bossm8/portfoli.go/utils"
)

const (
	cardsTpl = "cards.html"
)

// Card defines an element which will be rendered as a card
type Card interface {
	// CardTemplateName returns the name of the corresponding card html template
	CardTemplateName() string
}

// CardContentConfig defines a specific type of content config,
// which contains a list of elements to be renderer as bootstrap cards
type CardContentConfig interface {
	// Elements returns the elements loaded from the corresponding yaml file
	// Each of those will be rendered in its own card
	Elements() []Card
}

// CardBase contains shared attributes for all card content types
type CardBase struct {
	// Image to render in the card
	Image string `yaml:"image"`
	// Name to display in the heading
	Name string `yaml:"name"`
	// Link to external content
	Link string `yaml:"link"`
	// Description displayed in the card body
	Description template.HTML `yaml:"description"`
}

// CardDateRange specifies a range of two dates
type CardDateRange struct {
	// From a date
	From time.Time `yaml:"from"`
	// To, may be string or date format
	To interface{} `yaml:"to"`
}

// GetFromDateAsStr returns the date formatted as string
func (d *CardDateRange) GetFromDateAsStr() string {
	return d.From.Format("2006-01-02")
}

// GetToDateAsStr checks if to date is a date and formats if not a date
// the string content is returned, if not set this defaults to 'now'
func (d *CardDateRange) GetToDateAsStr() string {
	if date, ok := d.To.(time.Time); ok {
		return date.Format("2006-01-02")
	} else if str, ok := d.To.(string); ok {
		return str
	}
	return "now"
}

// castToCard casts the array of a specific kind to an array of Cards
// for further use
func castToCard[T Card](content []T) []Card {
	casted := make([]Card, len(content))
	for idx, crd := range content {
		casted[idx] = crd
	}
	return casted
}

// renderCard renders the passed content as html from its template
func renderCard(card Card) (template.HTML, error) {

	contentBaseTpl := filepath.Join(config.ContentTemplatesPath(), "base.html")
	htmlTpl := filepath.Join(config.ContentTemplatesPath(), card.CardTemplateName())

	rendered, err := apputils.RenderTemplate("content", card, contentBaseTpl, htmlTpl)
	if nil != err {
		log.Printf("[ERROR] Failed to parse template '%s': %s\n", htmlTpl, err)
		return "", err
	}

	return template.HTML(rendered), nil
}

func renderCards(obj CardContentConfig, cardType string) (*template.HTML, error) {
	// render the content read from yaml into the html models
	cards := obj.Elements()
	data := make([]template.HTML, 0)
	for _, crd := range cards {
		if tpl, err := renderCard(crd); nil != err {
			return nil, err
		} else {
			data = append(data, tpl)
		}
	}

	cardData := struct {
		Type  string
		Cards []template.HTML
	}{
		Type:  cardType,
		Cards: data,
	}

	baseTpl := filepath.Join(config.ContentTemplatesPath(), cardsTpl)
	rendered, err := apputils.RenderTemplate("cards", &cardData, baseTpl)
	if err != nil {
		log.Printf("[ERROR] Failed to render %s\n", baseTpl)
		return nil, err
	}
	html := template.HTML(rendered)

	return &html, nil
}
