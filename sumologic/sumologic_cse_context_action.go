package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSEContextAction(id string) (*CSEContextAction, error) {
	data, err := s.Get(fmt.Sprintf("sec/v1/context-actions/%s", id))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSEContextActionResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSEContextAction, nil
}

func (s *Client) DeleteCSEContextAction(id string) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/context-actions/%s", id))

	return err
}

func (s *Client) CreateCSEContextAction(CSEContextAction CSEContextAction) (string, error) {

	request := CSEContextAction

	var response CSEContextActionResponse

	responseBody, err := s.Post("sec/v1/context-actions", request)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return "", err
	}

	return response.CSEContextAction.ID, nil
}

func (s *Client) UpdateCSEContextAction(CSEContextAction CSEContextAction) error {
	url := fmt.Sprintf("sec/v1/context-actions/%s", CSEContextAction.ID)

	CSEContextAction.ID = ""

	request := CSEContextActionRequestUpdate{
		CSEContextAction,
	}

	_, err := s.Put(url, request)

	return err
}

type CSEContextActionRequestUpdate struct {
	CSEContextAction CSEContextAction `json:"fields"`
}

type CSEContextActionResponse struct {
	CSEContextAction CSEContextAction `json:"data"`
}

type CSEContextAction struct {
	ID              string   `json:"id,omitempty"`
	Name            string   `json:"name,omitempty"`
	Type            string   `json:"type"`
	Template        string   `json:"template"`
	IocTypes        []string `json:"iocTypes"`
	EntityTypes     []string `json:"entityTypes,omitempty"`
	RecordFields    []string `json:"recordFields,omitempty"`
	AllRecordFields bool     `json:"allRecordFields"`
	Enabled         bool     `json:"enabled"`
}
