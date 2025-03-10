package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSECustomInsight(id string) (*CSECustomInsight, error) {
	data, _, err := s.Get(fmt.Sprintf("sec/v1/custom-insights/%s", id))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSECustomInsightResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSECustomInsight, nil
}

func (s *Client) DeleteCSECustomInsight(id string) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/custom-insights/%s", id))

	return err
}

func (s *Client) CreateCSECustomInsight(CSECustomInsight CSECustomInsight) (string, error) {
	request := CSECustomInsightRequest{
		CSECustomInsight: CSECustomInsight,
	}

	var response CSECustomInsightResponse

	responseBody, err := s.Post("sec/v1/custom-insights", request)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return "", err
	}

	return response.CSECustomInsight.ID, nil
}

func (s *Client) UpdateCSECustomInsight(CSECustomInsight CSECustomInsight) error {
	url := fmt.Sprintf("sec/v1/custom-insights/%s", CSECustomInsight.ID)

	CSECustomInsight.ID = ""
	request := CSECustomInsightRequest{
		CSECustomInsight: CSECustomInsight,
	}

	_, err := s.Put(url, request)

	return err
}

type CSECustomInsightRequest struct {
	CSECustomInsight CSECustomInsight `json:"fields"`
}

type CSECustomInsightResponse struct {
	CSECustomInsight CSECustomInsight `json:"data"`
}

type DynamicSeverity struct {
	MinimumSignalSeverity int    `json:"minimumSignalSeverity"`
	InsightSeverity       string `json:"insightSeverity"`
}

type CSECustomInsight struct {
	ID                  string            `json:"id,omitempty"`
	Description         string            `json:"description"`
	Enabled             bool              `json:"enabled"`
	Name                string            `json:"name"`
	Ordered             bool              `json:"ordered"`
	RuleIds             []string          `json:"ruleIds"`
	Severity            string            `json:"severity"`
	DynamicSeverity     []DynamicSeverity `json:"dynamicSeverity"`
	SignalNames         []string          `json:"signalNames"`
	Tags                []string          `json:"tags"`
	SignalMatchStrategy string            `json:"signalMatchStrategy,omitempty"`
}
