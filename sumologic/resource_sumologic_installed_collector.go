package sumologic

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSumologicInstalledCollector() *schema.Resource {
	return &schema.Resource{
		Read:   resourceSumologicCollectorRead,
		Delete: resourceSumologicCollectorDelete,
		Update: resourceSumologicInstalledCollectorUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Etc/UTC",
			},
			"fields": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func resourceSumologicInstalledCollectorUpdate(d *schema.ResourceData, meta interface{}) error {

	collector := resourceToInstalledCollector(d)

	c := meta.(*Client)
	err := c.UpdateCollector(collector)

	if err != nil {
		return err
	}

	return resourceSumologicCollectorRead(d, meta)
}

func resourceToInstalledCollector(d *schema.ResourceData) Collector {
	id, _ := strconv.Atoi(d.Id())

	return Collector{
		ID:            id,
		CollectorType: "Installable",
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Category:      d.Get("category").(string),
		TimeZone:      d.Get("timezone").(string),
		Fields:        d.Get("fields").(map[string]interface{}),
	}
}
