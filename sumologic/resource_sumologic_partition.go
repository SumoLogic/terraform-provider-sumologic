package sumologic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSumologicPartition() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicPartitionCreate,
		Read:   resourceSumologicPartitionRead,
		Delete: resourceSumologicPartitionDelete,
		Update: resourceSumologicPartitionUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"routing_expression": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 16384),
			},
			"analytics_tier": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(-1),
				Default:      -1,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == "-1" && old != ""
				},
			},
			"is_compliantt": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
			"is_included_in_default_search": {
				Type:        schema.TypeBool,
				Description: "Indicates whether the partition is included in the default search scope. Configuring this property is exclusively permitted for flex partitions.",
				Optional:    true,
				Default:     true,
			},
		},
	}
}

func resourceSumologicPartitionCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	if d.Id() == "" {
		spartition := resourceToPartition(d)
		createdSpartition, err := c.CreatePartition(spartition)

		if err != nil {
			return err
		}

		d.SetId(createdSpartition.ID)
	}

	return resourceSumologicPartitionUpdate(d, meta)
}

func resourceSumologicPartitionRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	spartition, err := c.GetPartition(id)

	if err != nil {
		return err
	}

	if spartition == nil {
		log.Printf("[WARN] Partition not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	d.Set("routing_expression", spartition.RoutingExpression)
	d.Set("name", spartition.Name)
	d.Set("analytics_tier", spartition.AnalyticsTier)
	d.Set("retention_period", spartition.RetentionPeriod)
	d.Set("is_compliantt", spartition.IsCompliant)
	d.Set("data_forwarding_id", spartition.DataForwardingId)
	d.Set("is_active", spartition.IsActive)
	d.Set("total_bytes", spartition.TotalBytes)
	d.Set("index_type", spartition.IndexType)
	d.Set("is_included_in_default_search", spartition.IsIncludedInDefaultSearch)

	return nil
}
func resourceSumologicPartitionDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	return c.DecommissionPartition(d.Id())
}
func resourceSumologicPartitionUpdate(d *schema.ResourceData, meta interface{}) error {
	spartition := resourceToPartition(d)

	c := meta.(*Client)
	err := c.UpdatePartition(spartition)

	if err != nil {
		return err
	}

	return resourceSumologicPartitionRead(d, meta)
}

func resourceToPartition(d *schema.ResourceData) Partition {
	return Partition{
		ID:                               d.Id(),
		Name:                             d.Get("name").(string),
		RoutingExpression:                d.Get("routing_expression").(string),
		AnalyticsTier:                    d.Get("analytics_tier").(string),
		RetentionPeriod:                  d.Get("retention_period").(int),
		IsCompliant:                      d.Get("is_compliantt").(bool),
		DataForwardingId:                 d.Get("data_forwarding_id").(string),
		IsActive:                         d.Get("is_active").(bool),
		TotalBytes:                       d.Get("total_bytes").(int),
		IndexType:                        d.Get("index_type").(string),
		ReduceRetentionPeriodImmediately: d.Get("reduce_retention_period_immediately").(bool),
		IsIncludedInDefaultSearch:        d.Get("is_included_in_default_search").(bool),
	}
}
