package sumologic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSumologicScanBudget() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicScanBudgetCreate,
		Read:   resourceSumologicScanBudgetRead,
		Update: resourceSumologicScanBudgetUpdate,
		Delete: resourceSumologicScanBudgetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"orgId": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"capacity": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},

			"unit": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"budgetType": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"window": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"groupBy": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"applicableOn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"scope": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeList,
				},
			},

			"action": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"status": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourceSumologicScanBudgetCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		scanBudget := resourceToScanBudget(d)
		id, err := c.CreateScanBudget(scanBudget)
		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicScanBudgetRead(d, meta)
}

func resourceSumologicScanBudgetRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	ScanBudget, err := c.GetScanBudget(id)
	if err != nil {
		return err
	}

	if ScanBudget == nil {
		log.Printf("[WARN] ScanBudget not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("name", ScanBudget.Name)
	d.Set("orgId", ScanBudget.OrgID)
	d.Set("capacity", ScanBudget.Capacity)
	d.Set("unit", ScanBudget.Unit)
	d.Set("budgetType", ScanBudget.BudgetType)
	d.Set("window", ScanBudget.Window)
	d.Set("applicableOn", ScanBudget.Grouping)
	d.Set("groupBy", ScanBudget.GroupingEntity)
	d.Set("action", ScanBudget.Action)
	d.Set("scope", ScanBudget.Scope)
	d.Set("status", ScanBudget.Status)

	return nil
}

func resourceSumologicScanBudgetDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteScanBudget(d.Id())
}

func resourceSumologicScanBudgetUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	ScanBudget := resourceToScanBudget(d)
	err := c.UpdateScanBudget(ScanBudget)
	if err != nil {
		return err
	}

	return resourceSumologicScanBudgetRead(d, meta)
}

func resourceToScanBudget(d *schema.ResourceData) ScanBudget {
	return ScanBudget{
		ID:             d.Id(),
		OrgID:          d.Get("orgId").(string),
		Name:           d.Get("name").(string),
		Capacity:       d.Get("capacity").(int),
		Unit:           d.Get("unit").(string),
		BudgetType:     d.Get("budgetType").(string),
		Window:         d.Get("window").(string),
		Grouping:       d.Get("applicableOn").(string),
		GroupingEntity: d.Get("groupBy").(string),
		Action:         d.Get("action").(string),
		Scope:          d.Get("scope").(map[string]interface{}),
		Status:         d.Get("status").(string),
	}
}
