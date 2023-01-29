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

// package utils contains utility functions for static and dynamic builds
package utils

import (
	"bytes"
	"html/template"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/yosssi/gohtml"
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
		res := filepath.Join(basePath, path)
		if filepath.IsAbs(res) {
			return res
		}
		return string(filepath.Separator) + res
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

// checkFuncsInitializedOrAbort makes sure that the function maps were initialized
// by calling Init, if not is will abort the program, as it is a programmer error
func checkFuncsInitializedOrAbort() {
	if funcMap == nil {
		log.Fatal(
			"[ERROR] Please call util.Init at least once before rendering a template",
		)
	}
}

// RenderTemplate renders the baseTemplate containing an optional childTemplate with the data
// passed. (name) is passed to ExecuteTemplate
func RenderTemplate(
	tplName string,
	data interface{},
	templates ...string,
) ([]byte, error) {
	checkFuncsInitializedOrAbort()

	var tpl *template.Template
	var err error

	// Title is used in templates to title case content kind names
	if tpl, err = template.New(tplName).Funcs(funcMap).ParseFiles(templates...); nil != err {
		log.Printf("[ERROR] Failed to parse templates: %s with error %s\n", templates, err)
		return nil, err
	}

	resp := &bytes.Buffer{}
	if err = tpl.ExecuteTemplate(resp, tplName, data); nil != err {
		log.Printf("[ERROR] Failed to process template %s with error %s\n", tpl.Name(), err)
		return nil, err
	}
	pretty := gohtml.FormatBytes(resp.Bytes())

	return pretty, nil
}

// ProcessHTMLContent takes a html template which could contain some template
// pipelines and passes them through template.Execute.
// This makes it possible to have e.g. Assemble in the content configs.
func ProcessHTMLContent(html *template.HTML) (*template.HTML, error) {
	checkFuncsInitializedOrAbort()
	tpl, err := template.New("html").Funcs(funcMap).Parse(string(*html))
	if err != nil {
		log.Printf("[ERROR Parsing html content %s\n", err)
		return nil, err
	}

	res := &bytes.Buffer{}
	if err := tpl.Execute(res, nil); err != nil {
		log.Printf("[ERROR] Executing template on html failed: %s\n", err)
		return nil, err
	}
	newHTML := template.HTML(res.String())
	return &newHTML, nil
}
