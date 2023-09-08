package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicSCEAutomation_create_update(t *testing.T) {
	SkipCseTest(t)

	var Automation CSEAutomation
	nPlaybookId := "63ece953d5f0cb2ec4d5794e"
	nCseResourceType := "INSIGHT"
	nExecutionTypes := []string{"NEW_INSIGHT"}
	nEnabled := true
	uExecutionTypes := []string{"ON_DEMAND"}
	uEnabled := false
	resourceName := "sumologic_cse_automation.automation"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEAutomationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSEAutomationConfig(nPlaybookId, nCseResourceType, nExecutionTypes, nEnabled),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEAutomationExists(resourceName, &Automation),
					testCheckAutomationValues(&Automation, nPlaybookId, nCseResourceType, nExecutionTypes, nEnabled),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testCreateCSEAutomationConfig(nPlaybookId, nCseResourceType, uExecutionTypes, uEnabled),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEAutomationExists(resourceName, &Automation),
					testCheckAutomationValues(&Automation, nPlaybookId, nCseResourceType, uExecutionTypes, uEnabled),
				),
			},
		},
	})
}

func testAccCSEAutomationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_automation" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Automation destruction check: CSE Automation ID is not set")
		}

		s, err := client.GetCSEAutomation(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("automation still exists")
		}
	}
	return nil
}

func testCreateCSEAutomationConfig(nPlaybookId string, nCseResourceType string, nExecutionTypes []string, nEnabled bool) string {

	return fmt.Sprintf(`
resource "sumologic_cse_automation" "automation" {
	playbook_id = "%s"
	cse_resource_type = "%s"
	execution_types = ["%s"]
	enabled = "%t"	
}
`, nPlaybookId, nCseResourceType, nExecutionTypes[0], nEnabled)
}

func testCheckCSEAutomationExists(n string, Automation *CSEAutomation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("automation ID is not set")
		}

		c := testAccProvider.Meta().(*Client)
		AutomationResp, err := c.GetCSEAutomation(rs.Primary.ID)
		if err != nil {
			return err
		}

		*Automation = *AutomationResp

		return nil
	}
}

func testCheckAutomationValues(Automation *CSEAutomation, nPlaybookId string, nCseResourceType string, nExecutionTypes []string, nEnabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if Automation.PlaybookId != nPlaybookId {
			return fmt.Errorf("bad playbook_id, expected \"%s\", got: %#v", nPlaybookId, Automation.PlaybookId)
		}
		if Automation.CseResourceType != nCseResourceType {
			return fmt.Errorf("bad cse_resource_type, expected \"%s\", got: %#v", nCseResourceType, Automation.CseResourceType)
		}
		if Automation.ExecutionTypes != nil {
			if len(Automation.ExecutionTypes) != len(nExecutionTypes) {
				return fmt.Errorf("bad execution_types list, expected \"%d\", got: %d", len(nExecutionTypes), len(Automation.ExecutionTypes))
			}
			if Automation.ExecutionTypes[0] != nExecutionTypes[0] {
				return fmt.Errorf("bad execution_types in list, expected \"%s\", got: %s", nExecutionTypes[0], Automation.ExecutionTypes[0])
			}
		}
		if Automation.Enabled != nEnabled {
			return fmt.Errorf("bad enabled field, expected \"%t\", got: %#v", nEnabled, Automation.Enabled)
		}

		return nil
	}
}
