package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSEChainRule(id string) (*CSEChainRule, error) {
	data, err := s.Get(fmt.Sprintf("sec/v1/rules/%s", id))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSEChainRuleResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSEChainRule, nil
}

func (s *Client) DeleteCSEChainRule(id string) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/rules/%s", id))

	return err
}

func (s *Client) CreateCSEChainRule(CSEChainRule CSEChainRule) (string, error) {
	request := CSEChainRuleRequest{
		CSEChainRule: CSEChainRule,
	}

	var response CSEChainRuleResponse

	responseBody, err := s.Post("sec/v1/rules/chain", request)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return "", err
	}

	return response.CSEChainRule.ID, nil
}

func (s *Client) UpdateCSEChainRule(CSEChainRule CSEChainRule) error {
	url := fmt.Sprintf("sec/v1/rules/chain/%s", CSEChainRule.ID)

	CSEChainRule.ID = ""
	request := CSEChainRuleRequest{
		CSEChainRule: CSEChainRule,
	}

	_, err := s.Put(url, request)

	return err
}

type CSEChainRuleRequest struct {
	CSEChainRule CSEChainRule `json:"fields"`
}

type CSEChainRuleResponse struct {
	CSEChainRule CSEChainRule `json:"data"`
}

type ExpressionAndLimit struct {
	Expression string `json:"expression"`
	Limit      int    `json:"limit"`
}

type CSEChainRule struct {
	ID                     string               `json:"id,omitempty"`
	Description            string               `json:"description"`
	Enabled                bool                 `json:"enabled"`
	EntitySelectors        []EntitySelector     `json:"entitySelectors"`
	ExpressionsAndLimits   []ExpressionAndLimit `json:"expressionsAndLimits"`
	GroupByFields          []string             `json:"groupByFields"`
	IsPrototype            bool                 `json:"isPrototype"`
	Ordered                bool                 `json:"ordered"`
	Name                   string               `json:"name"`
	Severity               int                  `json:"score"`
	Stream                 string               `json:"stream"`
	SummaryExpression      string               `json:"summaryExpression"`
	Tags                   []string             `json:"tags"`
	WindowSize             windowSizeField      `json:"windowSize,omitempty"`
	WindowSizeName         string               `json:"windowSizeName,omitempty"`
	WindowSizeMilliseconds string               `json:"windowSizeMilliseconds,omitempty"`
	SuppressionWindowSize  *int                 `json:"suppressionWindowSize,omitempty"`
}
