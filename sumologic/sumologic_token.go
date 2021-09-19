package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) CreateToken(token Token) (string, error) {
	urlWithoutParams := "v1/tokens"

	data, err := s.Post(urlWithoutParams, token)
	if err != nil {
		return "", err
	}

	var createdToken Token

	err = json.Unmarshal(data, &createdToken)
	if err != nil {
		return "", err
	}

	return createdToken.ID, nil

}

func (s *Client) GetToken(id string) (*Token, error) {
	urlWithoutParams := "v1/tokens/%s"
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

	var token Token

	err = json.Unmarshal(data, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil

}

func (s *Client) DeleteToken(id string) error {
	urlWithoutParams := "v1/tokens/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	_, err := s.Delete(urlWithParams)

	return err
}

func (s *Client) UpdateToken(token Token) error {
	urlWithoutParams := "v1/tokens/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, token.ID)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	token.ID = ""

	_, err := s.Put(urlWithParams, token)

	return err

}

type Token struct {
	Name               string `json:"name"`
	Status             string `json:"status"`
	ID                 string `json:"id,omitempty"`
	Description        string `json:"description"`
	Type               string `json:"type"`
	Version            int    `json:"version"`
	EncodedTokenAndUrl string `json:"encodedTokenAndUrl"`
}
