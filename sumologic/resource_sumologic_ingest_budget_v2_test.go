package sumologic

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicIngestBudgetV2_basic(t *testing.T) {
	var ingestBudgetV2 IngestBudgetV2
	testName := "Developer Budget"
	testScope := "_sourceCategory=*prod*nginx*"
	testTimezone := "America/Los_Angeles"
	testResetTime := "23:30"
	testAuditThreshold := 85
	testDescription := "description-7hUwr"
	testAction := "stopCollecting"
	testCapacityBytes := 1000
	testBudgetType := "dailyVolume"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIngestBudgetV2Destroy(ingestBudgetV2),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSumologicIngestBudgetV2ConfigImported(testName, testScope, testTimezone, testResetTime, testAuditThreshold, testDescription, testAction, testCapacityBytes, testBudgetType),
			},
			{
				ResourceName:      "sumologic_ingest_budget_v2.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "sumologic_ingest_budget_v2.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     testName,
			},
		},
	})
}
func TestAccSumologicIngestBudgetV2_create(t *testing.T) {
	var ingestBudgetV2 IngestBudgetV2
	testName := "Developer Budget"
	testScope := "_sourceCategory=*prod*nginx*"
	testTimezone := "America/Los_Angeles"
	testResetTime := "23:30"
	testAuditThreshold := 85
	testDescription := "description-900AB"
	testAction := "stopCollecting"
	testCapacityBytes := 1000
	testBudgetType := "dailyVolume"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIngestBudgetV2Destroy(ingestBudgetV2),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicIngestBudgetV2(testName, testScope, testTimezone, testResetTime, testAuditThreshold, testDescription, testAction, testCapacityBytes, testBudgetType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIngestBudgetV2Exists("sumologic_ingest_budget_v2.test", &ingestBudgetV2, t),
					testAccCheckIngestBudgetV2Attributes("sumologic_ingest_budget_v2.test"),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "scope", testScope),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "timezone", testTimezone),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "reset_time", testResetTime),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "audit_threshold", strconv.Itoa(testAuditThreshold)),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "description", testDescription),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "action", testAction),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "capacity_bytes", strconv.Itoa(testCapacityBytes)),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "budget_type", testBudgetType),
				),
			},
		},
	})
}

func TestAccSumologicIngestBudgetV2_update(t *testing.T) {
	var ingestBudgetV2 IngestBudgetV2
	testName := "Developer Budget"
	testScope := "_sourceCategory=*prod*nginx*"
	testTimezone := "America/Los_Angeles"
	testResetTime := "23:30"
	testAuditThreshold := 85
	testDescription := "description-2tAk8"
	testAction := "stopCollecting"
	testCapacityBytes := 1000
	testBudgetType := "dailyVolume"

	testUpdatedName := "Developer BudgetUpdate"
	testUpdatedScope := "_sourceCategory=*prod*nginx*Update"
	testUpdatedTimezone := "America/Lima"
	testUpdatedResetTime := "22:05"
	testUpdatedAuditThreshold := 86
	testUpdatedDescription := "description-pY8kDUpdate"
	testUpdatedAction := "keepCollecting"
	testUpdatedCapacityBytes := 1001
	testUpdatedBudgetType := "dailyVolume"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIngestBudgetV2Destroy(ingestBudgetV2),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicIngestBudgetV2(testName, testScope, testTimezone, testResetTime, testAuditThreshold, testDescription, testAction, testCapacityBytes, testBudgetType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIngestBudgetV2Exists("sumologic_ingest_budget_v2.test", &ingestBudgetV2, t),
					testAccCheckIngestBudgetV2Attributes("sumologic_ingest_budget_v2.test"),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "scope", testScope),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "timezone", testTimezone),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "reset_time", testResetTime),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "audit_threshold", strconv.Itoa(testAuditThreshold)),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "description", testDescription),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "action", testAction),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "capacity_bytes", strconv.Itoa(testCapacityBytes)),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "budget_type", testBudgetType),
				),
			},
			{
				Config: testAccSumologicIngestBudgetV2Update(testUpdatedName, testUpdatedScope, testUpdatedTimezone, testUpdatedResetTime, testUpdatedAuditThreshold, testUpdatedDescription, testUpdatedAction, testUpdatedCapacityBytes, testUpdatedBudgetType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "name", testUpdatedName),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "scope", testUpdatedScope),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "timezone", testUpdatedTimezone),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "reset_time", testUpdatedResetTime),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "audit_threshold", strconv.Itoa(testUpdatedAuditThreshold)),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "description", testUpdatedDescription),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "action", testUpdatedAction),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "capacity_bytes", strconv.Itoa(testUpdatedCapacityBytes)),
					resource.TestCheckResourceAttr("sumologic_ingest_budget_v2.test", "budget_type", testUpdatedBudgetType),
				),
			},
		},
	})
}

