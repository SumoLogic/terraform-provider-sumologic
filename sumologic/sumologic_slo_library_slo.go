package sumologic

import (
	"encoding/json"
	"fmt"
	"log"
)

// ---------- ENDPOINTS ----------

func (s *Client) CreateSLO(slo SLOLibrarySLO, paramMap map[string]string) (string, error) {
	urlWithoutParams := SLOBaseApiUrl
	paramString := ""
	sprintfArgs := []interface{}{}

	paramString += "?"

	if val, ok := paramMap["parentId"]; ok {
		queryParam := fmt.Sprintf("parentId=%s&", val)
		paramString += queryParam
	}

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	data, err := s.Post(urlWithParams, slo)
	if err != nil {
		return "", err
	}
	log.Printf("created slo response: %v", data)

	var createdSLO SLOLibrarySLO

	err = json.Unmarshal(data, &createdSLO)
	if err != nil {
		return "", err
	}

	return createdSLO.ID, nil
}

func (s *Client) SLORead(id string, paramMap map[string]string) (*SLOLibrarySLO, error) {
	urlWithoutParams := SLOBaseApiUrl + "/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	data, _, err := s.Get(urlWithParams)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var sloRead SLOLibrarySLO

	err = json.Unmarshal(data, &sloRead)

	if err != nil {
		return nil, err
	}

	return &sloRead, nil
}

func (s *Client) DeleteSLO(id string) error {
	urlWithoutParams := SLOBaseApiUrl + "/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, id)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	_, err := s.Delete(urlWithParams)

	return err
}

func (s *Client) UpdateSLO(slo SLOLibrarySLO) error {
	urlWithoutParams := SLOBaseApiUrl + "/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, slo.ID)

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	slo.ID = ""

	_, err := s.Put(urlWithParams, slo)

	return err
}

func (s *Client) MoveSLOLibraryToFolder(slo SLOLibrarySLO) error {
	urlWithoutParams := SLOBaseApiUrl + "/%s"
	paramString := ""
	sprintfArgs := []interface{}{}
	sprintfArgs = append(sprintfArgs, slo.ID)

	paramString += "?"
	queryParam := fmt.Sprintf("parentId=%s&", slo.ParentID)
	paramString += queryParam

	urlWithParams := fmt.Sprintf(urlWithoutParams+paramString, sprintfArgs...)

	slo.ID = ""

	_, err := s.Put(urlWithParams, slo)

	return err
}

// ---------- TYPES ----------

type SLOLibrarySLO struct {
	ID          string        `json:"id,omitempty"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Version     int           `json:"version"`
	CreatedAt   string        `json:"createdAt"`
	CreatedBy   string        `json:"createdBy"`
	ModifiedAt  string        `json:"modifiedAt"`
	ModifiedBy  string        `json:"modifiedBy"`
	ParentID    string        `json:"parentId"`
	ContentType string        `json:"contentType"`
	Type        string        `json:"type"`
	IsSystem    bool          `json:"isSystem"`
	IsMutable   bool          `json:"isMutable"`
	IsLocked    bool          `json:"isLocked"`
	SignalType  string        `json:"signalType"` // string^(Latency|Error|Throughput|Availability|Other)$
	Compliance  SLOCompliance `json:"compliance"`
	Indicator   SLOIndicator  `json:"indicator"`
	Service     string        `json:"service"`
	Application string        `json:"application"`
}

type SLOCompliance struct {
	ComplianceType string  `json:"complianceType"`       // string^(Window|Request)$
	Target         float64 `json:"target"`               // [0..100]
	Timezone       string  `json:"timezone"`             // IANA Time Zone Database
	Size           string  `json:"size,omitempty"`       // Must be a multiple of days (minimum 1d, and maximum 14d)
	WindowType     string  `json:"windowType,omitempty"` // string^(Daily|Weekly|Monthly|Yearly)$
	StartFrom      string  `json:"startFrom,omitempty"`
}

type SLOIndicator struct {
	EvaluationType string          `json:"evaluationType"` // string^(Window|Request)$
	QueryType      string          `json:"queryType"`      // string^(Logs|Metrics)$
	Queries        []SLIQueryGroup `json:"queries"`
	Threshold      float64         `json:"threshold"`
	Op             string          `json:"op,omitempty"`
	Aggregation    string          `json:"aggregation,omitempty"`
	Size           string          `json:"size,omitempty"`
}

type SLI struct {
	EvaluationType string          `json:"evaluationType"` // string^(Window|Request)$
	QueryType      string          `json:"queryType"`      // string^(Logs|Metrics)$
	Queries        []SLIQueryGroup `json:"queries"`
}

type SLIQueryGroup struct {
	QueryGroupType string     `json:"queryGroupType"` // string^(Successful|Unsuccessful|Total|Threshold)$
	QueryGroup     []SLIQuery `json:"queryGroup"`
}

type SLIQuery struct {
	RowId       string `json:"rowId"`
	Query       string `json:"query"`
	UseRowCount bool   `json:"useRowCount"`
	Field       string `json:"field,omitempty"`
}

// SloBurnRateCondition struct for SloBurnRateCondition
type SloBurnRateCondition struct {
	TriggerCondition
	// The burn rate percentage.
	BurnRateThreshold float64 `json:"burnRateThreshold"`
	// The relative time range for the burn rate percentage evaluation.
	TimeRange string `json:"timeRange"`
}

// ---------- END ----------
