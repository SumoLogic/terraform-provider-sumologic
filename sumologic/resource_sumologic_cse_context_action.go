package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
)

func resourceSumologicCSEContextAction() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSEContextActionCreate,
		Read:   resourceSumologicCSEContextActionRead,
		Delete: resourceSumologicCSEContextActionDelete,
		Update: resourceSumologicCSEContextActionUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringIsNotEmpty,
					validation.StringInSlice([]string{"URL", "QUERY"}, false)),
			},
			"template": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ioc_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.All(validation.StringIsNotEmpty, validation.StringInSlice([]string{"ASN", "DOMAIN", "HASH", "IP_ADDRESS", "MAC_ADDRESS", "PORT", "RECORD_PROPERTY", "URL"}, false)),
				},
			},
			"entity_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
			"record_fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
			"all_record_fields": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceSumologicCSEContextActionRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSEContextAction *CSEContextAction
	id := d.Id()

	CSEContextAction, err := c.GetCSEContextAction(id)
	if err != nil {
		log.Printf("[WARN] CSE Context Action not found when looking by id: %s, err: %v", id, err)

	}

	if CSEContextAction == nil {
		log.Printf("[WARN] CSE Context Action not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("name", CSEContextAction.Name)
	d.Set("type", CSEContextAction.Type)
	d.Set("template", CSEContextAction.Template)
	d.Set("ioc_types", CSEContextAction.IocTypes)
	d.Set("entity_types", CSEContextAction.EntityTypes)
	d.Set("record_fields", CSEContextAction.RecordFields)
	d.Set("all_record_fields", CSEContextAction.AllRecordFields)
	d.Set("enabled", CSEContextAction.Enabled)

	return nil
}

func resourceSumologicCSEContextActionDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteCSEContextAction(d.Id())

}

func resourceSumologicCSEContextActionCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		id, err := c.CreateCSEContextAction(CSEContextAction{
			Name:            d.Get("name").(string),
			Type:            d.Get("type").(string),
			Template:        d.Get("template").(string),
			IocTypes:        resourceFieldsToStringArray(d.Get("ioc_types").([]interface{})),
			EntityTypes:     resourceFieldsToStringArray(d.Get("entity_types").([]interface{})),
			RecordFields:    resourceFieldsToStringArray(d.Get("record_fields").([]interface{})),
			AllRecordFields: d.Get("all_record_fields").(bool),
			Enabled:         d.Get("enabled").(bool),
		})

		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicCSEContextActionRead(d, meta)
}

func resourceSumologicCSEContextActionUpdate(d *schema.ResourceData, meta interface{}) error {
	CSEContextAction, err := resourceToCSEContextAction(d)
	if err != nil {
		return err
	}
	c := meta.(*Client)
	if err = c.UpdateCSEContextAction(CSEContextAction); err != nil {
		return err
	}

	return resourceSumologicCSEContextActionRead(d, meta)
}

func resourceToCSEContextAction(d *schema.ResourceData) (CSEContextAction, error) {
	id := d.Id()
	if id == "" {
		return CSEContextAction{}, nil
	}

	return CSEContextAction{
		ID:              id,
		Name:            d.Get("name").(string),
		Type:            d.Get("type").(string),
		Template:        d.Get("template").(string),
		IocTypes:        resourceFieldsToStringArray(d.Get("ioc_types").([]interface{})),
		EntityTypes:     resourceFieldsToStringArray(d.Get("entity_types").([]interface{})),
		RecordFields:    resourceFieldsToStringArray(d.Get("record_fields").([]interface{})),
		AllRecordFields: d.Get("all_record_fields").(bool),
		Enabled:         d.Get("enabled").(bool),
	}, nil
}
