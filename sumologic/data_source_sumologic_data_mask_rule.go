package sumologic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSumologicDataMaskRule() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceSumologicDataMaskRuleRead,
		Schema: dataSourceDataMaskRuleSchema(),
	}
}

func dataSourceDataMaskRuleSchema() map[string]*schema.Schema {
	schemaMap := dataSourceDataMaskRuleComputedSchema()
	schemaMap["id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	return schemaMap
}

func dataSourceDataMaskRuleComputedSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"regex_pattern": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"mask_string": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceSumologicDataMaskRuleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Get("id").(string)
	rule, err := c.GetDataMaskRule(id)
	if err != nil {
		return fmt.Errorf("error retrieving data mask rule with id %v: %v", id, err)
	}
	if rule == nil {
		return fmt.Errorf("data mask rule with id %v not found", id)
	}

	d.SetId(rule.ID)
	d.Set("name", rule.Name)
	d.Set("regex_pattern", rule.RegexPattern)
	d.Set("mask_string", rule.MaskString)
	d.Set("enabled", rule.Enabled)
	d.Set("description", rule.Description)

	return nil
}
