package sumologic

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func (s *Client) GetContent(id string, timeout time.Duration) (*Content, error) {
	url := fmt.Sprintf("v2/content/%s/export", id)
	log.Printf("[DEBUG] Exporting content with id: %s", id)

	// Begin the content export job
	rawJID, err := s.Post(url, nil, false)
	if err != nil {
		if strings.Contains(err.Error(), "Content with the given ID does not exist.") {
			return nil, nil
		}
		return nil, err
	}

	var jid JobId
	err = json.Unmarshal(rawJID, &jid)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Export job id: %s", jid.ID)

	// Wait for export job to finish
	url = fmt.Sprintf("v2/content/%s/export/%s/status", id, jid.ID)
	_, err = waitForJob(url, timeout, s)
	if err != nil {
		return nil, err
	}

	// Request the results of the job
	var content Content
	url = fmt.Sprintf("v2/content/%s/export/%s/result", id, jid.ID)
	rawContent, _, err := s.Get(url, false)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rawContent, &content)
	if err != nil {
		return nil, err
	}

	// Set the content.Config to be the raw results of the export.
	// This should be a complete, correctly formatted json representation of the content
	content.Config = string(rawContent)
	return &content, nil
}

func (s *Client) DeleteContent(id string, timeout time.Duration) error {
	log.Printf("[DEBUG] Deleting content with id: %s", id)
	url := fmt.Sprintf("v2/content/%s/delete", id)

	rawJID, err := s.Delete(url)
	if err != nil {
		return err
	}

	var jid JobId
	err = json.Unmarshal(rawJID, &jid)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Delete job id: %s", jid.ID)

	url = fmt.Sprintf("v2/content/%s/delete/%s/status", id, jid.ID)
	_, err = waitForJob(url, timeout, s)
	return err
}

func (s *Client) CreateOrUpdateContent(content Content, timeout time.Duration, overwrite bool) (string, error) {
	url := fmt.Sprintf("v2/content/folders/%s/import?overwrite=%s", content.ParentId, strconv.FormatBool(overwrite))
	log.Printf("[DEBUG] Import content in folder=%s, overwrite=%t", content.ParentId, overwrite)

	jobResponse, err := s.PostRawPayload(url, content.Config)
	if err != nil {
		return "", err
	}

	var jid JobId
	err = json.Unmarshal(jobResponse, &jid)
	if err != nil {
		return "", err
	}
	log.Printf("[DEBUG] Import content job id: %s", jid.ID)

	url = fmt.Sprintf("v2/content/folders/%s/import/%s/status", content.ParentId, jid.ID)
	status, err := waitForJob(url, timeout, s)
	if err != nil {
		return "", err
	}

	// extract id of newly created content
	contentId := strings.Split(status.StatusMessage, ":")[1]
	log.Printf("New content id: %s", contentId)
	return contentId, nil
}

func waitForJob(url string, timeout time.Duration, s *Client) (Status, error) {
	conf := &resource.StateChangeConf{
		Pending: []string{
			"InProgress",
		},
		Target: []string{
			"Success",
		},
		Refresh: func() (interface{}, string, error) {
			var status Status
			b, _, err := s.Get(url, false)
			if err != nil {
				return nil, "", err
			}

			err = json.Unmarshal(b, &status)
			if err != nil {
				return nil, "", err
			}

			if status.Status == "failed" {
				return status, status.Status, fmt.Errorf("job failed: %s", status.StatusMessage)
			}

			return status, status.Status, nil
		},
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	result, err := conf.WaitForState()
	log.Printf("[DEBUG] job result: %v", result)
	return result.(Status), err
}
