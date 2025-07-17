package sumologic

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSumologicInstalledCollector() *schema.Resource {
	return &schema.Resource{
		Read:   resourceSumologicCollectorRead,
		Delete: resourceSumologicCollectorDelete,
		Update: resourceSumologicInstalledCollectorUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cutoff_timestamp": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Etc/UTC",
			},
			"ephemeral": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"fields": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"alive": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"last_seen_alive": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source_sync_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "UI",
			},
			"target_cpu": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"collector_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSumologicInstalledCollectorUpdate(d *schema.ResourceData, meta interface{}) error {

	collector := resourceToInstalledCollector(d)

	c := meta.(*Client)
	err := c.UpdateCollector(collector)

	if err != nil {
		return err
	}

	return resourceSumologicCollectorRead(d, meta)
}

func resourceToInstalledCollector(d *schema.ResourceData) Collector {
	id, _ := strconv.Atoi(d.Id())

	return Collector{
		ID:               int64(id),
		CollectorType:    "Installable",
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Category:         d.Get("category").(string),
		TimeZone:         d.Get("timezone").(string),
		HostName:         d.Get("host_name").(string),
		Ephemeral:        d.Get("ephemeral").(bool),
		SourceSyncMode:   d.Get("source_sync_mode").(string),
		Targetcpu:        d.Get("target_cpu").(int),
		Fields:           d.Get("fields").(map[string]interface{}),
		CutoffTimestamp:  d.Get("cutoff_timestamp").(int),
		Alive:            d.Get("alive").(bool),
		LastSeenAlive:    d.Get("last_seen_alive").(int),
		CollectorVersion: d.Get("collector_version").(string),
	}
}
