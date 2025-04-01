package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetRoleV2(id string) (*RoleV2, error) {
	urlWithoutParams := "v2/roles/%s"
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

	var roleV2 RoleV2

	err = json.Unmarshal(data, &roleV2)
	if err != nil {
		return nil, err
	}

	return &roleV2, nil

}

func (s *Client) DeleteRoleV2(id string) error {
	urlWithoutParams := "v2/roles/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	_, err := s.Delete(urlWithParams)

	return err
}

func (s *Client) UpdateRoleV2(roleV2 RoleV2) error {
	urlWithoutParams := "v2/roles/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, roleV2.ID)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	roleV2.ID = ""

	_, err := s.Put(urlWithParams, roleV2)

	return err

}

func (s *Client) CreateRoleV2(roleV2 RoleV2) (string, error) {
	urlWithoutParams := "v2/roles"

	data, err := s.Post(urlWithoutParams, roleV2)
	if err != nil {
		return "", err
	}

	var createdRoleV2 RoleV2

	err = json.Unmarshal(data, &createdRoleV2)
	if err != nil {
		return "", err
	}

	return createdRoleV2.ID, nil

}

type RoleV2 struct {
	Capabilities       []string               `json:"capabilities,omitempty"`
	SecurityDataFilter string                 `json:"securityDataFilter"`
	Name               string                 `json:"name"`
	AuditDataFilter    string                 `json:"auditDataFilter"`
	ID                 string                 `json:"id,omitempty"`
	LogAnalyticsFilter string                 `json:"logAnalyticsFilter"`
	SelectionType      string                 `json:"selectionType,omitempty"`
	SelectedViews      []ViewFilterDefinition `json:"selectedViews,omitempty"`
	Description        string                 `json:"description"`
	Users              []string               `json:"users"`
}

type ViewFilterDefinition struct {
	ViewName   string `json:"viewName"`
	ViewFilter string `json:"viewFilter"`
}
