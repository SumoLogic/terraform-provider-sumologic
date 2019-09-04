package sumologic

import (
	"encoding/json"
	"time"
	"fmt"
)

func (s *Client) GetScheduledView(id string) (*ScheduledView, error) {
	data, _, err := s.Get(fmt.Sprintf("scheduledViews/%s", id))
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

	responseBody, err := s.Post("scheduledViews", sview)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &createdSview)

	if err != nil {
		return nil, err
	}

	return &createdSview, nil
}

type ScheduledView struct {
	ID               string    `json:"id,omitempty"`
	Query            string    `json:"query"`
	IndexName        string    `json:"indexName"`
	StartTime        time.Time `json:"startTime"`
	RetentionPeriod  int       `json:"retentionPeriod"`
	DataForwardingId string    `json:"dataForwardingId"`
}
