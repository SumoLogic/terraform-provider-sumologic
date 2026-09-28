package sumologic

import (
	"encoding/json"
	"fmt"
)

type DataPipeline struct {
	ID               string             `json:"id"`
	CustomerID       string             `json:"customerId"`
	Name             string             `json:"name"`
	Description      string             `json:"description,omitempty"`
	PipelineType     string             `json:"pipelineType"`
	IsEnabled        bool               `json:"isEnabled"`
	RouteExpression  string             `json:"routeExpression,omitempty"`
	Order            int                `json:"order"`
	Version          int                `json:"version"`
	CreatedAt        string             `json:"createdAt"`
	ModifiedAt       string             `json:"modifiedAt"`
	CreatedByUserID  string             `json:"createdByUserId"`
	ModifiedByUserID string             `json:"modifiedByUserId"`
	State            string             `json:"state,omitempty"`
	Nodes            []DataPipelineNode `json:"nodes"`
}

type DataPipelineNode struct {
	ID               string                  `json:"id,omitempty"`
	Name             string                  `json:"name"`
	NodeType         string                  `json:"nodeType"`
	FilterExpression string                  `json:"filterExpression,omitempty"`
	IsEnabled        *bool                   `json:"isEnabled,omitempty"`
	Processors       []DataPipelineProcessor `json:"processors,omitempty"`
	Outputs          []DataPipelineOutput    `json:"outputs,omitempty"`
}

type DataPipelineProcessor struct {
	ID            string          `json:"id,omitempty"`
	Name          string          `json:"name"`
	ProcessorType string          `json:"processorType"`
	Order         int             `json:"order"`
	IsEnabled     bool            `json:"isEnabled"`
	Config        json.RawMessage `json:"config,omitempty"`
}

type DataPipelineOutput struct {
	Target    string `json:"target"`
	Condition string `json:"condition,omitempty"`
	Order     *int   `json:"order,omitempty"`
}

type DataPipelineList struct {
	Pipelines []DataPipeline `json:"pipelines"`
}

func (s *Client) GetDataPipeline(id string) (*DataPipeline, error) {
	url := fmt.Sprintf("v1/pipelines/%s", id)

	data, err := s.Get(url)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var pipeline DataPipeline
	if err := json.Unmarshal(data, &pipeline); err != nil {
		return nil, err
	}

	return &pipeline, nil
}

func (s *Client) ListDataPipelines() ([]DataPipeline, error) {
	data, err := s.Get("v1/pipelines")
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var response DataPipelineList
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response.Pipelines, nil
}

func (s *Client) FindDataPipelineByName(name string) (*DataPipeline, error) {
	pipelines, err := s.ListDataPipelines()
	if err != nil {
		return nil, err
	}

	var matches []DataPipeline
	for _, pipeline := range pipelines {
		if pipeline.Name == name {
			matches = append(matches, pipeline)
		}
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf("found %d data pipelines named %q, expected at most 1", len(matches), name)
	}
	if len(matches) == 0 {
		return nil, nil
	}

	return s.GetDataPipeline(matches[0].ID)
}
