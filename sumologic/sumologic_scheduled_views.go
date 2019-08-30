package sumologic

import (
	"encoding/json"
)

func (s *Client) CreateScheduledView(sview ScheduledView) (string, error) {
	var createdSview ScheduledView

	responseBody, err := s.Post("scheduledViews", sview)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &createdSview)

	if err != nil {
		return "", err
	}

	return createdSview.ID, nil
}

type ScheduledView struct {
	ID               string   `json:"id,omitempty"`
	Query            string   `json:"query"`
	IndexName        string   `json:"indexName"`
	StartTime        string   `json:"startTime"`
	RetentionPeriod  int      `json:"retentionPeriod"`
	DataForwardingId string   `json:"dataForwardingId"`
}
