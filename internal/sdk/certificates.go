package sdk

import (
	"fmt"
	"time"
)

type CertificatesClient struct {
	maintClient *MaintenanceClient
}

func NewCertificatesClient(cfg *Config) *CertificatesClient {
	return &CertificatesClient{
		maintClient: NewMaintenanceClient(cfg),
	}
}

func (c *CertificatesClient) Authenticate(username, password string) (string, error) {
	return c.maintClient.Authenticate(username, password)
}

func (c *CertificatesClient) SetSID(sid string) {
	c.maintClient.sid = sid
}

func (c *CertificatesClient) IsAuthenticated() bool {
	return c.maintClient.IsAuthenticated()
}

type Certificate struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Domain      string    `json:"domain"`
	Issuer      string    `json:"issuer"`
	Subject     string    `json:"subject"`
	ValidFrom   time.Time `json:"valid_from"`
	ValidTo     time.Time `json:"valid_to"`
	Fingerprint string    `json:"fingerprint"`
	IsDefault   bool      `json:"is_default"`
	IsValid     bool      `json:"is_valid"`
}

func (c *CertificatesClient) ListCertificates() ([]Certificate, error) {
	if !c.maintClient.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	result, err := c.maintClient.Call("getcertificates", nil)
	if err != nil {
		return nil, err
	}

	certsData, ok := result["certificate"]
	if !ok {
		return []Certificate{}, nil
	}

	switch v := certsData.(type) {
	case []interface{}:
		certs := make([]Certificate, 0, len(v))
		for _, c := range v {
			if certMap, ok := c.(map[string]interface{}); ok {
				cert := Certificate{
					ID:          getCertString(certMap, "id"),
					Name:        getCertString(certMap, "name"),
					Domain:      getCertString(certMap, "domain"),
					Issuer:      getCertString(certMap, "issuer"),
					Subject:     getCertString(certMap, "subject"),
					Fingerprint: getCertString(certMap, "fingerprint"),
					IsDefault:   getCertString(certMap, "is_default") == "1",
					IsValid:     getCertString(certMap, "is_valid") == "1",
				}
				certs = append(certs, cert)
			}
		}
		return certs, nil
	}

	return []Certificate{}, nil
}

func (c *CertificatesClient) GetCertificateInfo(certID string) (*Certificate, error) {
	if !c.maintClient.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_certificate_id": certID,
	}

	result, err := c.maintClient.Call("getcertificateinfo", params)
	if err != nil {
		return nil, err
	}

	cert := &Certificate{
		ID:          getCertString(result, "id"),
		Name:        getCertString(result, "name"),
		Domain:      getCertString(result, "domain"),
		Issuer:      getCertString(result, "issuer"),
		Subject:     getCertString(result, "subject"),
		Fingerprint: getCertString(result, "fingerprint"),
		IsDefault:   getCertString(result, "is_default") == "1",
		IsValid:     getCertString(result, "is_valid") == "1",
	}

	return cert, nil
}

func (c *CertificatesClient) AddCertificate(name, domain, certData, keyData, password string) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_name":     name,
		"s_domain":   domain,
		"s_cert":     certData,
		"s_key":      keyData,
		"s_password": password,
	}

	_, err := c.maintClient.Call("addcertificate", params)
	return err
}

func (c *CertificatesClient) DeleteCertificate(certID string) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_certificate_id": certID,
	}

	_, err := c.maintClient.Call("deletecertificate", params)
	return err
}

func (c *CertificatesClient) SetDefaultCertificate(certID string) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_certificate_id": certID,
	}

	_, err := c.maintClient.Call("setdefaultcertificate", params)
	return err
}

func (c *CertificatesClient) CreateCSR(domain, commonName, org, orgUnit, city, state, country string) (string, error) {
	if !c.maintClient.IsAuthenticated() {
		return "", fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_domain":   domain,
		"s_cn":       commonName,
		"s_org":      org,
		"s_org_unit": orgUnit,
		"s_city":     city,
		"s_state":    state,
		"s_country":  country,
	}

	result, err := c.maintClient.Call("createcsr", params)
	if err != nil {
		return "", err
	}

	return getCertString(result, "csr"), nil
}

func (c *CertificatesClient) ImportCertificate(certData string) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_certificate": certData,
	}

	_, err := c.maintClient.Call("importcertificate", params)
	return err
}

func (c *CertificatesClient) ExportCertificate(certID string) (string, error) {
	if !c.maintClient.IsAuthenticated() {
		return "", fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_certificate_id": certID,
	}

	result, err := c.maintClient.Call("exportcertificate", params)
	if err != nil {
		return "", err
	}

	return getCertString(result, "certificate"), nil
}

func getCertString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
