package sumologic

import (
	"encoding/json"
	"fmt"
)

type CloudsyslogSource struct {
	Source
	Token string `json:"token,omitempty"`
}

func (s *Client) CreateCloudsyslogSource(name string, description string, collectorID int) (int, error) {

	type CloudsyslogSourceMessage struct {
		Source CloudsyslogSource `json:"source"`
	}

	request := CloudsyslogSourceMessage{
		Source: CloudsyslogSource{
			Source: Source{
				Type:        "Cloudsyslog",
				Name:        name,
				Description: description,
			},
		},
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

	newSource := response.Source

	return newSource.ID, nil
}

func (s *Client) GetCloudsyslogSource(collectorID, sourceID int) (*CloudsyslogSource, error) {

	urlPath := fmt.Sprintf("collectors/%d/sources/%d", collectorID, sourceID)
	body, _, err := s.Get(urlPath)

	if err != nil {
		return nil, err
	}

	type Response struct {
		Source CloudsyslogSource `json:"source"`
	}

	var response Response
	json.Unmarshal(body, &response)

	var source = response.Source

	return &source, nil
}

func (s *Client) UpdateCloudsyslogSource(source CloudsyslogSource, collectorID int) error {

	url := fmt.Sprintf("collectors/%d/sources/%d", collectorID, source.ID)

	type CloudsyslogSourceMessage struct {
		Source CloudsyslogSource `json:"source"`
	}

	request := CloudsyslogSourceMessage{
		Source: source,
	}

	_, err := s.Put(url, request)

	return err
}
