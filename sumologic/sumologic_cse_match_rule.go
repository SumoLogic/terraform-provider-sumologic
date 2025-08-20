package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSEMatchRule(id string) (*CSEMatchRule, error) {
	data, err := s.Get(fmt.Sprintf("sec/v1/rules/%s", id))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSEMatchRuleResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSEMatchRule, nil
}

func (s *Client) DeleteCSEMatchRule(id string) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/rules/%s", id))

	return err
}

func (s *Client) CreateCSEMatchRule(CSEMatchRule CSEMatchRule) (string, error) {

	request := CSEMatchRuleRequest{
		CSEMatchRule: CSEMatchRule,
	}

	var response CSEMatchRuleResponse

	responseBody, err := s.Post("sec/v1/rules/templated", request)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return "", err
	}

	return response.CSEMatchRule.ID, nil
}

func (s *Client) UpdateCSEMatchRule(CSEMatchRule CSEMatchRule) error {
	url := fmt.Sprintf("sec/v1/rules/templated/%s", CSEMatchRule.ID)

	CSEMatchRule.ID = ""
	request := CSEMatchRuleRequest{
		CSEMatchRule: CSEMatchRule,
	}

	_, err := s.Put(url, request)

	return err
}

func (s *Client) OverrideCSEMatchRule(CSEMatchRule CSEMatchRule) error {
	url := fmt.Sprintf("sec/v1/rules/templated/%s/override", CSEMatchRule.ID)

	CSEMatchRuleOverride := toOverride(CSEMatchRule)

	request := CSEMatchRuleOverrideRequest{
		CSEMatchRuleOverride: CSEMatchRuleOverride,
	}
	_, err := s.Put(url, request)

	return err
}

func toOverride(CSEMatchRule CSEMatchRule) CSEMatchRuleOverride {

	return CSEMatchRuleOverride{
		DescriptionExpression: CSEMatchRule.DescriptionExpression,
		EntitySelectors:       CSEMatchRule.EntitySelectors,
		IsPrototype:           CSEMatchRule.IsPrototype,
		Name:                  CSEMatchRule.Name,
		NameExpression:        CSEMatchRule.NameExpression,
		SeverityMapping:       CSEMatchRule.SeverityMapping,
		SummaryExpression:     CSEMatchRule.SummaryExpression,
		Tags:                  CSEMatchRule.Tags,
		SuppressionWindowSize: CSEMatchRule.SuppressionWindowSize,
	}
}

type CSEMatchRuleRequest struct {
	CSEMatchRule CSEMatchRule `json:"fields"`
}

type CSEMatchRuleOverrideRequest struct {
	CSEMatchRuleOverride CSEMatchRuleOverride `json:"fields"`
}

type CSEMatchRuleResponse struct {
	CSEMatchRule CSEMatchRule `json:"data"`
}

type SeverityMappingValueMapping struct {
	Type string `json:"type"`
	From string `json:"from"`
	To   int    `json:"to"`
}

type SeverityMapping struct {
	Type    string                        `json:"type"`
	Default int                           `json:"default"`
	Field   string                        `json:"field"`
	Mapping []SeverityMappingValueMapping `json:"mapping"`
}

type CSEMatchRule struct {
	ID                    string           `json:"id,omitempty"`
	DescriptionExpression string           `json:"descriptionExpression"`
	Enabled               bool             `json:"enabled"`
	EntitySelectors       []EntitySelector `json:"entitySelectors"`
	Expression            string           `json:"expression"`
	IsPrototype           bool             `json:"isPrototype"`
	Name                  string           `json:"name"`
	NameExpression        string           `json:"nameExpression"`
	SeverityMapping       SeverityMapping  `json:"scoreMapping"`
	Stream                string           `json:"stream"`
	SummaryExpression     string           `json:"summaryExpression"`
	Tags                  []string         `json:"tags"`
	SuppressionWindowSize *int             `json:"suppressionWindowSize,omitempty"`
}

type CSEMatchRuleOverride struct {
	DescriptionExpression string           `json:"descriptionExpression"`
	EntitySelectors       []EntitySelector `json:"entitySelectors"`
	IsPrototype           bool             `json:"isPrototype"`
	Name                  string           `json:"name"`
	NameExpression        string           `json:"nameExpression"`
	SeverityMapping       SeverityMapping  `json:"scoreMapping"`
	SummaryExpression     string           `json:"summaryExpression"`
	Tags                  []string         `json:"tags"`
	SuppressionWindowSize *int             `json:"suppressionWindowSize,omitempty"`
}
