package sumologic

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const parsersLibraryFolderContentType = "Folder"
const parsersLibraryFolderType = "ParsersLibraryFolder"

// The PUT body uses a different discriminator value than the POST body.
const parsersLibraryFolderUpdateType = "ParsersLibraryFolderUpdate"

// ---------- ENDPOINTS ----------

func (s *Client) CreateParsersLibraryFolder(folder ParsersLibraryFolder) (string, error) {
	urlWithParams := fmt.Sprintf("%s?parentId=%s", parsersLibraryBaseApiUrl, url.QueryEscape(folder.ParentID))

	data, err := s.Post(urlWithParams, folder.toCreateRequest())
	if err != nil {
		return "", err
	}

	var created ParsersLibraryFolder
	if err := json.Unmarshal(data, &created); err != nil {
		return "", err
	}
	return created.ID, nil
}

// GetParsersLibraryFolder reads a folder by id. Returns nil if it does not exist.
// An error is returned if the id refers to an item that is not a folder.
func (s *Client) GetParsersLibraryFolder(id string) (*ParsersLibraryFolder, error) {
	data, err := s.Get(fmt.Sprintf("%s/%s", parsersLibraryBaseApiUrl, id))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var folder ParsersLibraryFolder
	if err := json.Unmarshal(data, &folder); err != nil {
		return nil, err
	}
	if folder.ContentType != parsersLibraryFolderContentType {
		return nil, fmt.Errorf("parsers library item '%s' is not a folder (contentType: %s)", id, folder.ContentType)
	}
	return &folder, nil
}

// UpdateParsersLibraryFolder updates a folder. folder.Version must match the latest
// version of the folder in the library or the API rejects the update.
func (s *Client) UpdateParsersLibraryFolder(folder ParsersLibraryFolder) error {
	_, err := s.Put(fmt.Sprintf("%s/%s", parsersLibraryBaseApiUrl, folder.ID), folder.toUpdateRequest())
	return err
}

func (s *Client) MoveParsersLibraryFolder(id, parentID string) error {
	urlWithParams := fmt.Sprintf("%s/%s/move?parentId=%s", parsersLibraryBaseApiUrl, id, url.QueryEscape(parentID))
	_, err := s.Post(urlWithParams, nil)
	return err
}

func (s *Client) DeleteParsersLibraryFolder(id string) error {
	_, err := s.Delete(fmt.Sprintf("%s/%s", parsersLibraryBaseApiUrl, id))
	return err
}

// ---------- TYPES ----------

type ParsersLibraryFolder struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     int64  `json:"version"`
	CreatedAt   string `json:"createdAt,omitempty"`
	CreatedBy   string `json:"createdBy,omitempty"`
	ModifiedAt  string `json:"modifiedAt,omitempty"`
	ModifiedBy  string `json:"modifiedBy,omitempty"`
	ParentID    string `json:"parentId,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	Type        string `json:"type"`
	IsLocked    bool   `json:"isLocked"`
	IsSystem    bool   `json:"isSystem"`
	IsMutable   bool   `json:"isMutable"`

	// Children is only returned when a folder is read. It lists the direct children of the folder.
	Children []ParsersLibraryFolderChild `json:"children,omitempty"`
}

// ParsersLibraryFolderChild is a direct child (folder or parser) of a parsers library folder.
type ParsersLibraryFolderChild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ContentType string `json:"contentType"`
}

// parsersLibraryFolderCreateRequest is the body of POST /v1/parsers (ParsersLibraryFolder).
type parsersLibraryFolderCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// parsersLibraryFolderUpdateRequest is the body of PUT /v1/parsers/{id} (ParsersLibraryFolderUpdate).
type parsersLibraryFolderUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     int64  `json:"version"`
	Type        string `json:"type"`
}

func (f ParsersLibraryFolder) toCreateRequest() parsersLibraryFolderCreateRequest {
	return parsersLibraryFolderCreateRequest{
		Name:        f.Name,
		Description: f.Description,
		Type:        parsersLibraryFolderType,
	}
}

func (f ParsersLibraryFolder) toUpdateRequest() parsersLibraryFolderUpdateRequest {
	return parsersLibraryFolderUpdateRequest{
		Name:        f.Name,
		Description: f.Description,
		Version:     f.Version,
		Type:        parsersLibraryFolderUpdateType,
	}
}
