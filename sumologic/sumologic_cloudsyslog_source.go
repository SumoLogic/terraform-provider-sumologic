package sumologic

import (
	"encoding/json"
	"fmt"
)

type CloudSyslogSource struct {
	Source
	Token string `json:"token,omitempty"`
}

func (s *Client) CreateCloudsyslogSource(cloudSyslogSource CloudSyslogSource, collectorID int) (int, error) {

	type CloudSyslogSourceMessage struct {
		Source CloudSyslogSource `json:"source"`
	}

	request := CloudSyslogSourceMessage{
		Source: cloudSyslogSource,
	}

	urlPath := fmt.Sprintf("v1/collectors/%d/sources", collectorID)
	body, err := s.Post(urlPath, request, false)

	if err != nil {
		return -1, err
	}

	var response CloudSyslogSourceMessage

	err = json.Unmarshal(body, &response)
	if err != nil {
		return -1, err
	}

	return response.Source.ID, nil
}

func (s *Client) GetCloudSyslogSource(collectorID, sourceID int) (*CloudSyslogSource, error) {

	body, _, err := s.Get(fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, sourceID), false)
	if err != nil {
		return nil, err
	}

	if body == nil {
		return nil, nil
	}

	type Response struct {
		Source CloudSyslogSource `json:"source"`
	}

	var response Response

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Source, nil
}

func (s *Client) UpdateCloudSyslogSource(source CloudSyslogSource, collectorID int) error {

	type CloudSyslogSourceMessage struct {
		Source CloudSyslogSource `json:"source"`
	}

	request := CloudSyslogSourceMessage{
		Source: source,
	}

	urlPath := fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, source.ID)
	_, err := s.Put(urlPath, request, false)

	return err
}
