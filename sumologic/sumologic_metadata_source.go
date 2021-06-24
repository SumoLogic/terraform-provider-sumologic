package sumologic

import (
	"encoding/json"
	"fmt"
)

type MetadataSource struct {
	Source
	ContentType   string                `json:"contentType"`
	ScanInterval  int                   `json:"scanInterval"`
	Paused        bool                  `json:"paused"`
	URL           string                `json:"url"`
	ThirdPartyRef MetadataThirdPartyRef `json:"thirdPartyRef,omitempty"`
}

type MetadataThirdPartyRef struct {
	Resources []MetadataResource `json:"resources"`
}

type MetadataResource struct {
	ServiceType    string                 `json:"serviceType"`
	Authentication MetadataAuthentication `json:"authentication"`
	Path           MetadataPath           `json:"path"`
}

type MetadataAuthentication struct {
	Type    string `json:"type"`
	AwsID   string `json:"awsId"`
	AwsKey  string `json:"awsKey"`
	RoleARN string `json:"roleARN"`
}

type MetadataPath struct {
	Type              string   `json:"type"`
	LimitToRegions    []string `json:"limitToRegions,omitempty"`
	LimitToNamespaces []string `json:"limitToNamespaces,omitempty"`
	TagFilters        []string `json:"tagFilters,omitempty"`
}

func (s *Client) CreateMetadataSource(source MetadataSource, collectorID int) (int, error) {

	type MetadataSourceMessage struct {
		Source MetadataSource `json:"source"`
	}

	request := MetadataSourceMessage{
		Source: source,
	}

	urlPath := fmt.Sprintf("v1/collectors/%d/sources", collectorID)

	body, err := s.Post(urlPath, request, false)

	if err != nil {
		return -1, err
	}

	var response MetadataSourceMessage
	err = json.Unmarshal(body, &response)

	if err != nil {
		return -1, err
	}

	return response.Source.ID, nil
}

func (s *Client) GetMetadataSource(collectorID, sourceID int) (*MetadataSource, error) {
	urlPath := fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, sourceID)
	body, _, err := s.Get(urlPath, false)

	if err != nil {
		return nil, err
	}

	if body == nil {
		return nil, nil
	}

	type MetadataSourceResponse struct {
		Source MetadataSource `json:"source"`
	}

	var response MetadataSourceResponse
	err = json.Unmarshal(body, &response)

	if err != nil {
		return nil, err
	}

	return &response.Source, nil
}

func (s *Client) UpdateMetadataSource(source MetadataSource, collectorID int) error {
	url := fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, source.ID)

	type MetadataSourceMessage struct {
		Source MetadataSource `json:"source"`
	}

	request := MetadataSourceMessage{
		Source: source,
	}

	_, err := s.Put(url, request, false)

	return err
}
