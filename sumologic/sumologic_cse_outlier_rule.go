package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSEOutlierRule(id string) (*CSEOutlierRule, error) {
	data, err := s.Get(fmt.Sprintf("sec/v1/rules/%s", id))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSEOutlierRuleResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSEOutlierRule, nil
}

func (s *Client) DeleteCSEOutlierRule(id string) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/rules/%s", id))

	return err
}

func (s *Client) CreateCSEOutlierRule(CSEOutlierRule CSEOutlierRule) (string, error) {
	request := CSEOutlierRuleRequest{
		CSEOutlierRule: CSEOutlierRule,
	}

	var response CSEOutlierRuleResponse

	responseBody, err := s.Post("sec/v1/rules/outlier", request)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return "", err
	}

	return response.CSEOutlierRule.ID, nil
}

func (s *Client) UpdateCSEOutlierRule(CSEOutlierRule CSEOutlierRule) error {
	url := fmt.Sprintf("sec/v1/rules/outlier/%s", CSEOutlierRule.ID)

	CSEOutlierRule.ID = ""
	request := CSEOutlierRuleRequest{
		CSEOutlierRule: CSEOutlierRule,
	}

	_, err := s.Put(url, request)

	return err
}

func (s *Client) OverrideCSEOutlierRule(CSEOutlierRule CSEOutlierRule) error {
	url := fmt.Sprintf("sec/v1/rules/outlier/%s/override", CSEOutlierRule.ID)

	CSEOutlierRuleOverride := toOverrideOutlier(CSEOutlierRule)

	request := CSEOutlierRuleOverrideRequest{
		CSEOutlierRuleOverride: CSEOutlierRuleOverride,
	}
	_, err := s.Put(url, request)

	return err
}

func toOverrideOutlier(CSEOutlierRule CSEOutlierRule) CSEOutlierRuleOverride {
	return CSEOutlierRuleOverride{
		BaselineWindowSize:    CSEOutlierRule.BaselineWindowSize,
		DescriptionExpression: CSEOutlierRule.DescriptionExpression,
		DeviationThreshold:    CSEOutlierRule.DeviationThreshold,
		EntitySelectors:       CSEOutlierRule.EntitySelectors,
		FloorValue:            CSEOutlierRule.FloorValue,
		GroupByFields:         CSEOutlierRule.GroupByFields,
		IsPrototype:           CSEOutlierRule.IsPrototype,
		Name:                  CSEOutlierRule.Name,
		NameExpression:        CSEOutlierRule.NameExpression,
		RetentionWindowSize:   CSEOutlierRule.RetentionWindowSize,
		Severity:              CSEOutlierRule.Severity,
		SummaryExpression:     CSEOutlierRule.SummaryExpression,
		Tags:                  CSEOutlierRule.Tags,
		WindowSize:            CSEOutlierRule.WindowSize,
		SuppressionWindowSize: CSEOutlierRule.SuppressionWindowSize,
	}
}

type CSEOutlierRuleRequest struct {
	CSEOutlierRule CSEOutlierRule `json:"fields"`
}

type CSEOutlierRuleResponse struct {
	CSEOutlierRule CSEOutlierRule `json:"data"`
}

type CSEOutlierRuleOverrideRequest struct {
	CSEOutlierRuleOverride CSEOutlierRuleOverride `json:"fields"`
}

type CSEOutlierRule struct {
	ID                    string                `json:"id,omitempty"`
	AssetField            string                `json:"assetField"`
	AggregationFunctions  []AggregationFunction `json:"aggregationFunctions"`
	BaselineWindowSize    string                `json:"baselineWindowSize"`
	DescriptionExpression string                `json:"descriptionExpression"`
	DeviationThreshold    int                   `json:"deviationThreshold"`
	Enabled               bool                  `json:"enabled"`
	EntitySelectors       []EntitySelector      `json:"entitySelectors"`
	FloorValue            int                   `json:"floorValue"`
	GroupByFields         []string              `json:"groupByFields"`
	IsPrototype           bool                  `json:"isPrototype"`
	MatchExpression       string                `json:"matchExpression"`
	Name                  string                `json:"name"`
	NameExpression        string                `json:"nameExpression"`
	RetentionWindowSize   string                `json:"retentionWindowSize"`
	Severity              int                   `json:"score"`
	SummaryExpression     string                `json:"summaryExpression"`
	Tags                  []string              `json:"tags"`
	WindowSize            windowSizeField       `json:"windowSize,omitempty"`
	WindowSizeName        string                `json:"windowSizeName,omitempty"`
	SuppressionWindowSize *int                  `json:"suppressionWindowSize,omitempty"`
}

type CSEOutlierRuleOverride struct {
	BaselineWindowSize    string           `json:"baselineWindowSize"`
	DescriptionExpression string           `json:"descriptionExpression"`
	DeviationThreshold    int              `json:"deviationThreshold"`
	EntitySelectors       []EntitySelector `json:"entitySelectors"`
	FloorValue            int              `json:"floorValue"`
	GroupByFields         []string         `json:"groupByFields"`
	IsPrototype           bool             `json:"isPrototype"`
	Name                  string           `json:"name"`
	NameExpression        string           `json:"nameExpression"`
	RetentionWindowSize   string           `json:"retentionWindowSize"`
	Severity              int              `json:"score"`
	SummaryExpression     string           `json:"summaryExpression"`
	Tags                  []string         `json:"tags"`
	WindowSize            windowSizeField  `json:"windowSize,omitempty"`
	SuppressionWindowSize *int             `json:"suppressionWindowSize,omitempty"`
}
