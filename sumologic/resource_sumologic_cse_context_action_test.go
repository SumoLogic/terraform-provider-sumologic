package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicSCEContextAction_create_update(t *testing.T) {
	SkipCseTest(t)

	var ContextAction CSEContextAction
	nName := "Test Context Action"
	nType := "URL"
	nTemplate := "https://bar.com/?q={{value}}"
	nIocTypes := []string{"IP_ADDRESS"}
	nEntityTypes := []string{"_hostname"}
	nRecordFields := []string{"request_url"}
	nAllRecordFields := false
	nEnabled := true
	uIocTypes := []string{"MAC_ADDRESS"}
	uEnabled := false
	resourceName := "sumologic_cse_context_action.context_action"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEContextActionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSEContextActionConfig(nName, nType, nTemplate, nIocTypes, nEntityTypes, nRecordFields, nAllRecordFields, nEnabled),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEContextActionExists(resourceName, &ContextAction),
					testCheckContextActionValues(&ContextAction, nName, nType, nTemplate, nIocTypes, nEntityTypes, nRecordFields, nAllRecordFields, nEnabled),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testCreateCSEContextActionConfig(nName, nType, nTemplate, uIocTypes, nEntityTypes, nRecordFields, nAllRecordFields, uEnabled),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEContextActionExists(resourceName, &ContextAction),
					testCheckContextActionValues(&ContextAction, nName, nType, nTemplate, uIocTypes, nEntityTypes, nRecordFields, nAllRecordFields, uEnabled),
				),
			},
		},
	})
}

func testAccCSEContextActionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_context_action" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Context Action destruction check: CSE Context Action ID is not set")
		}

		s, err := client.GetCSEContextAction(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("Context Action still exists")
		}
	}
	return nil
}

func testCreateCSEContextActionConfig(nName string, nType string, nTemplate string, nIocTypes []string, nEntityTypes []string, nRecordFields []string, nAllRecordFields bool, nEnabled bool) string {

	return fmt.Sprintf(`
resource "sumologic_cse_context_action" "context_action" {
	name = "%s"
	type = "%s"
	template = "%s"
	ioc_types = ["%s"]
	entity_types = ["%s"]
	record_fields = ["%s"]
	all_record_fields = "%t"	
	enabled = "%t"	
}
`, nName, nType, nTemplate, nIocTypes[0], nEntityTypes[0], nRecordFields[0], nAllRecordFields, nEnabled)
}

func testCheckCSEContextActionExists(n string, ContextAction *CSEContextAction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Context Action ID is not set")
		}

		c := testAccProvider.Meta().(*Client)
		ContextActionResp, err := c.GetCSEContextAction(rs.Primary.ID)
		if err != nil {
			return err
		}

		*ContextAction = *ContextActionResp

		return nil
	}
}

func testCheckContextActionValues(ContextAction *CSEContextAction, nName string, nType string, nTemplate string, nIocTypes []string, nEntityTypes []string, nRecordFields []string, nAllRecordFields bool, nEnabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if ContextAction.Name != nName {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", nName, ContextAction.Name)
		}
		if ContextAction.Type != nType {
			return fmt.Errorf("bad type, expected \"%s\", got: %#v", nType, ContextAction.Type)
		}
		if ContextAction.Template != nTemplate {
			return fmt.Errorf("bad template, expected \"%s\", got: %#v", nTemplate, ContextAction.Template)
		}
		if ContextAction.IocTypes != nil {
			if len(ContextAction.IocTypes) != len(nIocTypes) {
				return fmt.Errorf("bad ioc_types list lenght, expected \"%d\", got: %d", len(nIocTypes), len(ContextAction.IocTypes))
			}
			if ContextAction.IocTypes[0] != nIocTypes[0] {
				return fmt.Errorf("bad ioc_types in list, expected \"%s\", got: %s", nIocTypes[0], ContextAction.IocTypes[0])
			}
		}
		if ContextAction.EntityTypes != nil {
			if len(ContextAction.EntityTypes) != len(nEntityTypes) {
				return fmt.Errorf("bad entity_types list lenght, expected \"%d\", got: %d", len(nEntityTypes), len(ContextAction.EntityTypes))
			}
			if ContextAction.EntityTypes[0] != nEntityTypes[0] {
				return fmt.Errorf("bad entity_types in list, expected \"%s\", got: %s", nEntityTypes[0], ContextAction.EntityTypes[0])
			}
		}
		if ContextAction.RecordFields != nil {
			if len(ContextAction.RecordFields) != len(nRecordFields) {
				return fmt.Errorf("bad record_fields list lenght, expected \"%d\", got: %d", len(nRecordFields), len(ContextAction.RecordFields))
			}
			if ContextAction.RecordFields[0] != nRecordFields[0] {
				return fmt.Errorf("bad record_fields in list, expected \"%s\", got: %s", nRecordFields[0], ContextAction.RecordFields[0])
			}
		}
		if ContextAction.AllRecordFields != nAllRecordFields {
			return fmt.Errorf("bad all_record_fields field, expected \"%t\", got: %#v", nAllRecordFields, ContextAction.AllRecordFields)
		}
		if ContextAction.Enabled != nEnabled {
			return fmt.Errorf("bad enabled field, expected \"%t\", got: %#v", nEnabled, ContextAction.Enabled)
		}

		return nil
	}
}
