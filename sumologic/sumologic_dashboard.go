package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) GetDashboard(id string) (*Dashboard, error) {
	url := fmt.Sprintf("v2/dashboards/%s", id)
	data, _, err := s.Get(url, false)
	if err != nil {
		return nil, err
	}

	var dashboard Dashboard
	err = json.Unmarshal(data, &dashboard)
	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}

func (s *Client) CreateDashboard(dashboardReq Dashboard) (*Dashboard, error) {
	responseBody, err := s.Post("v2/dashboards", dashboardReq, false)
	if err != nil {
		return nil, err
	}

	var dashboard Dashboard
	err = json.Unmarshal(responseBody, &dashboard)
	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}

func (s *Client) DeleteDashboard(id string) error {
	url := fmt.Sprintf("v2/dashboards/%s", id)
	_, err := s.Delete(url)
	return err
}

func (s *Client) UpdateDashboard(dashboard Dashboard) error {
	url := fmt.Sprintf("v2/dashboards/%s", dashboard.ID)
	_, err := s.Put(url, dashboard, false)
	return err
}

type Dashboard struct {
	ID               string         `json:"id,omitempty"`
	Title            string         `json:"title"`
	Description      string         `json:"description"`
	FolderId         string         `json:"folderId"`
	TopologyLabelMap *TopologyLabel `json:"topologyLabelMap"`
	RefreshInterval  int            `json:"refreshInterval"`
	TimeRange        interface{}    `json:"timeRange"`
	Panels           []interface{}  `json:"panels"`
	Layout           interface{}    `json:"layout"`
	Variables        []Variable     `json:"variables"`
	Theme            string         `json:"theme"`
	ColoringRules    []ColoringRule `json:"coloringRules"`
}

type TopologyLabel struct {
	Data map[string][]string `json:"data"`
}

// TimeRange related structs
type CompleteLiteralTimeRange struct {
	Type      string `json:"type"`
	RangeName string `json:"rangeName"`
}

type BeginBoundedTimeRange struct {
	Type string      `json:"type"`
	From interface{} `json:"from"`
	To   interface{} `json:"to"`
}

type RelativeTimeRangeBoundary struct {
	Type         string `json:"type"`
	RelativeTime string `json:"relativeTime"`
}

type EpochTimeRangeBoundary struct {
	Type        string `json:"type"`
	EpochMillis int64  `json:"epochMillis"`
}

type Iso8601TimeRangeBoundary struct {
	Type        string `json:"type"`
	Iso8601Time string `json:"iso8601Time"`
}

type LiteralTimeRangeBoundary struct {
	Type      string `json:"type"`
	RangeName string `json:"rangeName"`
}

// Panel related structs
type TextPanel struct {
	Id                                     string `json:"id,omitempty"`
	Key                                    string `json:"key"`
	Title                                  string `json:"title"`
	VisualSettings                         string `json:"visualSettings"`
	KeepVisualSettingsConsistentWithParent bool   `json:"keepVisualSettingsConsistentWithParent"`
	PanelType                              string `json:"panelType"`
	// Text panel related properties
	Text string `json:"text"`
}

type SumoSearchPanel struct {
	Id                                     string `json:"id,omitempty"`
	Key                                    string `json:"key"`
	Title                                  string `json:"title"`
	VisualSettings                         string `json:"visualSettings"`
	KeepVisualSettingsConsistentWithParent bool   `json:"keepVisualSettingsConsistentWithParent"`
	PanelType                              string `json:"panelType"`
	// Search panel related properties
	Queries          []SearchPanelQuery `json:"queries"`
	Description      string             `json:"description"`
	TimeRange        interface{}        `json:"timeRange"`
	ColoringRules    []ColoringRule     `json:"coloringRules"`
	LinkedDashboards []LinkedDashboard  `json:"linkedDashboards"`
}

type SearchPanelQuery struct {
	QueryString      string            `json:"queryString"`
	QueryType        string            `json:"queryType"`
	QueryKey         string            `json:"queryKey"`
	MetricsQueryMode string            `json:"metricsQueryMode,omitempty"`
	MetricsQueryData *MetricsQueryData `json:"metricsQueryData,omitempty"`
}

type MetricsQueryData struct {
	Metric          string                 `json:"metric"`
	AggregationType string                 `json:"aggregationType"`
	GroupBy         string                 `json:"groupBy,omitempty"`
	Filters         []MetricsQueryFilter   `json:"filters"`
	Operators       []MetricsQueryOperator `json:"operators"`
}

type MetricsQueryFilter struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Negation bool   `json:"negation,omitempty"`
}

type MetricsQueryOperator struct {
	Name       string                          `json:"operatorName"`
	Parameters []MetricsQueryOperatorParameter `json:"parameters"`
}

type MetricsQueryOperatorParameter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Layout related structs
type GridLayout struct {
	LayoutType       string            `json:"layoutType"`
	LayoutStructures []LayoutStructure `json:"layoutStructures"`
}

type LayoutStructure struct {
	Key       string `json:"key"`
	Structure string `json:"structure"`
}

// Variable related structs
type Variable struct {
	Id               string      `json:"id,omitempty"`
	Name             string      `json:"name"`
	DisplayName      string      `json:"displayName,omitempty"`
	DefaultValue     string      `json:"defaultValue,omitempty"`
	SourceDefinition interface{} `json:"sourceDefinition"`
	AllowMultiSelect bool        `json:"allowMultiSelect,omitempty"`
	IncludeAllOption bool        `json:"includeAllOption"`
	HideFromUI       bool        `json:"hideFromUI,omitempty"`
}

type MetadataVariableSourceDefinition struct {
	VariableSourceType string `json:"variableSourceType"`
	Filter             string `json:"filter"`
	Key                string `json:"key"`
}

type CsvVariableSourceDefinition struct {
	VariableSourceType string `json:"variableSourceType"`
	Values             string `json:"values"`
}

type LogQueryVariableSourceDefinition struct {
	VariableSourceType string `json:"variableSourceType"`
	Query              string `json:"query"`
	Field              string `json:"field"`
}

// Coloring Rule related structs
type ColoringRule struct {
	Scope                           string           `json:"scope"`
	SingleSeriesAggregateFunction   string           `json:"singleSeriesAggregateFunction"`
	MultipleSeriesAggregateFunction string           `json:"multipleSeriesAggregateFunction"`
	ColorThresholds                 []ColorThreshold `json:"colorThresholds"`
}

type ColorThreshold struct {
	Color string  `json:"color"`
	Min   float64 `json:"min,omitempty"`
	Max   float64 `json:"max,omitempty"`
}

type LinkedDashboard struct {
	Id               string `json:"id"`
	RelativePath     string `json:"relativePath,omitempty"`
	IncludeTimeRange bool   `json:"includeTimeRange"`
	IncludeVariables bool   `json:"includeVariables"`
}
