package sumologic

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSumologicScheduledView() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicScheduledViewCreate,
		Read:   resourceSumologicScheduledViewRead,
		Delete: resourceSumologicScheduledViewDelete,
		Update: resourceSumologicScheduledViewUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"query": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 16384),
			},
			"index_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"start_time": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},
			"index_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"data_forwarding_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parsing_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Manual",
				ValidateFunc: validation.StringInSlice([]string{"AutoParse", "Manual"}, false),
			},
			"reduce_retention_period_immediately": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"time_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Time zone for ingesting data in scheduled view. Follow the format in the [IANA Time Zone Database](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List).",
			},
			"auto_pause_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Default:     true,
				Description: "Auto Pause status of Scheduled View.",
			},
		},
	}
}

func resourceSumologicScheduledViewCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	if d.Id() == "" {
		sview := resourceToScheduledView(d)

		// add timeZone if specified
		if v, ok := d.GetOk("time_zone"); ok {
			sview.TimeZone = v.(string)
		}

		if v, ok := d.GetOk("auto_pause_enabled"); ok {
			sview.AutoPauseEnabled = v.(bool)
		}

		createdSview, err := c.CreateScheduledView(sview)

		if err != nil {
			return err
		}

		d.SetId(createdSview.ID)
		d.Set("index_id", createdSview.IndexId)
		d.Set("retention_period", createdSview.RetentionPeriod)
	}

	return resourceSumologicScheduledViewUpdate(d, meta)
}

func resourceSumologicScheduledViewRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	sview, err := c.GetScheduledView(id)

	if err != nil {
		return err
	}

	if sview == nil {
		log.Printf("[WARN] Scheduled view not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	d.Set("index_id", sview.IndexId)
	d.Set("query", sview.Query)
	d.Set("index_name", sview.IndexName)
	d.Set("start_time", sview.StartTime.Format(time.RFC3339))
	d.Set("retention_period", sview.RetentionPeriod)
	d.Set("data_forwarding_id", sview.DataForwardingId)
	d.Set("parsing_mode", sview.ParsingMode)
	d.Set("time_zone", sview.TimeZone)
	d.Set("auto_pause_enabled", sview.AutoPauseEnabled)

	return nil
}
func resourceSumologicScheduledViewDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	return c.DeleteScheduledView(d.Id())
}
func resourceSumologicScheduledViewUpdate(d *schema.ResourceData, meta interface{}) error {
	sview := resourceToScheduledView(d)

	if d.HasChange("time_zone") {
		sview.TimeZone = d.Get("time_zone").(string)
	}

	if d.HasChange("auto_pause_enabled") {
		sview.AutoPauseEnabled = d.Get("auto_pause_enabled").(bool)
	}

	c := meta.(*Client)
	err := c.UpdateScheduledView(sview)

	if err != nil {
		return err
	}

	return resourceSumologicScheduledViewRead(d, meta)
}

func resourceToScheduledView(d *schema.ResourceData) ScheduledView {
	var startTimeParsed, err = time.Parse(time.RFC3339, d.Get("start_time").(string))
	if err != nil {
		log.Fatal(err)
	}
	return ScheduledView{
		ID:                               d.Id(),
		Query:                            d.Get("query").(string),
		IndexName:                        d.Get("index_name").(string),
		StartTime:                        startTimeParsed,
		RetentionPeriod:                  d.Get("retention_period").(int),
		DataForwardingId:                 d.Get("data_forwarding_id").(string),
		ParsingMode:                      d.Get("parsing_mode").(string),
		ReduceRetentionPeriodImmediately: d.Get("reduce_retention_period_immediately").(bool),
	}
}
