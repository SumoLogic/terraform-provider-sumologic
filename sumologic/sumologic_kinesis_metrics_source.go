package sumologic

import (
	"encoding/json"
	"fmt"
)

type KinesisMetricsSource struct {
	Source
	MessagePerRequest bool                 `json:"messagePerRequest"`
	URL               string               `json:"url,omitempty"`
	ThirdPartyRef     PollingThirdPartyRef `json:"thirdPartyRef"`
}

func (s *Client) CreateKinesisMetricsSource(kinesisMetricsSource KinesisMetricsSource, collectorID int) (int, error) {

	type KinesisMetricsSourceMessage struct {
		Source KinesisMetricsSource `json:"source"`
	}

	request := KinesisMetricsSourceMessage{
		Source: kinesisMetricsSource,
	}

	urlPath := fmt.Sprintf("v1/collectors/%d/sources", collectorID)
	body, err := s.Post(urlPath, request, false)

	if err != nil {
		return -1, err
	}

	var response KinesisMetricsSourceMessage

	err = json.Unmarshal(body, &response)
	if err != nil {
		return -1, err
	}

	return response.Source.ID, nil
}

func (s *Client) GetKinesisMetricsSource(collectorID, sourceID int) (*KinesisMetricsSource, error) {

	body, _, err := s.Get(fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, sourceID), false)
	if err != nil {
		return nil, err
	}

	if body == nil {
		return nil, nil
	}

	type Response struct {
		Source KinesisMetricsSource `json:"source"`
	}

	var response Response

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Source, nil
}

func (s *Client) UpdateKinesisMetricsSource(source KinesisMetricsSource, collectorID int) error {

	type KinesisMetricsSourceMessage struct {
		Source KinesisMetricsSource `json:"source"`
	}

	request := KinesisMetricsSourceMessage{
		Source: source,
	}

	urlPath := fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, source.ID)
	_, err := s.Put(urlPath, request, false)

	return err
}
