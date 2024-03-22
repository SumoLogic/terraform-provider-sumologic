package sumologic

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ListPartitionResp struct {
	Data []Partition `json:"data"`
	Next string      `json:"next"`
}

func (s *Client) ListPartitions() ([]Partition, error) {
	var listPartitionResp ListPartitionResp

	data, _, err := s.Get("v1/partitions")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &listPartitionResp)
	if err != nil {
		return nil, err
	}

	spartitions := listPartitionResp.Data

	for listPartitionResp.Next != "" {
		data, _, err = s.Get("v1/partitions?token=" + listPartitionResp.Next)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(data, &listPartitionResp)
		if err != nil {
			return nil, err
		}

		spartitions = append(spartitions, listPartitionResp.Data...)
	}

	var activePartitions []Partition
	for _, partition := range spartitions {
		if partition.IsActive {
			activePartitions = append(activePartitions, partition)
		}
	}

	return activePartitions, nil
}

func (s *Client) GetPartition(id string) (*Partition, error) {
	data, _, err := s.Get(fmt.Sprintf("v1/partitions/%s", id))
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
	} else if !spartition.IsActive {
		return nil, nil
	}

	return &spartition, nil
}

func (s *Client) CreatePartition(spartition Partition) (*Partition, error) {
	var createdspartition Partition

	responseBody, err := s.Post("v1/partitions", spartition)
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
	_, err := s.Post(fmt.Sprintf("v1/partitions/%s/decommission", id), nil)

	return err
}

func (s *Client) UpdatePartition(spartition Partition) error {
	url := fmt.Sprintf("v1/partitions/%s", spartition.ID)
	_, err := s.Put(url, spartition)

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
