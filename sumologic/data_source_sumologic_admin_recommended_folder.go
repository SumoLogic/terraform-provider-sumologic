package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

func dataSourceSumologicAdminRecommendedFolder() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumologicAdminRecommendedFolderRead,

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
				Optional: true,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(1 * time.Minute),
		},
	}
}

func dataSourceSumologicAdminRecommendedFolderRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	adminRecommendedFolder, err := c.getAdminRecommendedFolder(d.Timeout(schema.TimeoutRead))

	if err != nil {
		return err
	}

	d.SetId(adminRecommendedFolder.ID)
	d.Set("name", adminRecommendedFolder.Name)
	d.Set("description", adminRecommendedFolder.Description)

	return nil
}
