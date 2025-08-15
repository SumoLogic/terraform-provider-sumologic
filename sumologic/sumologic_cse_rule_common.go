package sumologic

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getEntitySelectorsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"expression": {
					Type:     schema.TypeString,
					Required: true,
				},
				"entity_type": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func getSeverityMappingSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
				},
				"default": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"field": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"mapping": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"type": {
								Type:     schema.TypeString,
								Required: true,
							},
							"from": {
								Type:     schema.TypeString,
								Required: true,
							},
							"to": {
								Type:     schema.TypeInt,
								Required: true,
							},
						},
					},
				},
			},
		},
	}
}

func resourceToStringArray(resourceStrings []interface{}) []string {
	result := make([]string, len(resourceStrings))

	for i, resourceString := range resourceStrings {
		result[i] = resourceString.(string)
	}

	return result
}

func resourceToEntitySelectorArray(resourceEntitySelectors []interface{}) []EntitySelector {
	result := make([]EntitySelector, len(resourceEntitySelectors))

	for i, resourceEntitySelector := range resourceEntitySelectors {
		result[i] = EntitySelector{
			EntityType: resourceEntitySelector.(map[string]interface{})["entity_type"].(string),
			Expression: resourceEntitySelector.(map[string]interface{})["expression"].(string),
		}
	}

	return result
}

func entitySelectorArrayToResource(entitySelectors []EntitySelector) []map[string]interface{} {
	result := make([]map[string]interface{}, len(entitySelectors))

	for i, entitySelector := range entitySelectors {
		result[i] = map[string]interface{}{
			"entity_type": entitySelector.EntityType,
			"expression":  entitySelector.Expression,
		}
	}

	return result
}

func resourceToSeverityMappingValueMappingArray(resourceSeverityMappingValueMappings []interface{}) []SeverityMappingValueMapping {
	result := make([]SeverityMappingValueMapping, len(resourceSeverityMappingValueMappings))

	for i, resourceSeverityMappingValueMapping := range resourceSeverityMappingValueMappings {
		severityMappingValueMappingMap := resourceSeverityMappingValueMapping.(map[string]interface{})
		result[i] = SeverityMappingValueMapping{
			Type: severityMappingValueMappingMap["type"].(string),
			From: severityMappingValueMappingMap["from"].(string),
			To:   severityMappingValueMappingMap["to"].(int),
		}
	}

	return result
}

func resourceToSeverityMapping(resourceSeverityMapping interface{}) SeverityMapping {
	resourceSeverityMappingMap := resourceSeverityMapping.(map[string]interface{})
	return SeverityMapping{
		Type:    resourceSeverityMappingMap["type"].(string),
		Default: resourceSeverityMappingMap["default"].(int),
		Field:   resourceSeverityMappingMap["field"].(string),
		Mapping: resourceToSeverityMappingValueMappingArray(resourceSeverityMappingMap["mapping"].([]interface{})),
	}
}

func severityMappingValueMappingArrayToResource(severityMappingValueMappings []SeverityMappingValueMapping) []map[string]interface{} {
	result := make([]map[string]interface{}, len(severityMappingValueMappings))

	for i, severityMappingValueMapping := range severityMappingValueMappings {
		result[i] = map[string]interface{}{
			"type": severityMappingValueMapping.Type,
			"from": severityMappingValueMapping.From,
			"to":   severityMappingValueMapping.To,
		}
	}

	return result
}

func severityMappingToResource(severityMapping SeverityMapping) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"type":    severityMapping.Type,
			"default": severityMapping.Default,
			"field":   severityMapping.Field,
			"mapping": severityMappingValueMappingArrayToResource(severityMapping.Mapping),
		},
	}
}

type EntitySelector struct {
	Expression string `json:"expression"`
	EntityType string `json:"entityType"`
}

// Use explicit windowSizeField type so that we can unmarshall Int value as String
type windowSizeField string

func (r *windowSizeField) UnmarshalJSON(data []byte) error {
	cleanedData := strings.Trim(string(data), `"`)
	*r = windowSizeField(cleanedData)
	return nil
}

func getRuleSource(id string, meta interface{}) string {
	c := meta.(*Client)

	ruleSource, err := c.GetCSERuleSource(id)
	if err != nil {
		log.Printf("[WARN] CSE Match Rule not found when looking by id: %s, err: %v", id, err)

	}

	return ruleSource
}

func (s *Client) GetCSERuleSource(id string) (string, error) {
	data, err := s.Get(fmt.Sprintf("sec/v1/rules/%s", id))
	if err != nil {
		return "", err
	}

	if data == nil {
		return "", nil
	}

	var response CSERuleSourceResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return "", err
	}

	return response.CSERuleSource.RuleSource, nil
}

type CSERuleSource struct {
	RuleSource string `json:"ruleSource"`
}

type CSERuleSourceResponse struct {
	CSERuleSource CSERuleSource `json:"data"`
}

func suppressSpaceDiff(_, old, new string, _ *schema.ResourceData) bool {
	oldNormalized := strings.TrimSpace(old)
	newNormalized := strings.TrimSpace(new)

	return oldNormalized == newNormalized
}
