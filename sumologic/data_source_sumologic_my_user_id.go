package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSumologicMyUserId() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumologicMyUserIdRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceSumologicMyUserIdRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	personalFolder, err := c.getPersonalFolder()

	if err != nil {
		return err
	}

	d.SetId(personalFolder.CreatedBy)

	return nil
}
