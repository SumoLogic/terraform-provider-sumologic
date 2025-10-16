package sumologic

import (
	"fmt"
	"net/url"
)

func (s *Client) DeletePlaybook(name string) error {
	_, err := s.Delete(fmt.Sprintf("api/csoar/v3/playbook/?name=%s", url.QueryEscape(name)))

	return err
}

func (s *Client) UpdatePlaybook(playbook Playbook) error {
	_, err := s.Put("api/csoar/v3/playbook/", playbook)

	return err
}

type Playbook struct {
	Description string                   `json:"description,omitempty"`
	Name        string                   `json:"name"`
	UpdatedName string                   `json:"updated_name,omitempty"`
	Tags        string                   `json:"tags,omitempty"`
	IsDeleted   bool                     `json:"is_deleted,omitempty"`
	Draft       bool                     `json:"draft,omitempty"`
	IsPublished bool                     `json:"is_published,omitempty"`
	LastUpdated int64                    `json:"last_updated,omitempty"`
	Links       []map[string]interface{} `json:"links,omitempty"`
	Nodes       []map[string]interface{} `json:"nodes,omitempty"`
	CreatedBy   int64                    `json:"created_by,omitempty"`
	UpdatedBy   int64                    `json:"updated_by,omitempty"`
	Nested      bool                     `json:"nested,omitempty"`
	Type        string                   `json:"type,omitempty"`
	IsEnabled   bool                     `json:"is_enabled,omitempty"`
}
