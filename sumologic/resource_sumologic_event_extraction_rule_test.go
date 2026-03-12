package sumologic

import (
	"fmt"
	"strings"
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
					// Validate basic attributes
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "query", "_sourceCategory=deployments"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),

					// Validate configuration count
					resource.TestCheckResourceAttr(resourceName, "configuration.#", "4"),

					// Validate individual configuration fields (order may vary, so we check using TypeSetElemNestedAttrs)
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"configuration.*",
						map[string]string{
							"field_name":   "eventType",
							"value_source": "Deployment",
							"mapping_type": "HardCoded",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"configuration.*",
						map[string]string{
							"field_name":   "eventPriority",
							"value_source": "High",
							"mapping_type": "HardCoded",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"configuration.*",
						map[string]string{
							"field_name":   "eventSource",
							"value_source": "Jenkins",
							"mapping_type": "HardCoded",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"configuration.*",
						map[string]string{
							"field_name":   "eventName",
							"value_source": "monitor-manager deployed",
							"mapping_type": "HardCoded",
						},
					),
				),
			},
			{
				Config: testAccSumologicEventExtractionRuleUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					// Validate basic attributes
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "updated description"),
					resource.TestCheckResourceAttr(resourceName, "query", "_sourceCategory=deployments"),

					// Validate configuration count (now 5 fields)
					resource.TestCheckResourceAttr(resourceName, "configuration.#", "5"),

					// Validate the new eventDescription field was added
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"configuration.*",
						map[string]string{
							"field_name":   "eventDescription",
							"value_source": "2 containers upgraded",
							"mapping_type": "HardCoded",
						},
					),
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
			// If the error is "not found" or "record_not_found", the resource was successfully destroyed
			if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "record_not_found") {
				continue
			}
			return err
		}
		if found != nil {
			return fmt.Errorf("event extraction rule %s still exists", r.Primary.ID)
		}
	}
	return nil
}
