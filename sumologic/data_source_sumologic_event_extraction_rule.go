package sumologic

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func dataSourceEventExtractionRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEventExtractionRuleRead,
		Schema: map[string]*schema.Schema{
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":   {Type: schema.TypeString, Computed: true},
						"name": {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceEventExtractionRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	rules, err := client.GetAllEventExtractionRules()
	if err != nil {
		return err
	}

	out := []interface{}{}
	for _, r := range rules {
		out = append(out, map[string]interface{}{
			"id":   r.ID,
			"name": r.Name,
		})
	}

	d.SetId("event_extraction_rules")
	return d.Set("rules", out)
}
