package sumologic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSumologicHTTPSource() *schema.Resource {
	httpSource := resourceSumologicSource()
	httpSource.Create = resourceSumologicHTTPSourceCreate
	httpSource.Read = resourceSumologicHTTPSourceRead
	httpSource.Update = resourceSumologicHTTPSourceUpdate
	httpSource.Importer = &schema.ResourceImporter{
		State: resourceSumologicSourceImport,
	}

	httpSource.Schema["message_per_request"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
	httpSource.Schema["url"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	httpSource.Schema["token"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	httpSource.Schema["base_url"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	httpSource.Schema["json_unrolling"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
	httpSource.Schema["json_unroll_settings"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"path": {
					Type:     schema.TypeString,
					Required: true,
				},
				"field_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}

	return httpSource
}

func resourceSumologicHTTPSourceCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		source := resourceToHTTPSource(d)

		id, err := c.CreateHTTPSource(source, d.Get("collector_id").(int))

		if err != nil {
			return err
		}

		d.SetId(strconv.Itoa(id))
	}

	return resourceSumologicHTTPSourceRead(d, meta)
}

func resourceSumologicHTTPSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	source := resourceToHTTPSource(d)

	err := c.UpdateHTTPSource(source, d.Get("collector_id").(int))

	if err != nil {
		return err
	}

	return resourceSumologicHTTPSourceRead(d, meta)
}

func resourceToHTTPSource(d *schema.ResourceData) HTTPSource {
	source := resourceToSource(d)
	source.Type = "HTTP"

	httpSource := HTTPSource{
		Source:            source,
		MessagePerRequest: d.Get("message_per_request").(bool),
		JsonUnrolling:     d.Get("json_unrolling").(bool),
	}

	if settings := d.Get("json_unroll_settings").([]interface{}); len(settings) > 0 {
		s := settings[0].(map[string]interface{})
		js := &JsonUnrollSettings{
			Path: s["path"].(string),
		}
		if v := s["field_name"].(string); v != "" {
			js.FieldName = &v
		}
		httpSource.JsonUnrollSettings = js
	}

	return httpSource
}

func resourceSumologicHTTPSourceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	source, err := c.GetHTTPSource(d.Get("collector_id").(int), id)

	if err != nil {
		return err
	}

	if source == nil {
		log.Printf("[WARN] HTTP source not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	if err := resourceSumologicSourceRead(d, source.Source); err != nil {
		return fmt.Errorf("%s", err)
	}
	d.Set("message_per_request", source.MessagePerRequest)
	d.Set("url", source.URL)
	d.Set("token", source.Token)
	d.Set("base_url", source.BaseUrl)
	d.Set("json_unrolling", source.JsonUnrolling)
	settings := []map[string]interface{}{}
	if source.JsonUnrollSettings != nil {
		fieldName := ""
		if source.JsonUnrollSettings.FieldName != nil {
			fieldName = *source.JsonUnrollSettings.FieldName
		}
		settings = []map[string]interface{}{
			{
				"path":       source.JsonUnrollSettings.Path,
				"field_name": fieldName,
			},
		}
	}
	d.Set("json_unroll_settings", settings)

	return nil
}
