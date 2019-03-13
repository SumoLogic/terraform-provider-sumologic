package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetRole(id string) (*Role, error) {
	data, _, err := s.Get(fmt.Sprintf("roles/%s", id))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var role Role
	err = json.Unmarshal(data, &role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (s *Client) GetRoleName(name string) (*Role, error) {
	data, _, err := s.Get("roles")
	if err != nil {
		return nil, err
	}

	if data == nil {
		return &Role{}, nil
	}

	var response RoleList
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	for _, c := range response.Roles {
		if c.Name == name {
			return &c, nil
		}
	}

	return nil, nil
}

func (s *Client) DeleteRole(id string) error {
	_, err := s.Delete(fmt.Sprintf("roles/%s", id))

	return err
}

func (s *Client) CreateRole(role Role) (string, error) {

	var createdRole Role

	responseBody, err := s.Post("roles", role)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &createdRole)

	if err != nil {
		return "", err
	}

	return createdRole.ID, nil
}

func (s *Client) UpdateRole(role Role) error {
	url := fmt.Sprintf("roles/%s", role.ID)

	_, err := s.Put(url, role)

	return err
}

type RoleList struct {
	Roles []Role `json:"data"`
}

type Role struct {
	ID              string   `json:"id,omitempty"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	FilterPredicate string   `json:"filterPredicate"`
	Users           []string ` json:"users"`
	Capabilities    []string ` json:"capabilities"`
}
