package sumologic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSumologicRumSource() *schema.Resource {
	rumSource := resourceSumologicSource()
	rumSource.Create = resourceSumologicRumSourceCreate
	rumSource.Read = resourceSumologicRumSourceRead
	rumSource.Update = resourceSumologicRumSourceUpdate
	rumSource.Importer = &schema.ResourceImporter{
		State: resourceSumologicSourceImport,
	}

	rumSource.Schema["content_type"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"Rum"}, false),
	}

	rumSource.Schema["path"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: false,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"application_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"service_name": {
					Type:     schema.TypeString,
					Optional: false,
				},
				"deployment_environment": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"sampling_rate": {
					Type:     schema.TypeFloat,
					Optional: true,
				},
				"ignore_urls": {
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Optional: true,
				},
				"custom_tags": {
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Optional: true,
				},
				"propagate_trace_headers_cors_urls": {
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Optional: true,
				},
				"selected_country": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}

	return rumSource
}

func resourceSumologicRumSourceCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		source, err := resourceToRumSource(d)
		if err != nil {
			return err
		}

		id, err := c.CreateRumSource(source, d.Get("collector_id").(int))

		if err != nil {
			return err
		}

		d.SetId(strconv.Itoa(id))
	}

	return resourceSumologicRumSourceRead(d, meta)
}

func resourceSumologicRumSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	source, err := resourceToRumSource(d)
	if err != nil {
		return err
	}

	err = c.UpdateRumSource(source, d.Get("collector_id").(int))

	if err != nil {
		return err
	}

	return resourceSumologicRumSourceRead(d, meta)
}

func resourceSumologicRumSourceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	source, err := c.GetRumSource(d.Get("collector_id").(int), id)

	if err != nil {
		return err
	}

	if source == nil {
		log.Printf("[WARN] Rum source not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	if err := resourceSumologicSourceRead(d, source.Source); err != nil {
		return fmt.Errorf("%s", err)
	}
	d.Set("content_type", source.ContentType)

	return nil
}

func resourceToRumSource(d *schema.ResourceData) (RumSource, error) {
	source := resourceToSource(d)

	rumSource := RumSource{
		Source: source,
	}

	return rumSource, nil
}
