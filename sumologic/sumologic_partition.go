package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetPartition(id string) (*Partition, error) {
	data, _, err := s.Get(fmt.Sprintf("partitions/%s", id))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var spartition Partition
	err = json.Unmarshal(data, &spartition)
	if err != nil {
		return nil, err
	}

	return &spartition, nil
}

func (s *Client) CreatePartition(spartition Partition) (*Partition, error) {
	var createdspartition Partition

	responseBody, err := s.Post("partitions", spartition)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &createdspartition)

	if err != nil {
		return nil, err
	}

	return &createdspartition, nil
}

func (s *Client) DeletePartition(id string) error {
	_, err := s.Delete(fmt.Sprintf("partitions/%s/disable", id))

	return err
}

func (s *Client) UpdatePartition(spartition Partition) error {
	url := fmt.Sprintf("partitions/%s", spartition.ID)

	_, err := s.Put(url, spartition)

	return err
}

type Partition struct {
	ID               	string    `json:"id,omitempty"`
	Name             	string    `json:"name"`
	RoutingExpression	string    `json:"routingExpression"`
	AnalyticsTier       string    `json:"analyticsTier"`
	RetentionPeriod  	int       `json:"retentionPeriod"`
	IsCompliant			bool      `json:"isCompliant"`
	DataForwardingId 	string    `json:"dataForwardingId"`
}
