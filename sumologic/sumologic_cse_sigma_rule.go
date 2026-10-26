package sumologic

import (
	"encoding/json"
	"fmt"
)

const cseSigmaOriginalSource = "sigma"

// GetCSESigmaRule fetches a converted Sigma rule by its Sigma ID.
func (s *Client) GetCSESigmaRule(sigmaID string) (*CSESigmaRule, error) {
	data, err := s.GetWithErrOpt(fmt.Sprintf("sec/v1/rules/sigma/%s", sigmaID), true)
	if err != nil {
		return nil, err
	}

	var response CSESigmaRuleResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.CSESigmaRule, nil
}

// ListCSESigmaRules returns up to limit converted Sigma rules.
func (s *Client) ListCSESigmaRules(limit int) ([]CSESigmaRule, error) {
	data, err := s.GetWithErrOpt(fmt.Sprintf("sec/v1/rules/sigma?limit=%d", limit), true)
	if err != nil {
		return nil, err
	}

	var response CSESigmaRuleListResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return response.Data.Data, nil
}

// toMatchRule converts the Sigma rule fields into a templated match rule request.
func (r CSESigmaRule) toMatchRule() CSEMatchRule {
	selectors := make([]EntitySelector, len(r.EntitySelectors))
	for i, selector := range r.EntitySelectors {
		selectors[i] = EntitySelector{
			EntityType: selector.EntityType,
			Expression: selector.Expression,
		}
	}

	tags := r.Tags
	if tags == nil {
		tags = []string{}
	}

	return CSEMatchRule{
		DescriptionExpression: r.DescriptionExpression,
		Enabled:               r.Enabled,
		EntitySelectors:       selectors,
		Expression:            r.Expression,
		IsPrototype:           r.IsPrototype,
		Name:                  r.Name,
		NameExpression:        r.NameExpression,
		SeverityMapping: SeverityMapping{
			Type:    r.ScoreMapping.Type,
			Default: r.ScoreMapping.Default,
			Field:   r.ScoreMapping.Field,
			Mapping: r.ScoreMapping.valueMappings(),
		},
		Stream:            r.Stream,
		SummaryExpression: r.SummaryExpression,
		Tags:              tags,
		OriginalSource:    cseSigmaOriginalSource,
	}
}

// valueMappings parses mapping entries that are objects; other shapes are ignored.
func (m CSESigmaScoreMapping) valueMappings() []SeverityMappingValueMapping {
	result := []SeverityMappingValueMapping{}
	for _, raw := range m.Mapping {
		var mapping SeverityMappingValueMapping
		if err := json.Unmarshal(raw, &mapping); err == nil && mapping.From != "" {
			result = append(result, mapping)
		}
	}
	return result
}

type CSESigmaRuleResponse struct {
	CSESigmaRule CSESigmaRule `json:"data"`
}

type CSESigmaRuleListResponse struct {
	Data struct {
		Data []CSESigmaRule `json:"data"`
	} `json:"data"`
}

type CSESigmaEntitySelector struct {
	EntityType string `json:"entity_type"`
	Expression string `json:"expression"`
}

type CSESigmaScoreMapping struct {
	Type    string            `json:"type"`
	Default int               `json:"default"`
	Field   string            `json:"field"`
	Mapping []json.RawMessage `json:"mapping"`
}

// CSESigmaRule is a converted Sigma rule, as returned by the Sigma rules API.
type CSESigmaRule struct {
	Name                  string                   `json:"name"`
	Enabled               bool                     `json:"enabled"`
	IsPrototype           bool                     `json:"is_prototype"`
	NameExpression        string                   `json:"name_expression"`
	DescriptionExpression string                   `json:"description_expression"`
	SummaryExpression     string                   `json:"summary_expression"`
	Expression            string                   `json:"expression"`
	Stream                string                   `json:"stream"`
	EntitySelectors       []CSESigmaEntitySelector `json:"entity_selectors"`
	ScoreMapping          CSESigmaScoreMapping     `json:"score_mapping"`
	Tags                  []string                 `json:"tags"`
	Category              string                   `json:"category"`
	ConversionMetadata    struct {
		SigmaUID string `json:"sigma_uid"`
	} `json:"conversion_metadata"`
}
