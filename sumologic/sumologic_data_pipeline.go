package sumologic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

type DataPipelineRequest struct {
	Name            string             `json:"name"`
	Description     string             `json:"description,omitempty"`
	PipelineType    string             `json:"pipelineType"`
	IsEnabled       bool               `json:"isEnabled"`
	RouteExpression string             `json:"routeExpression,omitempty"`
	Nodes           []DataPipelineNode `json:"nodes"`
}

type DataPipelineUpdateRequest struct {
	DataPipelineRequest
	Version int `json:"version"`
}

type DataPipelineMetadataRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
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

func (s *Client) CreateDataPipeline(request DataPipelineRequest) (*DataPipeline, error) {
	data, err := s.Post("v1/pipelines", request)
	if err != nil {
		return nil, err
	}

	var pipeline DataPipeline
	if err := json.Unmarshal(data, &pipeline); err != nil {
		return nil, err
	}

	return &pipeline, nil
}

func (s *Client) UpdateDataPipeline(id string, request DataPipelineUpdateRequest) (*DataPipeline, error) {
	url := fmt.Sprintf("v1/pipelines/%s", id)

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := s.createSumoRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := s.doSumoRequest(req)
	if err != nil {
		return nil, err
	}

	data, err := s.handleSumoResponse(resp)
	if err != nil {
		return nil, err
	}

	var pipeline DataPipeline
	if err := json.Unmarshal(data, &pipeline); err != nil {
		return nil, err
	}

	return &pipeline, nil
}

func (s *Client) PublishDataPipeline(id string) (*DataPipeline, error) {
	url := fmt.Sprintf("v1/pipelines/%s/publish", id)

	data, err := s.PostRawPayload(url, "")
	if err != nil {
		return nil, err
	}

	var pipeline DataPipeline
	if err := json.Unmarshal(data, &pipeline); err != nil {
		return nil, err
	}

	return &pipeline, nil
}

func (s *Client) SetDataPipelineEnabled(id string, enabled bool) error {
	action := "disable"
	if enabled {
		action = "enable"
	}

	url := fmt.Sprintf("v1/pipelines/%s/%s", id, action)
	_, err := s.Patch(url, nil)
	return err
}

func (s *Client) UpdateDataPipelineMetadata(id string, request DataPipelineMetadataRequest) (*DataPipeline, error) {
	url := fmt.Sprintf("v1/pipelines/%s/metadata", id)

	data, err := s.Patch(url, request)
	if err != nil {
		return nil, err
	}

	var pipeline DataPipeline
	if err := json.Unmarshal(data, &pipeline); err != nil {
		return nil, err
	}

	return &pipeline, nil
}

func (s *Client) DeleteDataPipeline(id string) error {
	url := fmt.Sprintf("v1/pipelines/%s", id)
	_, err := s.Delete(url)
	return err
}
