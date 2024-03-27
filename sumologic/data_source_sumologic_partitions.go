package sumologic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSumologicPartitions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumologicPartitionsRead,

		Schema: map[string]*schema.Schema{
			"partitions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: dataSourcePartitionSchema(),
					Read:   dataSourceSumologicPartitionRead,
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

	partitions := make([]map[string]interface{}, 0, len(spartitions))

	for _, spartition := range spartitions {
		partition := map[string]interface{}{
			"id":                                  spartition.ID,
			"name":                                spartition.Name,
			"routing_expression":                  spartition.RoutingExpression,
			"analytics_tier":                      spartition.AnalyticsTier,
			"retention_period":                    spartition.RetentionPeriod,
			"is_compliant":                        spartition.IsCompliant,
			"total_bytes":                         spartition.TotalBytes,
			"data_forwarding_id":                  spartition.DataForwardingID,
			"is_active":                           spartition.IsActive,
			"index_type":                          spartition.IndexType,
			"reduce_retention_period_immediately": spartition.ReduceRetentionPeriodImmediately,
		}

		partitions = append(partitions, partition)
	}

	d.Set("partitions", partitions)
	d.SetId(c.BaseURL.Host)

	return nil
}
