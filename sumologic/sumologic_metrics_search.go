package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetMetricsSearch(id string) (*MetricsSearch, error) {
	urlWithoutParams := "v1/metricsSearches/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	data, _, err := s.Get(urlWithParams)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var metricsSearch MetricsSearch

	err = json.Unmarshal(data, &metricsSearch)
	if err != nil {
		return nil, err
	}

	return &metricsSearch, nil
}

func (s *Client) DeleteMetricsSearch(id string) error {
	urlWithoutParams := "v1/metricsSearches/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	_, err := s.Delete(urlWithParams)

	return err
}

func (s *Client) UpdateMetricsSearch(metricsSearch MetricsSearch) error {
	urlWithoutParams := "v1/metricsSearches/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, metricsSearch.ID)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	metricsSearch.ID = ""

	_, err := s.Put(urlWithParams, metricsSearch)

	return err
}

func (s *Client) CreateMetricsSearch(metricsSearch MetricsSearch) (string, error) {
	urlWithoutParams := "v1/metricsSearches"

	data, err := s.Post(urlWithoutParams, metricsSearch)
	if err != nil {
		return "", err
	}

	var createdMetricsSearch MetricsSearch

	err = json.Unmarshal(data, &createdMetricsSearch)
	if err != nil {
		return "", err
	}

	return createdMetricsSearch.ID, nil
}

type MetricsSearch struct {
	ID                        string               `json:"id,omitempty"`
	Title                     string               `json:"title"`
	Description               string               `json:"description"`
	ParentId                  string               `json:"parentId"`
	LogQuery                  string               `json:"logQuery"`
	DesiredQuantizationInSecs int                  `json:"desiredQuantizationInSecs"`
	TimeRange                 interface{}          `json:"timeRange"`
	MetricsQueries            []MetricsSearchQuery `json:"metricsQueries"`
}

type MetricsSearchQuery struct {
	RowId string `json:"rowId"`
	Query string `json:"query"`
}
