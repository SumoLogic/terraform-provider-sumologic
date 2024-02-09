package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSEMatchList(id string, getCustomTargetColumnName bool) (*CSEMatchListGet, error) {
	data, _, err := s.Get(fmt.Sprintf("sec/v1/match-lists/%s?getTargetCustomColumnName=%t", id, getCustomTargetColumnName))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSEMatchListResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSEMatchListGet, nil
}

func (s *Client) DeleteCSEMatchList(id string) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/match-lists/%s", id))

	return err
}

func (s *Client) CreateCSEMatchList(CSEMatchListPost CSEMatchListPost) (string, error) {

	request := CSEMatchListRequestPost{
		CSEMatchListPost: CSEMatchListPost,
	}

	var response CSEMatchListResponse

	responseBody, err := s.Post("sec/v1/match-lists", request)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return "", err
	}

	return response.CSEMatchListGet.ID, nil
}

func (s *Client) UpdateCSEMatchList(CSEMatchListPost CSEMatchListPost) error {
	url := fmt.Sprintf("sec/v1/match-lists/%s", CSEMatchListPost.ID)

	request := CSEMatchListRequestUpdate{
		CSEMatchListUpdate{
			Active:      CSEMatchListPost.Active,
			DefaultTtl:  CSEMatchListPost.DefaultTtl,
			Description: CSEMatchListPost.Description,
		},
	}

	_, err := s.Put(url, request)

	return err
}

type CSEMatchListRequestPost struct {
	CSEMatchListPost CSEMatchListPost `json:"fields"`
}

type CSEMatchListRequestUpdate struct {
	CSEMatchListUpdate CSEMatchListUpdate `json:"fields"`
}

type CSEMatchListResponse struct {
	CSEMatchListGet CSEMatchListGet `json:"data"`
}

type CSEMatchListPost struct {
	ID           string `json:"id,omitempty"`
	Active       bool   `json:"active,omitempty"`
	DefaultTtl   int    `json:"defaultTtl,omitempty"`
	Description  string `json:"description,omitempty"`
	Name         string `json:"name,omitempty"`
	TargetColumn string `json:"targetColumn"`
}

type CSEMatchListGet struct {
	ID            string `json:"id,omitempty"`
	Active        bool   `json:"active,omitempty"`
	DefaultTtl    int    `json:"defaultTtl,omitempty"`
	Description   string `json:"description,omitempty"`
	Name          string `json:"name,omitempty"`
	TargetColumn  string `json:"targetColumn,omitempty"`
	Created       string `json:"created,omitempty"`
	CreatedBy     string `json:"createdBy,omitempty"`
	LastUpdated   string `json:"lastUpdated,omitempty"`
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
}

type CSEMatchListUpdate struct {
	Active      bool   `json:"active,omitempty"`
	DefaultTtl  int    `json:"defaultTtl,omitempty"`
	Description string `json:"description,omitempty"`
}
