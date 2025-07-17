package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSEAggregationRule(id string) (*CSEAggregationRule, error) {
	data, err := s.Get(fmt.Sprintf("sec/v1/rules/%s", id))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSEAggregationRuleResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSEAggregationRule, nil
}

func (s *Client) DeleteCSEAggregationRule(id string) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/rules/%s", id))

	return err
}

func (s *Client) CreateCSEAggregationRule(CSEAggregationRule CSEAggregationRule) (string, error) {
	request := CSEAggregationRuleRequest{
		CSEAggregationRule: CSEAggregationRule,
	}

	var response CSEAggregationRuleResponse

	responseBody, err := s.Post("sec/v1/rules/aggregation", request)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return "", err
	}

	return response.CSEAggregationRule.ID, nil
}

func (s *Client) UpdateCSEAggregationRule(CSEAggregationRule CSEAggregationRule) error {
	url := fmt.Sprintf("sec/v1/rules/aggregation/%s", CSEAggregationRule.ID)

	CSEAggregationRule.ID = ""
	request := CSEAggregationRuleRequest{
		CSEAggregationRule: CSEAggregationRule,
	}

	_, err := s.Put(url, request)

	return err
}

type CSEAggregationRuleRequest struct {
	CSEAggregationRule CSEAggregationRule `json:"fields"`
}

type CSEAggregationRuleResponse struct {
	CSEAggregationRule CSEAggregationRule `json:"data"`
}

type AggregationFunction struct {
	Name      string   `json:"name"`
	Function  string   `json:"function"`
	Arguments []string `json:"arguments"`
}

type CSEAggregationRule struct {
	ID                     string                `json:"id,omitempty"`
	AggregationFunctions   []AggregationFunction `json:"aggregationFunctions"`
	DescriptionExpression  string                `json:"descriptionExpression"`
	Enabled                bool                  `json:"enabled"`
	EntitySelectors        []EntitySelector      `json:"entitySelectors"`
	GroupByEntity          bool                  `json:"groupByAsset"`
	GroupByFields          []string              `json:"groupByFields"`
	IsPrototype            bool                  `json:"isPrototype"`
	MatchExpression        string                `json:"matchExpression"`
	Name                   string                `json:"name"`
	NameExpression         string                `json:"nameExpression"`
	SeverityMapping        SeverityMapping       `json:"scoreMapping"`
	Stream                 string                `json:"stream"`
	SummaryExpression      string                `json:"summaryExpression"`
	TriggerExpression      string                `json:"triggerExpression"`
	Tags                   []string              `json:"tags"`
	WindowSize             windowSizeField       `json:"windowSize,omitempty"`
	WindowSizeName         string                `json:"windowSizeName,omitempty"`
	WindowSizeMilliseconds string                `json:"windowSizeMilliseconds,omitempty"`
	SuppressionWindowSize  *int                  `json:"suppressionWindowSize,omitempty"`
}
