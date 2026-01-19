package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicEventExtractionRule_crud(t *testing.T) {
	suffix := acctest.RandString(10)
	name := "terraform_event_rule_" + suffix
	resourceName := "sumologic_event_extraction_rule.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEventExtractionRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicEventExtractionRule(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "query", "_sourceCategory=deployments"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					// Note: Testing specific indices in a list.
					// If the order is non-deterministic, use TestCheckTypeSetElemAttr instead.
					resource.TestCheckResourceAttr(resourceName, "configuration.0.field_name", "eventType"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.value_source", "Deployment"),
					resource.TestCheckResourceAttr(resourceName, "configuration.1.field_name", "eventPriority"),
					resource.TestCheckResourceAttr(resourceName, "configuration.1.value_source", "High"),
				),
			},
			{
				Config: testAccSumologicEventExtractionRuleUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "updated description"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					// Verifying the new added block
					resource.TestCheckResourceAttr(resourceName, "configuration.4.field_name", "eventDescription"),
					resource.TestCheckResourceAttr(resourceName, "configuration.4.value_source", "2 containers upgraded"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSumologicEventExtractionRule(name string) string {
	return fmt.Sprintf(`
resource "sumologic_event_extraction_rule" "test" {
  name  = "%s"
  query = "_sourceCategory=deployments"

  configuration {
    field_name   = "eventType"
    value_source = "Deployment"
  }
  configuration {
    field_name   = "eventPriority"
    value_source = "High"
  }
  configuration {
    field_name   = "eventSource"
    value_source = "Jenkins"
  }
  configuration {
    field_name   = "eventName"
    value_source = "monitor-manager deployed"
  }
}
`, name)
}

func testAccSumologicEventExtractionRuleUpdate(name string) string {
	return fmt.Sprintf(`
resource "sumologic_event_extraction_rule" "test" {
  name        = "%s"
  description = "updated description"
  query       = "_sourceCategory=deployments"
  enabled     = false

  configuration {
    field_name   = "eventType"
    value_source = "Deployment"
  }
  configuration {
    field_name   = "eventPriority"
    value_source = "High"
  }
  configuration {
    field_name   = "eventSource"
    value_source = "Jenkins"
  }
  configuration {
    field_name   = "eventName"
    value_source = "monitor-manager deployed"
  }
  configuration {
    field_name   = "eventDescription"
    value_source = "2 containers upgraded"
  }
}
`, name)
}

func testAccCheckEventExtractionRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, r := range s.RootModule().Resources {
		if r.Type != "sumologic_event_extraction_rule" {
			continue
		}

		found, err := client.GetEventExtractionRule(r.Primary.ID)
		if err != nil {
			return err
		}
		if found != nil {
			return fmt.Errorf("event extraction rule %s still exists", r.Primary.ID)
		}
	}
	return nil
}
