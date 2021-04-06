package sumologic

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSumologicCollector() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumologicCollectorRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fields": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSumologicCollectorRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var collector *Collector
	var err error
	if cid, ok := d.GetOk("id"); ok {
		id := cid.(int)
		collector, err = c.GetCollector(id)
		if err != nil {
			return fmt.Errorf("collector with id %d not found: %v", id, err)
		}
	} else {
		if cname, ok := d.GetOk("name"); ok {
			name := cname.(string)
			collector, err = c.GetCollectorName(name)
			if err != nil {
				return fmt.Errorf("collector with name %s not found: %v", name, err)
			}
			if collector == nil {
				return fmt.Errorf("collector with name %s not found", name)
			}
		} else {
			return errors.New("please specify either id or name")
		}
	}

	d.SetId(strconv.FormatInt(collector.ID, 10))
	d.Set("name", collector.Name)
	d.Set("description", collector.Description)
	d.Set("category", collector.Category)
	d.Set("timezone", collector.TimeZone)
	if err := d.Set("fields", collector.Fields); err != nil {
		return fmt.Errorf("error setting fields for datasource %s: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] data_source_sumologic_collector: retrieved %v", collector)
	return nil
}
