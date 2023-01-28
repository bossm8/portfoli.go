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
	AboutMeTemplateName = "me"
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
