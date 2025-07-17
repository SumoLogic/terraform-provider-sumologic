package sumologic

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
)

func resourceSumologicCSECustomMatchListColumn() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSECustomMatchListColumnCreate,
		Read:   resourceSumologicCSECustomMatchListColumnRead,
		Delete: resourceSumologicCSECustomMatchListColumnDelete,
		Update: resourceSumologicCSECustomMatchListColumnUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"fields": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func resourceSumologicCSECustomMatchListColumnRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSECustomMatchListColumn *CSECustomMatchListColumn
	id := d.Id()

	CSECustomMatchListColumn, err := c.GetCSECustomMatchListColumn(id)
	if err != nil {
		log.Printf("[WARN] CSE Custom Match List Column not found when looking by id: %s, err: %v", id, err)

	}

	if CSECustomMatchListColumn == nil {
		log.Printf("[WARN] CSE Custom Match List Column not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("name", CSECustomMatchListColumn.Name)
	d.Set("fields", CSECustomMatchListColumn.Fields)

	return nil
}

func resourceSumologicCSECustomMatchListColumnDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	return c.DeleteCSECustomMatchListColumn(id)

}

func resourceSumologicCSECustomMatchListColumnCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		id, err := c.CreateCSECustomMatchListColumn(CSECustomMatchListColumn{
			Name:   d.Get("name").(string),
			Fields: fieldsToStringArray(d.Get("fields").([]interface{})),
		})

		if err != nil {
			return err
		}
		d.SetId(id)
	}

	return resourceSumologicCSECustomMatchListColumnUpdate(d, meta)
}

func fieldsToStringArray(resourceStrings []interface{}) []string {
	result := make([]string, len(resourceStrings))

	for i, resourceString := range resourceStrings {
		result[i] = resourceString.(string)
	}

	return result
}

func resourceSumologicCSECustomMatchListColumnUpdate(d *schema.ResourceData, meta interface{}) error {
	CSECustomMatchListColumn, err := resourceToCSECustomMatchListColumn(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSECustomMatchListColumn(CSECustomMatchListColumn); err != nil {
		return err
	}

	return resourceSumologicCSECustomMatchListColumnRead(d, meta)
}

func resourceToCSECustomMatchListColumn(d *schema.ResourceData) (CSECustomMatchListColumn, error) {
	id := d.Id()
	if id == "" {
		return CSECustomMatchListColumn{}, errors.New("Custom Match List Column id not specified")
	}

	return CSECustomMatchListColumn{
		ID:     id,
		Name:   d.Get("name").(string),
		Fields: fieldsToStringArray(d.Get("fields").([]interface{})),
	}, nil
}
