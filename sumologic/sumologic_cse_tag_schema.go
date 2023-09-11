package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetCSETagSchema(key string) (*CSETagSchema, error) {
	data, _, err := s.Get(fmt.Sprintf("sec/v1/tag-schemas/%s", key))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response CSETagSchemaResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSETagSchema, nil
}

func (s *Client) DeleteCSETagSchema(id string) error {
	_, err := s.Delete(fmt.Sprintf("sec/v1/tag-schemas/%s", id))

	return err
}

func (s *Client) CreateCSETagSchema(CSECreateUpdateTagSchema CSECreateUpdateTagSchema) (string, error) {

	request := CSETagSchemaRequestCreateUpdate{
		CSECreateUpdateTagSchema,
	}

	var response CSETagSchemaResponse

	responseBody, err := s.Post("sec/v1/tag-schemas", request)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return "", err
	}

	return response.CSETagSchema.Key, nil
}

func (s *Client) UpdateCSETagSchema(CSECreateUpdateTagSchema CSECreateUpdateTagSchema) error {
	url := fmt.Sprintf("sec/v1/tag-schemas/")

	request := CSETagSchemaRequestCreateUpdate{
		CSECreateUpdateTagSchema,
	}

	_, err := s.Put(url, request)

	return err
}

type CSETagSchemaRequestCreateUpdate struct {
	CSETagSchema CSECreateUpdateTagSchema `json:"fields"`
}

type CSETagSchemaResponse struct {
	CSETagSchema CSETagSchema `json:"data"`
}

type CSECreateUpdateTagSchema struct {
	Key          string        `json:"key,omitempty"`
	Label        string        `json:"label,omitempty"`
	ContentTypes []string      `json:"contentTypes"`
	FreeForm     bool          `json:"freeform"`
	ValueOptions []ValueOption `json:"valueOptions,omitempty"`
}

type CSETagSchema struct {
	Key                string        `json:"key,omitempty"`
	Label              string        `json:"label,omitempty"`
	ContentTypes       []string      `json:"contentTypeEnums"`
	FreeForm           bool          `json:"freeform"`
	ValueOptionObjects []ValueOption `json:"valueOptionObjects,omitempty"`
}

type ValueOption struct {
	Value string `json:"value"`
	Label string `json:"label,omitempty"`
	Link  string `json:"link,omitempty"`
}
