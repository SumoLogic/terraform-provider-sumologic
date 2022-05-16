package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSEntityGroupConfiguration(id string) (*CSEEntityGroupConfiguration, error) {
	data, _, err := s.Get(fmt.Sprintf("sec/v1/entity-group-configurations/%s", id))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSEEntityGroupConfigurationResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSEEntityGroupConfiguration, nil
}

func (s *Client) DeleteCSEEntityGroupConfiguration(id string) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/entity-group-configurations/%s", id))

	return err
}

func (s *Client) CreateCSEEntityEntityGroupConfiguration(CSEEntityGroupConfiguration CSEEntityGroupConfiguration) (string, error) {
	CSEEntityGroupConfiguration.ConfigurationType = "ENTITY_VALUE"
	request := CSEEntityGroupConfigurationRequest{
		CSEEntityGroupConfiguration: CSEEntityGroupConfiguration,
	}

	var response CSEEntityGroupConfigurationResponse

	responseBody, err := s.Post("sec/v1/entity-group-configurations/entity_value", request)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return "", err
	}

	return response.CSEEntityGroupConfiguration.ID, nil
}

func (s *Client) UpdateCSEEntityEntityGroupConfiguration(CSEEntityGroupConfiguration CSEEntityGroupConfiguration) error {
	url := fmt.Sprintf("sec/v1/entity-group-configurations/entity_value/%s", CSEEntityGroupConfiguration.ID)

	CSEEntityGroupConfiguration.ID = ""
	CSEEntityGroupConfiguration.ConfigurationType = "ENTITY_VALUE"
	request := CSEEntityGroupConfigurationRequest{
		CSEEntityGroupConfiguration: CSEEntityGroupConfiguration,
	}

	_, err := s.Put(url, request)

	return err
}

func (s *Client) CreateCSEInventoryEntityGroupConfiguration(CSEEntityGroupConfiguration CSEEntityGroupConfiguration) (string, error) {
	CSEEntityGroupConfiguration.ConfigurationType = "INVENTORY"
	request := CSEEntityGroupConfigurationRequest{
		CSEEntityGroupConfiguration: CSEEntityGroupConfiguration,
	}

	var response CSEEntityGroupConfigurationResponse

	responseBody, err := s.Post("sec/v1/entity-group-configurations/inventory", request)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return "", err
	}

	return response.CSEEntityGroupConfiguration.ID, nil
}

func (s *Client) UpdateCSEInventoryEntityGroupConfiguration(CSEEntityGroupConfiguration CSEEntityGroupConfiguration) error {
	url := fmt.Sprintf("sec/v1/entity-group-configurations/inventory/%s", CSEEntityGroupConfiguration.ID)

	CSEEntityGroupConfiguration.ID = ""
	CSEEntityGroupConfiguration.ConfigurationType = "INVENTORY"
	request := CSEEntityGroupConfigurationRequest{
		CSEEntityGroupConfiguration: CSEEntityGroupConfiguration,
	}

	_, err := s.Put(url, request)

	return err
}

type CSEEntityGroupConfigurationRequest struct {
	CSEEntityGroupConfiguration CSEEntityGroupConfiguration `json:"fields"`
}

type CSEEntityGroupConfigurationResponse struct {
	CSEEntityGroupConfiguration CSEEntityGroupConfiguration `json:"data"`
}

type CSEEntityGroupConfiguration struct {
	ID                string   `json:"id,omitempty"`
	ConfigurationType string   `json:"configurationType,omitempty"`
	Created           string   `json:"created,omitempty"`
	CreatedBy         string   `json:"createdBy,omitempty"`
	Criticality       string   `json:"criticality,omitempty"`
	Deleted           bool     `json:"deleted,omitempty"`
	Description       string   `json:"description,omitempty"`
	EntityNamespace   string   `json:"entityNamespace,omitempty"`
	EntityType        string   `json:"entityType,omitempty"`
	Groups            []string `json:"groups,omitempty"`
	InventoryType     string   `json:"inventoryType,omitempty"`
	LastUpdated       string   `json:"lastUpdated,omitempty,omitempty"`
	LastUpdatedBy     string   `json:"lastUpdatedBy,omitempty,omitempty"`
	Name              string   `json:"name,omitempty"`
	NetworkBlock      string   `json:"networkBlock,omitempty"`
	Prefix            string   `json:"prefix,omitempty"`
	Suffix            string   `json:"suffix,omitempty"`
	Suppressed        bool     `json:"suppressed,omitempty"`
	Tags              []string `json:"tags,omitempty"`
}
