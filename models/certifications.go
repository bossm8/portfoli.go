package models

import "html/template"

type CertificationConfig struct {
	Certifications []*Certification `yaml:"certifications"`
}

// Make sure the interface is implemented
var _ listConfig = &CertificationConfig{}

func (cc *CertificationConfig) GetRenderedElements() ([]template.HTML, error) {
	data := make([]template.HTML, len(cc.Certifications))
	for idx, crt := range cc.Certifications {
		rendered, err := renderCard(crt)
		if nil != err {
			return nil, err
		}
		data[idx] = rendered
	}
	return data, nil
}

func (cc *CertificationConfig) GetConfigName() string {
	return cc.GetContentKind() + ".yml"
}

func (cc *CertificationConfig) Load() error {
	return unmarshal(cc)
}

func (cc *CertificationConfig) GetContentKind() string {
	return kinds[kindCert]
}

type Certification struct {
	Base           `yaml:",inline"`
	Specialization string `yaml:"specialization"`
}

// Make sure the interface is implemented
var _ portfolioCard = &Certification{}

func (c *Certification) GetTemplateName() string {
	return "certification.html"
}

func LoadCertifications() (listConfig, error) {
	cert := &CertificationConfig{}
	err := cert.Load()
	return cert, err
}
