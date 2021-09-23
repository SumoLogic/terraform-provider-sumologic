package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSEEntityCriticalityConfig(id string) (*CSEEntityCriticalityConfig, error) {
	data, _, err := s.Get(fmt.Sprintf("sec/v1/entity-criticality-configs/%s", id))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSEEntityCriticalityConfigResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSEEntityCriticalityConfig, nil
}

func (s *Client) DeleteCSEEntityCriticalityConfig(id string) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/entity-criticality-configs/%s", id))

	return err
}

func (s *Client) CreateCSEEntityCriticalityConfig(CSEEntityCriticalityConfig CSEEntityCriticalityConfig) (string, error) {

	request := CSEEntityCriticalityConfigRequestPost{
		CSEEntityCriticalityConfig: CSEEntityCriticalityConfig,
	}

	var response CSEEntityCriticalityConfigResponse

	responseBody, err := s.Post("sec/v1/entity-criticality-configs", request)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return "", err
	}

	return response.CSEEntityCriticalityConfig.ID, nil
}

func (s *Client) UpdateCSEEntityCriticalityConfig(CSEEntityCriticalityConfig CSEEntityCriticalityConfig) error {
	url := fmt.Sprintf("sec/v1/entity-criticality-configs/%s", CSEEntityCriticalityConfig.ID)

	request := CSEEntityCriticalityConfigRequestUpdate{
		CSEEntityCriticalityConfigUpdate{
			SeverityExpression: CSEEntityCriticalityConfig.SeverityExpression,
		},
	}

	_, err := s.Put(url, request)

	return err
}

type CSEEntityCriticalityConfigRequestPost struct {
	CSEEntityCriticalityConfig CSEEntityCriticalityConfig `json:"fields"`
}

type CSEEntityCriticalityConfigRequestUpdate struct {
	CSEEntityCriticalityConfigUpdate CSEEntityCriticalityConfigUpdate `json:"fields"`
}

type CSEEntityCriticalityConfigResponse struct {
	CSEEntityCriticalityConfig CSEEntityCriticalityConfig `json:"data"`
}

type CSEEntityCriticalityConfig struct {
	ID                 string `json:"id,omitempty"`
	Name               string `json:"name"`
	SeverityExpression string `json:"severityExpression"`
}

type CSEEntityCriticalityConfigUpdate struct {
	SeverityExpression string `json:"severityExpression"`
}
