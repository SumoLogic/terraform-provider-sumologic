package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetLogSearch(id string) (*LogSearch, error) {
	urlWithoutParams := "v1/logSearches/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	data, err := s.Get(urlWithParams)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var logSearch LogSearch

	err = json.Unmarshal(data, &logSearch)
	if err != nil {
		return nil, err
	}

	return &logSearch, nil
}

func (s *Client) DeleteLogSearch(id string) error {
	urlWithoutParams := "v1/logSearches/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	_, err := s.Delete(urlWithParams)

	return err
}

func (s *Client) UpdateLogSearch(logSearch LogSearch) error {
	urlWithoutParams := "v1/logSearches/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, logSearch.ID)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	logSearch.ID = ""

	_, err := s.Put(urlWithParams, logSearch)

	return err
}

func (s *Client) CreateLogSearch(logSearch LogSearch) (string, error) {
	urlWithoutParams := "v1/logSearches"

	data, err := s.Post(urlWithoutParams, logSearch)
	if err != nil {
		return "", err
	}

	var createdLogSearch LogSearch

	err = json.Unmarshal(data, &createdLogSearch)
	if err != nil {
		return "", err
	}

	return createdLogSearch.ID, nil
}

type LogSearch struct {
	ID               string                    `json:"id,omitempty"`
	Name             string                    `json:"name"`
	Description      string                    `json:"description"`
	ParentId         string                    `json:"parentId"`
	QueryString      string                    `json:"queryString"`
	RunByReceiptTime bool                      `json:"runByReceiptTime"`
	TimeRange        interface{}               `json:"timeRange"`
	ParsingMode      string                    `json:"parsingMode"`
	QueryParameters  []LogSearchQueryParameter `json:"queryParameters,omitempty"`
	Schedule         *LogSearchSchedule        `json:"schedule,omitempty"`
}

type LogSearchQueryParameter struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DataType    string `json:"dataType"`
	Value       string `json:"value"`
}

type LogSearchSchedule struct {
	CronExpression     string                       `json:"cronExpression"`
	ParseableTimeRange interface{}                  `json:"parseableTimeRange"`
	TimeZone           string                       `json:"timeZone"`
	Threshold          *SearchNotificationThreshold `json:"threshold,omitempty"`
	Parameters         []ScheduleSearchParameter    `json:"parameters,omitempty"`
	MuteErrorEmails    bool                         `json:"muteErrorEmails"`
	Notification       interface{}                  `json:"notification"`
	ScheduleType       string                       `json:"scheduleType"`
}

type SearchNotificationThreshold struct {
	ThresholdType string `json:"thresholdType"`
	Operator      string `json:"operator"`
	Count         int    `json:"count"`
}

type ScheduleSearchParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type AlertSearchNotification struct {
	TaskType string `json:"taskType"`
	SourceId string `json:"sourceId"`
}

type CseSignalNotification struct {
	TaskType   string `json:"taskType"`
	RecordType string `json:"recordType"`
}

type EmailSearchNotification struct {
	TaskType             string   `json:"taskType"`
	ToList               []string `json:"toList"`
	SubjectTemplate      string   `json:"subjectTemplate"`
	IncludeQuery         bool     `json:"includeQuery"`
	IncludeResultSet     bool     `json:"includeResultSet"`
	IncludeHistogram     bool     `json:"includeHistogram"`
	IncludeCsvAttachment bool     `json:"includeCsvAttachment"`
}

type SaveToViewNotification struct {
	TaskType string `json:"taskType"`
	ViewName string `json:"viewName"`
}

type SaveToLookupNotification struct {
	TaskType               string `json:"taskType"`
	LookupFilePath         string `json:"lookupFilePath"`
	IsLookupMergeOperation bool   `json:"isLookupMergeOperation"`
}

type ServiceNowSearchNotification struct {
	TaskType   string           `json:"taskType"`
	ExternalId string           `json:"externalId"`
	Fields     ServiceNowFields `json:"fields"`
}

type WebhookSearchNotification struct {
	TaskType          string  `json:"taskType"`
	WebhookId         string  `json:"webhookId"`
	Payload           *string `json:"payload"`
	ItemizeAlerts     bool    `json:"itemizeAlerts"`
	MaxItemizedAlerts int     `json:"maxItemizedAlerts"`
}

type ServiceNowFields struct {
	EventType string `json:"eventType"`
	Severity  int    `json:"severity"`
	Resource  string `json:"resource"`
	Node      string `json:"node"`
}
