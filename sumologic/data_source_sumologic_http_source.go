package sumologic

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSumologicHTTPSource() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceSumologicHTTPSourceRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"collector_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceSumologicHTTPSourceRead(d *schema.ResourceData, meta interface{}) error {
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

	//d.Set("message_per_request", source.MessagePerRequest)
	//resourceSumologicSourceRead(d, source.Source)
	d.SetId(strconv.Itoa(source.ID))
	d.Set("name", "some name")
	d.Set("url", "some url")

	return nil
}
