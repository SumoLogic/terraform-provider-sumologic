package sumologic

import (
	"encoding/json"
	"fmt"
)

type GCPSource struct {
	Source
	MessagePerRequest bool             `json:"messagePerRequest"`
	URL               string           `json:"url,omitempty"`
	ThirdPartyRef     GCPThirdPartyRef `json:"thirdPartyRef,omitempty"`
}

type GCPThirdPartyRef struct {
	Resources []GCPResource `json:"resources"`
}

type GCPResource struct {
	ServiceType    string             `json:"serviceType"`
	Authentication *GCPAuthentication `json:"authentication,omitempty"`
	Path           *GCPPath           `json:"path,omitempty"`
}

type GCPAuthentication struct {
	Type string `json:"type"`
}

type GCPPath struct {
	Type string `json:"type"`
}

func (s *Client) CreateGCPSource(gcpSource GCPSource, collectorID int) (int, error) {

	type GCPSourceMessage struct {
		Source GCPSource `json:"source"`
	}

	request := GCPSourceMessage{
		Source: gcpSource,
	}

	urlPath := fmt.Sprintf("v1/collectors/%d/sources", collectorID)
	body, err := s.Post(urlPath, request, false)

	if err != nil {
		return -1, err
	}

	var response GCPSourceMessage

	err = json.Unmarshal(body, &response)
	if err != nil {
		return -1, err
	}

	return response.Source.ID, nil
}

func (s *Client) GetGCPSource(collectorID, sourceID int) (*GCPSource, error) {

	body, _, err := s.Get(fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, sourceID), false)
	if err != nil {
		return nil, err
	}

	if body == nil {
		return nil, nil
	}

	type Response struct {
		Source GCPSource `json:"source"`
	}

	var response Response

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Source, nil
}

func (s *Client) UpdateGCPSource(source GCPSource, collectorID int) error {

	type GCPSourceMessage struct {
		Source GCPSource `json:"source"`
	}

	request := GCPSourceMessage{
		Source: source,
	}

	urlPath := fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, source.ID)
	_, err := s.Put(urlPath, request, false)

	return err
}
