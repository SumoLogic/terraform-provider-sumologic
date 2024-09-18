package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) getDataForwardingDestination(id string) (*DataForwardingDestination, error) {

	data, _, err := s.Get(fmt.Sprintf("v1/logsDataForwarding/destinations/%s", id))
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
func (s *Client) CreateDataForwardingDestination(dataForwardingDestination DataForwardingDestination) (*DataForwardingDestination, error) {
	var createdDataForwardingDestination DataForwardingDestination

	responseBody, err := s.Post("v1/logsDataForwarding/destinations", dataForwardingDestination)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &createdDataForwardingDestination)

	if err != nil {
		return nil, err
	}

	return &createdDataForwardingDestination, nil
}

func (s *Client) UpdateDataForwardingDestination(dataForwardingDestination DataForwardingDestination) error {

	url := fmt.Sprintf("v1/logsDataForwarding/destinations/%s", dataForwardingDestination.ID)
	_, err := s.Put(url, dataForwardingDestination)

	return err
}
func (s *Client) DeleteDataForwardingDestination(id string) error {
	url := fmt.Sprintf("v1/logsDataForwarding/destinations/%s", id)

	_, err := s.Delete(url)

	return err
}

type DataForwardingDestination struct {
	ID                     string `json:"id,omitempty"`
	DestinationName        string `json:"destinationName"`
	Description            string `json:"description,omitempty"`
	BucketName             string `json:"bucketName"`
	AccessMethod           string `json:"authenticationMode"`
	AccessKey              string `json:"accessKeyId,omitempty"`
	SecretKey              string `json:"secretAccessKey,omitempty"`
	RoleArn                string `json:"roleArn,omitempty"`
	S3Region               string `json:"region,omitempty"`
	S3ServerSideEncryption bool   `json:"encrypted"`
}
