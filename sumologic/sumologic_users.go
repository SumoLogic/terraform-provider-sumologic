package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetUser(id string) (*User, error) {
	data, _, err := s.Get(fmt.Sprintf("users/%s", id))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var user User
	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Client) GetUserName(name string) (*User, error) {
	data, _, err := s.Get("users")
	if err != nil {
		return nil, err
	}

	if data == nil {
		return &User{}, nil
	}

	var response UserList
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	for _, c := range response.Users {
		if c.Name == name {
			return &c, nil
		}
	}

	return nil, nil
}

func (s *Client) DeleteUser(id string) error {
	_, err := s.Delete(fmt.Sprintf("users/%s", id))

	return err
}

func (s *Client) CreateUser(user User) (string, error) {

	var createdUser User

	responseBody, err := s.Post("users", user)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &createdUser)

	if err != nil {
		return "", err
	}

	return createdUser.ID, nil
}

func (s *Client) UpdateUser(user User) error {
	url := fmt.Sprintf("user/%s", user.ID)

	_, err := s.Put(url, user)

	return err
}

type UserList struct {
	Users []User `json:"data"`
}

type User struct {
	ID              string   `json:"id,omitempty"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	FilterPredicate string   `json:"filterPredicate"`
	Users           []string ` json:"users"`
	Capabilities    []string ` json:"capabilities"`
}
