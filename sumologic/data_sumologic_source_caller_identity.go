package sumologic

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
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

	d.SetId(time.Now().UTC().String())
	d.Set("access_id", c.AccessID)
	d.Set("environment", c.Environment)

	return nil
}
