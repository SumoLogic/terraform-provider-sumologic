package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicDataMaskRule_basic(t *testing.T) {
	var rule DataMaskRule
	resourceName := "sumologic_data_mask_rule.test"
	testName := fmt.Sprintf("terraform_data_mask_%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataMaskRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataMaskRuleConfig(testName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataMaskRuleExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "name", testName),
					resource.TestCheckResourceAttr(resourceName, "regex_pattern", "[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}"),
					resource.TestCheckResourceAttr(resourceName, "mask_string", "[REDACTED]"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
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

func TestAccSumologicDataMaskRule_update(t *testing.T) {
	resourceName := "sumologic_data_mask_rule.test"
	testName := fmt.Sprintf("terraform_data_mask_%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataMaskRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataMaskRuleConfig(testName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "mask_string", "[REDACTED]"),
				),
			},
			{
				Config: testAccDataMaskRuleConfigUpdated(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "mask_string", "***MASKED***"),
					resource.TestCheckResourceAttr(resourceName, "description", "updated by terraform acceptance tests"),
				),
			},
		},
	})
}

func testAccCheckDataMaskRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "sumologic_data_mask_rule" {
			continue
		}

		rule, err := client.GetDataMaskRule(r.Primary.ID)
		if err != nil {
			if isDataMaskRuleNotFoundErr(err) {
				continue
			}
			return fmt.Errorf("error retrieving data mask rule during destroy check: %w", err)
		}
		if rule != nil {
			return fmt.Errorf("data mask rule still exists: %s", r.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDataMaskRuleExists(n string, rule *DataMaskRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no data mask rule ID is set")
		}

		client := testAccProvider.Meta().(*Client)
		r, err := client.GetDataMaskRule(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error retrieving data mask rule: %w", err)
		}
		if r == nil {
			return fmt.Errorf("data mask rule not found: %s", rs.Primary.ID)
		}

		*rule = *r
		return nil
	}
}

func testAccDataMaskRuleConfig(testName string, enabled bool) string {
	return fmt.Sprintf(`
resource "sumologic_data_mask_rule" "test" {
  name          = "%s"
  regex_pattern = "[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}"
  mask_string   = "[REDACTED]"
  enabled       = %t
  description   = "managed by terraform acceptance tests"
}
`, testName, enabled)
}

func testAccDataMaskRuleConfigUpdated(testName string) string {
	return fmt.Sprintf(`
resource "sumologic_data_mask_rule" "test" {
  name          = "%s"
  regex_pattern = "[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}"
  mask_string   = "***MASKED***"
  enabled       = false
  description   = "updated by terraform acceptance tests"
}
`, testName)
}
