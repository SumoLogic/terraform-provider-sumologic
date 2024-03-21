package sumologic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSumologicPartitions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumologicPartitionsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

	ids := make([]string, 0, len(spartitions))

	for _, partition := range spartitions {
		ids = append(ids, partition.ID)
	}

	d.SetId(c.BaseURL.Host)
	d.Set("ids", ids)

	return nil
}
