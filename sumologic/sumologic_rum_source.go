package sumologic

import (
	"encoding/json"
	"fmt"
)

type RumSource struct {
	HTTPSource
	RumThirdPartyRef RumThirdPartyRef `json:"thirdPartyRef"`
}

type RumThirdPartyRef struct {
	Resources []RumThirdPartyResource `json:"resources"`
}

type RumThirdPartyResource struct {
	ServiceType string        `json:"serviceType"`
	Path        RumSourcePath `json:"path,omitempty"`
}

type RumSourcePath struct {
	Type                         string                 `json:"type"`
	ApplicationName              string                 `json:"applicationName,omitempty"`
	ServiceName                  string                 `json:"serviceName"`
	DeploymentEnvironment        string                 `json:"deploymentEnvironment,omitempty"`
	SamplingRate                 float32                `json:"samplingRate,omitempty"`
	IgnoreUrls                   []string               `json:"ignoreUrls,omitempty"`
	CustomTags                   map[string]interface{} `json:"customTags,omitempty"`
	PropagateTraceHeaderCorsUrls []string               `json:"propagateTraceHeaderCorsUrls,omitempty"`
	SelectedCountry              string                 `json:"selectedCountry,omitempty"`
}

func (s *Client) CreateRumSource(rumSource RumSource, collectorID int) (int, error) {

	type RumSourceMessage struct {
		Source RumSource `json:"source"`
	}

	request := RumSourceMessage{
		Source: rumSource,
	}

	urlPath := fmt.Sprintf("v1/collectors/%d/sources", collectorID)
	body, err := s.Post(urlPath, request)

	if err != nil {
		return -1, err
	}

	var response RumSourceMessage

	err = json.Unmarshal(body, &response)
	if err != nil {
		return -1, err
	}

	return response.Source.ID, nil
}

func (s *Client) GetRumSource(collectorID, sourceID int) (*RumSource, error) {

	body, _, err := s.Get(fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, sourceID))
	if err != nil {
		return nil, err
	}

	if body == nil {
		return nil, nil
	}

	type Response struct {
		Source RumSource `json:"source"`
	}

	var response Response

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Source, nil
}

func (s *Client) UpdateRumSource(source RumSource, collectorID int) error {

	type RumSourceMessage struct {
		Source RumSource `json:"source"`
	}

	request := RumSourceMessage{
		Source: source,
	}

	urlPath := fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, source.ID)
	_, err := s.Put(urlPath, request)

	return err
}
