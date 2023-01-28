package config

import (
	"fmt"
	"log"
	"os"
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
	if filepath.IsAbs(*path) {
		return *path
	}
	exe, err := os.Executable()
	if err != nil {
		log.Fatalf("[ERROR] parsing configuration directory: %s\n", err)
	}
	return filepath.Join(
		filepath.Dir(exe),
		*path,
	)
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

func HTMLTeplatesPath() string {
	return filepath.Join(templatesDir, "html")
}

func BaseTemplatePath() string {
	return filepath.Join(HTMLTeplatesPath(), BaseTemplateName+".html")
}

func MailTemplatePath() string {
	return filepath.Join(templatesDir, "mail", MailTemplate)
}

func ContentTemplatesPath() string {
	return filepath.Join(HTMLTeplatesPath(), "content")
}

func StaticContentPath() string {
	return staticDir
}

func DistDir() string {
	return distDir
}
