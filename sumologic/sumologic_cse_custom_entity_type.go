package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSECustomEntityType(id string) (*CSECustomEntityType, error) {
	data, err := s.Get(fmt.Sprintf("sec/v1/custom-entity-types/%s", id))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSECustomEntityTypeResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSECustomEntityType, nil
}

func (s *Client) DeleteCSECustomEntityType(id string) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/custom-entity-types/%s", id))

	return err
}

func (s *Client) CreateCSECustomEntityType(CSECustomEntityType CSECustomEntityType) (string, error) {

	request := CSECustomEntityTypeRequestPost{
		CSECustomEntityType: CSECustomEntityType,
	}

	var response CSECustomEntityTypeResponse

	responseBody, err := s.Post("sec/v1/custom-entity-types", request)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return "", err
	}

	return response.CSECustomEntityType.ID, nil
}

func (s *Client) UpdateCSECustomEntityType(CSECustomEntityType CSECustomEntityType) error {
	url := fmt.Sprintf("sec/v1/custom-entity-types/%s", CSECustomEntityType.ID)

	request := CSECustomEntityTypeRequestUpdate{
		CSECustomEntityTypeUpdate{
			Fields: CSECustomEntityType.Fields,
			Name:   CSECustomEntityType.Name,
		},
	}

	_, err := s.Put(url, request)

	return err
}

type CSECustomEntityTypeRequestPost struct {
	CSECustomEntityType CSECustomEntityType `json:"fields"`
}

type CSECustomEntityTypeRequestUpdate struct {
	CSECustomEntityTypeUpdate CSECustomEntityTypeUpdate `json:"fields"`
}

type CSECustomEntityTypeResponse struct {
	CSECustomEntityType CSECustomEntityType `json:"data"`
}

type CSECustomEntityType struct {
	ID         string   `json:"id,omitempty"`
	Fields     []string `json:"fields"`
	Name       string   `json:"name"`
	Identifier string   `json:"identifier"`
}

type CSECustomEntityTypeUpdate struct {
	Fields []string `json:"fields"`
	Name   string   `json:"name"`
}
