package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

func resourceSumologicCSETagSchema() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSETagSchemaCreate,
		Read:   resourceSumologicCSETagSchemaRead,
		Delete: resourceSumologicCSETagSchemaDelete,
		Update: resourceSumologicCSETagSchemaUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.All(
						validation.StringIsNotEmpty,
						validation.StringInSlice([]string{"customInsight", "entity", "rule", "threatIntelligence"}, false)),
				},
			},
			"free_form": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"value_options": getValueOptionsSchema(),
		},
	}
}

func getValueOptionsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     schema.TypeString,
					Required: true,
				},
				"label": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"link": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func resourceSumologicCSETagSchemaRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSETagSchema *CSETagSchema
	key := d.Id()

	CSETagSchema, err := c.GetCSETagSchema(key)
	if err != nil {
		log.Printf("[WARN] CSE Tag Schema not found when looking by key: %s, err: %v", key, err)

	}

	if CSETagSchema == nil {
		log.Printf("[WARN] CSE tag Schema not found, removing from state: %v - %v", key, err)
		d.SetId("")
		return nil
	}

	d.Set("key", CSETagSchema.Key)
	d.Set("label", CSETagSchema.Label)
	d.Set("content_types", CSETagSchema.ContentTypes)
	d.Set("value_options", valueOptionsArrayToResource(CSETagSchema.ValueOptionObjects))
	d.Set("free_form", CSETagSchema.FreeForm)

	return nil
}

func valueOptionsArrayToResource(valueOptions []ValueOption) []map[string]interface{} {
	result := make([]map[string]interface{}, len(valueOptions))

	for i, valueOption := range valueOptions {
		result[i] = map[string]interface{}{
			"value": valueOption.Value,
			"label": valueOption.Label,
			"link":  valueOption.Link,
		}
	}

	return result
}

func resourceSumologicCSETagSchemaDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteCSETagSchema(d.Id())

}

func resourceSumologicCSETagSchemaCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	key, err := c.CreateCSETagSchema(CSECreateUpdateTagSchema{
		Key:          d.Get("key").(string),
		Label:        d.Get("label").(string),
		ContentTypes: resourceFieldsToStringArray(d.Get("content_types").([]interface{})),
		ValueOptions: resourceToValueOptionArray(d.Get("value_options").([]interface{})),
		FreeForm:     d.Get("free_form").(bool),
	})

	if err != nil {
		return err
	}

	d.SetId(key)

	return resourceSumologicCSETagSchemaRead(d, meta)
}

func resourceToValueOptionArray(resouceValueOptions []interface{}) []ValueOption {
	result := make([]ValueOption, len(resouceValueOptions))

	for i, resourceValueOption := range resouceValueOptions {
		result[i] = ValueOption{
			Value: resourceValueOption.(map[string]interface{})["value"].(string),
			Label: resourceValueOption.(map[string]interface{})["label"].(string),
			Link:  resourceValueOption.(map[string]interface{})["link"].(string),
		}
	}

	return result
}

func resourceSumologicCSETagSchemaUpdate(d *schema.ResourceData, meta interface{}) error {
	CSETagSchema, err := resourceToCSETagSchema(d)
	if err != nil {
		return err
	}
	c := meta.(*Client)
	if err = c.UpdateCSETagSchema(CSETagSchema); err != nil {
		return err
	}

	return resourceSumologicCSETagSchemaRead(d, meta)
}

func resourceToCSETagSchema(d *schema.ResourceData) (CSECreateUpdateTagSchema, error) {
	key := d.Id()
	if key == "" {
		return CSECreateUpdateTagSchema{}, nil
	}

	return CSECreateUpdateTagSchema{
		Key:          key,
		Label:        d.Get("label").(string),
		ContentTypes: resourceFieldsToStringArray(d.Get("content_types").([]interface{})),
		ValueOptions: resourceToValueOptionArray(d.Get("value_options").([]interface{})),
		FreeForm:     d.Get("free_form").(bool),
	}, nil
}
