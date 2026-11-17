package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSumologicUsageAttributionConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumologicUsageAttributionConfigRead,

		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"custom_tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"standard_tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"breakdown_by_options": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dashboard_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSumologicUsageAttributionConfigRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	config, err := c.GetUsageAttributionConfig()
	if err != nil {
		return err
	}

	d.SetId(usageAttributionConfigId)
	return setUsageAttributionConfigResource(d, config)
}
