package sumologic

import (
	"encoding/json"
	"fmt"
)

type HTTPSource struct {
	Source
	MessagePerRequest bool   `json:"messagePerRequest"`
	URL               string `json:"url,omitempty"`
}

func (s *Client) CreateHTTPSource(httpSource HTTPSource, collectorID int) (int, error) {

	type HTTPSourceMessage struct {
		Source HTTPSource `json:"source"`
	}

	request := HTTPSourceMessage{
		Source: httpSource,
	}

	urlPath := fmt.Sprintf("v1/collectors/%d/sources", collectorID)
	body, err := s.Post(urlPath, request, false)

	if err != nil {
		return -1, err
	}

	var response HTTPSourceMessage

	err = json.Unmarshal(body, &response)
	if err != nil {
		return -1, err
	}

	return response.Source.ID, nil
}

func (s *Client) GetHTTPSource(collectorID, sourceID int) (*HTTPSource, error) {

	body, _, err := s.Get(fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, sourceID), false)
	if err != nil {
		return nil, err
	}

	if body == nil {
		return nil, nil
	}

	type Response struct {
		Source HTTPSource `json:"source"`
	}

	var response Response

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Source, nil
}

func (s *Client) UpdateHTTPSource(source HTTPSource, collectorID int) error {

	type HTTPSourceMessage struct {
		Source HTTPSource `json:"source"`
	}

	request := HTTPSourceMessage{
		Source: source,
	}

	urlPath := fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, source.ID)
	_, err := s.Put(urlPath, request, false)

	return err
}
