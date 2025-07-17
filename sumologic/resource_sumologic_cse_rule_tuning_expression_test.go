package sumologic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccSumologicSCERuleTuningExpression_create(t *testing.T) {
	SkipCseTest(t)

	var ruleTuningExpression CSERuleTuningExpression
	nName := acctest.RandomWithPrefix("New Rule Tuning Name")
	nDescription := "New Rule Tuning Description"
	nExpression := "accountId = 1234"
	nEnabled := true
	nExclude := true
	nIsGlobal := false
	nRuleIds := []string{"LEGACY-S00084", "THRESHOLD-S00514", "AGGREGATION-S00002"}
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
					testCheckRuleTuningExpressionValues(t, &ruleTuningExpression, nName, nDescription, nExpression, nEnabled, nExclude, nIsGlobal, nRuleIds),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
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
	quotedRuleIds := make([]string, len(nRuleIds))
	for i, id := range nRuleIds {
		quotedRuleIds[i] = fmt.Sprintf(`"%s"`, id)
	}

	return fmt.Sprintf(`
resource "sumologic_cse_rule_tuning_expression" "rule_tuning_expression" {
	name = "%s"
	description = "%s"
	expression = "%s"
	enabled = "%t"
	exclude = "%t"
	is_global = "%t"
	rule_ids = [%s]
}
`, nName, nDescription, nExpression, nEnabled, nExclude, nIsGlobal, strings.Join(quotedRuleIds, ", "))
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

func testCheckRuleTuningExpressionValues(t *testing.T, ruleTuningExpression *CSERuleTuningExpression, nName string, nDescription string, nExpression string, nEnabled bool, nExclude bool, nIsGlobal bool, nRuleIds []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		assert.Equal(t, nName, ruleTuningExpression.Name)
		assert.Equal(t, nDescription, ruleTuningExpression.Description)
		assert.Equal(t, nExpression, ruleTuningExpression.Expression)
		assert.Equal(t, nEnabled, ruleTuningExpression.Enabled)
		assert.Equal(t, nExclude, ruleTuningExpression.Exclude)
		assert.Equal(t, nIsGlobal, ruleTuningExpression.IsGlobal)
		assert.ElementsMatch(t, nRuleIds, ruleTuningExpression.RuleIds)
		return nil
	}
}
