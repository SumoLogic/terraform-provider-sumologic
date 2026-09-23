package sumologic

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"base_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"json_unrolling": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"json_unroll_settings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path":       {Type: schema.TypeString, Computed: true},
						"field_name": {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}

}

func dataSourceSumologicHTTPSourceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var collectorId int
	switch cid := d.Get("collector_id").(type) {
	case int:
		collectorId = cid
	case int64:
		collectorId = int(cid)
	default:
		return fmt.Errorf("unknown data type of collector_id: %T, value: %v", cid, cid)
	}

	var source *HTTPSource
	var err error

	// If we have both collector_id and name, look up by name first
	if d.Get("name").(string) != "" {
		baseSource, err := c.GetSourceName(int64(collectorId), d.Get("name").(string))
		if err != nil {
			return err
		}
		if baseSource == nil {
			d.SetId("")
			return fmt.Errorf("HTTP source not found")
		}
		// Now get the full HTTP source with all fields
		source, err = c.GetHTTPSource(collectorId, baseSource.ID)
		if err != nil {
			return err
		}
	} else if d.Id() != "" {
		// If we have an ID, get directly
		id, _ := strconv.Atoi(d.Id())
		source, err = c.GetHTTPSource(collectorId, id)
		if err != nil {
			return err
		}
	}

	if source == nil {
		d.SetId("")
		return fmt.Errorf("HTTP source not found")
	}

	d.SetId(strconv.Itoa(source.ID))
	d.Set("name", source.Name)
	d.Set("description", source.Description)
	d.Set("category", source.Category)
	d.Set("timezone", source.TimeZone)
	d.Set("multiline", source.MultilineProcessingEnabled)
	d.Set("url", source.Url)
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
