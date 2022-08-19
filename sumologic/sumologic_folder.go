package sumologic

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// READ
func (s *Client) GetFolder(id string) (*Folder, error) {
	url := fmt.Sprintf("v2/content/folders/%s", id)
	rawFolder, _, err := s.Get(url)
	if err != nil {
		return nil, err
	}
	if rawFolder == nil {
		return nil, nil
	}

	var folder Folder
	err = json.Unmarshal(rawFolder, &folder)
	if err != nil {
		return nil, err
	}

	return &folder, nil
}

func (s *Client) DeleteFolder(id string, timeout time.Duration) error {
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
	log.Printf("[DEBUG] Delete folder job id: %s", jid.ID)

	url = fmt.Sprintf("v2/content/%s/delete/%s/status", id, jid.ID)
	_, err = waitForJob(url, timeout, s)
	return err
}

func (s *Client) CreateFolder(folder Folder) (string, error) {
	url := "v2/content/folders"
	responseData, err := s.Post(url, folder)
	if err != nil {
		return "", err
	}

	var folderResponse Folder
	err = json.Unmarshal(responseData, &folderResponse)
	if err != nil {
		return "", err
	}

	log.Printf("New folder id: %s", folderResponse.ID)
	return folderResponse.ID, nil
}

func (s *Client) UpdateFolder(folder Folder) error {
	url := fmt.Sprintf("v2/content/folders/%s", folder.ID)
	_, err := s.Put(url, folder)
	return err
}

func (s *Client) getPersonalFolder() (*Folder, error) {
	url := "v2/content/folders/personal"
	rawFolder, _, err := s.Get(url)
	if err != nil {
		return nil, err
	}

	var personalFolder Folder
	err = json.Unmarshal(rawFolder, &personalFolder)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] personal folder id: %s", personalFolder.ID)
	return &personalFolder, nil
}

func (s *Client) getAdminRecommendedFolder(timeout time.Duration) (*Folder, error) {
	url := "v2/content/folders/adminRecommended"
	rawJID, _, err := s.Get(url)
	if err != nil {
		return nil, err
	}

	var jid JobId
	err = json.Unmarshal(rawJID, &jid)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Admin Recommended folder job id: %s", jid.ID)

	url = fmt.Sprintf("v2/content/folders/adminRecommended/%s/status", jid.ID)
	_, err = waitForJob(url, timeout, s)
	if err != nil {
		return nil, err
	}

	url = fmt.Sprintf("v2/content/folders/adminRecommended/%s/result", jid.ID)
	rawContent, _, err := s.Get(url)
	if err != nil {
		return nil, err
	}

	var adminRecommendedFolder Folder
	err = json.Unmarshal(rawContent, &adminRecommendedFolder)
	if err != nil {
		return nil, err
	}
	return &adminRecommendedFolder, nil
}
