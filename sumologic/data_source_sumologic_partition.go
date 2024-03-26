package sumologic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSumologicPartition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumologicPartitionRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"routing_expression": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"analytics_tier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"retention_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_compliant": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"total_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_forwarding_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"index_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reduce_retention_period_immediately": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func dataSourceSumologicPartitionRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var err error
	var spartition *Partition

	if rid, ok := d.GetOk("id"); ok {
		id := rid.(string)
		spartition, err = c.GetPartition(id)
		if err != nil {
			return fmt.Errorf("partition with id %v not found: %v", id, err)
		}
		if spartition == nil {
			return fmt.Errorf("partition with id %v not found", id)
		}
	}

	d.SetId(spartition.ID)
	d.Set("routing_expression", spartition.RoutingExpression)
	d.Set("name", spartition.Name)
	d.Set("analytics_tier", spartition.AnalyticsTier)
	d.Set("retention_period", spartition.RetentionPeriod)
	d.Set("is_compliant", spartition.IsCompliant)
	d.Set("data_forwarding_id", spartition.DataForwardingId)
	d.Set("is_active", spartition.IsActive)
	d.Set("total_bytes", spartition.TotalBytes)
	d.Set("index_type", spartition.IndexType)

	return nil
}
