package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetSourceTemplate(id string) (*SourceTemplate, error) {
	urlWithoutParams := "v1/sourceTemplate/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	data, err := s.Get(urlWithParams)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var sourceTemplate SourceTemplate

	err = json.Unmarshal(data, &sourceTemplate)
	if err != nil {
		return nil, err
	}
	return &sourceTemplate, nil

}

func (s *Client) UpdateSourceTemplate(sourceTemplate SourceTemplate) error {
	urlWithoutParams := "v1/sourceTemplate/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, sourceTemplate.ID)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)
	_, err := s.Post(urlWithParams, sourceTemplate)

	return err
}

func (s *Client) DeleteSourceTemplate(id string) error {
	urlWithoutParams := "v1/sourceTemplate/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	_, err := s.Delete(urlWithParams)

	return err
}

func (s *Client) CreateSourceTemplate(sourceTemplate SourceTemplate) (string, error) {
	urlWithoutParams := "v1/sourceTemplate"

	data, err := s.Post(urlWithoutParams, sourceTemplate)
	if err != nil {
		return "", err
	}

	var response SourceTemplate

	err = json.Unmarshal(data, &response)
	if err != nil {
		return "", err
	}
	return response.ID, nil

}

type SourceTemplate struct {
	TotalCollectorLinked int             `json:"totalCollectorLinked"`
	ID                   string          `json:"id,omitempty"`
	SchemaRef            SchemaRef       `json:"schemaRef"`
	Selector             Selector        `json:"selector,omitempty"`
	Config               string          `json:"config"`
	InputJson            json.RawMessage `json:"inputJson"`
	CreatedAt            string          `json:"createdAt"`
	CreatedBy            string          `json:"createdBy"`
	ModifiedBy           string          `json:"modifiedBy"`
	ModifiedAt           string          `json:"modifiedAt"`
}

type SchemaRef struct {
	Type          string `json:"type"`
	Version       string `json:"version"`
	LatestVersion string `json:"latestVersion"`
}

type Selector struct {
	Tags  [][]OtTag `json:"tags,omitempty"`
	Names []string  `json:"names,omitempty"`
}

type OtTag struct {
	Key    string   `json:"key"`
	Values []string `json:"values,omitempty"`
}
