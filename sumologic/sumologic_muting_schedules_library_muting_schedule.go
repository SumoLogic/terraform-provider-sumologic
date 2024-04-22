package sumologic

import (
	"encoding/json"
	"fmt"
)

// ---------- ENDPOINTS ----------

func (s *Client) CreateMutingSchedulesLibraryMutingSchedule(mutingSchedulesLibraryMutingSchedule MutingSchedulesLibraryMutingSchedule, paramMap map[string]string) (string, error) {
	urlWithoutParams := "v1/mutingSchedules"
	paramString := ""
	sprintfArgs := []interface{}{}

	paramString += "?"

	if val, ok := paramMap["parentId"]; ok {
		queryParam := fmt.Sprintf("parentId=%s&", val)
		paramString += queryParam
	}

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	data, err := s.Post(urlWithParams, mutingSchedulesLibraryMutingSchedule)
	if err != nil {
		return "", err
	}

	var createdMutingSchedulesLibraryMutingSchedule MutingSchedulesLibraryMutingSchedule

	err = json.Unmarshal(data, &createdMutingSchedulesLibraryMutingSchedule)
	if err != nil {
		return "", err
	}

	return createdMutingSchedulesLibraryMutingSchedule.ID, nil

}

func (s *Client) MutingSchedulesRead(id string) (*MutingSchedulesLibraryMutingSchedule, error) {
	urlWithoutParams := "v1/mutingSchedules/%s"
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams, sprintfArgs...)

	data, _, err := s.Get(urlWithParams)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var mutingSchedulesLibraryMutingSchedule MutingSchedulesLibraryMutingSchedule

	err = json.Unmarshal(data, &mutingSchedulesLibraryMutingSchedule)

	if err != nil {
		return nil, err
	}

	return &mutingSchedulesLibraryMutingSchedule, nil
}

func (s *Client) DeleteMutingSchedulesLibraryMutingSchedule(id string) error {
	urlWithoutParams := "v1/mutingSchedules/%s"
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams, sprintfArgs...)

	_, err := s.Delete(urlWithParams)

	return err
}

func (s *Client) UpdateMutingSchedulesLibraryMutingSchedule(mutingSchedulesLibraryMutingSchedule MutingSchedulesLibraryMutingSchedule) error {
	urlWithoutParams := "v1/mutingSchedules/%s"
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, mutingSchedulesLibraryMutingSchedule.ID)

	urlWithParams := fmt.Sprintf(urlWithoutParams, sprintfArgs...)

	_, err := s.Put(urlWithParams, mutingSchedulesLibraryMutingSchedule)

	return err
}

func (s *Client) GetMutingSchedulesLibraryFolder(id string) (*MutingSchedulesLibraryFolder, error) {
	urlWithoutParams := "v1/mutingSchedules/%s"
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams, sprintfArgs...)

	data, _, err := s.Get(urlWithParams)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var mutingSchedulesLibraryFolder MutingSchedulesLibraryFolder

	err = json.Unmarshal(data, &mutingSchedulesLibraryFolder)

	if err != nil {
		return nil, err
	}

	return &mutingSchedulesLibraryFolder, nil

}

// ---------- TYPES ----------
type MutingSchedulesLibraryMutingSchedule struct {
	ID                 string                        `json:"id"`
	Type               string                        `json:"type"`
	IsSystem           bool                          `json:"isSystem"`
	IsMutable          bool                          `json:"isMutable"`
	Schedule           ScheduleDefinition            `json:"schedule"`
	Monitor            *MonitorScope                 `json:"monitor"`
	ParentID           string                        `json:"parentId"`
	Name               string                        `json:"name"`
	Version            int                           `json:"version"`
	CreatedBy          string                        `json:"createdBy"`
	Description        string                        `json:"description"`
	CreatedAt          string                        `json:"createdAt"`
	ModifiedAt         string                        `json:"modifiedAt"`
	ContentType        string                        `json:"contentType"`
	ModifiedBy         string                        `json:"modifiedBy"`
	NotificationGroups []NotificationGroupDefinition `json:"notificationGroups"`
}

type ScheduleDefinition struct {
	TimeZone  string `json:"timezone"`
	StartDate string `json:"startDate"`
	StartTime string `json:"startTime"`
	Duration  int    `json:"duration"`
	RRule     string `json:"rrule,omitempty"`
}

type MonitorScope struct {
	Ids []string `json:"ids,omitempty"`
	All bool     `json:"all,omitempty"`
}

type NotificationGroupDefinition struct {
	GroupKey    string   `json:"groupKey"`
	GroupValues []string `json:"groupValues"`
}

type MutingSchedulesLibraryFolder struct {
	ID          string `json:"id,omitempty"`
	Type        string `json:"type"`
	ContentType string `json:"contentType"`
	ParentID    string `json:"parentId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedBy   string `json:"createdBy"`
	CreatedAt   string `json:"createdAt"`
	ModifiedBy  string `json:"modifiedBy"`
	ModifiedAt  string `json:"modifiedAt"`
	IsLocked    bool   `json:"isLocked"`
	IsMutable   bool   `json:"isMutable"`
	IsSystem    bool   `json:"isSystem"`
	Version     int    `json:"version"`
}

// ---------- END ----------
