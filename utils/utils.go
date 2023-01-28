// package utils contains utility functions for static and dynamic builds
package utils

import (
	"bytes"
	"html/template"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	funcMap template.FuncMap = nil
)

// assembleBasePath return a function which adds the server base path to
// any string passed to it
func assembleBasePath(basePath string) func(string) string {
	return func(path string) string {
		u, err := url.ParseRequestURI(path)
		if err == nil && u.Scheme != "" && u.Host != "" {
			return path
		}
		path = strings.TrimPrefix(path, basePath)
		return filepath.Join(basePath, path)
	}
}

// Init initializes the functions needed for rendering templates
func Init(serverBasePath string) {

	if !filepath.IsAbs(serverBasePath) {
		serverBasePath = "/" + serverBasePath
	}
	if !strings.HasSuffix(serverBasePath, "/") {
		serverBasePath = serverBasePath + "/"
	}
	log.Printf("[INFO] Using %s as server base path\n", serverBasePath)

	funcMap = template.FuncMap{
		"Title":    cases.Title(language.English).String,
		"Assemble": assembleBasePath(serverBasePath),
	}
}

// RenderTemplate renders the baseTemplate containing a childTemplate with the data
// passed. (name) is passed to ExecuteTemplate
func RenderTemplate(
	tplName string,
	baseTemplate string,
	childTemplate string,
	data interface{},
) (*bytes.Buffer, error) {

	if funcMap == nil {
		log.Fatal(
			"[ERROR] Please call util.Init at least once before rendering a template",
		)
	}

	var tpl *template.Template
	var err error

	// Title is used in templates to title case content kind names
	if tpl, err = template.New(tplName).Funcs(funcMap).ParseFiles(baseTemplate, childTemplate); nil != err {
		log.Printf("[ERROR] Failed to parse template: %s with error %s\n", childTemplate, err)
		return nil, err
	}

	resp := &bytes.Buffer{}
	if err = tpl.ExecuteTemplate(resp, tplName, data); nil != err {
		log.Printf("[ERROR] Failed to process template %s with error %s\n", tpl.Name(), err)
		return nil, err
	}

	return resp, nil
}
