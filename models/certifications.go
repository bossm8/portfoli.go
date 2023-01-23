package models

type CertificationConfig struct {
	Certifications []*Certification `yaml:"certifications"`
}

// Make sure the interface is implemented
var _ listConfig = &CertificationConfig{}

func (cc *CertificationConfig) GetElements() []portfolioCard {
	return castToCard(cc.Certifications)
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
