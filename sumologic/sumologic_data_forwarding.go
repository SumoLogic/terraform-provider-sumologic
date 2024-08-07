package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) getDataForwarding(id string) (*DataForwarding, error) {

	data, _, err := s.Get(fmt.Sprintf("v1/logsDataForwarding/destinations/%s", id))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var dataForwarding DataForwarding
	err = json.Unmarshal(data, &dataForwarding)

	if err != nil {
		return nil, err
	}

	return &dataForwarding, nil

}
func (s *Client) CreateDataForwarding(dataForwarding DataForwarding) (*DataForwarding, error) {
	var createdDataForwarding DataForwarding

	responseBody, err := s.Post("v1/logsDataForwarding/destinations", dataForwarding)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &createdDataForwarding)

	if err != nil {
		return nil, err
	}

	return &createdDataForwarding, nil
}

func (s *Client) UpdateDataForwarding(dataForwarding DataForwarding) error {

	url := fmt.Sprintf("v1/logsDataForwarding/destinations/%s", dataForwarding.ID)
	_, err := s.Put(url, dataForwarding)

	return err
}
func (s *Client) DeleteDataForwarding(id string) error {
	url := fmt.Sprintf("v1/logsDataForwarding/destinations/%s", id)

	_, err := s.Delete(url)

	return err
}

type DataForwarding struct {
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
