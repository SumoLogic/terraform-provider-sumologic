package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicCSEMatchRule_createAndUpdate(t *testing.T) {
	var matchRule CSEMatchRule
	descriptionExpression := "Test description"
	enabled := true
	entitySelectorEntityType := "_ip"
	entitySelectorExpression := "srcDevice_ip"
	expression := "foo = bar"
	isPrototype := false
	name := "Test Match Rule"
	nameExpression := "Signal Name"
	severityMappingType := "constant"
	severityMappingDefault := 5
	summaryExpression := "Signal Summary"
	tag := "foo"

	nameUpdated := "Updated Match Rule"
	tagUpdated := "bar"

	resourceName := "sumologic_cse_match_rule.match_rule"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEMatchRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSEMatchRuleConfig(descriptionExpression, enabled,
					entitySelectorEntityType, entitySelectorExpression, expression,
					isPrototype, name, nameExpression, severityMappingType, severityMappingDefault,
					summaryExpression, tag),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEMatchRuleExists(resourceName, &matchRule),
					testCheckMatchRuleValues(&matchRule, descriptionExpression, enabled,
						entitySelectorEntityType, entitySelectorExpression, expression,
						isPrototype, name, nameExpression, severityMappingType, severityMappingDefault,
						summaryExpression, tag),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testCreateCSEMatchRuleConfig(descriptionExpression, enabled,
					entitySelectorEntityType, entitySelectorExpression, expression,
					isPrototype, nameUpdated, nameExpression, severityMappingType, severityMappingDefault,
					summaryExpression, tagUpdated),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEMatchRuleExists(resourceName, &matchRule),
					testCheckMatchRuleValues(&matchRule, descriptionExpression, enabled,
						entitySelectorEntityType, entitySelectorExpression, expression,
						isPrototype, nameUpdated, nameExpression, severityMappingType, severityMappingDefault,
						summaryExpression, tagUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccCSEMatchRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_match_rule" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Match Rule destruction check: CSE Match Rule ID is not set")
		}

		s, err := client.GetCSEMatchRule(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("Match rule still exists")
		}
	}
	return nil
}

func testCreateCSEMatchRuleConfig(
	descriptionExpression string, enabled bool, entitySelectorEntityType string,
	entitySelectorExpression string, expression string, isPrototype bool, name string,
	nameExpression string, severityMappingType string, severityMappingDefault int,
	summaryExpression string, tag string) string {
	return fmt.Sprintf(`
resource "sumologic_cse_match_rule" "match_rule" {
	description_expression = "%s"
    enabled = %t
    entity_selectors {
    	entity_type = "%s"
    	expression = "%s"
    }
    expression = "%s"
    is_prototype = %t
    name = "%s"
    name_expression = "%s"
    severity_mapping {
    	type = "%s"
        default = %d
    }
    summary_expression = "%s"
    tags = ["%s"]
}
`, descriptionExpression, enabled, entitySelectorEntityType, entitySelectorExpression,
		expression, isPrototype, name, nameExpression, severityMappingType, severityMappingDefault,
		summaryExpression, tag)
}

func testCheckCSEMatchRuleExists(n string, matchRule *CSEMatchRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("match rule ID is not set")
		}

		c := testAccProvider.Meta().(*Client)
		matchRuleResp, err := c.GetCSEMatchRule(rs.Primary.ID)
		if err != nil {
			return err
		}

		*matchRule = *matchRuleResp

		return nil
	}
}

func testCheckMatchRuleValues(matchRule *CSEMatchRule, descriptionExpression string, enabled bool,
	entitySelectorEntityType string, entitySelectorExpression string, expression string,
	isPrototype bool, name string, nameExpression string, severityMappingType string, severityMappingDefault int,
	summaryExpression string, tag string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if matchRule.DescriptionExpression != descriptionExpression {
			return fmt.Errorf("bad descriptionExpression, expected \"%s\", got %#v", descriptionExpression, matchRule.DescriptionExpression)
		}
		if matchRule.Enabled != enabled {
			return fmt.Errorf("bad enabled, expected \"%t\", got %#v", enabled, matchRule.Enabled)
		}
		if matchRule.EntitySelectors[0].EntityType != entitySelectorEntityType {
			return fmt.Errorf("bad entitySelectorEntityType, expected \"%s\", got %#v", entitySelectorEntityType, matchRule.EntitySelectors[0].EntityType)
		}
		if matchRule.EntitySelectors[0].Expression != entitySelectorExpression {
			return fmt.Errorf("bad entitySelectorExpression, expected \"%s\", got %#v", entitySelectorExpression, matchRule.EntitySelectors[0].Expression)
		}
		if matchRule.Expression != expression {
			return fmt.Errorf("bad expression, expected \"%s\", got %#v", expression, matchRule.Expression)
		}
		if matchRule.IsPrototype != isPrototype {
			return fmt.Errorf("bad isPrototype, expected \"%t\", got %#v", isPrototype, matchRule.IsPrototype)
		}
		if matchRule.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got %#v", name, matchRule.Name)
		}
		if matchRule.NameExpression != nameExpression {
			return fmt.Errorf("bad nameExpression, expected \"%s\", got %#v", nameExpression, matchRule.NameExpression)
		}
		if matchRule.SeverityMapping.Type != severityMappingType {
			return fmt.Errorf("bad severityMappingType, expected \"%s\", got %#v", severityMappingType, matchRule.SeverityMapping.Type)
		}
		if matchRule.SeverityMapping.Default != severityMappingDefault {
			return fmt.Errorf("bad severityMappingDefault, expected \"%d\", got %#v", severityMappingDefault, matchRule.SeverityMapping.Default)
		}
		if matchRule.SummaryExpression != summaryExpression {
			return fmt.Errorf("bad summaryExpression, expected \"%s\", got %#v", summaryExpression, matchRule.SummaryExpression)
		}
		if matchRule.Tags[0] != tag {
			return fmt.Errorf("bad tag, expected \"%s\", got %#v", tag, matchRule.Tags[0])
		}

		return nil
	}
}
