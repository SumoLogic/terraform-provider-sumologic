package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicIngestBudget_create(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))
	fieldValue := fmt.Sprintf("tf-%s", acctest.RandString(5))
	description := fmt.Sprintf("tf-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIngestBudgetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicIngestBudgetConfig(name, fieldValue, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "name", name),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "field_value", fieldValue),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "capacity_bytes", "30000000000"),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "description", description),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "timezone", "Etc/UTC"),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "reset_time", "00:00"),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "action", "keepCollecting"),
				),
			},
		}})
}

func TestAccSumologicIngestBudget_assign(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))
	collectorName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	fieldValue := fmt.Sprintf("tf-%s", acctest.RandString(5))
	description := fmt.Sprintf("tf-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIngestBudgetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicIngestBudgetAssignmentConfig(collectorName, name, fieldValue, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_ingest_budget.testBudget", "name", name),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.testBudget", "field_value", fieldValue),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.testBudget", "capacity_bytes", "2"),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.testBudget", "timezone", "Etc/UTC"),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.testBudget", "reset_time", "00:00"),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.testBudget", "action", "keepCollecting"),
					resource.TestCheckResourceAttr("sumologic_collector.testCollector", "name", collectorName),
					resource.TestCheckResourceAttr("sumologic_collector.testCollector", "fields._budget", fieldValue),
				),
			},
		}})
}

func TestAccSumologicIngestBudget_basic(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))
	fieldValue := fmt.Sprintf("tf-%s", acctest.RandString(5))
	description := fmt.Sprintf("tf-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIngestBudgetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicIngestBudgetConfig(name, fieldValue, description),
			},
			{
				ResourceName:      "sumologic_ingest_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     name,
			},
		}})
}

func testAccCheckIngestBudgetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "sumologic_ingest_budget" {
			continue
		}
		id := r.Primary.ID
		u, err := client.GetIngestBudget(id)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if u != nil {
			return fmt.Errorf("Ingest Budget still exists")
		}
	}
	return nil

}

func testAccSumologicIngestBudgetConfig(name, fieldValue, description string) string {
	return fmt.Sprintf(`
resource "sumologic_ingest_budget" "budget" {
	name = "%s"
  	field_value = "%s"
  	capacity_bytes = 30000000000
  	description = "%s"
}
`, name, fieldValue, description)
}

func testAccSumologicIngestBudgetAssignmentConfig(collectorName, name, fieldValue, description string) string {
	return fmt.Sprintf(`
resource "sumologic_collector" "testCollector" {
	name = "%s"
	fields = {
	  "_budget" = "${sumologic_ingest_budget.testBudget.field_value}"
	}
}

resource "sumologic_ingest_budget" "testBudget" {
	name = "%s"
  	field_value = "%s"
  	capacity_bytes = 2
  	description = "%s"
}
`, collectorName, name, fieldValue, description)
}
