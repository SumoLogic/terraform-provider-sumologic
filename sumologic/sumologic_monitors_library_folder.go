package sumologic

import (
	"encoding/json"
	"fmt"
	"log"
)

// ---------- ENDPOINTS ----------

func (s *Client) CreateMonitorsLibraryFolder(monitorsLibraryFolder MonitorsLibraryFolder, paramMap map[string]string) (string, error) {
	urlWithoutParams := "v1/monitors"
	paramString := ""
	sprintfArgs := []interface{}{}

	paramString += "?"

	if val, ok := paramMap["parentId"]; ok {
		queryParam := fmt.Sprintf("parentId=%s&", val)
		paramString += queryParam
	}

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	data, err := s.Post(urlWithParams, monitorsLibraryFolder, false)
	if err != nil {
		return "", err
	}
	log.Printf("created monitor response: %v", data)

	var createdMonitorsLibraryFolder MonitorsLibraryFolder

	err = json.Unmarshal(data, &createdMonitorsLibraryFolder)
	if err != nil {
		return "", err
	}

	return createdMonitorsLibraryFolder.ID, nil

}

func (s *Client) GetMonitorsLibraryFolder(id string) (*MonitorsLibraryFolder, error) {
	urlWithoutParams := "v1/monitors/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	data, _, err := s.Get(urlWithParams, false)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var monitorsLibraryFolder MonitorsLibraryFolder

	err = json.Unmarshal(data, &monitorsLibraryFolder)

	if err != nil {
		return nil, err
	}

	return &monitorsLibraryFolder, nil

}

func (s *Client) DeleteMonitorsLibraryFolder(id string) error {
	urlWithoutParams := "v1/monitors/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	_, err := s.Delete(urlWithParams)

	return err
}

func (s *Client) UpdateMonitorsLibraryFolder(monitorsLibraryFolder MonitorsLibraryFolder) error {
	urlWithoutParams := "v1/monitors/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, monitorsLibraryFolder.ID)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	monitorsLibraryFolder.ID = ""

	_, err := s.Put(urlWithParams, monitorsLibraryFolder, false)

	return err

}

// ---------- TYPES ----------
type MonitorsLibraryFolder struct {
	ID          string `json:"id,omitempty"`
	Type        string `json:"type"`
	ContentType string `json:"contentType"`
	ParentID    string `json:"parentId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedBy   string `json:"createdBy"`
	CreatedAt   string `json:"createdAt"`
	ModifiedBy  string `json:"modifiedBy"`
	ModifiedAt  string `json:"modifiedAt"`
	IsLocked    bool   `json:"isLocked"`
	IsMutable   bool   `json:"isMutable"`
	IsSystem    bool   `json:"isSystem"`
	Version     int    `json:"version"`
}

// ---------- END ----------
