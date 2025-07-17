package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
)

func resourceSumologicCSEAutomation() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSEAutomationCreate,
		Read:   resourceSumologicCSEAutomationRead,
		Delete: resourceSumologicCSEAutomationDelete,
		Update: resourceSumologicCSEAutomationUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"playbook_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cse_resource_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cse_resource_sub_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"execution_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceSumologicCSEAutomationRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSEAutomation *CSEAutomation
	id := d.Id()

	CSEAutomation, err := c.GetCSEAutomation(id)
	if err != nil {
		log.Printf("[WARN] CSE Automation not found when looking by id: %s, err: %v", id, err)

	}

	if CSEAutomation == nil {
		log.Printf("[WARN] CSE Automation not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("name", CSEAutomation.Name)
	d.Set("description", CSEAutomation.Description)
	d.Set("playbook_id", CSEAutomation.PlaybookId)
	d.Set("cse_resource_type", CSEAutomation.CseResourceType)
	d.Set("cse_resource_sub_types", CSEAutomation.CseResourceSubTypes)
	d.Set("execution_types", CSEAutomation.ExecutionTypes)
	d.Set("enabled", CSEAutomation.Enabled)

	return nil
}

func resourceSumologicCSEAutomationDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteCSEAutomation(d.Id())

}

func resourceSumologicCSEAutomationCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		id, err := c.CreateCSEAutomation(CSEAutomation{
			PlaybookId:          d.Get("playbook_id").(string),
			CseResourceType:     d.Get("cse_resource_type").(string),
			CseResourceSubTypes: resourceFieldsToStringArray(d.Get("cse_resource_sub_types").([]interface{})),
			ExecutionTypes:      resourceFieldsToStringArray(d.Get("execution_types").([]interface{})),
			Enabled:             d.Get("enabled").(bool),
		})

		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicCSEAutomationRead(d, meta)
}

func resourceSumologicCSEAutomationUpdate(d *schema.ResourceData, meta interface{}) error {
	CSEAutomation, err := resourceToCSEAutomation(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSEAutomation(CSEAutomation); err != nil {
		return err
	}

	return resourceSumologicCSEAutomationRead(d, meta)
}

func resourceToCSEAutomation(d *schema.ResourceData) (CSEAutomation, error) {
	id := d.Id()
	if id == "" {
		return CSEAutomation{}, nil
	}

	return CSEAutomation{
		ID:                  id,
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		PlaybookId:          d.Get("playbook_id").(string),
		CseResourceType:     d.Get("cse_resource_type").(string),
		CseResourceSubTypes: resourceFieldsToStringArray(d.Get("cse_resource_sub_types").([]interface{})),
		ExecutionTypes:      resourceFieldsToStringArray(d.Get("execution_types").([]interface{})),
		Enabled:             d.Get("enabled").(bool),
	}, nil
}
