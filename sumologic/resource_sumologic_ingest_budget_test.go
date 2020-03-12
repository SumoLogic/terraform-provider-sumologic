package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicIngestBudget_create(t *testing.T) {
	var ingestBudget IngestBudget
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
					testAccCheckIngestBudgetExists("sumologic_ingest_budget.budget", &ingestBudget),
					testAccCheckIngestBudgetValues(&ingestBudget, name, fieldValue, description, 30000000000),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "name", name),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "field_value", fieldValue),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "capacity_bytes", "30000000000"),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "description", description),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "timezone", "Etc/UTC"),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "reset_time", "00:00"),
					resource.TestCheckResourceAttr("sumologic_ingest_budget.budget", "action", "keepCollecting"),
				),
			},
			{
				ResourceName:      "sumologic_ingest_budget.budget",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     name,
			},
		}})
}

func TestAccSumologicIngestBudget_assign(t *testing.T) {
	var ingestBudget IngestBudget
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
					testAccCheckIngestBudgetExists("sumologic_ingest_budget.testBudget", &ingestBudget),
					testAccCheckIngestBudgetValues(&ingestBudget, name, fieldValue, description, 2),
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

func testAccCheckIngestBudgetExists(n string, ingestBudget *IngestBudget) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		c := testAccProvider.Meta().(*Client)
		ingestBudgetResp, err := c.GetIngestBudget(rs.Primary.ID)
		if err != nil {
			return err
		}

		// If no error, assign the response ingest budget attribute to the ingest budget pointer
		*ingestBudget = *ingestBudgetResp

		return nil
	}
}

func testAccCheckIngestBudgetValues(ingestBudget *IngestBudget, name, fieldValue, description string, capacity int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if ingestBudget.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, ingestBudget.Name)
		}
		if ingestBudget.Capacity != capacity {
			return fmt.Errorf("bad name, expected \"%d\", got: %#v", capacity, ingestBudget.Capacity)
		}
		if ingestBudget.FieldValue != fieldValue {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", fieldValue, ingestBudget.FieldValue)
		}
		if ingestBudget.Description != description {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", description, ingestBudget.Description)
		}
		return nil
	}
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
