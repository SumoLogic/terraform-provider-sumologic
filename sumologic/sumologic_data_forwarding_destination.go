package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) CreateDataForwardingDestination(dataForwardingDestination DataForwardingDestination) (string, error) {
	urlWithoutParams := "v1/logsDataForwarding/destinations"

	data, err := s.Post(urlWithoutParams, dataForwardingDestination)
	if err != nil {
		return "", err
	}

	var createdDataForwardingDestination DataForwardingDestination

	err = json.Unmarshal(data, &createdDataForwardingDestination)
	if err != nil {
		return "", err
	}

	return createdDataForwardingDestination.ID, nil

}

func (s *Client) GetDataForwardingDestination(id string) (*DataForwardingDestination, error) {
	urlWithoutParams := "v1/logsDataForwarding/destinations/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	data, _, err := s.Get(urlWithParams)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var dataForwardingDestination DataForwardingDestination

	err = json.Unmarshal(data, &dataForwardingDestination)
	if err != nil {
		return nil, err
	}

	return &dataForwardingDestination, nil

}

func (s *Client) DeleteDataForwardingDestination(id string) error {
	urlWithoutParams := "v1/logsDataForwarding/destinations/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	_, err := s.Delete(urlWithParams)

	return err
}

func (s *Client) UpdateDataForwardingDestination(dataForwardingDestination DataForwardingDestination) error {
	urlWithoutParams := "v1/logsDataForwarding/destinations/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, dataForwardingDestination.ID)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	dataForwardingDestination.ID = ""

	_, err := s.Put(urlWithParams, dataForwardingDestination)

	return err

}

type DataForwardingDestination struct {
	AuthenticationMode string `json:"authenticationMode"`
	RoleArn            string `json:"roleArn"`
	Encrypted          bool   `json:"encrypted"`
	ID                 string `json:"id,omitempty"`
	AccessKeyId        string `json:"accessKeyId"`
	SecretAccessKey    string `json:"secretAccessKey"`
	Region             string `json:"region"`
	Enabled            bool   `json:"enabled"`
	DestinationName    string `json:"destinationName"`
	Description        string `json:"description"`
	BucketName         string `json:"bucketName"`
}
