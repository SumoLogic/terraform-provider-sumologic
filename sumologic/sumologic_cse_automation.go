package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSEAutomation(id string) (*CSEAutomation, error) {
	data, _, err := s.Get(fmt.Sprintf("sec/v1/automations/%s", id))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSEAutomationResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSEAutomation, nil
}

func (s *Client) DeleteCSEAutomation(id string) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/automations/%s", id))

	return err
}

func (s *Client) CreateCSEAutomation(CSEAutomation CSEAutomation) (string, error) {

	request := CSEAutomationRequestPost{
		CSEAutomation: CSEAutomation,
	}

	var response CSEAutomationResponse

	responseBody, err := s.Post("sec/v1/automations", request)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return "", err
	}

	return response.CSEAutomation.ID, nil
}

func (s *Client) UpdateCSEAutomation(CSEAutomation CSEAutomation) error {
	url := fmt.Sprintf("sec/v1/automations/%s", CSEAutomation.ID)

	request := CSEAutomationRequestUpdate{
		CSEAutomationUpdate{
			CseResourceSubTypes: CSEAutomation.CseResourceSubTypes,
			ExecutionTypes:      CSEAutomation.ExecutionTypes,
			Enabled:             CSEAutomation.Enabled,
		},
	}

	_, err := s.Put(url, request)

	return err
}

type CSEAutomationRequestPost struct {
	CSEAutomation CSEAutomation `json:"fields"`
}

type CSEAutomationRequestUpdate struct {
	CSEAutomationUpdate CSEAutomationUpdate `json:"fields"`
}

type CSEAutomationResponse struct {
	CSEAutomation CSEAutomation `json:"data"`
}

type CSEAutomation struct {
	ID                  string   `json:"id,omitempty"`
	PlaybookId          string   `json:"playbookId"`
	CseResourceType     string   `json:"cseResourceType"`
	CseResourceSubTypes []string `json:"cseResourceSubTypes,omitempty"`
	Name                string   `json:"name,omitempty"`
	Description         string   `json:"description,omitempty"`
	ExecutionTypes      []string `json:"executionTypes"`
	Enabled             bool     `json:"enabled"`
}

type CSEAutomationUpdate struct {
	CseResourceSubTypes []string `json:"cseResourceSubTypes"`
	ExecutionTypes      []string `json:"executionTypes"`
	Enabled             bool     `json:"enabled"`
}
