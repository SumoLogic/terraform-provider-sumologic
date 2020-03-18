package sumologic

import (
	"fmt"
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
			"source_name": {
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
	source, err := c.GetSourceName(d.Get("collector_id").(int), d.Get("source_name").(string))

	if err != nil {
		return err
	}

	if source == nil {
		d.SetId("")
		return fmt.Errorf("HTTP source not found, removing from state: %v - %v", id, err)

	}

	d.SetId(strconv.Itoa(source.ID))
	d.Set("source_name", source.Name)
	d.Set("url", source.Url)

	return nil
}
