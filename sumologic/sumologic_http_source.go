package sumologic

import (
	"encoding/json"
	"fmt"
)

type HTTPSource struct {
	Source
	MessagePerRequest bool              `json:"messagePerRequest"`
	URL               string            `json:"url,omitempty"`
	Token             string            `json:"token,omitempty"`
	BaseUrl           string            `json:"baseUrl,omitempty"`
	ThirdPartyRef     HTTPThirdPartyRef `json:"thirdPartyRef,omitempty"`
}

type HTTPThirdPartyRef struct {
	Resources []HTTPResource `json:"resources"`
}

type HTTPResource struct {
	ServiceType    string             `json:"serviceType"`
	Path           HTTPPath           `json:"path"`
	Authentication HTTPAuthentication `json:"authentication"`
}

type HTTPPath struct {
	Type     string `json:"type"`
	Workload string `json:"workload,omitempty"`
	Region   string `json:"region,omitempty"`
}

type HTTPAuthentication struct {
	Type         string `json:"type"`
	TenantId     string `json:"tenantId,omitempty"`
	ClientId     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
}

func (s *Client) CreateHTTPSource(httpSource HTTPSource, collectorID int) (int, error) {

	type HTTPSourceMessage struct {
		Source HTTPSource `json:"source"`
	}

	request := HTTPSourceMessage{
		Source: httpSource,
	}

	urlPath := fmt.Sprintf("v1/collectors/%d/sources", collectorID)
	body, err := s.Post(urlPath, request)

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

	body, err := s.Get(fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, sourceID))
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
	_, err := s.Put(urlPath, request)

	return err
}
