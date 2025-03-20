package sumologic

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceSumologicOTCollector() *schema.Resource {
	return &schema.Resource{
		Read:   resourceSumologicOTCollectorRead,
		Delete: resourceSumologicOTCollectorDelete,
		Update: resourceSumologicOTCollectorUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"modified_by": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"alive": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"is_remotely_managed": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"modified_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"ephemeral": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"tags": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func resourceSumologicOTCollectorRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	otCollector, err := c.GetOTCollector(id)
	if err != nil {
		return err
	}

	if otCollector == nil {
		log.Printf("[WARN] OTCollector not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	if err := d.Set("tags", otCollector.Tags); err != nil {
		return fmt.Errorf("error setting fields for resource %s: %s", d.Id(), err)
	}

	d.Set("alive", otCollector.IsAlive)
	d.Set("created_at", otCollector.CreatedAt)
	d.Set("description", otCollector.Description)
	d.Set("time_zone", otCollector.TimeZone)
	d.Set("category", otCollector.Category)
	d.Set("modified_at", otCollector.ModifiedAt)
	d.Set("is_remotely_managed", otCollector.IsRemotelyManaged)
	d.Set("created_by", otCollector.CreatedBy)
	d.Set("modified_by", otCollector.ModifiedBy)
	d.Set("ephemeral", otCollector.Ephemeral)
	d.Set("name", otCollector.Name)

	return nil
}

func resourceSumologicOTCollectorDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteOTCollector(d.Id())
}

func resourceSumologicOTCollectorUpdate(d *schema.ResourceData, meta interface{}) error {
	return errors.New("Terraform does not support OTel collector updates")
}

func resourceToOTCollector(d *schema.ResourceData) OTCollector {

	return OTCollector{
		IsAlive:           d.Get("alive").(bool),
		CreatedBy:         d.Get("created_by").(string),
		Ephemeral:         d.Get("ephemeral").(bool),
		CreatedAt:         d.Get("created_at").(string),
		Description:       d.Get("description").(string),
		ModifiedBy:        d.Get("modified_by").(string),
		ModifiedAt:        d.Get("modified_at").(string),
		Category:          d.Get("category").(string),
		IsRemotelyManaged: d.Get("is_remotely_managed").(bool),
		TimeZone:          d.Get("time_zone").(string),
		ID:                d.Id(),
		Name:              d.Get("name").(string),
		Tags:              d.Get("fields").(map[string]interface{}),
	}
}
