package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCollector(id int) (*Collector, error) {
	data, _, err := s.Get(fmt.Sprintf("v1/collectors/%d", id), false)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CollectorResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.Collector, nil
}

func (s *Client) GetCollectorName(name string) (*Collector, error) {
	data, _, err := s.Get(fmt.Sprintf("v1/collectors/name/%s", name), false)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("collector with name '%s' does not exist", name)
	}

	var response CollectorResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.Collector, nil
}

func (s *Client) DeleteCollector(id int) error {
	_, err := s.Delete(fmt.Sprintf("v1/collectors/%d", id))

	return err
}

func (s *Client) CreateCollector(collector Collector) (int64, error) {

	request := CollectorRequest{
		Collector: collector,
	}

	var response CollectorResponse

	responseBody, err := s.Post("v1/collectors", request, false)
	if err != nil {
		return -1, err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return -1, err
	}

	return response.Collector.ID, nil
}

func (s *Client) UpdateCollector(collector Collector) error {
	url := fmt.Sprintf("v1/collectors/%d", collector.ID)

	request := CollectorRequest{
		Collector: collector,
	}

	_, err := s.Put(url, request, false)

	return err
}

type CollectorRequest struct {
	Collector Collector `json:"collector"`
}

type CollectorResponse struct {
	Collector Collector `json:"collector"`
}

type CollectorList struct {
	Collectors []Collector `json:"collectors"`
}

type CollectorLink struct {
	Rel  string `json:"rel,omitempty"`
	Href string `json:"href,omitempty"`
}

type Collector struct {
	ID               int64                  `json:"id,omitempty"`
	CollectorType    string                 `json:"collectorType,omitempty"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description,omitempty"`
	Category         string                 `json:"category,omitempty"`
	TimeZone         string                 `json:"timeZone,omitempty"`
	Fields           map[string]interface{} `json:"fields,omitempty"`
	Links            []CollectorLink        `json:"links,omitempty"`
	CollectorVersion string                 `json:"collectorVersion,omitempty"`
	LastSeenAlive    int64                  `json:"lastSeenAlive,omitempty"`
	Alive            bool                   `json:"alive,omitempty"`
}
