package models

import "html/template"

const certificationConfigName = "certifications.yml"

type CertificationConfig struct {
	Certifications []*PortfolioCard `yaml:"certifications"`
}

// Make sure the interface is implemented
var _ listConfig = &CertificationConfig{}

func (cc *CertificationConfig) GetElements() []*PortfolioCard {
	return cc.Certifications
}

func (cc *CertificationConfig) GetConfigName() string {
	return certificationConfigName
}

type Certification struct {
	Base           `yaml:",inline"`
	Specialization string `yaml:"specialization"`
}

// Make sure the interface is implemented
var _ PortfolioCard = &Certification{}

func GetCertifications() []*PortfolioCard {
	return unmarshal(&CertificationConfig{})
}

func (c *Certification) GetHTMLTemplate() template.HTML {
	return ""
}
