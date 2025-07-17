package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceSumologicCSECustomEntityType() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSECustomEntityTypeCreate,
		Read:   resourceSumologicCSECustomEntityTypeRead,
		Delete: resourceSumologicCSECustomEntityTypeDelete,
		Update: resourceSumologicCSECustomEntityTypeUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"identifier": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fields": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSumologicCSECustomEntityTypeRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSECustomEntityType *CSECustomEntityType
	id := d.Id()

	CSECustomEntityType, err := c.GetCSECustomEntityType(id)
	if err != nil {
		log.Printf("[WARN] CSE Custom Entity Type not found when looking by id: %s, err: %v", id, err)

	}

	if CSECustomEntityType == nil {
		log.Printf("[WARN] CSE Custom Entity Type not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("name", CSECustomEntityType.Name)
	d.Set("identifier", CSECustomEntityType.Identifier)
	d.Set("fields", CSECustomEntityType.Fields)

	return nil
}

func resourceSumologicCSECustomEntityTypeDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteCSECustomEntityType(d.Id())

}

func resourceSumologicCSECustomEntityTypeCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		id, err := c.CreateCSECustomEntityType(CSECustomEntityType{
			Name:       d.Get("name").(string),
			Identifier: d.Get("identifier").(string),
			Fields:     resourceFieldsToStringArray(d.Get("fields").([]interface{})),
		})

		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicCSECustomEntityTypeRead(d, meta)
}

func resourceSumologicCSECustomEntityTypeUpdate(d *schema.ResourceData, meta interface{}) error {
	CSECustomEntityType, err := resourceToCSECustomEntityType(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSECustomEntityType(CSECustomEntityType); err != nil {
		return err
	}

	return resourceSumologicCSECustomEntityTypeRead(d, meta)
}

func resourceToCSECustomEntityType(d *schema.ResourceData) (CSECustomEntityType, error) {
	id := d.Id()
	if id == "" {
		return CSECustomEntityType{}, nil
	}

	return CSECustomEntityType{
		ID:         id,
		Name:       d.Get("name").(string),
		Identifier: d.Get("identifier").(string),
		Fields:     resourceFieldsToStringArray(d.Get("fields").([]interface{})),
	}, nil
}

func resourceFieldsToStringArray(resourceFields []interface{}) []string {
	fields := make([]string, len(resourceFields))

	for i, field := range resourceFields {
		fields[i] = field.(string)
	}

	return fields
}
