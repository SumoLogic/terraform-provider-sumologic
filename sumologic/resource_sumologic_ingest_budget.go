package sumologic

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSumologicIngestBudget() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Use resource_sumologic_ingest_budget_v2 instead.",
		Create: func(d *schema.ResourceData, meta any) error {
			return errors.New("Use resource_sumologic_ingest_budget_v2 instead.")
		},
		Read: func(d *schema.ResourceData, meta any) error {
			return errors.New("Use resource_sumologic_ingest_budget_v2 instead.")
		},
		Update: func(d *schema.ResourceData, meta any) error {
			return errors.New("Use resource_sumologic_ingest_budget_v2 instead.")
		},
		Delete: func(d *schema.ResourceData, meta any) error {
			return errors.New("Use resource_sumologic_ingest_budget_v2 instead.")
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"field_value": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringLenBetween(1, 1024),
			},
			"capacity_bytes": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "Etc/UTC",
			},
			"reset_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				Default:      "00:00",
				ValidateFunc: validation.StringLenBetween(5, 5),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "",
			},
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringInSlice([]string{"stopCollecting", "keepCollecting"}, false),
				Default:      "keepCollecting",
			},
		},
	}
}
