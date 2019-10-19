package sumologic

import (
	"encoding/json"
	"fmt"
	"time"
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
	err = json.Unmarshal(data, &sview)
	if err != nil {
		return nil, err
	}

	return &sview, nil
}

func (s *Client) CreatesPartition(spartition Partition) (*Partition, error) {
	var createdspartition Partition

	responseBody, err := s.Post("partitions", sview)
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

func (s *Client) UpdatePartition(sview ScheduledView) error {
	url := fmt.Sprintf("partitions/%s", sview.ID)

	_, err := s.Put(url, sview)

	return err
}

type ScheduledView struct {
	ID               	string    `json:"id,omitempty"`
	Name             	string    `json:"name"`
	RoutingExpression	string    `json:"routingExpression"`
	AnalyticsTier       string    `json:"analyticsTier"`
	RetentionPeriod  	int       `json:"retentionPeriod"`
	IsCompliant			bool      `json:"isCompliant"`
	DataForwardingId 	string    `json:"dataForwardingId"`
}
