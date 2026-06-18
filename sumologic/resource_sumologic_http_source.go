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
		log.Printf("[WARN] HTTP sources not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	if err := resourceSumologicSourceRead(d, source.Source); err != nil {
		return fmt.Errorf("%s", err)
	}
	d.Set("message_per_request", source.MessagePerRequest)
	d.Set("url", source.URL)

	return nil
}
