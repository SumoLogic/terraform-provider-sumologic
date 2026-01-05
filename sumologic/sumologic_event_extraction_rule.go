package sumologic

import (
	"encoding/json"
	"fmt"
)

/*
========================
Models
========================
*/

type EventExtractionRule struct {
	ID                    string                  `json:"id,omitempty"`
	Name                  string                  `json:"name"`
	Description           string                  `json:"description,omitempty"`
	Query                 string                  `json:"query"`
	CorrelationExpression *CorrelationExpression  `json:"correlationExpression,omitempty"`
	Configuration         map[string]FieldMapping `json:"configuration"`
	Enabled               bool                    `json:"enabled,omitempty"`
	DisableReason         string                  `json:"disableReason,omitempty"`
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

type EventExtractionRuleListResponse struct {
	Data []EventExtractionRule `json:"data"`
}

/*
========================
API Calls
========================
*/

func (s *Client) CreateEventExtractionRule(rule EventExtractionRule) (*EventExtractionRule, error) {
	resp, err := s.Post("v1/eventExtractionRules", rule)
	if err != nil {
		return nil, err
	}

	var created EventExtractionRule
	if err := json.Unmarshal(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

func (s *Client) GetEventExtractionRule(id string) (*EventExtractionRule, error) {
	resp, err := s.Get(fmt.Sprintf("v1/eventExtractionRules/%s", id))
	if err != nil || resp == nil {
		return nil, err
	}

	var rule EventExtractionRule
	if err := json.Unmarshal(resp, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

func (s *Client) GetAllEventExtractionRules() ([]EventExtractionRule, error) {
	resp, err := s.Get("v1/eventExtractionRules")
	if err != nil {
		return nil, err
	}

	var list EventExtractionRuleListResponse
	if err := json.Unmarshal(resp, &list); err != nil {
		return nil, err
	}
	return list.Data, nil
}

func (s *Client) UpdateEventExtractionRule(id string, rule EventExtractionRule) error {
	_, err := s.Put(fmt.Sprintf("v1/eventExtractionRules/%s", id), rule)
	return err
}

func (s *Client) DeleteEventExtractionRule(id string) error {
	_, err := s.Delete(fmt.Sprintf("v1/eventExtractionRules/%s", id))
	return err
}
