package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSumologicCallerIdentity() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumologicCallerIdentityRead,

		Schema: map[string]*schema.Schema{
			"access_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"environment": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSumologicCallerIdentityRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	d.SetId("access_id")
	d.Set("access_id", c.AccessID)
	d.Set("environment", c.Environment)

	return nil
}
