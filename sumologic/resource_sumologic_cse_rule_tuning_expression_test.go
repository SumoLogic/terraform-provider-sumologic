package sumologic

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicSCERuleTuningExpression_create(t *testing.T) {
	var ruleTuningExpression CSERuleTuningExpression
	nName := "New Rule Tuning Name"
	nDescription := "New Rule Tuning Description"
	nExpression := "expression"
	nEnabled := true
	nExclude := true
	nIsGlobal := false
	nRuleIds := []string{"LEGACY-S00084"}
	resourceName := "sumologic_cse_rule_tuning_expression.rule_tuning_expression"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSERuleTuningExpressionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSERuleTuningExpressionConfig(nName, nDescription, nExpression, nEnabled, nExclude, nIsGlobal, nRuleIds),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSERuleTuningExpressionExists(resourceName, &ruleTuningExpression),
					testCheckRuleTuningExpressionValues(&ruleTuningExpression, nName, nDescription, nExpression, nEnabled, nExclude, nIsGlobal),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccCSERuleTuningExpressionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_rule_tuning_expression" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Rule Tuning Expression destruction check: CSE Rule Tuning Expresion ID is not set")
		}

		s, err := client.GetCSERuleTuningExpression(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("rule tuning expression still exists")
		}
	}
	return nil
}

func testCreateCSERuleTuningExpressionConfig(nName string, nDescription string, nExpression string, nEnabled bool, nExclude bool, nIsGlobal bool, nRuleIds []string) string {

	log.Printf(`
resource "sumologic_cse_rule_tuning_expression" "rule_tuning_expression" {
	name = "%s"
	description = "%s"
	expression = "%s"
	enabled = "%t"
	exclude = "%t"
	is_global = "%t"
	rule_ids = ["%s"]
}
`, nName, nDescription, nExpression, nEnabled, nExclude, nIsGlobal, nRuleIds[0])

	return fmt.Sprintf(`
resource "sumologic_cse_rule_tuning_expression" "rule_tuning_expression" {
	name = "%s"
	description = "%s"
	expression = "%s"
	enabled = "%t"
	exclude = "%t"
	is_global = "%t"
	rule_ids = ["%s"]
}
`, nName, nDescription, nExpression, nEnabled, nExclude, nIsGlobal, nRuleIds[0])
}

func testCheckCSERuleTuningExpressionExists(n string, ruleTuningExpression *CSERuleTuningExpression) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("rule tuning expression ID is not set")
		}

		c := testAccProvider.Meta().(*Client)
		ruleTuningExpressionResp, err := c.GetCSERuleTuningExpression(rs.Primary.ID)
		if err != nil {
			return err
		}

		*ruleTuningExpression = *ruleTuningExpressionResp

		return nil
	}
}

func testCheckRuleTuningExpressionValues(ruleTuningExpression *CSERuleTuningExpression, nName string, nDescription string, nExpression string, nEnabled bool, nExclude bool, nIsGlobal bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if ruleTuningExpression.Name != nName {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", nName, ruleTuningExpression.Name)
		}
		if ruleTuningExpression.Description != nDescription {
			return fmt.Errorf("bad description, expected \"%s\", got: %#v", nDescription, ruleTuningExpression.Description)
		}
		if ruleTuningExpression.Expression != nExpression {
			return fmt.Errorf("bad expression, expected \"%s\", got: %#v", nExpression, ruleTuningExpression.Expression)
		}
		if ruleTuningExpression.Enabled != nEnabled {
			return fmt.Errorf("bad enabled flag, expected \"%t\", got: %#v", nEnabled, ruleTuningExpression.Enabled)
		}
		if ruleTuningExpression.Exclude != nExclude {
			return fmt.Errorf("bad exclude flag, expected \"%t\", got: %#v", nExclude, ruleTuningExpression.Exclude)
		}
		if ruleTuningExpression.IsGlobal != nIsGlobal {
			return fmt.Errorf("bad isGlobal flag, expected \"%t\", got: %#v", nIsGlobal, ruleTuningExpression.IsGlobal)
		}

		return nil
	}
}
