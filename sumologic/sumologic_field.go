package sumologic

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (s *Client) GetField(id string) (*Field, error) {
	urlWithoutParams := "v1/fields/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	data, _, err := s.Get(urlWithParams, false)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var field Field

	err = json.Unmarshal(data, &field)
	if err != nil {
		return nil, err
	}

	return &field, nil

}

func (s *Client) DeleteField(id string) error {
	urlWithoutParams := "v1/fields/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	_, err := s.Delete(urlWithParams)

	return err
}

func (s *Client) CreateField(field Field) (string, error) {
	urlWithoutParams := "v1/fields"

	data, err := s.Post(urlWithoutParams, field, false)
	if err != nil {
		return "", err
	}

	var createdField Field

	err = json.Unmarshal(data, &createdField)
	if err != nil {
		return "", err
	}

	return createdField.FieldId, nil
}

func (s *Client) FindFieldId(name string) (string, error) {
	urlWithoutParams := "v1/fields"

	body, _, err := s.Get(urlWithoutParams, false)
	if err != nil {
		return "", err
	}

	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, body, "", "\t")
	if error != nil {
		return "", error
	}

	var fields FieldsResponse

	err = json.Unmarshal(body, &fields)
	if err != nil {
		return "", err
	}

	for _, f := range fields.Data {
		if f.FieldName == name {
			return f.FieldId, nil
		}
	}

	return "", fmt.Errorf("Field \"%s\" not found", name)
}

type Field struct {
	FieldId   string `json:"fieldId"`
	DataType  string `json:"dataType"`
	State     string `json:"state"`
	FieldName string `json:"fieldName"`
}

type FieldsResponse struct {
	Data []Field `json:"data"`
}
