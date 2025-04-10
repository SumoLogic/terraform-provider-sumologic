package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetMetricsSearchV2(id string) (*MetricsSearchV2, error) {
	urlWithoutParams := "v2/metricsSearches/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	data, err := s.Get(urlWithParams)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var metricsSearchV2 MetricsSearchV2

	err = json.Unmarshal(data, &metricsSearchV2)
	if err != nil {
		return nil, err
	}

	return &metricsSearchV2, nil
}

func (s *Client) DeleteMetricsSearchV2(id string) error {
	urlWithoutParams := "v2/metricsSearches/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	_, err := s.Delete(urlWithParams)

	return err
}

func (s *Client) UpdateMetricsSearchV2(metricsSearchV2 MetricsSearchV2) error {
	urlWithoutParams := "v2/metricsSearches/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, metricsSearchV2.ID)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	metricsSearchV2.ID = ""

	_, err := s.Put(urlWithParams, metricsSearchV2)

	return err
}

func (s *Client) CreateMetricsSearchV2(metricsSearchV2 MetricsSearchV2) (string, error) {
	urlWithoutParams := "v2/metricsSearches"
	data, err := s.Post(urlWithoutParams, metricsSearchV2)
	if err != nil {
		return "", err
	}

	var createdMetricsSearchV2 MetricsSearchV2

	err = json.Unmarshal(data, &createdMetricsSearchV2)
	if err != nil {
		return "", err
	}

	return createdMetricsSearchV2.ID, nil
}

type MetricsSearchV2 struct {
	Title          string                 `json:"title"`
	TimeRange      interface{}            `json:"timeRange"`
	Description    string                 `json:"description"`
	VisualSettings string                 `json:"visualSettings"`
	FolderID       string                 `json:"folderId"`
	ID             string                 `json:"id,omitempty"`
	Queries        []MetricsSearchQueryV2 `json:"queries"`
}

type MetricsSearchQueryV2 struct {
	QueryKey         string            `json:"queryKey"`
	QueryString      string            `json:"queryString"`
	QueryType        string            `json:"queryType"`
	MetricsQueryMode string            `json:"metricsQueryMode,omitempty"`
	ParseMode        string            `json:"parseMode,omitempty"`
	TimeSource       string            `json:"timeSource,omitempty"`
	Transient        bool              `json:"transient,omitempty"`
	MetricsQueryData *MetricsQueryData `json:"metricsQueryData,omitempty"`
	TracesQueryData  interface{}       `json:"tracesQueryData,omitempty"`
	SpansQueryData   interface{}       `json:"spansQueryData,omitempty"`
}
