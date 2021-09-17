package sumologic

import (
	"encoding/json"
)

func (s *Client) GetCSEInsightsConfiguration() (*CSEInsightsConfiguration, error) {
	data, _, err := s.Get("sec/v1/insights-configuration", false)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSEInsightsConfigurationResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSEInsightsConfiguration, nil
}

func (s *Client) UpdateCSEInsightsConfiguration(CSEInsightsConfiguration CSEInsightsConfiguration) error {
	url := "sec/v1/insights-configuration"

	request := CSEInsightsConfigurationRequest{
		CSEInsightsConfiguration: CSEInsightsConfiguration,
	}

	_, err := s.Put(url, request, false)

	return err
}

type CSEInsightsConfigurationRequest struct {
	CSEInsightsConfiguration CSEInsightsConfiguration `json:"config"`
}

type CSEInsightsConfigurationResponse struct {
	CSEInsightsConfiguration CSEInsightsConfiguration `json:"data"`
}

type CSEInsightsConfiguration struct {
	LookbackDays float64 `json:"lookbackDays,omitempty"`
	Threshold    float64 `json:"threshold,omitempty"`
}
