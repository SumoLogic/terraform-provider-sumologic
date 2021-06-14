package sumologic

import (
	"encoding/json"
	"fmt"
)

type CloudToCloudSource struct {
	ID        int             `json:"id,omitempty"`
	Type      string          `json:"sourceType"`
	Config    json.RawMessage `json:"config"`
	SchemaRef SchemaReference `json:"schemaRef"`
}

type SchemaReference struct {
	Type string `json:"type"`
}

func (s *Client) CreateCloudToCloudSource(source CloudToCloudSource, collectorID int) (int, error) {

	type CloudToCloudSourceMessage struct {
		Source CloudToCloudSource `json:"source"`
	}

	request := CloudToCloudSourceMessage{
		Source: source,
	}

	urlPath := fmt.Sprintf("v1/collectors/%d/sources", collectorID)

	body, err := s.Post(urlPath, request, false)

	if err != nil {
		return -1, err
	}

	var response CloudToCloudSourceMessage
	err = json.Unmarshal(body, &response)

	if err != nil {
		return -1, err
	}

	return response.Source.ID, nil
}

func (s *Client) GetCloudToCloudSource(collectorID, sourceID int) (*CloudToCloudSource, error) {
	urlPath := fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, sourceID)
	body, _, err := s.Get(urlPath, false)

	if err != nil {
		return nil, err
	}

	if body == nil {
		return nil, nil
	}

	type CloudToCloudSourceResponse struct {
		Source CloudToCloudSource `json:"source"`
	}

	var response CloudToCloudSourceResponse
	err = json.Unmarshal(body, &response)

	if err != nil {
		return nil, err
	}

	return &response.Source, nil
}

func (s *Client) UpdateCloudToCloudSource(source CloudToCloudSource, collectorID int) error {
	url := fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, source.ID)

	type CloudToCloudSourceMessage struct {
		Source CloudToCloudSource `json:"source"`
	}

	request := CloudToCloudSourceMessage{
		Source: source,
	}

	_, err := s.Put(url, request, false)

	return err
}
