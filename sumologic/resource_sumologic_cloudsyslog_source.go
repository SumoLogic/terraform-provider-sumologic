package sumologic

import (
	"strconv"

	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSumologicCloudsyslogSource() *schema.Resource {
	cloudsyslogSource := resourceSumologicSource()
	cloudsyslogSource.Create = resourceSumologicCloudsyslogSourceCreate
	cloudsyslogSource.Read = resourceSumologicCloudsyslogSourceRead
	cloudsyslogSource.Update = resourceSumologicCloudsyslogSourceUpdate

	return cloudsyslogSource
}

func resourceSumologicCloudsyslogSourceCreate(d *schema.ResourceData, meta interface{}) error {
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
		source := resourceToCloudsyslogSource(d)

		id, err := c.CreateCloudsyslogSource(source, d.Get("collector_id").(int))

		if err != nil {
			return err
		}

		d.SetId(strconv.Itoa(id))
	}

	return resourceSumologicCloudsyslogSourceRead(d, meta)
}

func resourceSumologicCloudsyslogSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	source := resourceToCloudsyslogSource(d)

	err := c.UpdateCloudsyslogSource(source, d.Get("collector_id").(int))

	if err != nil {
		return err
	}

	return resourceSumologicCloudsyslogSourceRead(d, meta)
}

func resourceToCloudsyslogSource(d *schema.ResourceData) CloudsyslogSource {
	source := resourceToSource(d)
	source.Type = "Cloudsyslog"

	cloudsyslogSource := CloudsyslogSource{
		Source:            source,
	}

	return cloudsyslogSource
}

func resourceSumologicCloudsyslogSourceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	source, err := c.GetCloudsyslogSource(d.Get("collector_id").(int), id)

	if err != nil {
		log.Printf("[WARN] Cloudsyslog source not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}
	//TODO: Create a common READ and then do the unique reads specific to each resource.
	d.Set("name", source.Name)
	d.Set("description", source.Description)
	d.Set("category", source.Category)

	return nil
}
