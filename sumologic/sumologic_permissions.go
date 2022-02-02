package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetPermissions(id string) (*PermissionsResponse, error) {
	url := fmt.Sprintf("v2/content/%s/permissions", id)
	data, _, err := s.Get(url)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var contentPermssionsResponse PermissionsResponse
	err = json.Unmarshal(data, &contentPermssionsResponse)
	if err != nil {
		return nil, err
	}
	return &contentPermssionsResponse, nil
}

func (s *Client) UpdatePermissions(contentPermissionsRequest PermissionsRequest, id string) (string, error) {
	url := fmt.Sprintf("v2/content/%s/permissions/add", id)
	_, err := s.Put(url, contentPermissionsRequest)
	return id, err
}

func (s *Client) DeletePermissions(contentPermissionsRequest PermissionsRequest, id string) error {
	url := fmt.Sprintf("v2/content/%s/permissions/remove", id)
	_, err := s.Put(url, contentPermissionsRequest)
	return err
}

func (s *Client) GetContentPath(id string) (string, error) {
	url := fmt.Sprintf("v2/content/%s/path", id)
	data, _, err := s.Get(url)
	if err != nil {
		return "", err
	}
	if data == nil {
		return "", nil
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(data, &m)
	if err != nil {
		return "", err
	}
	return m["path"].(string), nil
}

func (s *Client) GetCreatorId(path string) (string, error) {
	url := fmt.Sprintf("v2/content/path?path=%s", path)
	data, _, err := s.Get(url)
	if err != nil {
		return "", err
	}
	if data == nil {
		return "", nil
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(data, &m)
	if err != nil {
		return "", err
	}
	return m["createdBy"].(string), nil
}

type PermissionsResponse struct {
	ExplicitPermissions []Permission `json:"explicitPermissions"`
	ImplicitPermissions []Permission `json:"implicitPermissions"` // ignored for now.
}

type PermissionsRequest struct {
	PermissionAssignmentype []Permission `json:"contentPermissionAssignments"`
	NotifyRecipients        bool         `json:"notifyRecipients"`
	NotificationMessage     string       `json:"notificationMessage"`
}

type Permission struct {
	PermissionName string `json:"permissionName"`
	SourceType     string `json:"sourceType"`
	SourceId       string `json:"sourceId"`
	ContentId      string `json:"contentId"`
}