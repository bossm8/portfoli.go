package content

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"bossm8.ch/portfolio/models/utils"
)

const (
	typeExperience = iota
	typeEducation
	typeProject
	typeCertification
)

var (
	ContentTypes    = []string{"experience", "education", "projects", "certifications"}
	contentMappings = map[string]ContentConfig{
		ContentTypes[typeExperience]:    &ExperienceConfig{},
		ContentTypes[typeEducation]:     &EducationConfig{},
		ContentTypes[typeProject]:       &ProjectConfig{},
		ContentTypes[typeCertification]: &CertificationConfig{},
	}
	rex = regexp.MustCompile(fmt.Sprintf("(%s)", strings.Join(ContentTypes, "|")))
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

var (
	templatesDir = filepath.Join("models", "content", "templates", "html")
)

// GetRenderedContent reads the content kind passed from its yaml configuration
// and returns all configured elements as html to be placed in the main
// template directly
func GetRenderedContent(contentType string) ([]template.HTML, error) {
	// Get the correct object to load
	// TODO validate so we do not have null values
	obj := contentMappings[contentType]

	err := loadContentConfig(obj)
	if nil != err {
		log.Printf("[ERROR] Generating content failed: %s\n", err)
		return nil, err
	}

	// render the content read from yaml into the html models
	cards := obj.GetElements()
	data := make([]template.HTML, 0)
	for _, crd := range cards {
		if tpl, err := renderContent(crd); nil != err {
			return nil, err
		} else {
			data = append(data, tpl)
		}

	}
	return data, err
}

func GetRoutingRegexString() string {
	return rex.String()
}

func renderContent(content Content) (template.HTML, error) {
	contentBaseTpl := filepath.Join(templatesDir, "base.html")
	htmlTpl := filepath.Join(templatesDir, content.GetTemplateName())
	tpl, err := template.ParseFiles(contentBaseTpl, htmlTpl)
	if nil != err {
		log.Printf("[ERROR] Failed to parse template '%s': %s\n", htmlTpl, err)
		return "", err
	}
	rendered := bytes.Buffer{}
	if err := tpl.ExecuteTemplate(&rendered, "content", content); nil != err {
		log.Printf("[ERROR] Failed to process template %s with error %s\n", tpl.Name(), err)
		return "", err
	}

	return template.HTML(rendered.String()), nil
}

func castToContent[T Content](content []T) []Content {
	casted := make([]Content, len(content))
	for idx, crd := range content {
		casted[idx] = crd
	}
	return casted
}

func unmarshalContentConfig(content ContentConfig) error {
	return utils.LoadFromYAMLFile(content.GetConfigName(), content)
}

func loadContentConfig(content ContentConfig) error {
	err := unmarshalContentConfig(content)
	return err
}

func IsValidContentType(contentType string) bool {
	isValid := true
	// Check if all content kinds specified in the yaml config are valid
	if !rex.MatchString(contentType) {
		log.Printf("[ERROR] Invalid content kind '%s', allowed values are: %s\n", contentType, ContentTypes)
		isValid = false
	}
	return isValid
}
