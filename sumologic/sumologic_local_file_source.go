package sumologic

import (
	"encoding/json"
	"fmt"
)

type LocalFileSource struct {
	Source
}

func (s *Client) CreateLocalFileSource(source LocalFileSource, collectorID int) (int, error) {

	type LocalFileSourceMessage struct {
		Source LocalFileSource `json:"source"`
	}

	request := LocalFileSourceMessage{
		Source: source,
	}

	urlPath := fmt.Sprintf("v1/collectors/%d/sources", collectorID)
	body, err := s.Post(urlPath, request)

	if err != nil {
		return -1, err
	}

	var response LocalFileSourceMessage

	err = json.Unmarshal(body, &response)
	if err != nil {
		return -1, err
	}

	return response.Source.ID, nil
}

func (s *Client) GetLocalFileSource(collectorID, sourceID int) (*LocalFileSource, error) {

	body, _, err := s.Get(fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, sourceID))
	if err != nil {
		return nil, err
	}

	if body == nil {
		return nil, nil
	}

	type Response struct {
		Source LocalFileSource `json:"source"`
	}

	var response Response

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Source, nil
}

func (s *Client) UpdateLocalFileSource(source LocalFileSource, collectorID int) error {

	type LocalFileSourceMessage struct {
		Source LocalFileSource `json:"source"`
	}

	request := LocalFileSourceMessage{
		Source: source,
	}

	urlPath := fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, source.ID)
	_, err := s.Put(urlPath, request)

	return err
}
