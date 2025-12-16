package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetOTCollector(id string) (*OTCollector, error) {
	urlWithoutParams := "v1/otCollectors/%s"
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

	var colResponse OTCollector

	err = json.Unmarshal(data, &colResponse)
	if err != nil {
		return nil, err
	}

	return &colResponse, nil
}

func (s *Client) DeleteOTCollector(id string) error {
	urlWithoutParams := "v1/otCollectors/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	_, err := s.Delete(urlWithParams)
	return err
}

func (s *Client) UpdateOTCollector(id string, req map[string]interface{}) (*UpdateOTCollectorResponse, error) {
	url := fmt.Sprintf("v1/otCollectors?id=%s", id)

	data, err := s.Put(url, req)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	if len(data) == 0 {
		return &UpdateOTCollectorResponse{}, nil
	}
	var resp UpdateOTCollectorResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type UpdateOTCollectorResponse struct {
	ID     string                   `json:"id"`
	Errors []map[string]interface{} `json:"errors,omitempty"`
}

type OTCollector struct {
	Category          string                 `json:"category"`
	CreatedAt         string                 `json:"createdAt"`
	CreatedBy         string                 `json:"createdBy"`
	ModifiedBy        string                 `json:"modifiedBy"`
	ModifiedAt        string                 `json:"modifiedAt"`
	TimeZone          string                 `json:"timeZone"`
	IsAlive           bool                   `json:"alive"`
	IsRemotelyManaged bool                   `json:"isRemotelyManaged"`
	ID                string                 `json:"id,omitempty"`
	Name              string                 `json:"name"`
	Ephemeral         bool                   `json:"ephemeral"`
	Description       string                 `json:"description"`
	Tags              map[string]interface{} `json:"tags,omitempty"`
}
