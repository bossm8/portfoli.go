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

// Package config contains the application configuration
package config

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
)

const (
	BaseTemplateName    = "base"
	ContentTemplateName = "content"
	StatusTemplateName  = "status"
	ContactTemplateName = "contact"
	MailTemplate        = "mail.html"
)

var (
	wwwDir = filepath.Join("/", "var", "www", "portfoli.go")

	templatesDir = filepath.Join(wwwDir, "templates")
	staticDir    = filepath.Join(wwwDir, "public")
	distDir      = filepath.Join(wwwDir, "dist")
)

// ConvertToAbsPath takes the path of a directory
// (either relative or absolute) and converts it to an absolute path
func ConvertToAbsPath(path *string) string {
	abs, err := filepath.Abs(*path)
	if err != nil {
		log.Fatalf("[ERROR] parsing configuration directory: %s\n", err)
	}
	return abs
}

// SetPaths sets the paths used by the application
// It will only set values that are not nil, otherwise
// the default values are used
func SetPaths(templates *string, static *string, dist *string) {
	if nil != templates {
		templatesDir = ConvertToAbsPath(templates)
		log.Printf("[INFO] using templates path %s\n", templatesDir)
	}
	if nil != static {
		staticDir = ConvertToAbsPath(static)
		log.Printf("[INFO] using static path %s\n", staticDir)
	}
	if nil != dist {
		distDir = ConvertToAbsPath(dist)
		log.Printf("[INFO] using dist path %s\n", distDir)
	}
}

// StaticIgnoreRegex returns a regex which contains the names of all templates
// which cannot be rendered on their own when building the static website
func StaticIgnoreRegex() *regexp.Regexp {
	return regexp.MustCompile(
		fmt.Sprintf("(%s|%s|%s|%s)",
			BaseTemplateName,
			ContentTemplateName,
			StatusTemplateName,
			ContactTemplateName,
		),
	)
}

// HTMLTemplatesPath returns the path where the html templates are located
func HTMLTemplatesPath() string {
	return filepath.Join(templatesDir, "html")
}

// BaseTeplatePath returns the complete path to the base html template
func BaseTemplatePath() string {
	return filepath.Join(HTMLTemplatesPath(), BaseTemplateName+".html")
}

// MailTemplatePath returns the complete path to the mail html template
func MailTemplatePath() string {
	return filepath.Join(templatesDir, "mail", MailTemplate)
}

// ContentTemplatesPath returns the path to the directory containg the content
// html templates
func ContentTemplatesPath() string {
	return filepath.Join(HTMLTemplatesPath(), "content")
}

// TemplatesPath returns the path to the directory containing all templates
func TemplatesPath() string {
	return templatesDir
}

// StaticContentPath returns the path to the directory containing static
// content, such as js, css, and images
func StaticContentPath() string {
	return staticDir
}

// DistDir returns the path to where the dist build should be output
func DistDir() string {
	return distDir
}
