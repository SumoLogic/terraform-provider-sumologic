package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCollector(id int) (*Collector, error) {
	data, _, err := s.Get(fmt.Sprintf("collectors/%d", id))
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
	// TODO: check default limit count of 1000 and paginate
	data, _, err := s.Get("collectors")
	if err != nil {
		return nil, err
	}

	if data == nil {
		return &Collector{}, nil
	}

	var response CollectorList
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	for _, c := range response.Collectors {
		if c.Name == name {
			return &c, nil
		}
	}

	return nil, nil
}

func (s *Client) DeleteCollector(id int) error {
	_, err := s.Delete(fmt.Sprintf("collectors/%d", id))

	return err
}

func (s *Client) CreateCollector(collector Collector) (int, error) {

	request := CollectorRequest{
		Collector: collector,
	}

	var response CollectorResponse

	responseBody, err := s.Post("collectors", request)
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
	url := fmt.Sprintf("collectors/%d", collector.ID)

	request := CollectorRequest{
		Collector: collector,
	}

	_, err := s.Put(url, request)

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
	ID               int             `json:"id,omitempty"`
	CollectorType    string          `json:"collectorType,omitempty"`
	Name             string          `json:"name"`
	Description      string          `json:"description,omitempty"`
	Category         string          `json:"category,omitempty"`
	TimeZone         string          `json:"timeZone,omitempty"`
	Links            []CollectorLink `json:"links,omitempty"`
	CollectorVersion string          `json:"collectorVersion,omitempty"`
	LastSeenAlive    int             `json:"lastSeenAlive,omitempty"`
	Alive            bool            `json:"alive,omitempty"`
}
