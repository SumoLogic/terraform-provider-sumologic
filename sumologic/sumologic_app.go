package sumologic

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func (s *Client) GetAppInstance(id string) (*AppInstance, error) {
	url := fmt.Sprintf("v2/apps/instances/%s", id)
	data, _, err := s.Get(url)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var appInstance AppInstance
	err = json.Unmarshal(data, &appInstance)
	if err != nil {
		return nil, err
	}
	log.Printf("[GetAppInstance] response: %+v\n", appInstance)
	return &appInstance, nil
}

func (s *Client) CreateAppInstance(uuid string, appInstallPayload AppInstallPayload) (*AppInstallResponse, error) {
	url := fmt.Sprintf("v2/apps/%s/install", uuid)
	jobId, err := s.Post(url, appInstallPayload)
	if err != nil {
		return nil, err
	}

	// Wait for install job to finish
	url = fmt.Sprintf("v2/apps/install/%s/status", jobId)
	_, err = waitForJob(url, time.Minute, s)
	if err != nil {
		return nil, err
	}

	var appInstallResponse AppInstallResponse
	b, _, _ := s.Get(url)
	err = json.Unmarshal(b, &appInstallResponse)
	if err != nil {
		return nil, err
	}
	log.Printf("[CreateAppInstance] response: %+v\n", appInstallResponse)
	return &appInstallResponse, nil
}

func (s *Client) DeleteAppInstance(uuid string) error {
	url := fmt.Sprintf("v2/apps/%s/uninstall", uuid)
	jobId, err := s.Post(url, nil)
	if err != nil {
		return err
	}

	// Wait for install job to finish
	url = fmt.Sprintf("v2/apps/uninstall/%s/status", jobId)
	_, err = waitForJob(url, time.Minute, s)
	return err
}

func (s *Client) UpdateAppInstance(uuid string, appInstallPayload AppInstallPayload) (*AppInstallResponse, error) {
	url := fmt.Sprintf("v2/apps/%s/upgrade", uuid)
	jobId, err := s.Post(url, appInstallPayload)
	if err != nil {
		return nil, err
	}

	// Wait for install job to finish
	url = fmt.Sprintf("v2/apps/upgrade/%s/status", jobId)
	_, err = waitForJob(url, time.Minute, s)
	if err != nil {
		return nil, err
	}

	var appInstallResponse AppInstallResponse
	b, _, _ := s.Get(url)
	err = json.Unmarshal(b, &appInstallResponse)
	if err != nil {
		return nil, err
	}
	log.Printf("[UpdateAppInstance] response: %+v\n", appInstallResponse)
	return &appInstallResponse, nil
}

type AppInstallPayload struct {
	VERSION    string            `json:"version"`
	PARAMETERS map[string]string `json:"parameters"`
}

type AppInstallResponse struct {
	INSTANCEID string `json:"instanceId"`
	PATH       string `json:"path"`
	FOLDERID   string `json:"folderId"`
}

type AppInstance struct {
	ID                string          `json:"id"`
	UUID              string          `json:"uuid"`
	VERSION           string          `json:"version"`
	NAME              string          `json:"name"`
	DESCRIPTION       string          `json:"description"`
	CONFIGURATIONBLOB string          `json:"configurationBlob"`
	PREVIOUSVERSION   string          `json:"previousVersion"`
	LATESTVERSION     string          `json:"latestVersion"`
	PATH              string          `json:"path"`
	MANAGEDOBJECTS    []ManagedObject `json:"managedObjects"`
	FOLDERID          string          `json:"folderId"`
	CREATEDAT         string          `json:"createdAt"`
	CREATEDBY         string          `json:"createdBy"`
	MODIFIEDAT        string          `json:"modifiedAt"`
	MODIFIEDBY        string          `json:"modifiedBy"`
}

type ManagedObject struct {
	ID          string `json:"id"`
	NAME        string `json:"name"`
	TYPE        string `json:"type"`
	DESCRIPTION string `json:"description"`
}
