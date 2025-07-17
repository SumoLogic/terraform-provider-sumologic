package sumologic

import (
	"encoding/json"
	"fmt"
)

type KinesisLogSource struct {
	Source
	MessagePerRequest bool                    `json:"messagePerRequest"`
	URL               string                  `json:"url,omitempty"`
	ThirdPartyRef     KinesisLogThirdPartyRef `json:"thirdPartyRef"`
}

type KinesisLogThirdPartyRef struct {
	Resources []KinesisLogResource `json:"resources"`
}

type KinesisLogResource struct {
	ServiceType    string                `json:"serviceType"`
	Authentication PollingAuthentication `json:"authentication,omitempty"`
	Path           KinesisLogPath        `json:"path,omitempty"`
}

type KinesisLogPath struct {
	Type           string `json:"type"`
	BucketName     string `json:"bucketName,omitempty"`
	PathExpression string `json:"pathExpression,omitempty"`
	ScanInterval   int    `json:"scanInterval,omitempty"`
}

func (s *Client) CreateKinesisLogSource(kinesisLogSource KinesisLogSource, collectorID int) (int, error) {

	type KinesisLogSourceMessage struct {
		Source KinesisLogSource `json:"source"`
	}

	request := KinesisLogSourceMessage{
		Source: kinesisLogSource,
	}

	urlPath := fmt.Sprintf("v1/collectors/%d/sources", collectorID)
	body, err := s.Post(urlPath, request)

	if err != nil {
		return -1, err
	}

	var response KinesisLogSourceMessage

	err = json.Unmarshal(body, &response)
	if err != nil {
		return -1, err
	}

	return response.Source.ID, nil
}

func (s *Client) GetKinesisLogSource(collectorID, sourceID int) (*KinesisLogSource, error) {

	body, err := s.Get(
		fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, sourceID),
	)
	if err != nil {
		return nil, err
	}

	if body == nil {
		return nil, nil
	}

	type Response struct {
		Source KinesisLogSource `json:"source"`
	}

	var response Response

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Source, nil
}

func (s *Client) UpdateKinesisLogSource(source KinesisLogSource, collectorID int) error {

	type KinesisLogSourceMessage struct {
		Source KinesisLogSource `json:"source"`
	}

	request := KinesisLogSourceMessage{
		Source: source,
	}

	urlPath := fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, source.ID)
	_, err := s.Put(urlPath, request)

	return err
}
