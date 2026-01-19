package sumologic

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSumologicEventExtractionRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicEventExtractionRuleCreate,
		Read:   resourceSumologicEventExtractionRuleRead,
		Update: resourceSumologicEventExtractionRuleUpdate,
		Delete: resourceSumologicEventExtractionRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of event extraction rule.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of event extraction rule.",
			},
			"query": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Query string for the Event Extraction Rule.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Flag indicating whether the event extraction rule is enabled or disabled.",
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
								[]string{"ExactMatch"},
								false,
							),
						},
					},
				},
			},

			"configuration": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Structured field mappings for the extraction rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the field being mapped.",
						},
						"value_source": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mapping_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "HardCoded",
							ValidateFunc: validation.StringInSlice([]string{"HardCoded"}, false),
						},
					},
				},
			},
		},
	}
}

func resourceSumologicEventExtractionRuleCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	rule := resourceToEventExtractionRule(d)
	id, err := c.CreateEventExtractionRule(rule)
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceSumologicEventExtractionRuleRead(d, meta)
}

func resourceSumologicEventExtractionRuleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	rule, err := c.GetEventExtractionRule(d.Id())
	if err != nil {
		return err
	}

	if rule == nil {
		log.Printf("[WARN] Event Extraction Rule not found: %s", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("name", rule.Name)
	d.Set("description", rule.Description)
	d.Set("query", rule.Query)
	d.Set("enabled", rule.Enabled)

	if rule.CorrelationExpression != nil {
		d.Set("correlation_expression", flattenCorrelationExpression(rule.CorrelationExpression))
	}

	if err := d.Set("configuration", flattenConfiguration(rule.Configuration)); err != nil {
		return fmt.Errorf("error setting configuration: %w", err)
	}

	return nil
}

func resourceSumologicEventExtractionRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	rule := resourceToEventExtractionRule(d)
	rule.ID = d.Id()

	if err := c.UpdateEventExtractionRule(rule); err != nil {
		return err
	}

	return resourceSumologicEventExtractionRuleRead(d, meta)
}

func resourceSumologicEventExtractionRuleDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	return c.DeleteEventExtractionRule(d.Id())
}

func resourceToEventExtractionRule(d *schema.ResourceData) EventExtractionRule {
	rule := EventExtractionRule{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Query:       d.Get("query").(string),
		Enabled:     d.Get("enabled").(bool),
	}

	if v, ok := d.GetOk("correlation_expression"); ok {
		item := v.([]interface{})[0].(map[string]interface{})
		rule.CorrelationExpression = &CorrelationExpression{
			QueryFieldName:          item["query_field_name"].(string),
			EventFieldName:          item["event_field_name"].(string),
			StringMatchingAlgorithm: item["string_matching_algorithm"].(string),
		}
	}

	config := make(map[string]FieldMapping)
	if v, ok := d.GetOk("configuration"); ok {
		configs := v.([]interface{})
		for _, raw := range configs {
			val := raw.(map[string]interface{})
			fieldName := val["field_name"].(string)
			config[fieldName] = FieldMapping{
				ValueSource: val["value_source"].(string),
				MappingType: val["mapping_type"].(string),
			}
		}
	}
	rule.Configuration = config

	return rule
}

func flattenCorrelationExpression(ce *CorrelationExpression) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"query_field_name":          ce.QueryFieldName,
			"event_field_name":          ce.EventFieldName,
			"string_matching_algorithm": ce.StringMatchingAlgorithm,
		},
	}
}

func flattenConfiguration(config map[string]FieldMapping) []interface{} {
	var result []interface{}
	for key, val := range config {
		result = append(result, map[string]interface{}{
			"field_name":   key,
			"value_source": val.ValueSource,
			"mapping_type": val.MappingType,
		})
	}
	return result
}
