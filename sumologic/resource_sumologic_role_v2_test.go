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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicRoleV2_basic(t *testing.T) {
	var roleV2 RoleV2
	testName := acctest.RandomWithPrefix("tf-acc-test")
	testAuditDataFilter := "info"
	testSelectionType := "All"
	testCapabilities := []string{"\"manageContent\""}
	testDescription := "Manage data of the org."
	testSecurityDataFilter := "error"
	testLogAnalyticsFilter := "!_sourceCategory=collector"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRoleV2Destroy(roleV2),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSumologicRoleV2ConfigImported(testName, testAuditDataFilter, testSelectionType, testCapabilities, testDescription, testSecurityDataFilter, testLogAnalyticsFilter),
			},
			{
				ResourceName:      "sumologic_role_v2.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccSumologicRoleV2_create(t *testing.T) {
	var roleV2 RoleV2
	testName := acctest.RandomWithPrefix("tf-acc-test")
	testAuditDataFilter := "info"
	testSelectionType := "Allow"
	testCapabilities := []string{"\"manageContent\""}
	testDescription := "Manage data of the org."
	testSecurityDataFilter := "error"
	testLogAnalyticsFilter := "!_sourceCategory=collector"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRoleV2Destroy(roleV2),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicRoleV2(testName, testAuditDataFilter, testSelectionType, testCapabilities, testDescription, testSecurityDataFilter, testLogAnalyticsFilter),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoleV2Exists("sumologic_role_v2.test", &roleV2, t),
					testAccCheckRoleV2Attributes("sumologic_role_v2.test"),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "audit_data_filter", testAuditDataFilter),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "selection_type", testSelectionType),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "capabilities.0", strings.Replace(testCapabilities[0], "\"", "", 2)),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "description", testDescription),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "security_data_filter", testSecurityDataFilter),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "log_analytics_filter", testLogAnalyticsFilter),
				),
			},
		},
	})
}

func TestAccSumologicRoleV2_update(t *testing.T) {
	var roleV2 RoleV2
	testName := acctest.RandomWithPrefix("tf-acc-test")
	testAuditDataFilter := "info"
	testSelectionType := "Allow"
	testCapabilities := []string{"\"manageContent\""}
	testDescription := "Manage data of the org."
	testSecurityDataFilter := "error"
	testLogAnalyticsFilter := "!_sourceCategory=collector"

	testUpdatedName := acctest.RandomWithPrefix("tf-acc-test")
	testUpdatedAuditDataFilter := "infoUpdate"
	testUpdatedSelectionType := "All"
	testUpdatedCapabilities := []string{"\"manageContent\""}
	testUpdatedDescription := "Manage data of the org.Update"
	testUpdatedSecurityDataFilter := "errorUpdate"
	testUpdatedLogAnalyticsFilter := "!_sourceCategory=collectorUpdate"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRoleV2Destroy(roleV2),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicRoleV2(testName, testAuditDataFilter, testSelectionType, testCapabilities, testDescription, testSecurityDataFilter, testLogAnalyticsFilter),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoleV2Exists("sumologic_role_v2.test", &roleV2, t),
					testAccCheckRoleV2Attributes("sumologic_role_v2.test"),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "audit_data_filter", testAuditDataFilter),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "selection_type", testSelectionType),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "capabilities.0", strings.Replace(testCapabilities[0], "\"", "", 2)),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "description", testDescription),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "security_data_filter", testSecurityDataFilter),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "log_analytics_filter", testLogAnalyticsFilter),
				),
			},
			{
				Config: testAccSumologicRoleV2Update(testUpdatedName, testUpdatedAuditDataFilter, testUpdatedSelectionType, testUpdatedCapabilities, testUpdatedDescription, testUpdatedSecurityDataFilter, testUpdatedLogAnalyticsFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "name", testUpdatedName),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "audit_data_filter", testUpdatedAuditDataFilter),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "selection_type", testUpdatedSelectionType),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "capabilities.0", strings.Replace(testUpdatedCapabilities[0], "\"", "", 2)),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "description", testUpdatedDescription),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "security_data_filter", testUpdatedSecurityDataFilter),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "log_analytics_filter", testUpdatedLogAnalyticsFilter),
				),
			},
		},
	})
}

func testAccCheckRoleV2Destroy(roleV2 RoleV2) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			u, err := client.GetRoleV2(id)
			if err != nil {
				return fmt.Errorf("Encountered an error: " + err.Error())
			}
			if u != nil {
				return fmt.Errorf("RoleV2 %s still exists", id)
			}
		}
		return nil
	}
}

func testAccCheckRoleV2Exists(name string, roleV2 *RoleV2, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			//need this so that we don't get an unused import error for strconv in some cases
			return fmt.Errorf("Error = %s. RoleV2 not found: %s", strconv.FormatBool(ok), name)
		}

		//need this so that we don't get an unused import error for strings in some cases
		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("RoleV2 ID is not set")
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)
		newRoleV2, err := client.GetRoleV2(id)
		if err != nil {
			return fmt.Errorf("RoleV2 %s not found", id)
		}
		roleV2 = newRoleV2
		return nil
	}
}
func testAccCheckSumologicRoleV2ConfigImported(name string, auditDataFilter string, selectionType string, capabilities []string, description string, securityDataFilter string, logAnalyticsFilter string) string {
	return fmt.Sprintf(`
resource "sumologic_role_v2" "foo" {
      name = "%s"
      audit_data_filter = "%s"
      selection_type = "%s"
      capabilities = %v
      description = "%s"
      security_data_filter = "%s"
      log_analytics_filter = "%s"
}
`, name, auditDataFilter, selectionType, capabilities, description, securityDataFilter, logAnalyticsFilter)
}

func testAccSumologicRoleV2(name string, auditDataFilter string, selectionType string, capabilities []string, description string, securityDataFilter string, logAnalyticsFilter string) string {
	return fmt.Sprintf(`
resource "sumologic_role_v2" "test" {
    selected_views {
		view_name = "sumologic_default"
	}
    name = "%s"
    audit_data_filter = "%s"
    selection_type = "%s"
    capabilities = %v
    description = "%s"
    security_data_filter = "%s"
    log_analytics_filter = "%s"
}
`, name, auditDataFilter, selectionType, capabilities, description, securityDataFilter, logAnalyticsFilter)
}

func testAccSumologicRoleV2Update(name string, auditDataFilter string, selectionType string, capabilities []string, description string, securityDataFilter string, logAnalyticsFilter string) string {
	return fmt.Sprintf(`
resource "sumologic_role_v2" "test" {
      name = "%s"
      audit_data_filter = "%s"
      selection_type = "%s"
      capabilities = %v
      description = "%s"
      security_data_filter = "%s"
      log_analytics_filter = "%s"
}
`, name, auditDataFilter, selectionType, capabilities, description, securityDataFilter, logAnalyticsFilter)
}

func testAccCheckRoleV2Attributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(name, "name"),
			resource.TestCheckResourceAttrSet(name, "audit_data_filter"),
			resource.TestCheckResourceAttrSet(name, "selection_type"),
			resource.TestCheckResourceAttrSet(name, "description"),
			resource.TestCheckResourceAttrSet(name, "security_data_filter"),
			resource.TestCheckResourceAttrSet(name, "log_analytics_filter"),
		)
		return f(s)
	}
}
