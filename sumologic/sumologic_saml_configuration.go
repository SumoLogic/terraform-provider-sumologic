package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetSamlConfiguration(id string) (*SamlConfiguration, error) {
	// We don't have a get SAML configuration by id endpoint, only a list endpoint
	url := "v1/saml/identityProviders"

	data, _, err := s.Get(url, false)
	if err != nil || data == nil {
		return nil, err
	}

	var samlConfigurations []SamlConfiguration
	err = json.Unmarshal(data, &samlConfigurations)
	if err != nil {
		return nil, err
	}

	for _, samlConfiguration := range samlConfigurations {
		if samlConfiguration.ID == id {
			return &samlConfiguration, nil
		}
	}

	return nil, nil
}

func (s *Client) CreateSamlConfiguration(samlConfiguration SamlConfiguration) (*SamlConfiguration, error) {
	url := "v1/saml/identityProviders"

	data, err := s.Post(url, samlConfiguration, false)
	if err != nil {
		return nil, err
	}

	var createdSamlConfiguration SamlConfiguration
	err = json.Unmarshal(data, &createdSamlConfiguration)
	if err != nil {
		return nil, err
	}

	return &createdSamlConfiguration, nil
}

func (s *Client) DeleteSamlConfiguration(id string) error {
	url := fmt.Sprintf("v1/saml/identityProviders/%s", id)

	_, err := s.Delete(url)
	return err
}

func (s *Client) UpdateSamlConfiguration(id string, samlConfiguration SamlConfiguration) error {
	url := fmt.Sprintf("v1/saml/identityProviders/%s", id)

	_, err := s.Put(url, samlConfiguration, false)
	return err
}

type OnDemandProvisioningEnabled struct {
	FirstNameAttribute        string   `json:"firstNameAttribute"`
	LastNameAttribute         string   `json:"lastNameAttribute"`
	OnDemandProvisioningRoles []string `json:"onDemandProvisioningRoles"`
}

type SamlConfiguration struct {
	SpInitiatedLoginPath         string                       `json:"spInitiatedLoginPath"`
	ConfigurationName            string                       `json:"configurationName"`
	Issuer                       string                       `json:"issuer"`
	SpInitiatedLoginEnabled      bool                         `json:"spInitiatedLoginEnabled"`
	AuthnRequestUrl              string                       `json:"authnRequestUrl"`
	X509cert1                    string                       `json:"x509cert1"`
	X509cert2                    string                       `json:"x509cert2"`
	X509cert3                    string                       `json:"x509cert3"`
	OnDemandProvisioningEnabled  *OnDemandProvisioningEnabled `json:"onDemandProvisioningEnabled,omitempty"`
	RolesAttribute               string                       `json:"rolesAttribute"`
	LogoutEnabled                bool                         `json:"logoutEnabled"`
	LogoutUrl                    string                       `json:"logoutUrl"`
	EmailAttribute               string                       `json:"emailAttribute"`
	DebugMode                    bool                         `json:"debugMode"`
	SignAuthnRequest             bool                         `json:"signAuthnRequest"`
	DisableRequestedAuthnContext bool                         `json:"disableRequestedAuthnContext"`
	IsRedirectBinding            bool                         `json:"isRedirectBinding"`
	Certificate                  string                       `json:"certificate,omitempty"`
	ID                           string                       `json:"id,omitempty"`
	AssertionConsumerUrl         string                       `json:"assertionConsumerUrl,omitempty"`
}
