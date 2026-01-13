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
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of event extraction rule.",
				ValidateFunc: validation.StringLenBetween(1, 256),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Description of event extraction rule.",
				ValidateFunc: validation.StringLenBetween(0, 1024),
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
							ValidateFunc: validation.StringInSlice([]string{
								"ExactMatch",
							}, false),
						},
					},
				},
			},
			"field_mapping": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value_source": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 256),
						},
						"mapping_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "HardCoded",
							ValidateFunc: validation.StringInSlice([]string{
								"HardCoded",
							}, false),
						},
					},
				},
			},
		},
	}
}

func resourceSumologicEventExtractionRuleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	rule, err := c.GetEventExtractionRule(id)
	if err != nil {
		return err
	}

	if rule == nil {
		log.Printf("[WARN] Event Extraction Rule not found, removing from state: %v", id)
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

	if err := d.Set("field_mapping", flattenConfigurationMap(rule.Configuration)); err != nil {
		return fmt.Errorf("error setting field_mapping for resource %s: %s", id, err)
	}

	return nil
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

func resourceSumologicEventExtractionRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	rule := resourceToEventExtractionRule(d)
	rule.ID = d.Id()

	err := c.UpdateEventExtractionRule(rule)
	if err != nil {
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
		list := v.([]interface{})
		if len(list) > 0 {
			item := list[0].(map[string]interface{})
			rule.CorrelationExpression = &CorrelationExpression{
				QueryFieldName:          item["query_field_name"].(string),
				EventFieldName:          item["event_field_name"].(string),
				StringMatchingAlgorithm: item["string_matching_algorithm"].(string),
			}
		}
	}

	if v, ok := d.GetOk("field_mapping"); ok {
		mappingsSet := v.(*schema.Set).List()
		configMap := make(map[string]FieldMapping)
		for _, m := range mappingsSet {
			item := m.(map[string]interface{})
			key := item["field_name"].(string)
			configMap[key] = FieldMapping{
				ValueSource: item["value_source"].(string),
				MappingType: item["mapping_type"].(string),
			}
		}
		rule.Configuration = configMap
	}

	return rule
}

func flattenCorrelationExpression(ce *CorrelationExpression) []interface{} {
	if ce == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"query_field_name":          ce.QueryFieldName,
			"event_field_name":          ce.EventFieldName,
			"string_matching_algorithm": ce.StringMatchingAlgorithm,
		},
	}
}

func flattenConfigurationMap(config map[string]FieldMapping) []interface{} {
	result := []interface{}{}
	for key, val := range config {
		result = append(result, map[string]interface{}{
			"field_name":   key,
			"value_source": val.ValueSource,
			"mapping_type": val.MappingType,
		})
	}
	return result
}
