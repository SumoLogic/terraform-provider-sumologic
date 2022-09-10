package sumologic

import (
	"encoding/json"
)

func (s *Client) GetCSEEntityNormalizationConfiguration() (*CSEEntityNormalizationConfiguration, error) {
	data, _, err := s.Get("sec/v1/entity-normalization/domain-configuration")
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSEEntityNormalizationConfigurationResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSEEntityNormalizationConfiguration, nil
}

func (s *Client) UpdateCSEEntityNormalizationConfiguration(CSEEntityNormalizationConfiguration CSEEntityNormalizationConfiguration) error {
	url := "sec/v1/entity-normalization/domain-configuration"

	request := CSEEntityNormalizationConfigurationRequest{
		CSEEntityNormalizationConfiguration: CSEEntityNormalizationConfiguration,
	}

	_, err := s.Put(url, request)

	return err
}

type CSEEntityNormalizationConfigurationRequest struct {
	CSEEntityNormalizationConfiguration CSEEntityNormalizationConfiguration `json:"fields"`
}

type CSEEntityNormalizationConfigurationResponse struct {
	CSEEntityNormalizationConfiguration CSEEntityNormalizationConfiguration `json:"data"`
}

type CSEEntityNormalizationConfiguration struct {
	WindowsNormalizationEnabled bool            `json:"adDomainNormalizationEnabled"`
	FqdnNormalizationEnabled    bool            `json:"fqdnNormalizationEnabled"`
	AwsNormalizationEnabled     bool            `json:"awsNormalizationEnabled"`
	DefaultNormalizedDomain     string          `json:"defaultNormalizedDomain,omitempty"`
	NormalizeHostnames          bool            `json:"normalizeHostnames"`
	NormalizeUsernames          bool            `json:"normalizeUsernames"`
	DomainMappings              []DomainMapping `json:"domainMappings"`
}

type DomainMapping struct {
	NormalizedDomain string `json:"normalizedDomain"`
	RawDomain        string `json:"rawDomain"`
}
