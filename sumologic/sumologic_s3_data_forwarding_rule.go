package sumologic

import (
	"encoding/json"
	"fmt"
	"strings"
)

type S3DataForwardingRule struct {
	ID            string `json:"id,omitempty"`
	IndexID       string `json:"indexId,omitempty"`
	DestinationID string `json:"destinationId,omitempty"`
	Enabled       bool   `json:"enabled"`
	FileFormat    string `json:"fileFormat,omitempty"`
	Format        string `json:"format,omitempty"`
}

func (s *Client) GetS3DataForwardingRule(id string) (*S3DataForwardingRule, error) {
	data, _, err := s.Get(fmt.Sprintf("v1/logsDataForwarding/rules/%s", id))
	if err != nil {
		if strings.Contains(err.Error(), "S3DataForwardingRule Not Found") {
			if data == nil {
				return nil, nil
			} else {
				return nil, err
			}
		}
	} else {
		if data == nil {
			return nil, nil
		}
	}

	var dfr S3DataForwardingRule
	err = json.Unmarshal(data, &dfr)
	if err != nil {
		return nil, err
	}

	return &dfr, nil
}

func (s *Client) CreateS3DataForwardingRule(dfr S3DataForwardingRule) (*S3DataForwardingRule, error) {
	var createdDfr S3DataForwardingRule

	responseBody, err := s.Post("v1/logsDataForwarding/rules", dfr)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &createdDfr)

	if err != nil {
		return nil, err
	}

	return &createdDfr, nil
}

func (s *Client) DeleteS3DataForwardingRule(indexId string) error {
	_, err := s.Delete(fmt.Sprintf("v1/logsDataForwarding/rules/%s", indexId))

	return err
}

func (s *Client) UpdateS3DataForwardingRule(dfr S3DataForwardingRule) error {
	url := fmt.Sprintf("v1/logsDataForwarding/rules/%s", dfr.IndexID)
	_, err := s.Put(url, dfr)

	return err
}
