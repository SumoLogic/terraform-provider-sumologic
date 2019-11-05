package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) CreateFolder(folder FolderCreate) (Folder, error) {
	response, err := s.Post("v2/content/folders", folder)
	if err != nil {
		return Folder{}, err
	}
	var folderResponse Folder
	err = json.Unmarshal(response, &folderResponse)
	if err != nil {
		return Folder{}, err
	}
	return folderResponse, nil
}

func (s *Client) GetFolder(id string) (Folder, error) {
	var folderResponse Folder
	response, _, err := s.Get(fmt.Sprintf("v2/content/folders/%s", id))
	if err != nil {
		return Folder{}, err
	}

	err = json.Unmarshal(response, &folderResponse)
	if err != nil {
		return Folder{}, err
	}
	return folderResponse, nil
}

func (s *Client) UpdateFolder(id string, folder FolderUpdate) (Folder, error) {
	var folderResponse Folder
	response, err := s.Put(fmt.Sprintf("v2/content/folders/%s", id), folder)
	if err != nil {
		return Folder{}, err
	}

	err = json.Unmarshal(response, &folderResponse)
	if err != nil {
		return Folder{}, err
	}
	return folderResponse, nil
}

func (s *Client) StartDeleteFolder(id string) (string, error) {
	var deletionJob BeginDeletionJobResponse
	response, err := s.Delete(fmt.Sprintf("v2/content/%s/delete", id))
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(response, &deletionJob)
	if err != nil {
		return "", err
	}
	return deletionJob.ID, nil
}

func (s *Client) DeleteFolderStatus(id string, job_id string) (string, error) {
	var deletionStatus DeletionJobStatus
	response, _, err := s.Get(fmt.Sprintf("v2/content/%s/delete/%s/status", id, job_id))
	err = json.Unmarshal(response, &deletionStatus)
	if err != nil {
		return "", err
	}
	return deletionStatus.Status, nil
}

type Folder struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ParentId    string     `json:"parentId"`
	ItemType    string     `json:"itemType"`
	Permissions []string   `json:"permissions"`
	CreatedAt   string     `json:"createdAt"`
	CreatedBy   string     `json:"createdBy"`
	ModifiedAt  string     `json:"modifiedAt"`
	ModifiedBy  string     `json:"modifiedBy"`
	ID          string     `json:"id"`
	Children    []Children `json:"children"`
}

type Children struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ParentId    string   `json:"parentId"`
	ItemType    string   `json:"itemType"`
	Permissions []string `json:"permissions"`
	CreatedAt   string   `json:"createdAt"`
	CreatedBy   string   `json:"createdBy"`
	ModifiedAt  string   `json:"modifiedAt"`
	ModifiedBy  string   `json:"modifiedBy"`
	ID          string   `json:"id"`
}

type FolderCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentId    string `json:"parentId"`
}

type FolderUpdate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type BeginDeletionJobResponse struct {
	ID string `json:"id"`
}

type DeletionJobStatus struct {
	Status        string           `json:"status"`
	StatusMessage string           `json:"statusMessage"`
	Error         ErrorDescription `json:"error"`
}

type ErrorDescription struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}
