package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicCSEAggregationRule_createAndUpdate(t *testing.T) {
	var AggregationRule CSEAggregationRule
	aggregationFunctionName := "distinct_eventid_count"
	aggregationFunction := "count_distinct"
	aggregationFunctionArgument := "metadata_deviceEventId"
	descriptionExpression := "Test description"
	enabled := true
	entitySelectorEntityType := "_ip"
	entitySelectorExpression := "srcDevice_ip"
	groupByEntity := true
	groupByField := "dstDevice_hostname"
	isPrototype := false
	matchExpression := "foo = bar"
	name := "Test Aggregation Rule"
	nameExpression := "Signal Name"
	severityMappingType := "constant"
	severityMappingDefault := 5
	summaryExpression := "Signal Summary"
	triggerExpression := "foo = bar"
	tag := "foo"
	windowSize := "T30M"

	nameUpdated := "Updated Aggregation Rule"
	tagUpdated := "bar"

	resourceName := "sumologic_cse_aggregation_rule.aggregation_rule"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEAggregationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSEAggregationRuleConfig(aggregationFunctionName, aggregationFunction,
					aggregationFunctionArgument, descriptionExpression, enabled, entitySelectorEntityType,
					entitySelectorExpression, groupByEntity, groupByField, isPrototype, matchExpression,
					name, nameExpression, severityMappingType, severityMappingDefault, summaryExpression,
					triggerExpression, tag, windowSize),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEAggregationRuleExists(resourceName, &AggregationRule),
					testCheckAggregationRuleValues(&AggregationRule, aggregationFunctionName,
						aggregationFunction, aggregationFunctionArgument, descriptionExpression, enabled,
						entitySelectorEntityType, entitySelectorExpression, groupByEntity, groupByField,
						isPrototype, matchExpression, name, nameExpression, severityMappingType,
						severityMappingDefault, summaryExpression, triggerExpression, tag, windowSize),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testCreateCSEAggregationRuleConfig(aggregationFunctionName, aggregationFunction,
					aggregationFunctionArgument, descriptionExpression, enabled, entitySelectorEntityType,
					entitySelectorExpression, groupByEntity, groupByField, isPrototype, matchExpression,
					nameUpdated, nameExpression, severityMappingType, severityMappingDefault, summaryExpression,
					triggerExpression, tagUpdated, windowSize),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEAggregationRuleExists(resourceName, &AggregationRule),
					testCheckAggregationRuleValues(&AggregationRule, aggregationFunctionName,
						aggregationFunction, aggregationFunctionArgument, descriptionExpression, enabled,
						entitySelectorEntityType, entitySelectorExpression, groupByEntity, groupByField,
						isPrototype, matchExpression, nameUpdated, nameExpression, severityMappingType,
						severityMappingDefault, summaryExpression, triggerExpression, tagUpdated, windowSize),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccCSEAggregationRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_aggregation_rule" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Aggregation Rule destruction check: CSE Aggregation Rule ID is not set")
		}

		s, err := client.GetCSEAggregationRule(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("Aggregation rule still exists")
		}
	}
	return nil
}

func testCreateCSEAggregationRuleConfig(
	aggregationFunctionName string, aggregationFunction string,
	aggregationFunctionArgument string, descriptionExpression string, enabled bool,
	entitySelectorEntityType string, entitySelectorExpression string, groupByEntity bool,
	groupByField string, isPrototype bool, matchExpression string, name string,
	nameExpression string, severityMappingType string, severityMappingDefault int,
	summaryExpression string, triggerExpression string, tag string, windowSize string) string {
	return fmt.Sprintf(`
resource "sumologic_cse_aggregation_rule" "aggregation_rule" {
	aggregation_functions {
		name = "%s"
		function = "%s"
		arguments = ["%s"]
	}
	description_expression = "%s"
    enabled = %t
    entity_selectors {
    	entity_type = "%s"
    	expression = "%s"
    }
    group_by_entity = %t
    group_by_fields = ["%s"]
    is_prototype = %t
    match_expression = "%s"
    name = "%s"
    name_expression = "%s"
    severity_mapping {
    	type = "%s"
        default = %d
    }
    summary_expression = "%s"
    trigger_expression = "%s"
    tags = ["%s"]
    window_size = "%s"
}
`, aggregationFunctionName, aggregationFunction, aggregationFunctionArgument,
		descriptionExpression, enabled, entitySelectorEntityType, entitySelectorExpression,
		groupByEntity, groupByField, isPrototype, matchExpression, name, nameExpression,
		severityMappingType, severityMappingDefault, summaryExpression, triggerExpression,
		tag, windowSize)
}

