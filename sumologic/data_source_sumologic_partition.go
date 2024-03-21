package sumologic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSumologicPartition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumologicPartitionRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
			},
			"routing_expression": {
				Type: schema.TypeString,
			},
			"analytics_tier": {
				Type: schema.TypeString,
			},
			"retention_period": {
				Type: schema.TypeInt,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == "-1" && old != ""
				},
			},
			"is_compliant": {
				Type: schema.TypeBool,
			},
			"total_bytes": {
				Type: schema.TypeInt,
			},
			"data_forwarding_id": {
				Type: schema.TypeString,
			},
			"is_active": {
				Type: schema.TypeBool,
			},
			"index_type": {
				Type: schema.TypeString,
			},
			"reduce_retention_period_immediately": {
				Type: schema.TypeBool,
			},
		},
	}
}

func dataSourceSumologicPartitionRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	spartition, err := c.GetPartition(id)

	if err != nil {
		return fmt.Errorf("error retrieving partition: %v", err)
	}

	if spartition == nil {
		d.SetId("")
		return nil
	}

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
