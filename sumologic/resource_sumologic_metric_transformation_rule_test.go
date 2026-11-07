package sumologic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicMetricTransformationRule_basic(t *testing.T) {
	var rule MetricTransformationRule
	testName := "terraform_test_mtr_" + acctest.RandString(16)
	resourceName := "sumologic_metric_transformation_rule.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetricTransformationRuleDestroy,
		Steps: []resource.TestStep{
			{
				// create without dimension transformations
				Config: testAccSumologicMetricTransformationRuleBasic(testName, "_sourceCategory=terraform_test_mtr", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetricTransformationRuleExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "name", testName),
					resource.TestCheckResourceAttr(resourceName, "selector", "_sourceCategory=terraform_test_mtr"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "retention", "400"),
					resource.TestCheckResourceAttr(resourceName, "transformed_metrics_retention", "0"),
					resource.TestCheckResourceAttr(resourceName, "dimension_transformation.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
				),
			},
			{
				// update: always test this; create can pass while update fails
				Config: testAccSumologicMetricTransformationRuleBasic(testName, "_sourceCategory=terraform_test_mtr_edited", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetricTransformationRuleExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "selector", "_sourceCategory=terraform_test_mtr_edited"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				// import
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSumologicMetricTransformationRule_dimensionTransformations(t *testing.T) {
	var rule MetricTransformationRule
	testName := "terraform_test_mtr_" + acctest.RandString(16)
	resourceName := "sumologic_metric_transformation_rule.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetricTransformationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMetricTransformationRuleWithTransformations(testName, `"metric"`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetricTransformationRuleExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "retention", "0"),
					resource.TestCheckResourceAttr(resourceName, "transformed_metrics_retention", "8"),
					resource.TestCheckResourceAttr(resourceName, "dimension_transformation.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "dimension_transformation.0.transformation_type", "AggregateOnTransformation"),
					resource.TestCheckResourceAttr(resourceName, "dimension_transformation.0.aggregate_on.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dimension_transformation.0.aggregate_on.0", "metric"),
					resource.TestCheckResourceAttr(resourceName, "dimension_transformation.1.transformation_type", "AddOrReplaceTransformation"),
					resource.TestCheckResourceAttr(resourceName, "dimension_transformation.1.dimension_to_replace", "metric"),
					resource.TestCheckResourceAttr(resourceName, "dimension_transformation.1.value", "{{metric}}_aggregated"),
				),
			},
			{
				// update a nested value. "metric" must stay in aggregate_on because the
				// AddOrReplace value references {{metric}}.
				Config: testAccSumologicMetricTransformationRuleWithTransformations(testName, `"metric", "cluster"`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetricTransformationRuleExists(resourceName, &rule),
					resource.TestCheckResourceAttr(resourceName, "dimension_transformation.0.aggregate_on.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "dimension_transformation.0.aggregate_on.1", "cluster"),
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

func TestAccSumologicMetricTransformationRule_invalid(t *testing.T) {
	testName := "terraform_test_mtr_" + acctest.RandString(16)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetricTransformationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "sumologic_metric_transformation_rule" "test" {
  name                          = "%s"
  selector                      = "_sourceCategory=terraform_test_mtr"
  transformed_metrics_retention = 8
}
`, testName),
				ExpectError: regexp.MustCompile("transformed_metrics_retention must be 0"),
			},
			{
				Config: fmt.Sprintf(`
resource "sumologic_metric_transformation_rule" "test" {
  name      = "%s"
  selector  = "_sourceCategory=terraform_test_mtr"
  retention = 30
}
`, testName),
				ExpectError: regexp.MustCompile("expected retention to be one of"),
			},
		},
	})
}

func TestAccSumologicMetricTransformationRule_disappears(t *testing.T) {
	var rule MetricTransformationRule
	testName := "terraform_test_mtr_" + acctest.RandString(16)
	resourceName := "sumologic_metric_transformation_rule.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetricTransformationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMetricTransformationRuleBasic(testName, "_sourceCategory=terraform_test_mtr", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetricTransformationRuleExists(resourceName, &rule),
					// Delete out of band; the next plan must see it as gone and want to recreate it.
					func(s *terraform.State) error {
						client := testAccProvider.Meta().(*Client)
						return client.DeleteMetricTransformationRule(rule.ID)
					},
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckMetricTransformationRuleExists(name string, rule *MetricTransformationRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("metric transformation rule not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("metric transformation rule ID is not set")
		}

		client := testAccProvider.Meta().(*Client)
		found, err := client.GetMetricTransformationRule(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error reading metric transformation rule %s: %w", rs.Primary.ID, err)
		}
		if found == nil {
			return fmt.Errorf("metric transformation rule %s does not exist", rs.Primary.ID)
		}
		*rule = *found
		return nil
	}
}

func testAccCheckMetricTransformationRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "sumologic_metric_transformation_rule" {
			continue
		}
		found, err := client.GetMetricTransformationRule(r.Primary.ID)
		if err != nil {
			return fmt.Errorf("encountered an error: %w", err)
		}
		if found != nil {
			return fmt.Errorf("metric transformation rule %s still exists", r.Primary.ID)
		}
	}
	return nil
}

func testAccSumologicMetricTransformationRuleBasic(name, selector string, enabled bool) string {
	return fmt.Sprintf(`
resource "sumologic_metric_transformation_rule" "test" {
  name     = "%s"
  selector = "%s"
  enabled  = %t
}
`, name, selector, enabled)
}

func testAccSumologicMetricTransformationRuleWithTransformations(name, aggregateOn string) string {
	return fmt.Sprintf(`
resource "sumologic_metric_transformation_rule" "test" {
  name                          = "%s"
  selector                      = "_sourceCategory=terraform_test_mtr"
  retention                     = 0
  transformed_metrics_retention = 8

  dimension_transformation {
    transformation_type = "AggregateOnTransformation"
    aggregate_on        = [%s]
  }

  dimension_transformation {
    transformation_type  = "AddOrReplaceTransformation"
    dimension_to_replace = "metric"
    value                = "{{metric}}_aggregated"
  }
}
`, name, aggregateOn)
}
