package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicCSEThresholdRule_createAndUpdate(t *testing.T) {
	SkipCseTest(t)

	var thresholdRule CSEThresholdRule
	countDistinct := true
	countField := "dstDevice_hostname"
	description := "Test description"
	enabled := true
	entitySelectorEntityType := "_ip"
	entitySelectorExpression := "srcDevice_ip"
	expression := "foo = bar"
	groupByField := "destPort"
	isPrototype := false
	limit := 20
	name := "Test Threshold Rule"
	severity := 5
	summaryExpression := "Signal Summary"
	tag := "foo"
	windowSize := "T30M"

	nameUpdated := "Updated Threshold Rule"
	windowSizeUpdated := "T12H"
	resourceName := "sumologic_cse_threshold_rule.threshold_rule"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEThresholdRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSEThresholdRuleConfig(countDistinct, countField,
					description, enabled, entitySelectorEntityType,
					entitySelectorExpression, expression, groupByField, isPrototype,
					limit, name, severity, summaryExpression, tag, windowSize),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEThresholdRuleExists(resourceName, &thresholdRule),
					testCheckThresholdRuleValues(&thresholdRule, countDistinct, countField,
						description, enabled, entitySelectorEntityType,
						entitySelectorExpression, expression, groupByField, isPrototype,
						limit, name, severity, summaryExpression, tag, windowSize),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testCreateCSEThresholdRuleConfig(countDistinct, countField,
					description, enabled, entitySelectorEntityType,
					entitySelectorExpression, expression, groupByField, isPrototype,
					limit, nameUpdated, severity, summaryExpression, tag, windowSizeUpdated),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEThresholdRuleExists(resourceName, &thresholdRule),
					testCheckThresholdRuleValues(&thresholdRule, countDistinct, countField,
						description, enabled, entitySelectorEntityType,
						entitySelectorExpression, expression, groupByField, isPrototype,
						limit, nameUpdated, severity, summaryExpression, tag, windowSizeUpdated),
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

func testAccCSEThresholdRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_threshold_rule" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Threshold Rule destruction check: CSE Threshold Rule ID is not set")
		}

		s, err := client.GetCSEThresholdRule(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("Threshold rule still exists")
		}
	}
	return nil
}

func testCreateCSEThresholdRuleConfig(
	countDistinct bool, countField string, description string, enabled bool,
	entitySelectorEntityType string, entitySelectorExpression string,
	expression string, groupByField string, isPrototype bool, limit int,
	name string, severity int, summaryExpression string, tag string,
	windowSize string) string {
	return fmt.Sprintf(`
resource "sumologic_cse_threshold_rule" "threshold_rule" {
	count_distinct = %t
	count_field = "%s"
	description = "%s"
    enabled = %t
    entity_selectors {
    	entity_type = "%s"
    	expression = "%s"
    }
    expression = "%s"
    group_by_fields = ["%s"]
    is_prototype = %t
    limit = %d
    name = "%s"
    severity = %d
    summary_expression = "%s"
    tags = ["%s"]
    window_size = "%s"
}
`, countDistinct, countField, description, enabled, entitySelectorEntityType,
		entitySelectorExpression, expression, groupByField, isPrototype, limit, name,
		severity, summaryExpression, tag, windowSize)
}

func testCheckCSEThresholdRuleExists(n string, thresholdRule *CSEThresholdRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("threshold rule ID is not set")
		}

		c := testAccProvider.Meta().(*Client)
		thresholdRuleResp, err := c.GetCSEThresholdRule(rs.Primary.ID)
		if err != nil {
			return err
		}

		*thresholdRule = *thresholdRuleResp

		return nil
	}
}

func testCheckThresholdRuleValues(thresholdRule *CSEThresholdRule, countDistinct bool,
	countField string, description string, enabled bool, entitySelectorEntityType string,
	entitySelectorExpression string, expression string, groupByField string,
	isPrototype bool, limit int, name string, severity int, summaryExpression string,
	tag string, windowSize string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if thresholdRule.CountDistinct != countDistinct {
			return fmt.Errorf("bad countDistinct, expected \"%t\", got %#v", countDistinct, thresholdRule.CountDistinct)
		}
		if thresholdRule.CountField != countField {
			return fmt.Errorf("bad countField, expected \"%s\", got %#v", countField, thresholdRule.CountField)
		}
		if thresholdRule.Description != description {
			return fmt.Errorf("bad description, expected \"%s\", got %#v", description, thresholdRule.Description)
		}
		if thresholdRule.Enabled != enabled {
			return fmt.Errorf("bad enabled, expected \"%t\", got %#v", enabled, thresholdRule.Enabled)
		}
		if thresholdRule.EntitySelectors[0].EntityType != entitySelectorEntityType {
			return fmt.Errorf("bad entitySelectorEntityType, expected \"%s\", got %#v", entitySelectorEntityType, thresholdRule.EntitySelectors[0].EntityType)
		}
		if thresholdRule.EntitySelectors[0].Expression != entitySelectorExpression {
			return fmt.Errorf("bad entitySelectorExpression, expected \"%s\", got %#v", entitySelectorExpression, thresholdRule.EntitySelectors[0].Expression)
		}
		if thresholdRule.Expression != expression {
			return fmt.Errorf("bad expression, expected \"%s\", got %#v", expression, thresholdRule.Expression)
		}
		if thresholdRule.GroupByFields[0] != groupByField {
			return fmt.Errorf("bad groupByField, expected \"%s\", got %#v", groupByField, thresholdRule.GroupByFields[0])
		}
		if thresholdRule.IsPrototype != isPrototype {
			return fmt.Errorf("bad isPrototype, expected \"%t\", got %#v", isPrototype, thresholdRule.IsPrototype)
		}
		if thresholdRule.Limit != limit {
			return fmt.Errorf("bad limit, expected \"%d\", got %#v", limit, thresholdRule.Limit)
		}
		if thresholdRule.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got %#v", name, thresholdRule.Name)
		}
		if thresholdRule.Severity != severity {
			return fmt.Errorf("bad severity, expected \"%d\", got %#v", severity, thresholdRule.Severity)
		}
		if thresholdRule.SummaryExpression != summaryExpression {
			return fmt.Errorf("bad summaryExpression, expected \"%s\", got %#v", summaryExpression, thresholdRule.SummaryExpression)
		}
		if thresholdRule.Tags[0] != tag {
			return fmt.Errorf("bad tag, expected \"%s\", got %#v", tag, thresholdRule.Tags[0])
		}
		if thresholdRule.WindowSizeName != windowSize {
			return fmt.Errorf("bad windowSize, expected \"%s\", got %#v", windowSize, thresholdRule.WindowSize)
		}

		return nil
	}
}
