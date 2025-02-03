package sumologic

import (
	"encoding/json"
	"fmt"
)

type PollingSource struct {
	Source
	ContentType   string               `json:"contentType"`
	ScanInterval  int                  `json:"scanInterval"`
	Paused        bool                 `json:"paused"`
	URL           string               `json:"url"`
	ThirdPartyRef PollingThirdPartyRef `json:"thirdPartyRef,omitempty"`
}

type PollingThirdPartyRef struct {
	Resources []PollingResource `json:"resources"`
}

type PollingResource struct {
	ServiceType    string                `json:"serviceType"`
	Authentication PollingAuthentication `json:"authentication"`
	Path           PollingPath           `json:"path"`
}

type PollingAuthentication struct {
	Type                    string `json:"type"`
	AwsID                   string `json:"awsId"`
	AwsKey                  string `json:"awsKey"`
	RoleARN                 string `json:"roleARN"`
	Region                  string `json:"region,omitempty"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUrl                 string `json:"auth_uri"`
	TokenUrl                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
	SharedAccessPolicyName  string `json:"sharedAccessPolicyName"`
	SharedAccessPolicyKey   string `json:"sharedAccessPolicyKey"`
}

type PollingPath struct {
	Type                      string                           `json:"type"`
	BucketName                string                           `json:"bucketName,omitempty"`
	PathExpression            string                           `json:"pathExpression,omitempty"`
	LimitToRegions            []string                         `json:"limitToRegions,omitempty"`
	LimitToNamespaces         []string                         `json:"limitToNamespaces,omitempty"`
	LimitToServices           []string                         `json:"limitToServices,omitempty"`
	CustomServices            []string                         `json:"customServices,omitempty"`
	TagFilters                []TagFilter                      `json:"tagFilters,omitempty"`
	SnsTopicOrSubscriptionArn PollingSnsTopicOrSubscriptionArn `json:"snsTopicOrSubscriptionArn,omitempty"`
	UseVersionedApi           *bool                            `json:"useVersionedApi,omitempty"`
	Namespace                 string                           `json:"namespace,omitempty"`
	EventHubName              string                           `json:"eventHubName,omitempty"`
	ConsumerGroup             string                           `json:"consumerGroup,omitempty"`
	Region                    string                           `json:"region,omitempty"`
}

type TagFilter struct {
	Type      string   `json:"type"`
	Namespace string   `json:"namespace"`
	Tags      []string `json:"tags"`
}

type PollingSnsTopicOrSubscriptionArn struct {
	IsSuccess bool   `json:"isSuccess"`
	Arn       string `json:"arn"`
}

func (s *Client) CreatePollingSource(source PollingSource, collectorID int) (int, error) {

	type PollingSourceMessage struct {
		Source PollingSource `json:"source"`
	}

	request := PollingSourceMessage{
		Source: source,
	}

	urlPath := fmt.Sprintf("v1/collectors/%d/sources", collectorID)

	body, err := s.Post(urlPath, request)

	if err != nil {
		return -1, err
	}

	var response PollingSourceMessage
	err = json.Unmarshal(body, &response)

	if err != nil {
		return -1, err
	}

	return response.Source.ID, nil
}

func (s *Client) GetPollingSource(collectorID, sourceID int) (*PollingSource, error) {
	urlPath := fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, sourceID)
	body, _, err := s.Get(urlPath)

	if err != nil {
		return nil, err
	}

	if body == nil {
		return nil, nil
	}

	type PollingSourceResponse struct {
		Source PollingSource `json:"source"`
	}

	var response PollingSourceResponse
	err = json.Unmarshal(body, &response)

	if err != nil {
		return nil, err
	}

	return &response.Source, nil
}

func (s *Client) UpdatePollingSource(source PollingSource, collectorID int) error {
	url := fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, source.ID)

	type PollingSourceMessage struct {
		Source PollingSource `json:"source"`
	}

	request := PollingSourceMessage{
		Source: source,
	}

	_, err := s.Put(url, request)

	return err
}
