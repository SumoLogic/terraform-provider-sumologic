package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSERuleTuningExpression(id string) (*CSERuleTuningExpression, error) {
	data, err := s.Get(fmt.Sprintf("sec/v1/rule-tuning-expressions/%s", id))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSERuleTuningExpressionResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSERuleTuningExpression, nil
}

func (s *Client) DeleteCSERuleTuningExpression(id string) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/rule-tuning-expressions/%s", id))

	return err
}

func (s *Client) CreateCSERuleTuningExpression(CSERuleTuningExpression CSERuleTuningExpression) (string, error) {

	request := CSERuleTuningExpressionRequest{
		CSERuleTuningExpression: CSERuleTuningExpression,
	}

	var response CSERuleTuningExpressionResponse

	responseBody, err := s.Post("sec/v1/rule-tuning-expressions", request)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return "", err
	}

	return response.CSERuleTuningExpression.ID, nil
}

func (s *Client) UpdateCSERuleTuningExpression(CSERuleTuningExpression CSERuleTuningExpression) error {
	url := fmt.Sprintf("sec/v1/rule-tuning-expressions/%s", CSERuleTuningExpression.ID)

	CSERuleTuningExpression.ID = ""
	request := CSERuleTuningExpressionRequest{
		CSERuleTuningExpression: CSERuleTuningExpression,
	}

	_, err := s.Put(url, request)

	return err
}

type CSERuleTuningExpressionRequest struct {
	CSERuleTuningExpression CSERuleTuningExpression `json:"fields"`
}

type CSERuleTuningExpressionResponse struct {
	CSERuleTuningExpression CSERuleTuningExpression `json:"data"`
}

type CSERuleTuningExpression struct {
	ID          string   `json:"id,omitempty"`
	Description string   `json:"description"`
	Name        string   `json:"name"`
	Expression  string   `json:"expression"`
	Enabled     bool     `json:"enabled"`
	Exclude     bool     `json:"exclude"`
	IsGlobal    bool     `json:"isGlobal"`
	RuleIds     []string `json:"ruleIds"`
}
