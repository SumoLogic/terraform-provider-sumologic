package sumologic

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (s *Client) GetPartition(id string) (*Partition, error) {
	data, _, err := s.Get(fmt.Sprintf("v1/partitions/%s", id), false)
	if err != nil {
		if strings.Contains(err.Error(), "Partition Not Found") {
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

	var spartition Partition
	err = json.Unmarshal(data, &spartition)
	if err != nil {
		return nil, err
	} else if spartition.IsActive == false {
		return nil, nil
	}

	return &spartition, nil
}

func (s *Client) CreatePartition(spartition Partition) (*Partition, error) {
	var createdspartition Partition

	responseBody, err := s.Post("v1/partitions", spartition, false)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &createdspartition)

	if err != nil {
		return nil, err
	}

	return &createdspartition, nil
}

func (s *Client) DecommissionPartition(id string) error {
	_, err := s.Post(fmt.Sprintf("v1/partitions/%s/decommission", id), nil, false)

	return err
}

func (s *Client) UpdatePartition(spartition Partition) error {
	url := fmt.Sprintf("v1/partitions/%s", spartition.ID)
	_, err := s.Put(url, spartition, false)

	return err
}

type Partition struct {
	ID                               string `json:"id,omitempty"`
	Name                             string `json:"name"`
	RoutingExpression                string `json:"routingExpression,omitempty"`
	AnalyticsTier                    string `json:"analyticsTier"`
	RetentionPeriod                  int    `json:"retentionPeriod"`
	IsCompliant                      bool   `json:"isCompliant"`
	DataForwardingId                 string `json:"dataForwardingId"`
	IsActive                         bool   `json:"isActive"`
	TotalBytes                       int    `json:"totalBytes"`
	IndexType                        string `json:"indexType"`
	ReduceRetentionPeriodImmediately bool   `json:"reduceRetentionPeriodImmediately,omitempty"`
}
