package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceSumologicCSEEntityCriticalityConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSEEntityCriticalityConfigCreate,
		Read:   resourceSumologicCSEEntityCriticalityConfigRead,
		Delete: resourceSumologicCSEEntityCriticalityConfigDelete,
		Update: resourceSumologicCSEEntityCriticalityConfigUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"severity_expression": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourceSumologicCSEEntityCriticalityConfigRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSEEntityCriticalityConfig *CSEEntityCriticalityConfig
	id := d.Id()

	CSEEntityCriticalityConfig, err := c.GetCSEEntityCriticalityConfig(id)
	if err != nil {
		log.Printf("[WARN] CSE Entity Criticality Config not found when looking by id: %s, err: %v", id, err)

	}

	if CSEEntityCriticalityConfig == nil {
		log.Printf("[WARN] CSE Entity Criticality Config not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("name", CSEEntityCriticalityConfig.Name)
	d.Set("severity_expression", CSEEntityCriticalityConfig.SeverityExpression)

	return nil
}

func resourceSumologicCSEEntityCriticalityConfigDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteCSEEntityCriticalityConfig(d.Id())

}

func resourceSumologicCSEEntityCriticalityConfigCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		id, err := c.CreateCSEEntityCriticalityConfig(CSEEntityCriticalityConfig{
			Name:               d.Get("name").(string),
			SeverityExpression: d.Get("severity_expression").(string),
		})

		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicCSEEntityCriticalityConfigUpdate(d, meta)
}

func resourceSumologicCSEEntityCriticalityConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	CSEEntityCriticalityConfig, err := resourceToCSEEntityCriticalityConfig(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSEEntityCriticalityConfig(CSEEntityCriticalityConfig); err != nil {
		return err
	}

	return resourceSumologicCSEEntityCriticalityConfigRead(d, meta)
}

func resourceToCSEEntityCriticalityConfig(d *schema.ResourceData) (CSEEntityCriticalityConfig, error) {
	id := d.Id()
	if id == "" {
		return CSEEntityCriticalityConfig{}, nil
	}

	return CSEEntityCriticalityConfig{
		ID:                 id,
		SeverityExpression: d.Get("severity_expression").(string),
		Name:               d.Get("name").(string),
	}, nil
}
