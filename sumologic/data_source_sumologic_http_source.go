package sumologic

import (
	"fmt"
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
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"multiline": {
				Type:     schema.TypeBool,
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
	var collectorId int64
	switch cid := d.Get("collector_id").(type) {
	case int:
		collectorId = int64(cid)
	case int64:
		collectorId = cid
	default:
		return fmt.Errorf("unknown data type of collector_id: %T, value: %v", cid, cid)
	}
	source, err := c.GetSourceName(collectorId, d.Get("name").(string))

	if err != nil {
		return err
	}

	if source == nil {
		d.SetId("")
		return fmt.Errorf("HTTP source not found, removing from state: %v - %v", id, err)
	}

	d.SetId(strconv.Itoa(source.ID))
	d.Set("name", source.Name)
	d.Set("description", source.Description)
	d.Set("category", source.Category)
	d.Set("timezone", source.TimeZone)
	d.Set("multiline", source.MultilineProcessingEnabled)
	d.Set("url", source.Url)

	return nil
}
