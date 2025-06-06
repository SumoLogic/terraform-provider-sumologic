// ----------------------------------------------------------------------------
//
//	***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//	This file is automatically generated by Sumo Logic and manual
//	changes will be clobbered when the file is regenerated. Do not submit
//	changes to this file.
//
// ----------------------------------------------------------------------------
package sumologic

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicScanBudget_basic(t *testing.T) {
	var scanBudget ScanBudget

	testName := "Test Budget"
	testCapacity := 10
	testUnit := "GB"
	testBudgetType := "ScanBudget"
	testWindow := "Query"
	testApplicableOn := "PerEntity"
	testGroupBy := "User"
	testAction := "StopForeGroundScan"
	testScope := ScanBudgetScope{
		IncludedUsers: []string{"000000000000011C"},
		ExcludedUsers: []string{},
		IncludedRoles: []string{},
		ExcludedRoles: []string{"0000000000000196"},
	}
	testStatus := "active"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScanBudgetDestroy(scanBudget),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSumologicScanBudgetConfigImported(testName, testCapacity, testUnit, testBudgetType, testWindow, testApplicableOn, testGroupBy, testAction, testScope, testStatus),
			},
			{
				ResourceName:      "sumologic_scan_budget.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccSumologicScanBudget_create(t *testing.T) {
	var scanBudget ScanBudget

	testName := "Test Budget"
	testCapacity := 10
	testUnit := "GB"
	testBudgetType := "ScanBudget"
	testWindow := "Query"
	testApplicableOn := "PerEntity"
	testGroupBy := "User"
	testAction := "StopForeGroundScan"
	testScope := ScanBudgetScope{
		IncludedUsers: []string{"000000000000011C"},
		ExcludedUsers: []string{},
		IncludedRoles: []string{},
		ExcludedRoles: []string{"0000000000000196"},
	}
	testStatus := "active"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScanBudgetDestroy(scanBudget),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicScanBudget(testName, testCapacity, testUnit, testBudgetType, testWindow, testApplicableOn, testGroupBy, testAction, testScope, testStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScanBudgetExists("sumologic_scan_budget.test", &scanBudget, t),
					testAccCheckScanBudgetAttributes("sumologic_scan_budget.test"),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "capacity", strconv.Itoa(testCapacity)),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "unit", testUnit),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "budget_type", testBudgetType),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "window", testWindow),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "applicable_on", testApplicableOn),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "group_by", testGroupBy),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "action", testAction),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "status", testStatus),
				),
			},
		},
	})
}

func TestAccSumologicScanBudget_update(t *testing.T) {
	var scanBudget ScanBudget

	testName := "Test Budget"
	testCapacity := 10
	testUnit := "GB"
	testBudgetType := "ScanBudget"
	testWindow := "Query"
	testApplicableOn := "PerEntity"
	testGroupBy := "User"
	testAction := "StopForeGroundScan"
	testScope := ScanBudgetScope{
		IncludedUsers: []string{"000000000000011C"},
		ExcludedUsers: []string{},
		IncludedRoles: []string{},
		ExcludedRoles: []string{"0000000000000196"},
	}
	testStatus := "active"

	testUpdatedName := "Test Budget"
	testUpdatedCapacity := 20
	testUpdatedUnit := "GB"
	testUpdatedBudgetType := "ScanBudget"
	testUpdatedWindow := "Daily"
	testUpdatedApplicableOn := "PerEntity"
	testUpdatedGroupBy := "User"
	testUpdatedAction := "Warn"
	testUpdatedScope := ScanBudgetScope{
		IncludedUsers: []string{"000000000000011C"},
		ExcludedUsers: []string{},
		IncludedRoles: []string{},
		ExcludedRoles: []string{"0000000000000196"},
	}
	testUpdatedStatus := "active"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScanBudgetDestroy(scanBudget),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicScanBudget(testName, testCapacity, testUnit, testBudgetType, testWindow, testApplicableOn, testGroupBy, testAction, testScope, testStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScanBudgetExists("sumologic_scan_budget.test", &scanBudget, t),
					testAccCheckScanBudgetAttributes("sumologic_scan_budget.test"),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "capacity", strconv.Itoa(testCapacity)),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "unit", testUnit),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "budget_type", testBudgetType),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "window", testWindow),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "applicable_on", testApplicableOn),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "group_by", testGroupBy),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "action", testAction),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "status", testStatus),
				),
			},
			{
				Config: testAccSumologicScanBudgetUpdate(testUpdatedName, testUpdatedCapacity, testUpdatedUnit, testUpdatedBudgetType, testUpdatedWindow, testUpdatedApplicableOn, testUpdatedGroupBy, testUpdatedAction, testUpdatedScope, testUpdatedStatus),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "name", testUpdatedName),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "capacity", strconv.Itoa(testUpdatedCapacity)),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "unit", testUpdatedUnit),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "budget_type", testUpdatedBudgetType),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "window", testUpdatedWindow),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "applicable_on", testUpdatedApplicableOn),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "group_by", testUpdatedGroupBy),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "action", testUpdatedAction),
					resource.TestCheckResourceAttr("sumologic_scan_budget.test", "status", testUpdatedStatus),
				),
			},
		},
	})
}

