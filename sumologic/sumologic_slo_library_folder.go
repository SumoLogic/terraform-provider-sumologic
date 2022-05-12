package sumologic

import (
	"encoding/json"
	"fmt"
	"log"
)

// ---------- ENDPOINTS ----------

func (s *Client) CreateSLOLibraryFolder(sloLibraryFolder SLOLibraryFolder, paramMap map[string]string) (string, error) {

	urlWithoutParams := SLOBaseApiUrl
	paramString := ""
	sprintfArgs := []interface{}{}

	paramString += "?"

	if val, ok := paramMap["parentId"]; ok {
		queryParam := fmt.Sprintf("parentId=%s&", val)
		paramString += queryParam
	}

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	data, err := s.Post(urlWithParams, sloLibraryFolder)
	if err != nil {
		return "", err
	}
	log.Printf("created SLO folder response: %v", data)

	var createdSLOLibraryFolder SLOLibraryFolder

	err = json.Unmarshal(data, &createdSLOLibraryFolder)
	if err != nil {
		return "", err
	}

	return createdSLOLibraryFolder.ID, nil
}

func (s *Client) GetSLOLibraryFolder(id string) (*SLOLibraryFolder, error) {

	urlWithoutParams := SLOBaseApiUrl + "/%s"
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

	var sloLibraryFolder SLOLibraryFolder

	err = json.Unmarshal(data, &sloLibraryFolder)

	if err != nil {
		return nil, err
	}

	return &sloLibraryFolder, nil
}

func (s *Client) DeleteSLOLibraryFolder(id string) error {
	urlWithoutParams := SLOBaseApiUrl + "/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	_, err := s.Delete(urlWithParams)

	return err
}

func (s *Client) UpdateSLOLibraryFolder(sloLibraryFolder SLOLibraryFolder) error {
	urlWithoutParams := SLOBaseApiUrl + "/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, sloLibraryFolder.ID)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	sloLibraryFolder.ID = ""

	_, err := s.Put(urlWithParams, sloLibraryFolder)

	return err
}

// ---------- TYPES ----------

type SLOLibraryFolder struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     int      `json:"version"`
	CreatedAt   string   `json:"createdAt"`
	CreatedBy   string   `json:"createdBy"`
	ModifiedAt  string   `json:"modifiedAt"`
	ModifiedBy  string   `json:"modifiedBy"`
	ParentID    string   `json:"parentId"`
	ContentType string   `json:"contentType"`
	Type        string   `json:"type"`
	IsSystem    bool     `json:"isSystem"`
	IsMutable   bool     `json:"isMutable"`
	IsLocked    bool     `json:"isLocked"`
	Permissions []string `json:"permissions"`
}

// ---------- END ----------
