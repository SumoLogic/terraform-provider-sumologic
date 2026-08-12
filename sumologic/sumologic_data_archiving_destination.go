package sumologic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (s *Client) CreateDataArchivingDestination(dest DataArchivingDestination) (*DataArchivingDestination, error) {
	body, err := json.Marshal(dest)
	if err != nil {
		return nil, err
	}

	urlPath := "v1/dataarchiving/destinations"
	req, err := s.createSumoRequest(http.MethodPost, urlPath, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := s.doSumoRequest(req)
	if err != nil {
		return nil, err
	}

	d, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] CreateDataArchivingDestination POST %s - Status: %d", urlPath, resp.StatusCode)

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("POST %s failed with status %d: %s", urlPath, resp.StatusCode, string(d))
	}

	var created DataArchivingDestination
	err = json.Unmarshal(d, &created)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (s *Client) GetDataArchivingDestination(id string) (*DataArchivingDestination, error) {
	urlPath := fmt.Sprintf("v1/dataarchiving/destinations/%s", id)
	req, err := s.createSumoRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.doSumoRequest(req)
	if err != nil {
		return nil, err
	}

	d, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] GetDataArchivingDestination GET %s - Status: %d", urlPath, resp.StatusCode)

	if resp.StatusCode == 404 {
		return nil, nil
	}
	if resp.StatusCode == 400 && bytes.Contains(d, []byte("destination_name_not_exists")) {
		return nil, nil
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("GET %s failed with status %d: %s", urlPath, resp.StatusCode, string(d))
	}

	var dest DataArchivingDestination
	err = json.Unmarshal(d, &dest)
	if err != nil {
		return nil, err
	}

	return &dest, nil
}

func (s *Client) UpdateDataArchivingDestination(dest DataArchivingDestination) error {
	urlPath := fmt.Sprintf("v1/dataarchiving/destinations/%s", dest.ID)

	// The update contract accepts neither the id nor, for S3, the bucket name.
	config := dest.DestinationConfig
	config.BucketName = ""

	body, err := json.Marshal(UpdateDataArchivingDestinationRequest{
		DestinationName:   dest.DestinationName,
		DestinationConfig: config,
	})
	if err != nil {
		return err
	}

	req, err := s.createSumoRequest(http.MethodPut, urlPath, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := s.doSumoRequest(req)
	if err != nil {
		return err
	}

	d, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] UpdateDataArchivingDestination PUT %s - Status: %d", urlPath, resp.StatusCode)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("PUT %s failed with status %d: %s", urlPath, resp.StatusCode, string(d))
	}

	return nil
}

func (s *Client) DeleteDataArchivingDestination(id string) error {
	urlPath := fmt.Sprintf("v1/dataarchiving/destinations/%s", id)
	req, err := s.createSumoRequest(http.MethodDelete, urlPath, nil)
	if err != nil {
		return err
	}

	resp, err := s.doSumoRequest(req)
	if err != nil {
		return err
	}

	d, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] DeleteDataArchivingDestination DELETE %s - Status: %d", urlPath, resp.StatusCode)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("DELETE %s failed with status %d: %s", urlPath, resp.StatusCode, string(d))
	}

	return nil
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

type UpdateDataArchivingDestinationRequest struct {
	DestinationName   string                         `json:"destinationName"`
	DestinationConfig DataArchivingDestinationConfig `json:"destinationConfig"`
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