func testAccCheckScanBudgetDestroy(scanBudget ScanBudget) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			_, err := client.GetScanBudget(id)
			if err == nil {
				return fmt.Errorf("ScanBudget %s still exists", id)
			}
		}
		return nil
	}
}

func testAccCheckScanBudgetExists(name string, scanBudget *ScanBudget, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			//need this so that we don't get an unused import error for strconv in some cases
			return fmt.Errorf("Error = %s. ScanBudget not found: %s", strconv.FormatBool(ok), name)
		}

		//need this so that we don't get an unused import error for strings in some cases
		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("ScanBudget ID is not set")
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)
		newScanBudget, err := client.GetScanBudget(id)
		if err != nil {
			return fmt.Errorf("ScanBudget %s not found", id)
		}
		scanBudget = newScanBudget
		return nil
	}
}
func testAccCheckSumologicScanBudgetConfigImported(name string, capacity int, unit string, budgetType string, window string, applicableOn string, groupBy string, action string, scope ScanBudgetScope, status string) string {
	return fmt.Sprintf(`
resource "sumologic_scan_budget" "foo" {
      name = "%s"
      capacity = %d
      unit = "%s"
      budget_type = "%s"
      window = "%s"
      applicable_on = "%s"
	  group_by = "%s"
      action = "%s"
      scope {
	  	included_users = ["%s"]
	  	excluded_users = []
	  	included_roles = []
	  	excluded_roles = ["%s"]
	  }
	  status = "%s"
}
`, name, capacity, unit, budgetType, window, applicableOn, groupBy, action, scope.IncludedUsers[0], scope.ExcludedRoles[0], status)
}

func testAccSumologicScanBudget(name string, capacity int, unit string, budgetType string, window string, applicableOn string, groupBy string, action string, scope ScanBudgetScope, status string) string {
	return fmt.Sprintf(`
resource "sumologic_scan_budget" "test" {
    name = "%s"
    capacity = %d
    unit = "%s"
    budget_type = "%s"
    window = "%s"
    applicable_on = "%s"
	group_by = "%s"
    action = "%s"
    scope {
		included_users = ["%s"]
	  	excluded_users = []
	  	included_roles = []
	  	excluded_roles = ["%s"]
	}
	status = "%s"
}
`, name, capacity, unit, budgetType, window, applicableOn, groupBy, action, scope.IncludedUsers[0], scope.ExcludedRoles[0], status)
}

func testAccSumologicScanBudgetUpdate(name string, capacity int, unit string, budgetType string, window string, applicableOn string, groupBy string, action string, scope ScanBudgetScope, status string) string {
	return fmt.Sprintf(`
resource "sumologic_scan_budget" "test" {
      name = "%s"
      capacity = %d
      unit = "%s"
      budget_type = "%s"
      window = "%s"
      applicable_on = "%s"
      group_by = "%s"
      action = "%s"
      scope {
	  	included_users = ["%s"]
	  	excluded_users = []
	  	included_roles = []
	  	excluded_roles = ["%s"]
	  }
      status = "%s"
}
`, name, capacity, unit, budgetType, window, applicableOn, groupBy, action, scope.IncludedUsers[0], scope.ExcludedRoles[0], status)
}

func testAccCheckScanBudgetAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(name, "name"),
			resource.TestCheckResourceAttrSet(name, "capacity"),
			resource.TestCheckResourceAttrSet(name, "unit"),
			resource.TestCheckResourceAttrSet(name, "budget_type"),
			resource.TestCheckResourceAttrSet(name, "window"),
			resource.TestCheckResourceAttrSet(name, "applicable_on"),
			resource.TestCheckResourceAttrSet(name, "group_by"),
			resource.TestCheckResourceAttrSet(name, "action"),
			resource.TestCheckResourceAttrSet(name, "status"),
		)
		return f(s)
	}
}
