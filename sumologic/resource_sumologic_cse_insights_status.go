package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceSumologicCSEInsightsStatus() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSEInsightsStatusCreate,
		Read:   resourceSumologicCSEInsightsStatusRead,
		Delete: resourceSumologicCSEInsightsStatusDelete,
		Update: resourceSumologicCSEInsightsStatusUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSumologicCSEInsightsStatusRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSEInsightsStatusGet *CSEInsightsStatusGet
	id := d.Id()

	CSEInsightsStatusGet, err := c.GetCSEInsightsStatus(id)
	if err != nil {
		log.Printf("[WARN] CSE Insights Status not found when looking by id: %s, err: %v", id, err)

	}

	if CSEInsightsStatusGet == nil {
		log.Printf("[WARN] CSE Insights Status not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("name", CSEInsightsStatusGet.Name)
	d.Set("description", CSEInsightsStatusGet.Description)
	d.Set("display_name", CSEInsightsStatusGet.DisplayName)

	return nil
}

func resourceSumologicCSEInsightsStatusDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteCSEInsightsStatus(d.Id())

}

func resourceSumologicCSEInsightsStatusCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		id, err := c.CreateCSEInsightsStatus(CSEInsightsStatusPost{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		})

		if err != nil {
			return err
		}
		d.SetId(id)
	}

	return resourceSumologicCSEInsightsStatusUpdate(d, meta)
}

func resourceSumologicCSEInsightsStatusUpdate(d *schema.ResourceData, meta interface{}) error {
	CSEInsightsStatusPost, err := resourceToCSEInsightsStatus(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSEInsightsStatus(CSEInsightsStatusPost); err != nil {
		return err
	}

	return resourceSumologicCSEInsightsStatusRead(d, meta)
}

func resourceToCSEInsightsStatus(d *schema.ResourceData) (CSEInsightsStatusPost, error) {
	id := d.Id()
	if id == "" {
		return CSEInsightsStatusPost{}, nil
	}

	return CSEInsightsStatusPost{
		ID:          id,
		Description: d.Get("description").(string),
		Name:        d.Get("name").(string),
	}, nil
}