func testAccCheckIngestBudgetV2Destroy(ingestBudgetV2 IngestBudgetV2) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			u, err := client.GetIngestBudgetV2(id)
			if err != nil {
				return fmt.Errorf("Encountered an error: " + err.Error())
			}
			if u != nil {
				return fmt.Errorf("IngestBudgetV2 %s still exists", id)
			}
		}
		return nil
	}
}

func testAccCheckIngestBudgetV2Exists(name string, ingestBudgetV2 *IngestBudgetV2, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			//need this so that we don't get an unused import error for strconv in some cases
			return fmt.Errorf("Error = %s. IngestBudgetV2 not found: %s", strconv.FormatBool(ok), name)
		}

		//need this so that we don't get an unused import error for strings in some cases
		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("IngestBudgetV2 ID is not set")
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)
		newIngestBudgetV2, err := client.GetIngestBudgetV2(id)
		if err != nil {
			return fmt.Errorf("IngestBudgetV2 %s not found", id)
		}
		ingestBudgetV2 = newIngestBudgetV2
		return nil
	}
}
func testAccCheckSumologicIngestBudgetV2ConfigImported(name string, scope string, timezone string, resetTime string, auditThreshold int, description string, action string, capacityBytes int, budgetType string) string {
	return fmt.Sprintf(`
resource "sumologic_ingest_budget_v2" "foo" {
      name = "%s"
      scope = "%s"
      timezone = "%s"
      reset_time = "%s"
      audit_threshold = %d
      description = "%s"
      action = "%s"
      capacity_bytes = %d
	  budget_type = "%s"
}
`, name, scope, timezone, resetTime, auditThreshold, description, action, capacityBytes, budgetType)
}

func testAccSumologicIngestBudgetV2(name string, scope string, timezone string, resetTime string, auditThreshold int, description string, action string, capacityBytes int, budgetType string) string {
	return fmt.Sprintf(`
resource "sumologic_ingest_budget_v2" "test" {
    name = "%s"
    scope = "%s"
    timezone = "%s"
    reset_time = "%s"
    audit_threshold = %d
    description = "%s"
    action = "%s"
    capacity_bytes = %d
	budget_type = "%s"
}
`, name, scope, timezone, resetTime, auditThreshold, description, action, capacityBytes, budgetType)
}

func testAccSumologicIngestBudgetV2Update(name string, scope string, timezone string, resetTime string, auditThreshold int, description string, action string, capacityBytes int, budgetType string) string {
	return fmt.Sprintf(`
resource "sumologic_ingest_budget_v2" "test" {
      name = "%s"
      scope = "%s"
      timezone = "%s"
      reset_time = "%s"
      audit_threshold = %d
      description = "%s"
      action = "%s"
      capacity_bytes = %d
	  budget_type = "%s"
}
`, name, scope, timezone, resetTime, auditThreshold, description, action, capacityBytes, budgetType)
}

func testAccCheckIngestBudgetV2Attributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(name, "name"),
			resource.TestCheckResourceAttrSet(name, "scope"),
			resource.TestCheckResourceAttrSet(name, "timezone"),
			resource.TestCheckResourceAttrSet(name, "reset_time"),
			resource.TestCheckResourceAttrSet(name, "audit_threshold"),
			resource.TestCheckResourceAttrSet(name, "description"),
			resource.TestCheckResourceAttrSet(name, "action"),
			resource.TestCheckResourceAttrSet(name, "budget_type"),
			resource.TestCheckResourceAttrSet(name, "capacity_bytes"),
		)
		return f(s)
	}
}
