package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSEThresholdRule(id string) (*CSEThresholdRule, error) {
	data, _, err := s.Get(fmt.Sprintf("sec/v1/rules/%s", id))
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
	fmt.Printf("Window Size: %s\n", CSEThresholdRule.WindowSize)
	fmt.Printf("Window Size Name: %s\n", CSEThresholdRule.WindowSizeName)
	request := CSEThresholdRuleRequest{
		CSEThresholdRule: CSEThresholdRule,
	}

	_, err := s.Put(url, request)

	return err
}

type CSEThresholdRuleRequest struct {
	CSEThresholdRule CSEThresholdRule `json:"fields"`
}

type CSEThresholdRuleResponse struct {
	CSEThresholdRule CSEThresholdRule `json:"data"`
}

type CSEThresholdRule struct {
	ID                string           `json:"id,omitempty"`
	CountDistinct     bool             `json:"countDistinct"`
	CountField        string           `json:"countField"`
	Description       string           `json:"description"`
	Enabled           bool             `json:"enabled"`
	EntitySelectors   []EntitySelector `json:"entitySelectors"`
	Expression        string           `json:"expression"`
	GroupByFields     []string         `json:"groupByFields"`
	IsPrototype       bool             `json:"isPrototype"`
	Limit             int              `json:"limit"`
	Name              string           `json:"name"`
	Severity          int              `json:"score"`
	Stream            string           `json:"stream"`
	SummaryExpression string           `json:"summaryExpression"`
	Tags              []string         `json:"tags"`
	Version           int              `json:"version"`
	WindowSize        windowSizeField  `json:"windowSize,omitempty"`
	WindowSizeName    string           `json:"windowSizeName,omitempty"`
}
