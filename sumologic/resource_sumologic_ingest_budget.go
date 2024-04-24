package sumologic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSumologicIngestBudget() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Use resource_sumologic_ingest_budget_v2 instead.",
		Create:             resourceSumologicIngestBudgetCreate,
		Read:               resourceSumologicIngestBudgetRead,
		Delete:             resourceSumologicIngestBudgetDelete,
		Update:             resourceSumologicIngestBudgetUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceSumologicIngestBudgetImport,
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

func resourceSumologicIngestBudgetCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, err := c.CreateIngestBudget(resourceToIngestBudget(d))
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceSumologicIngestBudgetRead(d, meta)
}

func resourceSumologicIngestBudgetRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()

	budget, err := c.GetIngestBudget(id)
	if err != nil {
		return err
	}

	if budget == nil {
		log.Printf("[WARN] Ingest budget not found, removing from state: %v", id)
		d.SetId("")

		return nil
	}

	d.Set("name", budget.Name)
	d.Set("field_value", budget.FieldValue)
	d.Set("capacity_bytes", budget.Capacity)
	d.Set("timezone", budget.Timezone)
	d.Set("reset_time", budget.ResetTime)
	d.Set("description", budget.Description)
	d.Set("action", budget.Action)

	return nil
}

func resourceSumologicIngestBudgetUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	err := c.UpdateIngestBudget(resourceToIngestBudget(d))

	return err
}

func resourceSumologicIngestBudgetDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	err := c.DeleteIngestBudget(d.Id())

	return err
}

func resourceSumologicIngestBudgetImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	c := meta.(*Client)

	name := d.Id()

	budget, err := c.FindIngestBudget(name)

	if err != nil {
		return nil, err
	}

	d.SetId(budget.ID)

	return []*schema.ResourceData{d}, nil
}

func resourceToIngestBudget(d *schema.ResourceData) IngestBudget {
	return IngestBudget{
		ID:          d.Id(),
		Name:        d.Get("name").(string),
		FieldValue:  d.Get("field_value").(string),
		Capacity:    d.Get("capacity_bytes").(int),
		Timezone:    d.Get("timezone").(string),
		ResetTime:   d.Get("reset_time").(string),
		Description: d.Get("description").(string),
		Action:      d.Get("action").(string),
	}
}
