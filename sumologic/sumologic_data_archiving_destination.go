package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) CreateDataArchivingDestination(dest DataArchivingDestination) (*DataArchivingDestination, error) {
	data, err := s.Post("v1/dataarchiving/destinations", dest)
	if err != nil {
		return nil, err
	}

	var created DataArchivingDestination
	err = json.Unmarshal(data, &created)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (s *Client) GetDataArchivingDestination(id string) (*DataArchivingDestination, error) {
	data, err := s.Get(fmt.Sprintf("v1/dataarchiving/destinations/%s", id))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var dest DataArchivingDestination
	err = json.Unmarshal(data, &dest)
	if err != nil {
		return nil, err
	}

	return &dest, nil
}

func (s *Client) UpdateDataArchivingDestination(dest DataArchivingDestination) error {
	_, err := s.Put(fmt.Sprintf("v1/dataarchiving/destinations/%s", dest.ID), dest)
	return err
}

func (s *Client) DeleteDataArchivingDestination(id string) error {
	_, err := s.Delete(fmt.Sprintf("v1/dataarchiving/destinations/%s", id))
	return err
}

type DataArchivingDestination struct {
	ID                string                         `json:"id,omitempty"`
	DestinationName   string                         `json:"destinationName"`
	DestinationConfig DataArchivingDestinationConfig `json:"destinationConfig"`
	CreatedAt         string                         `json:"createdAt,omitempty"`
	CreatedBy         string                         `json:"createdBy,omitempty"`
	ModifiedAt        string                         `json:"modifiedAt,omitempty"`
	ModifiedBy        string                         `json:"modifiedBy,omitempty"`
}

type DataArchivingDestinationConfig struct {
	DestinationType string `json:"destinationType"`

	// S3 fields
	Description         string                 `json:"description,omitempty"`
	BucketName          string                 `json:"bucketName,omitempty"`
	Region              string                 `json:"region,omitempty"`
	Encrypted           *bool                  `json:"encrypted,omitempty"`
	Enabled             *bool                  `json:"enabled,omitempty"`
	InvalidatedBySystem *bool                  `json:"invalidatedBySystem,omitempty"`
	AuthConfig          *S3ArchivingAuthConfig `json:"authConfig,omitempty"`

	// Syslog fields
	Protocol string `json:"protocol,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Token    string `json:"token,omitempty"`

	// Hitachi and RestAPI shared fields
	URL      string `json:"url,omitempty"`
	ObjectId string `json:"objectId,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type S3ArchivingAuthConfig struct {
	AuthenticationMode string `json:"authenticationMode"`

	// AccessKey fields
	AccessKeyId     string `json:"accessKeyId,omitempty"`
	AccessKeySecret string `json:"accessKeySecret,omitempty"`

	// RoleBased fields
	RoleArn string `json:"roleArn,omitempty"`
}
