package sumologic

import (
	"encoding/json"
	"fmt"
	"time"
)

func (s *Client) GetScheduledView(id string) (*ScheduledView, error) {
	data, _, err := s.Get(fmt.Sprintf("v1/scheduledViews/%s", id), false)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var sview ScheduledView
	err = json.Unmarshal(data, &sview)
	if err != nil {
		return nil, err
	}

	return &sview, nil
}

func (s *Client) CreateScheduledView(sview ScheduledView) (*ScheduledView, error) {
	var createdSview ScheduledView

	responseBody, err := s.Post("v1/scheduledViews", sview, false)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &createdSview)

	if err != nil {
		return nil, err
	}

	return &createdSview, nil
}

func (s *Client) DeleteScheduledView(id string) error {
	_, err := s.Delete(fmt.Sprintf("v1/scheduledViews/%s/disable", id))

	return err
}

func (s *Client) UpdateScheduledView(sview ScheduledView) error {
	url := fmt.Sprintf("v1/scheduledViews/%s", sview.ID)

	_, err := s.Put(url, sview, false)

	return err
}

type ScheduledView struct {
	ID                               string    `json:"id,omitempty"`
	Query                            string    `json:"query"`
	IndexName                        string    `json:"indexName"`
	StartTime                        time.Time `json:"startTime"`
	RetentionPeriod                  int       `json:"retentionPeriod"`
	DataForwardingId                 string    `json:"dataForwardingId"`
	ParsingMode                      string    `json:"parsingMode"`
	ReduceRetentionPeriodImmediately bool      `json:"reduceRetentionPeriodImmediately"`
}
