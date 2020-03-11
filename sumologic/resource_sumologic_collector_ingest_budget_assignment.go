package sumologic

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"strconv"
)

func resourceSumologicCollectorIngestBudgetAssignment() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "use fields attribute of collector resource instead to assign an ingest budget",
		Create:             resourceSumologicCollectorIngestBudgetAssignmentCreate,
		Read:               resourceSumologicCollectorIngestBudgetAssignmentRead,
		Delete:             resourceSumologicCollectorIngestBudgetAssignmentDelete,

		Schema: map[string]*schema.Schema{
			"collector_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ingest_budget_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSumologicCollectorIngestBudgetAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	collectorId, err := strconv.Atoi(d.Get("collector_id").(string))
	if err != nil {
		return err
	}

	ingestBudgetId := d.Get("ingest_budget_id").(string)

	err = c.AssignCollectorToIngestBudget(ingestBudgetId, collectorId)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d-%s", collectorId, ingestBudgetId))

	return nil
}

func resourceSumologicCollectorIngestBudgetAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	collectorId, err := strconv.Atoi(d.Get("collector_id").(string))
	if err != nil {
		return err
	}

	ingestBudgetId := d.Get("ingest_budget_id").(string)

	assigned, err := c.CollectorAssignedToIngestBudget(ingestBudgetId, collectorId)
	if err != nil {
		return err
	}

	if !assigned {
		log.Printf("[WARN] Collector %d no longer assigned to ingest budget %s", collectorId, ingestBudgetId)
		d.SetId("")

		return nil
	}

	return nil
}

func resourceSumologicCollectorIngestBudgetAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	collectorId, err := strconv.Atoi(d.Get("collector_id").(string))
	if err != nil {
		return err
	}

	ingestBudgetId := d.Get("ingest_budget_id").(string)

	err = c.UnAssignCollectorToIngestBudget(ingestBudgetId, collectorId)

	return err
}
