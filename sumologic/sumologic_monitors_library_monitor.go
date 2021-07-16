package sumologic

import (
	"encoding/json"
	"fmt"
	"log"
)

// ---------- ENDPOINTS ----------

func (s *Client) CreateMonitorsLibraryMonitor(monitorsLibraryMonitor MonitorsLibraryMonitor, paramMap map[string]string) (string, error) {
	urlWithoutParams := "v1/monitors"
	paramString := ""
	sprintfArgs := []interface{}{}

	paramString += "?"

	if val, ok := paramMap["parentId"]; ok {
		queryParam := fmt.Sprintf("parentId=%s&", val)
		paramString += queryParam
	}

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	data, err := s.Post(urlWithParams, monitorsLibraryMonitor, false)
	if err != nil {
		return "", err
	}
	log.Printf("created monitor response: %v", data)

	var createdMonitorsLibraryMonitor MonitorsLibraryMonitor

	err = json.Unmarshal(data, &createdMonitorsLibraryMonitor)
	if err != nil {
		return "", err
	}

	return createdMonitorsLibraryMonitor.ID, nil

}

func (s *Client) MonitorsRead(id string) (*MonitorsLibraryMonitor, error) {
	urlWithoutParams := "v1/monitors/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	data, _, err := s.Get(urlWithParams, false)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var monitorsLibraryMonitor MonitorsLibraryMonitor

	err = json.Unmarshal(data, &monitorsLibraryMonitor)

	if err != nil {
		return nil, err
	}

	return &monitorsLibraryMonitor, nil
}

func (s *Client) DeleteMonitorsLibraryMonitor(id string) error {
	urlWithoutParams := "v1/monitors/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	_, err := s.Delete(urlWithParams)

	return err
}

func (s *Client) UpdateMonitorsLibraryMonitor(monitorsLibraryMonitor MonitorsLibraryMonitor) error {
	urlWithoutParams := "v1/monitors/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, monitorsLibraryMonitor.ID)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	monitorsLibraryMonitor.ID = ""

	_, err := s.Put(urlWithParams, monitorsLibraryMonitor, false)

	return err
}

func (s *Client) MoveMonitorsLibraryMonitor(monitorsLibraryMonitor MonitorsLibraryMonitor) error {
	urlWithoutParams := "v1/monitors/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, monitorsLibraryMonitor.ID)

	paramString += "?"
	queryParam := fmt.Sprintf("parentId=%s&", monitorsLibraryMonitor.ParentID)
	paramString += queryParam

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	monitorsLibraryMonitor.ID = ""

	_, err := s.Put(urlWithParams, monitorsLibraryMonitor, false)

	return err
}

// ---------- TYPES ----------
type MonitorsLibraryMonitor struct {
	ID                 string                `json:"id,omitempty"`
	IsSystem           bool                  `json:"isSystem"`
	Type               string                `json:"type"`
	Queries            []MonitorQuery        `json:"queries,omitempty"`
	ParentID           string                `json:"parentId"`
	Name               string                `json:"name"`
	IsMutable          bool                  `json:"isMutable"`
	Version            int                   `json:"version"`
	Notifications      []MonitorNotification `json:"notifications,omitempty"`
	CreatedBy          string                `json:"createdBy"`
	MonitorType        string                `json:"monitorType"`
	IsLocked           bool                  `json:"isLocked"`
	Description        string                `json:"description"`
	CreatedAt          string                `json:"createdAt"`
	Triggers           []TriggerCondition    `json:"triggers,omitempty"`
	ModifiedAt         string                `json:"modifiedAt"`
	ContentType        string                `json:"contentType"`
	ModifiedBy         string                `json:"modifiedBy"`
	IsDisabled         bool                  `json:"isDisabled"`
	Status             []string              `json:"status"`
	GroupNotifications bool                  `json:"groupNotifications"`
}

type MonitorQuery struct {
	RowID string `json:"rowId"`
	Query string `json:"query"`
}

type TriggerCondition struct {
	TimeRange       string  `json:"timeRange"`
	TriggerType     string  `json:"triggerType"`
	Threshold       float64 `json:"threshold,omitempty"`
	ThresholdType   string  `json:"thresholdType,omitempty"`
	OccurrenceType  string  `json:"occurrenceType"`
	TriggerSource   string  `json:"triggerSource"`
	DetectionMethod string  `json:"detectionMethod"`
	Field           string  `json:"field,omitempty"`
	Window          int     `json:"window,omitempty"`
	BaselineWindow  string  `json:"baselineWindow,omitempty"`
	Consecutive     int     `json:"consecutive,omitempty"`
	Direction       string  `json:"direction,omitempty"`
}

type MonitorNotification struct {
	Notification       interface{}   `json:"notification"`
	RunForTriggerTypes []interface{} `json:"runForTriggerTypes"`
}

type EmailNotification struct {
	ActionType     string        `json:"actionType,omitempty"`
	ConnectionType string        `json:"connectionType,omitempty"`
	Subject        string        `json:"subject"`
	Recipients     []interface{} `json:"recipients"`
	MessageBody    string        `json:"messageBody"`
	TimeZone       string        `json:"timeZone"`
}

type WebhookNotificiation struct {
	ActionType      string `json:"actionType,omitempty"`
	ConnectionType  string `json:"connectionType,omitempty"`
	ConnectionID    string `json:"connectionId"`
	PayloadOverride string `json:"payloadOverride,omitempty"`
}

// ---------- END ----------
