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

	// NOTE: Before sending out, denormalize TriggerConditions
	for i, t := range monitorsLibraryMonitor.Triggers {
		monitorsLibraryMonitor.Triggers[i] = DenormalizeTriggerCondition(t)
	}

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

	// NOTE: Normalize all TriggerConditions
	for i, t := range monitorsLibraryMonitor.Triggers {
		if normalized, err := NormalizeTriggerCondition(t); err == nil  {
			monitorsLibraryMonitor.Triggers[i] = normalized
		} else {
			return &monitorsLibraryMonitor, err
		}
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

	// NOTE: Before sending out, denormalize TriggerConditions
	for i, t := range monitorsLibraryMonitor.Triggers {
		monitorsLibraryMonitor.Triggers[i] = DenormalizeTriggerCondition(t)
	}
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

func NormalizeTriggerCondition(condition TriggerCondition) (TriggerCondition, error) {
	returnValue := TriggerCondition{Trigger: &Trigger{}}
	if condition.Trigger != nil {
		if isEmptyTrigger(*condition.Trigger) {
			return returnValue, fmt.Errorf("Bad Trigger: empty")
		}
		// Already in normal form. Return without the legacy fields
		return TriggerCondition{Trigger: condition.Trigger}, nil
	}
	switch condition.DetectionMethod {
	case "StaticCondition":
		returnValue.Trigger.StaticCondition = toStaticCondition(condition)
	case "LogsStaticCondition":
		returnValue.Trigger.LogsStaticCondition = toLogsStaticCondition(condition)
	case "MetricsStaticCondition":
		returnValue.Trigger.MetricsStaticCondition = toMetricsStaticCondition(condition)
	case "LogsOutlierCondition":
		returnValue.Trigger.LogsOutlierCondition = toLogsOutlierCondition(condition)
	case "MetricsOutlierCondition":
		returnValue.Trigger.MetricsOutlierCondition = toMetricsOutlierCondition(condition)
	case "LogsMissingDataCondition":
		returnValue.Trigger.LogsMissingDataCondition = toLogsMissingDataCondition(condition)
	case "MetricsMissingDataCondition":
		returnValue.Trigger.MetricsMissingDataCondition = toMetricsMissingDataCondition(condition)
	default:
		fmt.Println("NormalizeTriggerCondition: Invalid detection method:", condition.DetectionMethod)
	}
	return returnValue, nil
}

func DenormalizeTriggerCondition(condition TriggerCondition) TriggerCondition {
	switch trigger := condition.Trigger; {
	case trigger != nil:
		switch {
		case trigger.StaticCondition != nil:
			return fromStaticCondition(trigger.StaticCondition)
		case trigger.LogsStaticCondition != nil:
			return fromLogsStaticCondition(trigger.LogsStaticCondition)
		case trigger.MetricsStaticCondition != nil:
			return fromMetricsStaticCondition(trigger.MetricsStaticCondition)
		case trigger.LogsOutlierCondition != nil:
			return fromLogsOutlierCondition(trigger.LogsOutlierCondition)
		case trigger.MetricsOutlierCondition != nil:
			return fromMetricsOutlierCondition(trigger.MetricsOutlierCondition)
		case trigger.LogsMissingDataCondition != nil:
			return fromLogsMissingDataCondition(trigger.LogsMissingDataCondition)
		case trigger.MetricsMissingDataCondition != nil:
			return fromMetricsMissingDataCondition(trigger.MetricsMissingDataCondition)
		default:
			fmt.Println("DenormalizeTriggerCondition failed because none of the required attributes is present", trigger)
			return condition
		}
	default:
		// Already denormalized. Return as-is
		return condition
	}
}

func toStaticCondition(condition TriggerCondition) *StaticCondition {
	return &StaticCondition{
		TimeRange:      condition.TimeRange,
		TriggerType:    condition.TriggerType,
		Threshold:      condition.Threshold,
		ThresholdType:  condition.ThresholdType,
		OccurrenceType: condition.OccurrenceType,
		TriggerSource:  condition.TriggerSource,
		Field:          condition.Field,
	}
}

func toLogsStaticCondition(condition TriggerCondition) *LogsStaticCondition {
	return &LogsStaticCondition{
		TimeRange:     condition.TimeRange,
		TriggerType:   condition.TriggerType,
		Threshold:     condition.Threshold,
		ThresholdType: condition.ThresholdType,
		Field:         condition.Field,
	}
}

func toMetricsStaticCondition(condition TriggerCondition) *MetricsStaticCondition {
	return &MetricsStaticCondition{
		TimeRange:      condition.TimeRange,
		TriggerType:    condition.TriggerType,
		Threshold:      condition.Threshold,
		ThresholdType:  condition.ThresholdType,
		OccurrenceType: condition.OccurrenceType,
	}
}

func toLogsOutlierCondition(condition TriggerCondition) *LogsOutlierCondition {
	return &LogsOutlierCondition{
		TriggerType: condition.TriggerType,
		Threshold:   condition.Threshold,
		Field:       condition.Field,
		Window:      condition.Window,
		Consecutive: condition.Consecutive,
		Direction:   condition.Direction,
	}
}

func toMetricsOutlierCondition(condition TriggerCondition) *MetricsOutlierCondition {
	return &MetricsOutlierCondition{
		TriggerType:    condition.TriggerType,
		Threshold:      condition.Threshold,
		BaselineWindow: condition.BaselineWindow,
		Direction:      condition.Direction,
	}
}

func toLogsMissingDataCondition(condition TriggerCondition) *LogsMissingDataCondition {
	return &LogsMissingDataCondition{
		TimeRange:   condition.TimeRange,
		TriggerType: condition.TriggerType,
	}
}

func toMetricsMissingDataCondition(condition TriggerCondition) *MetricsMissingDataCondition {
	return &MetricsMissingDataCondition{
		TimeRange:     condition.TimeRange,
		TriggerType:   condition.TriggerType,
		TriggerSource: condition.TriggerSource,
	}
}

func fromStaticCondition(condition *StaticCondition) TriggerCondition {
	return TriggerCondition{
		TimeRange:       condition.TimeRange,
		TriggerType:     condition.TriggerType,
		Threshold:       condition.Threshold,
		ThresholdType:   condition.ThresholdType,
		OccurrenceType:  condition.OccurrenceType,
		TriggerSource:   condition.TriggerSource,
		Field:           condition.Field,
		DetectionMethod: "StaticCondition",
	}
}

func fromLogsStaticCondition(condition *LogsStaticCondition) TriggerCondition {
	return TriggerCondition{
		TimeRange:       condition.TimeRange,
		TriggerType:     condition.TriggerType,
		Threshold:       condition.Threshold,
		ThresholdType:   condition.ThresholdType,
		Field:           condition.Field,
		DetectionMethod: "LogsStaticCondition",
	}
}

func fromMetricsStaticCondition(condition *MetricsStaticCondition) TriggerCondition {
	return TriggerCondition{
		TimeRange:       condition.TimeRange,
		TriggerType:     condition.TriggerType,
		Threshold:       condition.Threshold,
		ThresholdType:   condition.ThresholdType,
		OccurrenceType:  condition.OccurrenceType,
		DetectionMethod: "MetricsStaticCondition",
	}
}

func fromLogsOutlierCondition(condition *LogsOutlierCondition) TriggerCondition {
	return TriggerCondition{
		TriggerType:     condition.TriggerType,
		Threshold:       condition.Threshold,
		Field:           condition.Field,
		Window:          condition.Window,
		Consecutive:     condition.Consecutive,
		Direction:       condition.Direction,
		DetectionMethod: "LogsOutlierCondition",
	}
}

func fromMetricsOutlierCondition(condition *MetricsOutlierCondition) TriggerCondition {
	return TriggerCondition{
		TriggerType:     condition.TriggerType,
		Threshold:       condition.Threshold,
		BaselineWindow:  condition.BaselineWindow,
		Direction:       condition.Direction,
		DetectionMethod: "MetricsOutlierCondition",
	}
}

func fromLogsMissingDataCondition(condition *LogsMissingDataCondition) TriggerCondition {
	return TriggerCondition{
		TimeRange:       condition.TimeRange,
		TriggerType:     condition.TriggerType,
		DetectionMethod: "LogsMissingDataCondition",
	}
}

func fromMetricsMissingDataCondition(condition *MetricsMissingDataCondition) TriggerCondition {
	return TriggerCondition{
		TimeRange:       condition.TimeRange,
		TriggerType:     condition.TriggerType,
		TriggerSource:   condition.TriggerSource,
		DetectionMethod: "MetricsMissingDataCondition",
	}
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
	// NOTE: The following fields are here for encoding/decoding JSON.
	Field           string  `json:"field,omitempty"`
	Window          int     `json:"window,omitempty"`
	BaselineWindow  string  `json:"baselineWindow,omitempty"`
	Consecutive     int     `json:"consecutive,omitempty"`
	Direction       string  `json:"direction,omitempty"`
	// NOTE: 'Trigger' is an internal-only optional field--it is not translated to/from infrastructure JSON.
	// Methods NormalizeTriggerCondition and DenormalizeTriggerCondition move TriggerCondition's attributes to (resp. from) Trigger.
	Trigger *Trigger `json:"-"`
}

func isEmptyTrigger(t Trigger) bool {
	return t.StaticCondition == nil &&
		t.LogsStaticCondition == nil &&
		t.MetricsStaticCondition == nil &&
		t.LogsOutlierCondition == nil &&
		t.MetricsOutlierCondition == nil &&
		t.LogsMissingDataCondition == nil &&
		t.MetricsMissingDataCondition == nil
}

type Trigger struct {
	StaticCondition             *StaticCondition
	LogsStaticCondition         *LogsStaticCondition
	MetricsStaticCondition      *MetricsStaticCondition
	LogsOutlierCondition        *LogsOutlierCondition
	MetricsOutlierCondition     *MetricsOutlierCondition
	LogsMissingDataCondition    *LogsMissingDataCondition
	MetricsMissingDataCondition *MetricsMissingDataCondition
}

type StaticCondition struct {
	TimeRange      string
	TriggerType    string
	Threshold      float64
	ThresholdType  string
	OccurrenceType string
	TriggerSource  string
	Field          string
}

type LogsStaticCondition struct {
	TimeRange     string
	TriggerType   string
	Threshold     float64
	ThresholdType string
	Field         string
}

type MetricsStaticCondition struct {
	TimeRange      string
	TriggerType    string
	Threshold      float64
	ThresholdType  string
	OccurrenceType string
}

type LogsOutlierCondition struct {
	TriggerType string
	Threshold   float64
	Field       string
	Window      int
	Consecutive int
	Direction   string
}

type MetricsOutlierCondition struct {
	TriggerType    string
	Threshold      float64
	BaselineWindow string
	Direction      string
}

type LogsMissingDataCondition struct {
	TimeRange   string
	TriggerType string
}

type MetricsMissingDataCondition struct {
	TimeRange     string
	TriggerType   string
	TriggerSource string
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
