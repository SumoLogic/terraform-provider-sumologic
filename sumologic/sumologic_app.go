package sumologic

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func (s *Client) GetAppInstance(id string) (*AppInstance, error) {
	url := fmt.Sprintf("v2/apps/instances/%s", id)
	data, err := s.Get(url)
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

func (s *Client) CreateAppInstance(uuid string, appInstallPayload AppInstallPayload) (string, error) {
	url := fmt.Sprintf("v2/apps/%s/install", uuid)
	response, err := s.Post(url, appInstallPayload)
	if err != nil {
		return "", err
	}

	var jobId AppInstallJobId
	err = json.Unmarshal(response, &jobId)
	if err != nil {
		return "", err
	}

	// Wait for install job to finish
	url = fmt.Sprintf("v2/apps/install/%s/status", jobId.JOBID)
	_, err = waitForJob(url, time.Minute, s)
	if err != nil {
		return "", err
	}

	var appInstallResponse AppInstallResponse
	b, _ := s.Get(url)
	err = json.Unmarshal(b, &appInstallResponse)
	if err != nil {
		return "", err
	}
	log.Printf("[CreateAppInstance] response: %+v\n", appInstallResponse)
	return appInstallResponse.INSTANCEID, nil
}

func (s *Client) DeleteAppInstance(uuid string) error {
	url := fmt.Sprintf("v2/apps/%s/uninstall", uuid)
	response, err := s.Post(url, nil)
	if err != nil {
		return err
	}

	var jobId AppInstallJobId
	err = json.Unmarshal(response, &jobId)
	if err != nil {
		return err
	}

	// Wait for install job to finish
	url = fmt.Sprintf("v2/apps/uninstall/%s/status", jobId.JOBID)
	_, err = waitForJob(url, time.Minute, s)
	return err
}

func (s *Client) UpdateAppInstance(uuid string, appInstallPayload AppInstallPayload) (string, error) {
	url := fmt.Sprintf("v2/apps/%s/upgrade", uuid)
	response, err := s.Post(url, appInstallPayload)
	if err != nil {
		return "", err
	}

	var jobId AppInstallJobId
	err = json.Unmarshal(response, &jobId)
	if err != nil {
		return "", err
	}

	// Wait for install job to finish
	url = fmt.Sprintf("v2/apps/upgrade/%s/status", jobId.JOBID)
	_, err = waitForJob(url, time.Minute, s)
	if err != nil {
		return "", err
	}

	var appInstallResponse AppInstallResponse
	b, _ := s.Get(url)
	err = json.Unmarshal(b, &appInstallResponse)
	if err != nil {
		return "", err
	}
	log.Printf("[UpdateAppInstance] response: %+v\n", appInstallResponse)
	return appInstallResponse.INSTANCEID, nil
}

type AppInstallPayload struct {
	VERSION    string                 `json:"version"`
	PARAMETERS map[string]interface{} `json:"parameters"`
}

type AppInstallJobId struct {
	JOBID string `json:"jobId"`
}

type AppInstallResponse struct {
	INSTANCEID string `json:"instanceId"`
	PATH       string `json:"path"`
	FOLDERID   string `json:"folderId"`
}

type AppInstance struct {
	ID                string `json:"id"`
	UUID              string `json:"uuid"`
	VERSION           string `json:"version"`
	CONFIGURATIONBLOB string `json:"configurationBlob"`
}
