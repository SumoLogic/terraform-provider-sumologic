package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicCSEChainRule_createAndUpdate(t *testing.T) {
	var ChainRule CSEChainRule
	description := "Test description"
	enabled := true
	entitySelectorEntityType := "_ip"
	entitySelectorExpression := "srcDevice_ip"
	expression1 := "foo = bar"
	limit1 := 5
	expression2 := "baz = qux"
	limit2 := 1
	groupByField := "destPort"
	isPrototype := false
	ordered := true
	name := "Test Chain Rule"
	severity := 5
	summaryExpression := "Signal Summary"
	tag := "foo"
	windowSize := "T30M"

	nameUpdated := "Updated Chain Rule"
	windowSizeUpdated := "T12H"

	resourceName := "sumologic_cse_chain_rule.chain_rule"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEChainRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSEChainRuleConfig(description, enabled,
					entitySelectorEntityType, entitySelectorExpression, expression1,
					limit1, expression2, limit2, groupByField, isPrototype, ordered,
					name, severity, summaryExpression, tag, windowSize),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEChainRuleExists(resourceName, &ChainRule),
					testCheckChainRuleValues(&ChainRule, description, enabled,
						entitySelectorEntityType, entitySelectorExpression, expression1,
						limit1, expression2, limit2, groupByField, isPrototype, ordered,
						name, severity, summaryExpression, tag, windowSize),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testCreateCSEChainRuleConfig(description, enabled,
					entitySelectorEntityType, entitySelectorExpression, expression1,
					limit1, expression2, limit2, groupByField, isPrototype, ordered,
					nameUpdated, severity, summaryExpression, tag, windowSizeUpdated),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEChainRuleExists(resourceName, &ChainRule),
					testCheckChainRuleValues(&ChainRule, description, enabled,
						entitySelectorEntityType, entitySelectorExpression, expression1,
						limit1, expression2, limit2, groupByField, isPrototype, ordered,
						nameUpdated, severity, summaryExpression, tag, windowSizeUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccCSEChainRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_chain_rule" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Chain Rule destruction check: CSE Chain Rule ID is not set")
		}

		s, err := client.GetCSEChainRule(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("Chain rule still exists")
		}
	}
	return nil
}

func testCreateCSEChainRuleConfig(
	description string, enabled bool, entitySelectorEntityType string,
	entitySelectorExpression string, expression1 string, limit1 int,
	expression2 string, limit2 int, groupByField string, isPrototype bool,
	ordered bool, name string, severity int, summaryExpression string, tag string,
	windowSize string) string {
	return fmt.Sprintf(`
resource "sumologic_cse_chain_rule" "chain_rule" {
	description = "%s"
    enabled = %t
    entity_selectors {
    	entity_type = "%s"
    	expression = "%s"
    }
    expressions_and_limits {
    	expression = "%s"
    	limit = %d
    }
    expressions_and_limits {
    	expression = "%s"
    	limit = %d
    }
    group_by_fields = ["%s"]
    is_prototype = %t
    ordered = %t
    name = "%s"
    severity = %d
    summary_expression = "%s"
    tags = ["%s"]
    window_size = "%s"
}
`, description, enabled, entitySelectorEntityType, entitySelectorExpression,
		expression1, limit1, expression2, limit2, groupByField, isPrototype, ordered,
		name, severity, summaryExpression, tag, windowSize)
}

func testCheckCSEChainRuleExists(n string, ChainRule *CSEChainRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("chain rule ID is not set")
		}

		c := testAccProvider.Meta().(*Client)
		ChainRuleResp, err := c.GetCSEChainRule(rs.Primary.ID)
		if err != nil {
			return err
		}

		*ChainRule = *ChainRuleResp

		return nil
	}
}

func testCheckChainRuleValues(ChainRule *CSEChainRule, description string,
	enabled bool, entitySelectorEntityType string, entitySelectorExpression string,
	expression1 string, limit1 int, expression2 string, limit2 int, groupByField string,
	isPrototype bool, ordered bool, name string, severity int, summaryExpression string,
	tag string, windowSize string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if ChainRule.Description != description {
			return fmt.Errorf("bad description, expected \"%s\", got %#v", description, ChainRule.Description)
		}
		if ChainRule.Enabled != enabled {
			return fmt.Errorf("bad enabled, expected \"%t\", got %#v", enabled, ChainRule.Enabled)
		}
		if ChainRule.EntitySelectors[0].EntityType != entitySelectorEntityType {
			return fmt.Errorf("bad entitySelectorEntityType, expected \"%s\", got %#v", entitySelectorEntityType, ChainRule.EntitySelectors[0].EntityType)
		}
		if ChainRule.EntitySelectors[0].Expression != entitySelectorExpression {
			return fmt.Errorf("bad entitySelectorExpression, expected \"%s\", got %#v", entitySelectorExpression, ChainRule.EntitySelectors[0].Expression)
		}
		if ChainRule.ExpressionsAndLimits[0].Expression != expression1 {
			return fmt.Errorf("bad expression1, expected \"%s\", got %#v", expression1, ChainRule.ExpressionsAndLimits[0].Expression)
		}
		if ChainRule.ExpressionsAndLimits[0].Limit != limit1 {
			return fmt.Errorf("bad limit1, expected \"%d\", got %#v", limit1, ChainRule.ExpressionsAndLimits[0].Limit)
		}
		if ChainRule.ExpressionsAndLimits[1].Expression != expression2 {
			return fmt.Errorf("bad expression2, expected \"%s\", got %#v", expression2, ChainRule.ExpressionsAndLimits[1].Expression)
		}
		if ChainRule.ExpressionsAndLimits[1].Limit != limit2 {
			return fmt.Errorf("bad limit2, expected \"%d\", got %#v", limit2, ChainRule.ExpressionsAndLimits[1].Limit)
		}
		if ChainRule.GroupByFields[0] != groupByField {
			return fmt.Errorf("bad groupByField, expected \"%s\", got %#v", groupByField, ChainRule.GroupByFields[0])
		}
		if ChainRule.IsPrototype != isPrototype {
			return fmt.Errorf("bad isPrototype, expected \"%t\", got %#v", isPrototype, ChainRule.IsPrototype)
		}
		if ChainRule.Ordered != ordered {
			return fmt.Errorf("bad ordered, expected \"%t\", got %#v", ordered, ChainRule.Ordered)
		}
		if ChainRule.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got %#v", name, ChainRule.Name)
		}
		if ChainRule.Severity != severity {
			return fmt.Errorf("bad severity, expected \"%d\", got %#v", severity, ChainRule.Severity)
		}
		if ChainRule.SummaryExpression != summaryExpression {
			return fmt.Errorf("bad summaryExpression, expected \"%s\", got %#v", summaryExpression, ChainRule.SummaryExpression)
		}
		if ChainRule.Tags[0] != tag {
			return fmt.Errorf("bad tag, expected \"%s\", got %#v", tag, ChainRule.Tags[0])
		}
		if ChainRule.WindowSizeName != windowSize {
			return fmt.Errorf("bad windowSize, expected \"%s\", got %#v", windowSize, ChainRule.WindowSizeName)
		}

		return nil
	}
}
