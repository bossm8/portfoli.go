package content

type CertificationConfig struct {
	Certifications []*Certification `yaml:"certifications"`
}

// Make sure the interface is implemented
var _ ContentConfig = &CertificationConfig{}

func (cc *CertificationConfig) GetElements() []Content {
	return castToContent(cc.Certifications)
}

func (cc *CertificationConfig) GetConfigName() string {
	return cc.GetContentKind() + ".yml"
}

func (cc *CertificationConfig) GetContentKind() string {
	return ContentTypes[typeCertification]
}

type Certification struct {
	ContentBase      `yaml:",inline"`
	ContentDateRange `yaml:",inline"`
}

// Make sure the interface is implemented
var _ Content = &Certification{}

func (c *Certification) GetTemplateName() string {
	return "certification.html"
}
