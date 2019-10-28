package sumologic

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
)

func resourceSumologicPartition() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicPartitionCreate,
		Read:   resourceSumologicPartitionRead,
		Delete: resourceSumologicPartitionDelete,
		Update: resourceSumologicPartitionUpdate,
		Exists: resourceSumologicPartitionExists,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"routing_expression": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringLenBetween(0, 16384),
			},
			"analytics_tier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validation.IntAtLeast(-1),
				Default:      -1,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// taken from https://stackoverflow.com/a/57785476/118587
					return old == "-1"
				},
			},
			"is_compliant": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: false,
			},
			"data_forwarding_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
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
		d.Set("retention_period", createdSpartition.RetentionPeriod)
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
	d.Set("is_compliant", spartition.RetentionPeriod)
	d.Set("data_forwarding_id", spartition.DataForwardingId)

	return nil
}
func resourceSumologicPartitionDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	return c.DeletePartition(d.Id())
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
func resourceSumologicPartitionExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	c := meta.(*Client)

	spartition, err := c.GetPartition(d.Id())
	if err != nil {
		return false, err
	}

	return spartition != nil, nil
}

func resourceToPartition(d *schema.ResourceData) Partition {
	return Partition{
		ID:                d.Id(),
		Name:              d.Get("name").(string),
		RoutingExpression: d.Get("routing_expression").(string),
		AnalyticsTier:     d.Get("analytics_tier").(string),
		RetentionPeriod:   d.Get("retention_period").(int),
		IsCompliant:       d.Get("is_compliant").(bool),
		DataForwardingId:  d.Get("data_forwarding_id").(string),
	}
}