func testCheckCSEAggregationRuleExists(n string, AggregationRule *CSEAggregationRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("match rule ID is not set")
		}

		c := testAccProvider.Meta().(*Client)
		AggregationRuleResp, err := c.GetCSEAggregationRule(rs.Primary.ID)
		if err != nil {
			return err
		}

		*AggregationRule = *AggregationRuleResp

		return nil
	}
}

func testCheckAggregationRuleValues(AggregationRule *CSEAggregationRule, aggregationFunctionName string,
	aggregationFunction string, aggregationFunctionArgument string, descriptionExpression string,
	enabled bool, entitySelectorEntityType string, entitySelectorExpression string, groupByEntity bool,
	groupByField string, isPrototype bool, matchExpression string, name string, nameExpression string,
	severityMappingType string, severityMappingDefault int, summaryExpression string,
	triggerExpression string, tag string, windowSize string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if AggregationRule.AggregationFunctions[0].Name != aggregationFunctionName {
			return fmt.Errorf("bad aggregationFunctionName, expected \"%s\", got %#v", aggregationFunctionName, AggregationRule.AggregationFunctions[0].Name)
		}
		if AggregationRule.AggregationFunctions[0].Function != aggregationFunction {
			return fmt.Errorf("bad aggregationFunction, expected \"%s\", got %#v", aggregationFunction, AggregationRule.AggregationFunctions[0].Function)
		}
		if AggregationRule.AggregationFunctions[0].Arguments[0] != aggregationFunctionArgument {
			return fmt.Errorf("bad aggregationFunctionArgument, expected \"%s\", got %#v", aggregationFunctionArgument, AggregationRule.AggregationFunctions[0].Arguments[0])
		}
		if AggregationRule.DescriptionExpression != descriptionExpression {
			return fmt.Errorf("bad descriptionExpression, expected \"%s\", got %#v", descriptionExpression, AggregationRule.DescriptionExpression)
		}
		if AggregationRule.Enabled != enabled {
			return fmt.Errorf("bad enabled, expected \"%t\", got %#v", enabled, AggregationRule.Enabled)
		}
		if AggregationRule.EntitySelectors[0].EntityType != entitySelectorEntityType {
			return fmt.Errorf("bad entitySelectorEntityType, expected \"%s\", got %#v", entitySelectorEntityType, AggregationRule.EntitySelectors[0].EntityType)
		}
		if AggregationRule.EntitySelectors[0].Expression != entitySelectorExpression {
			return fmt.Errorf("bad entitySelectorExpression, expected \"%s\", got %#v", entitySelectorExpression, AggregationRule.EntitySelectors[0].Expression)
		}
		if AggregationRule.GroupByEntity != groupByEntity {
			return fmt.Errorf("bad groupByEntity, expected \"%t\", got %#v", groupByEntity, AggregationRule.GroupByEntity)
		}
		if AggregationRule.GroupByFields[0] != groupByField {
			return fmt.Errorf("bad groupByField, expected \"%s\", got %#v", groupByField, AggregationRule.GroupByFields[0])
		}
		if AggregationRule.IsPrototype != isPrototype {
			return fmt.Errorf("bad isPrototype, expected \"%t\", got %#v", isPrototype, AggregationRule.IsPrototype)
		}
		if AggregationRule.MatchExpression != matchExpression {
			return fmt.Errorf("bad matchExpression, expected \"%s\", got %#v", matchExpression, AggregationRule.MatchExpression)
		}
		if AggregationRule.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got %#v", name, AggregationRule.Name)
		}
		if AggregationRule.NameExpression != nameExpression {
			return fmt.Errorf("bad nameExpression, expected \"%s\", got %#v", nameExpression, AggregationRule.NameExpression)
		}
		if AggregationRule.SeverityMapping.Type != severityMappingType {
			return fmt.Errorf("bad severityMappingType, expected \"%s\", got %#v", severityMappingType, AggregationRule.SeverityMapping.Type)
		}
		if AggregationRule.SeverityMapping.Default != severityMappingDefault {
			return fmt.Errorf("bad severityMappingDefault, expected \"%d\", got %#v", severityMappingDefault, AggregationRule.SeverityMapping.Default)
		}
		if AggregationRule.SummaryExpression != summaryExpression {
			return fmt.Errorf("bad summaryExpression, expected \"%s\", got %#v", summaryExpression, AggregationRule.SummaryExpression)
		}
		if AggregationRule.Tags[0] != tag {
			return fmt.Errorf("bad tag, expected \"%s\", got %#v", tag, AggregationRule.Tags[0])
		}
		if AggregationRule.TriggerExpression != triggerExpression {
			return fmt.Errorf("bad triggerExpression, expected \"%s\", got %#v", triggerExpression, AggregationRule.TriggerExpression)
		}
		if AggregationRule.WindowSizeName != windowSize {
			return fmt.Errorf("bad windowSize, expected \"%s\", got %#v", windowSize, AggregationRule.WindowSize)
		}

		return nil
	}
}
