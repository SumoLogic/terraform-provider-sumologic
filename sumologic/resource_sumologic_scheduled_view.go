package sumologic

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"time"
)

func resourceSumologicScheduledView() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicScheduledViewCreate,
		Read:   resourceSumologicScheduledViewRead,
		Delete: resourceSumologicScheduledViewDelete,
		Update: resourceSumologicScheduledViewUpdate,
		Exists: resourceSumologicScheduledViewExists,

		Schema: map[string]*schema.Schema{
			"query": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
                ValidateFunc: validation.StringLenBetween(1, 16384),
			},
			"index_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"start_time": {
				Type:         schema.TypeString, // TODO type should be different
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.ValidateRFC3339TimeString,
			},
			"retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     false,
				Default:      -1,
				ValidateFunc: validation.IntAtLeast(-1),
			},
			"data_forwarding_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceSumologicScheduledViewCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	if d.Id() == "" {
		sview := resourceToScheduledView(d)
		id, err := c.CreateScheduledView(sview)

		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicScheduledViewUpdate(d, meta)
}
func resourceSumologicScheduledViewRead(d *schema.ResourceData, meta interface{}) error {
    return nil // TODO implement
}
func resourceSumologicScheduledViewDelete(d *schema.ResourceData, meta interface{}) error {
    return nil // TODO implement
}
func resourceSumologicScheduledViewUpdate(d *schema.ResourceData, meta interface{}) error {
    return nil // TODO implement
}
func resourceSumologicScheduledViewExists(d *schema.ResourceData, meta interface{}) (bool, error) {
    return false, nil // TODO implement
}

func resourceToScheduledView(d *schema.ResourceData) ScheduledView {
    var startTimeParsed, err = time.Parse(time.RFC3339, d.Get("start_time").(string))
    if err != nil {
        panic(err)
    }
	return ScheduledView{
		ID:        d.Id(),
		Query:     d.Get("query").(string),
		IndexName: d.Get("index_name").(string),
		StartTime: startTimeParsed,
		RetentionPeriod:  d.Get("retention_period").(int),
		DataForwardingId: d.Get("data_forwarding_id").(string),
	}
}
