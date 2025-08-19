package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSEFirstSeenRule(id string) (*CSEFirstSeenRule, error) {
	data, err := s.Get(fmt.Sprintf("sec/v1/rules/%s", id))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSEFirstSeenRuleResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSEFirstSeenRule, nil
}

func (s *Client) DeleteCSEFirstSeenRule(id string) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/rules/%s", id))

	return err
}

func (s *Client) CreateCSEFirstSeenRule(CSEFirstSeenRule CSEFirstSeenRule) (string, error) {
	request := CSEFirstSeenRuleRequest{
		CSEFirstSeenRule: CSEFirstSeenRule,
	}

	var response CSEFirstSeenRuleResponse

	responseBody, err := s.Post("sec/v1/rules/first-seen", request)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return "", err
	}

	return response.CSEFirstSeenRule.ID, nil
}

func (s *Client) UpdateCSEFirstSeenRule(CSEFirstSeenRule CSEFirstSeenRule) error {
	url := fmt.Sprintf("sec/v1/rules/first-seen/%s", CSEFirstSeenRule.ID)

	CSEFirstSeenRule.ID = ""
	request := CSEFirstSeenRuleRequest{
		CSEFirstSeenRule: CSEFirstSeenRule,
	}

	_, err := s.Put(url, request)

	return err
}

func (s *Client) OverrideCSEFirstSeenRule(CSEFirstSeenRule CSEFirstSeenRule) error {
	url := fmt.Sprintf("sec/v1/rules/first-seen/%s/override", CSEFirstSeenRule.ID)

	CSEFirstSeenRuleOverride := toOverrideFirstSeen(CSEFirstSeenRule)
	request := CSEFirstSeenRuleOverrideRequest{
		CSEFirstSeenRuleOverride: CSEFirstSeenRuleOverride,
	}

	_, err := s.Put(url, request)

	return err
}

func toOverrideFirstSeen(CSEFirstSeenRule CSEFirstSeenRule) CSEFirstSeenRuleOverride {
	return CSEFirstSeenRuleOverride{
		BaselineWindowSize:    CSEFirstSeenRule.BaselineWindowSize,
		DescriptionExpression: CSEFirstSeenRule.DescriptionExpression,
		GroupByFields:         CSEFirstSeenRule.GroupByFields,
		IsPrototype:           CSEFirstSeenRule.IsPrototype,
		Name:                  CSEFirstSeenRule.Name,
		NameExpression:        CSEFirstSeenRule.NameExpression,
		RetentionWindowSize:   CSEFirstSeenRule.RetentionWindowSize,
		Severity:              CSEFirstSeenRule.Severity,
		SummaryExpression:     CSEFirstSeenRule.SummaryExpression,
		Tags:                  CSEFirstSeenRule.Tags,
		SuppressionWindowSize: CSEFirstSeenRule.SuppressionWindowSize,
	}
}

type CSEFirstSeenRuleRequest struct {
	CSEFirstSeenRule CSEFirstSeenRule `json:"fields"`
}

type CSEFirstSeenRuleResponse struct {
	CSEFirstSeenRule CSEFirstSeenRule `json:"data"`
}

type CSEFirstSeenRuleOverrideRequest struct {
	CSEFirstSeenRuleOverride CSEFirstSeenRuleOverride `json:"fields"`
}

type CSEFirstSeenRule struct {
	ID                    string           `json:"id,omitempty"`
	AssetField            string           `json:"assetField"`
	BaselineType          string           `json:"baselineType"`
	BaselineWindowSize    string           `json:"baselineWindowSize"`
	DescriptionExpression string           `json:"descriptionExpression"`
	Enabled               bool             `json:"enabled"`
	EntitySelectors       []EntitySelector `json:"entitySelectors"`
	FilterExpression      string           `json:"filterExpression"`
	GroupByFields         []string         `json:"groupByFields"`
	IsPrototype           bool             `json:"isPrototype"`
	Name                  string           `json:"name"`
	NameExpression        string           `json:"nameExpression"`
	RetentionWindowSize   string           `json:"retentionWindowSize"`
	Severity              int              `json:"score"`
	SummaryExpression     string           `json:"summaryExpression"`
	Tags                  []string         `json:"tags"`
	ValueFields           []string         `json:"valueFields"`
	Version               int              `json:"version"`
	SuppressionWindowSize *int             `json:"suppressionWindowSize,omitempty"`
}

type CSEFirstSeenRuleOverride struct {
	BaselineWindowSize    string   `json:"baselineWindowSize"`
	DescriptionExpression string   `json:"descriptionExpression"`
	GroupByFields         []string `json:"groupByFields"`
	IsPrototype           bool     `json:"isPrototype"`
	Name                  string   `json:"name"`
	NameExpression        string   `json:"nameExpression"`
	RetentionWindowSize   string   `json:"retentionWindowSize"`
	Severity              int      `json:"score"`
	SummaryExpression     string   `json:"summaryExpression"`
	Tags                  []string `json:"tags"`
	SuppressionWindowSize *int     `json:"suppressionWindowSize,omitempty"`
}
