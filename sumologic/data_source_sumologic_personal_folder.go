package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSumologicPersonalFolder() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumologicPersonalFolderRead,

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
	}
}

func dataSourceSumologicPersonalFolderRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	personalFolder, err := c.getPersonalFolder()

	if err != nil {
		return err
	}

	d.SetId(personalFolder.ID)
	d.Set("name", personalFolder.Name)
	d.Set("description", personalFolder.Description)

	return nil
}
