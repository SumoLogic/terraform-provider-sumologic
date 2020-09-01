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

	data, err := s.Post(urlWithParams, monitorsLibraryFolder)
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

	data, _, err := s.Get(urlWithParams)
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

	_, err := s.Put(urlWithParams, monitorsLibraryFolder)

	return err

}

// ---------- TYPES ----------
type MonitorsLibraryFolder struct {
	ID          string `json:"id,omitempty"`
	Type        string `json:"type"`
	ParentID    string `json:"parentId"`
	Name        string `json:"name"`
	Version     int    `json:"version"`
	CreatedBy   string `json:"createdBy"`
	IsLocked    bool   `json:"isLocked"`
	IsMutable   bool   `json:"isMutable"`
	IsSystem    bool   `json:"isSystem"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	ModifiedAt  string `json:"modifiedAt"`
	ContentType string `json:"contentType"`
	ModifiedBy  string `json:"modifiedBy"`
}

// ---------- END ----------
