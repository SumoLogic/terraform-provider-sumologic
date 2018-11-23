package sumologic

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSumologicHTTPSource() *schema.Resource {
	httpSource := resourceSumologicSource()
	httpSource.Create = resourceSumologicHTTPSourceCreate
	httpSource.Read = resourceSumologicHTTPSourceRead
	httpSource.Update = resourceSumologicHTTPSourceUpdate

	httpSource.Schema["message_per_request"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		ForceNew: false,
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

	if d.Get("lookup_by_name").(bool) {
		source, err := c.GetSourceName(d.Get("collector_id").(int), d.Get("name").(string))

		if err != nil {
			return err
		}

		if source != nil {
			d.SetId(strconv.Itoa(source.ID))
		}
	}

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
		log.Printf("[WARN] HTTP source not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	resourceSumologicSourceRead(d, source.Source)
	d.Set("message_per_request", source.MessagePerRequest)
	d.Set("url", source.URL)

	return nil
}
