package sumologic

import (
	"encoding/json"
	"fmt"
)

//TODO: Build Struct
type CloudsyslogSource struct {
	Source
}

func (s *Client) CreateCloudsyslogSource(cloudsyslogSource CloudsyslogSource, collectorID int) (int, error) {

	type CloudsyslogSourceMessage struct {
		Source CloudsyslogSource `json:"source"`
	}

	request := CloudsyslogSourceMessage{
		Source: cloudsyslogSource,
	}

	urlPath := fmt.Sprintf("collectors/%d/sources", collectorID)
	body, err := s.Post(urlPath, request)

	if err != nil {
		return -1, err
	}

	var response CloudsyslogSourceMessage

	err = json.Unmarshal(body, &response)
	if err != nil {
		return -1, err
	}

	return response.Source.ID, nil
}

func (s *Client) GetCloudsyslogSource(collectorID, sourceID int) (*CloudsyslogSource, error) {

	body, _, err := s.Get(
		fmt.Sprintf("collectors/%d/sources/%d", collectorID, sourceID),
	)

	if err != nil {
		return nil, err
	}

	type Response struct {
		Source CloudsyslogSource `json:"source"`
	}

	var response Response

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Source, nil
}

func (s *Client) UpdateCloudsyslogSource(source CloudsyslogSource, collectorID int) error {

	type CloudsyslogSourceMessage struct {
		Source CloudsyslogSource `json:"source"`
	}

	request := CloudsyslogSourceMessage{
		Source: source,
	}

	urlPath := fmt.Sprintf("collectors/%d/sources/%d", collectorID, source.ID)
	_, err := s.Put(urlPath, request)

	return err
}
