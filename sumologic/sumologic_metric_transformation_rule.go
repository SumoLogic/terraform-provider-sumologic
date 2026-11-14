package sumologic

import (
	"encoding/json"
	"fmt"
)

const metricTransformationRuleBaseApiUrl = "v1/transformationRules"

// ---------- ENDPOINTS ----------

// CreateMetricTransformationRule creates a rule and returns its ID.
func (s *Client) CreateMetricTransformationRule(rule MetricTransformationRule) (string, error) {
	data, err := s.Post(metricTransformationRuleBaseApiUrl, rule.toRequest())
	if err != nil {
		return "", err
	}

	// Decode only to read the ID. Do not log the response body.
	var created MetricTransformationRule
	if err := json.Unmarshal(data, &created); err != nil {
		return "", err
	}
	return created.ID, nil
}

// GetMetricTransformationRule reads a rule by id. It returns nil if the rule does not exist.
func (s *Client) GetMetricTransformationRule(id string) (*MetricTransformationRule, error) {
	data, err := s.Get(fmt.Sprintf("%s/%s", metricTransformationRuleBaseApiUrl, id))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var rule MetricTransformationRule
	if err := json.Unmarshal(data, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

// UpdateMetricTransformationRule replaces the rule's properties. The API does not use a version.
func (s *Client) UpdateMetricTransformationRule(rule MetricTransformationRule) error {
	_, err := s.Put(fmt.Sprintf("%s/%s", metricTransformationRuleBaseApiUrl, rule.ID), rule.toRequest())
	return err
}

func (s *Client) DeleteMetricTransformationRule(id string) error {
	_, err := s.Delete(fmt.Sprintf("%s/%s", metricTransformationRuleBaseApiUrl, id))
	return err
}

// ---------- TYPES ----------

const (
	metricTransformationTypeAggregateOn  = "AggregateOnTransformation"
	metricTransformationTypeAddOrReplace = "AddOrReplaceTransformation"
)

// MetricTransformationRule is the response shape of GET/POST/PUT /v1/transformationRules.
type MetricTransformationRule struct {
	ID             string                             `json:"id,omitempty"`
	Enabled        bool                               `json:"enabled"`
	RuleDefinition MetricTransformationRuleDefinition `json:"ruleDefinition"`
	CreatedAt      string                             `json:"createdAt,omitempty"`
	CreatedBy      string                             `json:"createdBy,omitempty"`
	ModifiedAt     string                             `json:"modifiedAt,omitempty"`
	ModifiedBy     string                             `json:"modifiedBy,omitempty"`
}

type MetricTransformationRuleDefinition struct {
	Name                        string                          `json:"name"`
	Selector                    string                          `json:"selector"`
	Retention                   int64                           `json:"retention"`
	TransformedMetricsRetention int64                           `json:"transformedMetricsRetention"`
	DimensionTransformations    []MetricDimensionTransformation `json:"dimensionTransformations"`
}

// MetricDimensionTransformation covers both subtypes, discriminated by TransformationType.
// AggregateOnTransformation uses AggregateOn; AddOrReplaceTransformation uses
// DimensionToReplace and Value.
type MetricDimensionTransformation struct {
	TransformationType string   `json:"transformationType"`
	AggregateOn        []string `json:"aggregateOn,omitempty"`
	DimensionToReplace string   `json:"dimensionToReplace,omitempty"`
	Value              string   `json:"value,omitempty"`
}

// metricTransformationRuleRequest is the body for POST and PUT. Read-only fields (id, metadata)
// are left out.
type metricTransformationRuleRequest struct {
	Enabled        bool                               `json:"enabled"`
	RuleDefinition MetricTransformationRuleDefinition `json:"ruleDefinition"`
}

func (r MetricTransformationRule) toRequest() metricTransformationRuleRequest {
	def := r.RuleDefinition
	if def.DimensionTransformations == nil {
		def.DimensionTransformations = []MetricDimensionTransformation{}
	}
	return metricTransformationRuleRequest{
		Enabled:        r.Enabled,
		RuleDefinition: def,
	}
}
