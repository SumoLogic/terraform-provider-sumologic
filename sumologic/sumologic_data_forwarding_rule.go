package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) getDataForwardingRule(indexId string) (*DataForwardingRule, error) {
	data, _, err := s.Get(fmt.Sprintf("v1/logsDataForwarding/rules/%s", indexId))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var dataForwardingRule DataForwardingRule
	err = json.Unmarshal(data, &dataForwardingRule)

	if err != nil {
		return nil, err
	}

	return &dataForwardingRule, nil

}

func (s *Client) CreateDataForwardingRule(dataForwardingRule DataForwardingRule) (*DataForwardingRule, error) {
	var createdDataForwardingRule DataForwardingRule

	responseBody, err := s.Post("v1/logsDataForwarding/rules", dataForwardingRule)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &createdDataForwardingRule)

	if err != nil {
		return nil, err
	}

	return &createdDataForwardingRule, nil
}

func (s *Client) UpdateDataForwardingRule(dataForwardingRule DataForwardingRule) error {

	url := fmt.Sprintf("v1/logsDataForwarding/rules/%s", dataForwardingRule.IndexId)
	_, err := s.Put(url, dataForwardingRule)

	return err
}

func (s *Client) DeleteDataForwardingRule(indexId string) error {
	url := fmt.Sprintf("v1/logsDataForwarding/rules/%s", indexId)

	_, err := s.Delete(url)

	return err
}

type DataForwardingRule struct {
	IndexId       string `json:"indexId"`
	DestinationId string `json:"destinationId"`
	Enabled       bool   `json:"enabled"`
	FileFormat    string `json:"fileFormat,omitempty"`
	PayloadSchema string `json:"payloadSchema,omitempty"`
	Format        string `json:"format,omitempty"`
}
