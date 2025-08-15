package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSEThresholdRule(id string) (*CSEThresholdRule, error) {
	data, err := s.Get(fmt.Sprintf("sec/v1/rules/%s", id))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSEThresholdRuleResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSEThresholdRule, nil
}

func (s *Client) DeleteCSEThresholdRule(id string) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/rules/%s", id))

	return err
}

func (s *Client) CreateCSEThresholdRule(CSEThresholdRule CSEThresholdRule) (string, error) {
	request := CSEThresholdRuleRequest{
		CSEThresholdRule: CSEThresholdRule,
	}

	var response CSEThresholdRuleResponse

	responseBody, err := s.Post("sec/v1/rules/threshold", request)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return "", err
	}

	return response.CSEThresholdRule.ID, nil
}

func (s *Client) UpdateCSEThresholdRule(CSEThresholdRule CSEThresholdRule) error {
	url := fmt.Sprintf("sec/v1/rules/threshold/%s", CSEThresholdRule.ID)

	CSEThresholdRule.ID = ""
	request := CSEThresholdRuleRequest{
		CSEThresholdRule: CSEThresholdRule,
	}

	_, err := s.Put(url, request)

	return err
}

func (s *Client) OverrideCSEThresholdRule(CSEThresholdRule CSEThresholdRule) error {
	url := fmt.Sprintf("sec/v1/rules/threshold/%s/override", CSEThresholdRule.ID)

	CSEThresholdRuleOverride := toOverrideThreshold(CSEThresholdRule)

	request := CSEThresholdRuleOverrideRequest{
		CSEThresholdRuleOverride: CSEThresholdRuleOverride,
	}
	_, err := s.Put(url, request)

	return err
}

func toOverrideThreshold(CSEThresholdRule CSEThresholdRule) CSEThresholdRuleOverride {

	return CSEThresholdRuleOverride{
		Description:            CSEThresholdRule.Description,
		EntitySelectors:        CSEThresholdRule.EntitySelectors,
		GroupByFields:          CSEThresholdRule.GroupByFields,
		IsPrototype:            CSEThresholdRule.IsPrototype,
		Limit:                  CSEThresholdRule.Limit,
		Name:                   CSEThresholdRule.Name,
		Severity:               CSEThresholdRule.Severity,
		SummaryExpression:      CSEThresholdRule.SummaryExpression,
		Tags:                   CSEThresholdRule.Tags,
		WindowSize:             CSEThresholdRule.WindowSize,
		WindowSizeMilliseconds: CSEThresholdRule.WindowSizeMilliseconds,
		SuppressionWindowSize:  CSEThresholdRule.SuppressionWindowSize,
	}
}

type CSEThresholdRuleRequest struct {
	CSEThresholdRule CSEThresholdRule `json:"fields"`
}

type CSEThresholdRuleResponse struct {
	CSEThresholdRule CSEThresholdRule `json:"data"`
}

type CSEThresholdRuleOverrideRequest struct {
	CSEThresholdRuleOverride CSEThresholdRuleOverride `json:"fields"`
}

type CSEThresholdRule struct {
	ID                     string           `json:"id,omitempty"`
	CountDistinct          bool             `json:"countDistinct"`
	CountField             string           `json:"countField"`
	Description            string           `json:"description"`
	Enabled                bool             `json:"enabled"`
	EntitySelectors        []EntitySelector `json:"entitySelectors"`
	Expression             string           `json:"expression"`
	GroupByFields          []string         `json:"groupByFields"`
	IsPrototype            bool             `json:"isPrototype"`
	Limit                  int              `json:"limit"`
	Name                   string           `json:"name"`
	Severity               int              `json:"score"`
	Stream                 string           `json:"stream"`
	SummaryExpression      string           `json:"summaryExpression"`
	Tags                   []string         `json:"tags"`
	Version                int              `json:"version"`
	WindowSize             windowSizeField  `json:"windowSize,omitempty"`
	WindowSizeName         string           `json:"windowSizeName,omitempty"`
	WindowSizeMilliseconds string           `json:"windowSizeMilliseconds,omitempty"`
	SuppressionWindowSize  *int             `json:"suppressionWindowSize,omitempty"`
}

type CSEThresholdRuleOverride struct {
	Description            string           `json:"description"`
	EntitySelectors        []EntitySelector `json:"entitySelectors"`
	GroupByFields          []string         `json:"groupByFields"`
	IsPrototype            bool             `json:"isPrototype"`
	Limit                  int              `json:"limit"`
	Name                   string           `json:"name"`
	Severity               int              `json:"score"`
	SummaryExpression      string           `json:"summaryExpression"`
	Tags                   []string         `json:"tags"`
	WindowSize             windowSizeField  `json:"windowSize,omitempty"`
	WindowSizeMilliseconds string           `json:"windowSizeMilliseconds,omitempty"`
	SuppressionWindowSize  *int             `json:"suppressionWindowSize,omitempty"`
}
