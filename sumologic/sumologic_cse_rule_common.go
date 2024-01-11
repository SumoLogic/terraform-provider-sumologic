package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

// Use explicit windowSizeField type so that we can ignore it when unmarshalling
type windowSizeField string

func (r windowSizeField) UnmarshalJSON(data []byte) error {
	r = ""
	return nil
}
