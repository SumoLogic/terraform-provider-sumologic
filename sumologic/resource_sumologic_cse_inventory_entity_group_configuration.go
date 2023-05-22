package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

func resourceSumologicCSEInventoryEntityGroupConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSEInventoryEntityGroupConfigurationCreate,
		Read:   resourceSumologicCSEInventoryEntityGroupConfigurationRead,
		Delete: resourceSumologicCSEInventoryEntityGroupConfigurationDelete,
		Update: resourceSumologicCSEInventoryEntityGroupConfigurationUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"criticality": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"group": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"inventory_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"user", "computer"}, false),
			},
			"inventory_source": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"suppressed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSumologicCSEInventoryEntityGroupConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSEInventoryEntityGroupConfigurationGet *CSEEntityGroupConfiguration
	id := d.Id()

	CSEInventoryEntityGroupConfigurationGet, err := c.GetCSEntityGroupConfiguration(id)
	if err != nil {
		log.Printf("[WARN] CSE Inventory Entity Group Configuration not found when looking by id: %s, err: %v", id, err)
	}

	if CSEInventoryEntityGroupConfigurationGet == nil {
		log.Printf("[WARN] CSE Inventory Entity Group Configuration not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("criticality", CSEInventoryEntityGroupConfigurationGet.Criticality)
	d.Set("description", CSEInventoryEntityGroupConfigurationGet.Description)
	d.Set("group", CSEInventoryEntityGroupConfigurationGet.Group)
	d.Set("inventory_type", CSEInventoryEntityGroupConfigurationGet.InventoryType)
	d.Set("inventory_source", CSEInventoryEntityGroupConfigurationGet.InventorySource)
	d.Set("name", CSEInventoryEntityGroupConfigurationGet.Name)
	d.Set("suppressed", CSEInventoryEntityGroupConfigurationGet.Suppressed)
	d.Set("tags", CSEInventoryEntityGroupConfigurationGet.Tags)

	return nil
}

func resourceSumologicCSEInventoryEntityGroupConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteCSEEntityGroupConfiguration(d.Id())

}

func resourceSumologicCSEInventoryEntityGroupConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		id, err := c.CreateCSEInventoryEntityGroupConfiguration(CSEEntityGroupConfiguration{
			Criticality:     d.Get("criticality").(string),
			Description:     d.Get("description").(string),
			Group:           d.Get("group").(string),
			InventoryType:   d.Get("inventory_type").(string),
			InventorySource: d.Get("inventory_source").(string),
			Name:            d.Get("name").(string),
			Suppressed:      d.Get("suppressed").(bool),
			Tags:            resourceToStringArray(d.Get("tags").([]interface{})),
		})

		if err != nil {
			return err
		}
		d.SetId(id)
	}

	return resourceSumologicCSEInventoryEntityGroupConfigurationRead(d, meta)
}

func resourceSumologicCSEInventoryEntityGroupConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	CSEInventoryEntityGroupConfiguration, err := resourceToCSEInventoryEntityGroupConfiguration(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSEInventoryEntityGroupConfiguration(CSEInventoryEntityGroupConfiguration); err != nil {
		return err
	}

	return resourceSumologicCSEInventoryEntityGroupConfigurationRead(d, meta)
}

func resourceToCSEInventoryEntityGroupConfiguration(d *schema.ResourceData) (CSEEntityGroupConfiguration, error) {
	id := d.Id()
	if id == "" {
		return CSEEntityGroupConfiguration{}, nil
	}

	return CSEEntityGroupConfiguration{
		ID:              id,
		Criticality:     d.Get("criticality").(string),
		Description:     d.Get("description").(string),
		Group:           d.Get("group").(string),
		InventoryType:   d.Get("inventory_type").(string),
		InventorySource: d.Get("inventory_source").(string),
		Name:            d.Get("name").(string),
		Suppressed:      d.Get("suppressed").(bool),
		Tags:            resourceToStringArray(d.Get("tags").([]interface{})),
	}, nil
}
