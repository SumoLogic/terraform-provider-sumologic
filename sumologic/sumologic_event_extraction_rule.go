package sumologic

import (
	"encoding/json"
	"fmt"
)

type EventExtractionRule struct {
	ID                    string                  `json:"id,omitempty"`
	Name                  string                  `json:"name"`
	Description           string                  `json:"description,omitempty"`
	Query                 string                  `json:"query"`
	CorrelationExpression *CorrelationExpression  `json:"correlationExpression,omitempty"`
	Configuration         map[string]FieldMapping `json:"configuration"`
	Enabled               bool                    `json:"enabled"`
	DisableReason         string                  `json:"disableReason,omitempty"`
	CreatedAt             string                  `json:"createdAt,omitempty"`
	CreatedBy             string                  `json:"createdBy,omitempty"`
	ModifiedAt            string                  `json:"modifiedAt,omitempty"`
	ModifiedBy            string                  `json:"modifiedBy,omitempty"`
}

type CorrelationExpression struct {
	QueryFieldName          string `json:"queryFieldName"`
	EventFieldName          string `json:"eventFieldName"`
	StringMatchingAlgorithm string `json:"stringMatchingAlgorithm"`
}

type FieldMapping struct {
	ValueSource string `json:"valueSource"`
	MappingType string `json:"mappingType,omitempty"`
}

type ListEventExtractionRulesResponse struct {
	Data []EventExtractionRule `json:"data"`
}

func (s *Client) GetEventExtractionRule(id string) (*EventExtractionRule, error) {
	urlPath := fmt.Sprintf("v1/eventExtractionRules/%s", id)
	data, err := s.Get(urlPath)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var rule EventExtractionRule
	err = json.Unmarshal(data, &rule)
	if err != nil {
		return nil, err
	}

	return &rule, nil
}

func (s *Client) GetEventExtractionRuleByName(name string) (*EventExtractionRule, error) {
	data, err := s.Get("v1/eventExtractionRules")
	if err != nil {
		return nil, err
	}

	var response ListEventExtractionRulesResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	for _, rule := range response.Data {
		if rule.Name == name {
			return &rule, nil
		}
	}

	return nil, fmt.Errorf("event extraction rule with name '%s' not found", name)
}

func (s *Client) CreateEventExtractionRule(rule EventExtractionRule) (string, error) {
	urlPath := "v1/eventExtractionRules"
	data, err := s.Post(urlPath, rule)
	if err != nil {
		return "", err
	}

	var createdRule EventExtractionRule
	err = json.Unmarshal(data, &createdRule)
	if err != nil {
		return "", err
	}

	return createdRule.ID, nil
}

func (s *Client) UpdateEventExtractionRule(rule EventExtractionRule) error {
	urlPath := fmt.Sprintf("v1/eventExtractionRules/%s", rule.ID)

	idBackup := rule.ID
	rule.ID = ""

	_, err := s.Put(urlPath, rule)

	rule.ID = idBackup
	return err
}

func (s *Client) DeleteEventExtractionRule(id string) error {
	urlPath := fmt.Sprintf("v1/eventExtractionRules/%s", id)
	_, err := s.Delete(urlPath)
	return err
}
