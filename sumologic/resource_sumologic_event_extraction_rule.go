package sumologic

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceEventExtractionRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceEventExtractionRuleCreate,
		Read:   resourceEventExtractionRuleRead,
		Update: resourceEventExtractionRuleUpdate,
		Delete: resourceEventExtractionRuleDelete,

		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			cfg := d.Get("configuration").(map[string]interface{})

			requiredKeys := []string{
				"eventType",
				"eventPriority",
				"eventSource",
				"eventName",
			}

			for _, k := range requiredKeys {
				if _, ok := cfg[k]; !ok {
					return fmt.Errorf("configuration.%s is required", k)
				}
			}
			return nil
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"query": {
				Type:     schema.TypeString,
				Required: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"correlation_expression": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query_field_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"event_field_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"string_matching_algorithm": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								[]string{"ExactMatch"}, false,
							),
						},
					},
				},
			},

			"configuration": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value_source": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mapping_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice(
								[]string{"HardCoded"}, false,
							),
						},
					},
				},
			},
		},
	}
}

func resourceEventExtractionRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	rule := expandEventExtractionRule(d)
	created, err := client.CreateEventExtractionRule(rule)
	if err != nil {
		return err
	}

	d.SetId(created.ID)
	return resourceEventExtractionRuleRead(d, meta)
}

func resourceEventExtractionRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	rule, err := client.GetEventExtractionRule(d.Id())
	if err != nil || rule == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", rule.Name)
	d.Set("description", rule.Description)
	d.Set("query", rule.Query)
	d.Set("enabled", rule.Enabled)
	d.Set("correlation_expression", flattenCorrelationExpression(rule.CorrelationExpression))
	d.Set("configuration", flattenConfiguration(rule.Configuration))

	return nil
}

func resourceEventExtractionRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	if err := client.UpdateEventExtractionRule(d.Id(), expandEventExtractionRule(d)); err != nil {
		return err
	}

	return resourceEventExtractionRuleRead(d, meta)
}

func resourceEventExtractionRuleDelete(d *schema.ResourceData, meta interface{}) error {
	return meta.(*Client).DeleteEventExtractionRule(d.Id())
}

/*
========================
Expand / Flatten
========================
*/

func expandEventExtractionRule(d *schema.ResourceData) EventExtractionRule {
	return EventExtractionRule{
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		Query:                 d.Get("query").(string),
		Enabled:               d.Get("enabled").(bool),
		CorrelationExpression: expandCorrelationExpression(d),
		Configuration:         expandConfiguration(d),
	}
}

func expandCorrelationExpression(d *schema.ResourceData) *CorrelationExpression {
	v := d.Get("correlation_expression").([]interface{})
	if len(v) == 0 {
		return nil
	}

	m := v[0].(map[string]interface{})
	return &CorrelationExpression{
		QueryFieldName:          m["query_field_name"].(string),
		EventFieldName:          m["event_field_name"].(string),
		StringMatchingAlgorithm: m["string_matching_algorithm"].(string),
	}
}

func expandConfiguration(d *schema.ResourceData) map[string]FieldMapping {
	cfg := map[string]FieldMapping{}

	for k, v := range d.Get("configuration").(map[string]interface{}) {
		m := v.(map[string]interface{})

		fm := FieldMapping{
			ValueSource: m["value_source"].(string),
		}

		if mt, ok := m["mapping_type"]; ok {
			fm.MappingType = mt.(string)
		}

		cfg[k] = fm
	}

	return cfg
}

func flattenCorrelationExpression(c *CorrelationExpression) []interface{} {
	if c == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"query_field_name":          c.QueryFieldName,
			"event_field_name":          c.EventFieldName,
			"string_matching_algorithm": c.StringMatchingAlgorithm,
		},
	}
}

func flattenConfiguration(cfg map[string]FieldMapping) map[string]interface{} {
	out := map[string]interface{}{}

	for k, v := range cfg {
		m := map[string]interface{}{
			"value_source": v.ValueSource,
		}

		if v.MappingType != "" {
			m["mapping_type"] = v.MappingType
		}

		out[k] = m
	}

	return out
}
