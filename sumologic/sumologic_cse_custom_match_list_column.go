package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSECustomMatchListColumn(id string) (*CSECustomMatchListColumn, error) {
	data, err := s.Get(fmt.Sprintf("sec/v1/custom-match-list-columns/%s", id))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSECustomMatchListColumnResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSECustomMatchListColumn, nil
}

func (s *Client) DeleteCSECustomMatchListColumn(id string) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/custom-match-list-columns/%s", id))

	return err
}

func (s *Client) CreateCSECustomMatchListColumn(CSECustomMatchListColumn CSECustomMatchListColumn) (string, error) {

	request := CSECustomMatchListColumnRequest{
		CSECustomMatchListColumn: CSECustomMatchListColumn,
	}

	var response CSECustomMatchListColumnResponse

	responseBody, err := s.Post("sec/v1/custom-match-list-columns", request)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return "", err
	}

	return response.CSECustomMatchListColumn.ID, nil
}

func (s *Client) UpdateCSECustomMatchListColumn(CSECustomMatchListColumn CSECustomMatchListColumn) error {
	url := fmt.Sprintf("sec/v1/custom-match-list-columns/%s", CSECustomMatchListColumn.ID)
	CSECustomMatchListColumn.ID = ""
	request := CSECustomMatchListColumnRequest{
		CSECustomMatchListColumn: CSECustomMatchListColumn,
	}

	_, err := s.Put(url, request)

	return err
}

type CSECustomMatchListColumnRequest struct {
	CSECustomMatchListColumn CSECustomMatchListColumn `json:"fields"`
}

type CSECustomMatchListColumnResponse struct {
	CSECustomMatchListColumn CSECustomMatchListColumn `json:"data"`
}

type CSECustomMatchListColumn struct {
	ID     string   `json:"id,omitempty"`
	Name   string   `json:"name"`
	Fields []string `json:"fields"`
}
