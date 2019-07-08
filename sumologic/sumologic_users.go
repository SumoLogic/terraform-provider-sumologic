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
	url := fmt.Sprintf("users/%s", user.ID)

	user.ID = ""
	user.Email = ""
	_, err := s.Put(url, user)

	return err
}

type UserList struct {
	Users []User `json:"data"`
}

type User struct {
	ID        string   `json:"id,omitempty"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Email     string   `json:"email,omitempty"`
	Active    bool     `json:"isActive"`
	RoleIds   []string `json:"roleIds"`
}
