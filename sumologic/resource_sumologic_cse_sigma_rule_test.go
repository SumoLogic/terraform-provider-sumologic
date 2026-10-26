package sumologic

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicCSESigmaRule_createAndUpdate(t *testing.T) {
	SkipCseTest(t)

	resourceName := "sumologic_cse_sigma_rule.sigma_rule"
	name := acctest.RandomWithPrefix("Test Sigma Rule")
	nameUpdated := acctest.RandomWithPrefix("Updated Sigma Rule")

	testAccPreCheck(t)
	sigmaID := testAccCSESigmaRuleID(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSESigmaRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCSESigmaRuleConfig(sigmaID, name, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "sigma_id", sigmaID),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "original_source", "sigma"),
					resource.TestCheckResourceAttrSet(resourceName, "category"),
					resource.TestCheckResourceAttrSet(resourceName, "expression"),
					resource.TestCheckResourceAttrSet(resourceName, "name_expression"),
					resource.TestCheckResourceAttrSet(resourceName, "description_expression"),
					resource.TestCheckResourceAttr(resourceName, "severity_mapping.#", "1"),
				),
			},
			{
				Config: testCSESigmaRuleConfig(sigmaID, nameUpdated, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					// The original source must survive updates.
					resource.TestCheckResourceAttr(resourceName, "original_source", "sigma"),
				),
			},
		},
	})
}

func TestAccSumologicCSESigmaRule_invalidSigmaID(t *testing.T) {
	SkipCseTest(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSESigmaRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
resource "sumologic_cse_sigma_rule" "sigma_rule" {
	sigma_id = "not-a-uuid"
	enabled  = true
}
`,
				ExpectError: regexp.MustCompile("expected \"sigma_id\" to be a valid UUID"),
			},
		},
	})
}

// testAccCSESigmaRuleID returns the ID of an available Sigma rule.
func testAccCSESigmaRuleID(t *testing.T) string {
	c, err := NewClient(
		os.Getenv("SUMOLOGIC_ACCESSID"),
		os.Getenv("SUMOLOGIC_ACCESSKEY"),
		"",
		os.Getenv("SUMOLOGIC_ENVIRONMENT"),
		os.Getenv("SUMOLOGIC_BASE_URL"),
		false,
	)
	if err != nil {
		t.Fatalf("unable to create client: %v", err)
	}
	rules, err := c.ListCSESigmaRules(1)
	if err != nil {
		t.Fatalf("unable to list Sigma rules: %v", err)
	}
	if len(rules) == 0 || rules[0].ConversionMetadata.SigmaUID == "" {
		t.Skip("no Sigma rules available")
	}
	return rules[0].ConversionMetadata.SigmaUID
}

func testAccCSESigmaRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_sigma_rule" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Sigma Rule destruction check: ID is not set")
		}

		rule, err := client.GetCSEMatchRule(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: %w", err)
		}
		if rule != nil {
			return fmt.Errorf("Sigma rule still exists")
		}
	}
	return nil
}

func testCSESigmaRuleConfig(sigmaID string, name string, enabled bool) string {
	return fmt.Sprintf(`
resource "sumologic_cse_sigma_rule" "sigma_rule" {
	sigma_id = "%s"
	name     = "%s"
	enabled  = %t
}
`, sigmaID, name, enabled)
}
