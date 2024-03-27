package sumologic

import (
	"encoding/json"
	"fmt"
	"strings"
)

type S3DataForwardingDestination struct {
	ID                 string `json:"id,omitempty"`
	Name               string `json:"destinationName,omitempty"`
	Description        string `json:"description,omitempty"`
	AuthenticationMode string `json:"authenticationMode"`
	AccessKeyID        string `json:"accessKeyId,omitempty"`
	SecretAccessKey    string `json:"secretAccessKey,omitempty"`
	RoleARN            string `json:"roleArn,omitempty"`
	Region             string `json:"region,omitempty"`
	Encrypted          bool   `json:"encrypted"`
	Enabled            bool   `json:"enabled"`
	BucketName         string `json:"bucketName,omitempty"`
}

func (s *Client) GetS3DataForwardingDestination(id string) (*S3DataForwardingDestination, error) {
	data, _, err := s.Get(fmt.Sprintf("v1/logsDataForwarding/destinations/%s", id))

	if err != nil {
		if strings.Contains(err.Error(), "gn:destination_name_not_exists") {
			return nil, nil
		}
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var dfd S3DataForwardingDestination
	err = json.Unmarshal(data, &dfd)
	if err != nil {
		return nil, err
	}

	return &dfd, nil
}

func (s *Client) CreateS3DataForwardingDestination(dfd S3DataForwardingDestination) (*S3DataForwardingDestination, error) {
	var createdDfd S3DataForwardingDestination

	responseBody, err := s.Post("v1/logsDataForwarding/destinations", dfd)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &createdDfd)

	if err != nil {
		return nil, err
	}

	return &createdDfd, nil
}

func (s *Client) DeleteS3DataForwardingDestination(id string) error {
	_, err := s.Delete(fmt.Sprintf("v1/logsDataForwarding/destinations/%s", id))

	return err
}

func (s *Client) UpdateS3DataForwardingDestination(dfd S3DataForwardingDestination) error {
	url := fmt.Sprintf("v1/logsDataForwarding/destinations/%s", dfd.ID)
	_, err := s.Put(url, dfd)

	return err
}
