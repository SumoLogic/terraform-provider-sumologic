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
	httpSource.Delete = resourceSumologicHTTPSourceDelete

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
		source := resourceToSource(d)
		source.Type = "HTTP"
		httpSource := HTTPSource{
			Source:            source,
			MessagePerRequest: d.Get("message_per_request").(bool),
		}

		id, err := c.CreateHTTPSource(httpSource, d.Get("collector_id").(int))

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

	id, _ := strconv.Atoi(d.Id())

	source := HTTPSource{}
	source.ID = id
	source.Type = "HTTP"
	source.Name = d.Get("name").(string)
	source.Description = d.Get("description").(string)
	source.Category = d.Get("category").(string)
	source.MessagePerRequest = d.Get("message_per_request").(bool)

	return source
}

func resourceSumologicHTTPSourceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	source, err := c.GetHTTPSource(d.Get("collector_id").(int), id)

	// Source is gone, remove it from state
	if err != nil {
		log.Printf("[WARN] HTTP source not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	d.Set("name", source.Name)
	d.Set("description", source.Description)
	d.Set("category", source.Category)
	d.Set("message_per_request", source.MessagePerRequest)
	d.Set("url", source.URL)

	return nil
}

func resourceSumologicHTTPSourceDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	// Destroy source if `destroy` is true, otherwise ignore
	if d.Get("destroy").(bool) {
		id, _ := strconv.Atoi(d.Id())
		collectorID, _ := d.Get("collector_id").(int)

		return c.DestroySource(id, collectorID)
	}

	return nil
}
