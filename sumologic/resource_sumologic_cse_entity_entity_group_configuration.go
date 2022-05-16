package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceSumologicCSEEntityEntityGroupConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSEEntityEntityGroupConfigurationCreate,
		Read:   resourceSumologicCSEEntityEntityGroupConfigurationRead,
		Delete: resourceSumologicCSEEntityEntityGroupConfigurationDelete,
		Update: resourceSumologicCSEEntityEntityGroupConfigurationUpdate,
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
			"entity_namespace": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"entity_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_block": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"suffix": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
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

func resourceSumologicCSEEntityEntityGroupConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSEEntityEntityGroupConfigurationGet *CSEEntityGroupConfiguration
	id := d.Id()

	CSEEntityEntityGroupConfigurationGet, err := c.GetCSEntityGroupConfiguration(id)
	if err != nil {
		log.Printf("[WARN] CSE Entity Entity Group Configuration not found when looking by id: %s, err: %v", id, err)
	}

	if CSEEntityEntityGroupConfigurationGet == nil {
		log.Printf("[WARN] CSE Entity Entity Group Configuration not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("criticality", CSEEntityEntityGroupConfigurationGet.Criticality)
	d.Set("description", CSEEntityEntityGroupConfigurationGet.Description)
	d.Set("entity_namespace", CSEEntityEntityGroupConfigurationGet.EntityNamespace)
	d.Set("entity_type", CSEEntityEntityGroupConfigurationGet.EntityType)
	d.Set("name", CSEEntityEntityGroupConfigurationGet.Name)
	d.Set("network_block", CSEEntityEntityGroupConfigurationGet.NetworkBlock)
	d.Set("prefix", CSEEntityEntityGroupConfigurationGet.Prefix)
	d.Set("suffix", CSEEntityEntityGroupConfigurationGet.Suffix)
	d.Set("suppressed", CSEEntityEntityGroupConfigurationGet.Suppressed)
	d.Set("tags", CSEEntityEntityGroupConfigurationGet.Tags)

	return nil
}

func resourceSumologicCSEEntityEntityGroupConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteCSEEntityGroupConfiguration(d.Id())

}

func resourceSumologicCSEEntityEntityGroupConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		id, err := c.CreateCSEEntityEntityGroupConfiguration(CSEEntityGroupConfiguration{
			Criticality:     d.Get("criticality").(string),
			Description:     d.Get("description").(string),
			EntityNamespace: d.Get("entity_namespace").(string),
			EntityType:      d.Get("entity_type").(string),
			Name:            d.Get("name").(string),
			NetworkBlock:    d.Get("network_block").(string),
			Prefix:          d.Get("prefix").(string),
			Suffix:          d.Get("suffix").(string),
			Suppressed:      d.Get("suppressed").(bool),
			Tags:            resourceToStringArray(d.Get("tags").([]interface{})),
		})

		if err != nil {
			return err
		}
		d.SetId(id)
	}

	return resourceSumologicCSEEntityEntityGroupConfigurationRead(d, meta)
}

func resourceSumologicCSEEntityEntityGroupConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	CSEEntityEntityGroupConfiguration, err := resourceToCSEEntityEntityGroupConfiguration(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSEEntityEntityGroupConfiguration(CSEEntityEntityGroupConfiguration); err != nil {
		return err
	}

	return resourceSumologicCSEEntityEntityGroupConfigurationRead(d, meta)
}

func resourceToCSEEntityEntityGroupConfiguration(d *schema.ResourceData) (CSEEntityGroupConfiguration, error) {
	id := d.Id()
	if id == "" {
		return CSEEntityGroupConfiguration{}, nil
	}

	return CSEEntityGroupConfiguration{
		ID:              id,
		Criticality:     d.Get("criticality").(string),
		Description:     d.Get("description").(string),
		EntityNamespace: d.Get("entity_namespace").(string),
		EntityType:      d.Get("entity_type").(string),
		Name:            d.Get("name").(string),
		NetworkBlock:    d.Get("network_block").(string),
		Prefix:          d.Get("prefix").(string),
		Suffix:          d.Get("suffix").(string),
		Suppressed:      d.Get("suppressed").(bool),
		Tags:            resourceToStringArray(d.Get("tags").([]interface{})),
	}, nil
}
