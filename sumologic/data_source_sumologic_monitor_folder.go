package sumologic

import (
    "fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSumologicMonitorFolder() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumologicMonitorFolderRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
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

func dataSourceSumologicMonitorFolderRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	monitorsLibraryFolder, err := c.GetMonitorsLibraryFolderByPath(d.Get("path").(string))

	if err != nil {
		return err
	}

    if monitorsLibraryFolder == nil || monitorsLibraryFolder.ID == "" {
        return fmt.Errorf("folder with path '%s' does not exist", d.Get("path").(string))
    }

    d.SetId(monitorsLibraryFolder.ID)
	d.Set("name", monitorsLibraryFolder.Name)
	d.Set("description", monitorsLibraryFolder.Description)

	return nil
}
