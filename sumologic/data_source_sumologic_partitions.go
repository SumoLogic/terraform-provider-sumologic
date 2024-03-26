package sumologic

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceSumologicPartitions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumologicPartitionsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"name", "routingExpression", "id"}, false),
						},
						"operator": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "equals",
							ValidateFunc: validation.StringInSlice([]string{"Equals", "Contains", "HasPrefix", "HasSuffix"}, false),
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSumologicPartitionsRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	spartitions, err := c.ListPartitions()
	if err != nil {
		return fmt.Errorf("error retrieving partitions: %v", err)
	}

	ids := []string{}

	filters, ok := d.GetOk("filters")
	if !ok {
		return fmt.Errorf("no filters provided")
	}

	for _, partition := range spartitions {
		if matchesAllFilters(partition, filters.([]interface{})) {
			ids = append(ids, partition.ID)
		}
	}

	d.SetId(c.BaseURL.Host)
	d.Set("ids", ids)

	return nil
}

func matchesAllFilters(partition Partition, filters []interface{}) bool {
	for _, filter := range filters {
		f := filter.(map[string]interface{})
		k := f["key"].(string)
		o := f["operator"].(string)
		v := f["value"].(string)

		var fv string
		switch k {
		case "name":
			fv = partition.Name
		case "routingExpression":
			fv = partition.RoutingExpression
		case "id":
			fv = partition.ID
		default:
			return false
		}

		switch o {
		case "Equals":
			if fv != v {
				return false
			}
		case "Contains":
			if !strings.Contains(fv, v) {
				return false
			}
		case "HasPrefix":
			if !strings.HasPrefix(fv, v) {
				return false
			}
		case "HasSuffix":
			if !strings.HasSuffix(fv, v) {
				return false
			}
		default:
			return false
		}
	}
	return true
}
