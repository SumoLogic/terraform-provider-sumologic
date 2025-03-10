package sumologic

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSumologicField() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicFieldCreate,
		Read:   resourceSumologicFieldRead,
		Update: resourceSumologicFieldUpdate,
		Delete: resourceSumologicFieldDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSumologicFieldImport,
		},

		Schema: map[string]*schema.Schema{

			"field_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"field_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: false,
			},

			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSumologicFieldRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Get("field_id").(string)
	name := d.Get("field_name").(string)
	if id == "" {
		newId, err := c.FindFieldId(name)
		if err != nil {
			return err
		}
		id = newId
	}

	field, err := c.GetField(id)
	if err != nil {
		return err
	}

	if field == nil {
		fmt.Printf("[WARN] Field not found, removing from state: %v - %v\n", id, err)
		d.SetId("")
		return nil
	}

	d.Set("field_name", field.FieldName)
	d.Set("field_id", field.FieldId)
	d.Set("state", field.State)

	return nil
}

func resourceSumologicFieldDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Get("field_id").(string)
	name := d.Get("field_name").(string)
	if id == "" {
		newId, err := c.FindFieldId(name)
		if err != nil {
			return err
		}
		id = newId
	}

	return c.DeleteField(id)
}

func resourceSumologicFieldCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		field := resourceToField(d)
		id, err := c.CreateField(field)
		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicFieldRead(d, meta)
}

func resourceSumologicFieldImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if d.Get("field_id").(string) == "" {
		d.Set("field_id", d.Id())
	}

	return []*schema.ResourceData{d}, nil
}

func resourceSumologicFieldUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Get("field_id").(string)
	name := d.Get("field_name").(string)
	state := d.Get("state").(string)
	if id == "" {
		newId, err := c.FindFieldId(name)
		if err != nil {
			return err
		}
		id = newId
	}
	_, err := c.GetField(id)

	if err != nil {
		return err
	}

	if state == "Enabled" {
		err := c.EnableField(id)
		if err != nil {
			return err
		}
	} else if state == "Disabled" {
		err := c.DisableField(id)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Invalid value of state field. Only Enabled or Disabled values are accepted")
	}

	return resourceSumologicFieldRead(d, meta)

}

func resourceToField(d *schema.ResourceData) Field {
	return Field{
		FieldId:   d.Get("field_id").(string),
		State:     d.Get("state").(string),
		FieldName: d.Get("field_name").(string),
	}
}
