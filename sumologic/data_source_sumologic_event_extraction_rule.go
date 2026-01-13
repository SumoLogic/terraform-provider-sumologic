package sumologic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSumologicEventExtractionRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumologicEventExtractionRuleRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"correlation_expression": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query_field_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_field_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"string_matching_algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"field_mapping": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value_source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mapping_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSumologicEventExtractionRuleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var rule *EventExtractionRule
	var err error

	if id, ok := d.GetOk("id"); ok {
		rule, err = c.GetEventExtractionRule(id.(string))
	} else if name, ok := d.GetOk("name"); ok {
		rule, err = c.GetEventExtractionRuleByName(name.(string))
	} else {
		return fmt.Errorf("please specify either id or name")
	}

	if err != nil {
		return err
	}
	if rule == nil {
		return fmt.Errorf("event extraction rule not found")
	}

	d.SetId(rule.ID)
	d.Set("name", rule.Name)
	d.Set("description", rule.Description)
	d.Set("query", rule.Query)
	d.Set("enabled", rule.Enabled)

	if rule.CorrelationExpression != nil {
		d.Set("correlation_expression", flattenCorrelationExpression(rule.CorrelationExpression))
	}
	if err := d.Set("field_mapping", flattenConfigurationMap(rule.Configuration)); err != nil {
		return fmt.Errorf("error setting field_mapping for datasource %s: %s", d.Id(), err)
	}

	return nil
}
