package sumologic

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func (s *Client) CreateServiceAccount(account ServiceAccount) (string, error) {
	data, err := s.Post("v1/serviceAccounts", account)
	if err != nil {
		return "", err
	}

	var created ServiceAccount
	if err := json.Unmarshal(data, &created); err != nil {
		return "", err
	}
	return created.ID, nil
}

// GetServiceAccount returns nil when the service account does not exist.
func (s *Client) GetServiceAccount(id string) (*ServiceAccount, error) {
	data, err := s.Get(fmt.Sprintf("v1/serviceAccounts/%s", url.PathEscape(id)))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var account ServiceAccount
	if err := json.Unmarshal(data, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

func (s *Client) ListServiceAccounts() ([]ServiceAccount, error) {
	data, err := s.Get("v1/serviceAccounts")
	if err != nil {
		return nil, err
	}

	var list struct {
		Data []ServiceAccount `json:"data"`
	}
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list.Data, nil
}

func (s *Client) UpdateServiceAccount(id string, update ServiceAccountUpdate) error {
	_, err := s.Put(fmt.Sprintf("v1/serviceAccounts/%s", url.PathEscape(id)), update)
	return err
}

// DeleteServiceAccount deletes the service account and its access keys. Only the options that are set are sent.
func (s *Client) DeleteServiceAccount(id string, transferTo string, deleteContent bool) error {
	query := url.Values{}
	if transferTo != "" {
		query.Set("transferTo", transferTo)
	}
	if deleteContent {
		query.Set("deleteContent", "true")
	}

	path := fmt.Sprintf("v1/serviceAccounts/%s", url.PathEscape(id))
	if len(query) > 0 {
		path += "?" + query.Encode()
	}
	_, err := s.Delete(path)
	return err
}

// isNotFoundError reports whether an API error body says the service account does not exist. It matches the error
// codes exactly, because other errors (such as an unknown transfer target) also mention a missing object.
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return HasErrorCode(err.Error(), []string{"serviceaccount:service_account_not_found", "user:user_not_found"}) != ""
}

// models
type ServiceAccount struct {
	ID string `json:"id,omitempty"`
	// Name of the service account.
	Name string `json:"name"`
	// Email address of the service account.
	Email string `json:"email"`
	// List of roleIds associated with the service account.
	RoleIds []string `json:"roleIds"`
	// True if the service account is active. Read-only on create.
	IsActive bool `json:"isActive,omitempty"`
}

// ServiceAccountUpdate holds the fields sent on update. Every field is always sent.
type ServiceAccountUpdate struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	RoleIds  []string `json:"roleIds"`
	IsActive bool     `json:"isActive"`
}
