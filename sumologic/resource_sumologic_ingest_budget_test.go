package sumologic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSumologicIngestBudget_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicIngestBudgetConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "name", "test"),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "field_value", "test"),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "capacity_bytes", "30000000000"),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "description", "For testing purposes"),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "timezone", "Etc/UTC"),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "reset_time", "00:00"),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "action", "keepCollecting"),
				),
			},
		}})
}

func TestAccSumologicIngestBudget_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicIngestBudgetConfig,
			},
			{
				ResourceName:      "sumologic_ingest_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "test",
			},
		}})
}

var testAccSumologicIngestBudgetConfig = `
resource "sumologic_ingest_budget" "budget" {
  name = "test"
  field_value = "test"
  capacity_bytes = 30000000000
  description = "For testing purposes"
}
`
