package sumologic

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccSumologicCSEOutlierRule_createAndUpdate(t *testing.T) {
	SkipCseTest(t)

	var payload = CSEOutlierRule{
		AggregationFunctions:  []AggregationFunction{{Name: "current", Function: "count", Arguments: []string{"true"}}},
		BaselineWindowSize:    "604800000",
		DescriptionExpression: "OutlierRuleTerraformTest - {{ user_username }}",
		Enabled:               true,
		EntitySelectors: []EntitySelector{
			{EntityType: "_username", Expression: "user_username"},
		},
		FloorValue:            0,
		DeviationThreshold:    3,
		GroupByFields:         []string{"user_username"},
		IsPrototype:           false,
		MatchExpression:       `objectType="Network"`,
		Name:                  fmt.Sprintf("OutlierRuleTerraformTest %s", uuid.New()),
		NameExpression:        "OutlierRuleTerraformTest - {{ user_username }}",
		RetentionWindowSize:   "1209600000",
		Severity:              1,
		SummaryExpression:     "OutlierRuleTerraformTest - {{ user_username }}",
		Tags:                  []string{"OutlierRuleTerraformTest"},
		WindowSize:            "T24H",
		SuppressionWindowSize: nil,
	}
	updatedPayload := payload
	updatedPayload.Enabled = false
	suppressionWindow := 25 * 60 * 60 * 1000
	updatedPayload.SuppressionWindowSize = &suppressionWindow

	var result CSEOutlierRule

	resourceName := "sumologic_cse_outlier_rule.outlier_rule"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEOutlierRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSEOutlierRuleConfig(t, &payload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEOutlierRuleExists(resourceName, &result),
					testCheckOutlierRuleValues(t, &payload, &result),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testCreateCSEOutlierRuleConfig(t, &updatedPayload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEOutlierRuleExists(resourceName, &result),
					testCheckOutlierRuleValues(t, &updatedPayload, &result),
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

func testAccCSEOutlierRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_outlier_rule" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Outlier Rule destruction check: CSE Outlier Rule ID is not set")
		}

		s, err := client.GetCSEOutlierRule(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("Outlier rule still exists")
		}
	}
	return nil
}

func testCreateCSEOutlierRuleConfig(t *testing.T, payload *CSEOutlierRule) string {
	configTemplate := template.Must(template.New("outlier_rule").Funcs(template.FuncMap{
		"quoteStringArray": func(arr []string) string {
			return `["` + strings.Join(arr, `","`) + `"]`
		}}).Parse(`
resource "sumologic_cse_outlier_rule" "outlier_rule" {
  {{ range .AggregationFunctions }}
  aggregation_functions {
  		name = "{{ .Name }}"
  		function = "{{ .Function }}"
  		arguments = {{ quoteStringArray .Arguments }}
  }
  {{ end }}
  baseline_window_size   = "{{ .BaselineWindowSize }}"
  description_expression = "{{ .DescriptionExpression }}"
  enabled                = {{ .Enabled }}
  {{ range .EntitySelectors }}
  entity_selectors {
        entity_type = "{{ .EntityType }}"
        expression = "{{ .Expression }}"
  }
  {{ end }}
  floor_value 		 	= {{ .FloorValue }}
  deviation_threshold 	= {{ .DeviationThreshold }}
  group_by_fields       = {{ quoteStringArray .GroupByFields }}
  is_prototype		  	= {{ .IsPrototype }}
  match_expression     	= "{{ js .MatchExpression }}"
  name                  = "{{ .Name }}"
  name_expression       = "{{ .NameExpression }}"
  retention_window_size = "{{ .RetentionWindowSize }}"
  severity              = {{ .Severity }}
  summary_expression	= "{{ .SummaryExpression }}"
  tags                  = {{ quoteStringArray .Tags }}
  window_size 		 	= "{{ .WindowSize }}"
	{{ if .SuppressionWindowSize }}
	suppression_window_size = {{ .SuppressionWindowSize }}
	{{ end }}
}
`))
	var buffer bytes.Buffer
	if err := configTemplate.Execute(&buffer, payload); err != nil {
		t.Error(err)
	}

	return buffer.String()
}

func testCheckCSEOutlierRuleExists(n string, outlierRule *CSEOutlierRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Outlier rule ID is not set")
		}

		c := testAccProvider.Meta().(*Client)
		outlierRuleResp, err := c.GetCSEOutlierRule(rs.Primary.ID)
		if err != nil {
			return err
		}

		*outlierRule = *outlierRuleResp

		return nil
	}
}

func testCheckOutlierRuleValues(t *testing.T, expected *CSEOutlierRule, actual *CSEOutlierRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		assert.Equal(t, expected.AggregationFunctions, actual.AggregationFunctions)
		assert.Equal(t, expected.BaselineWindowSize, actual.BaselineWindowSize)
		assert.Equal(t, expected.DescriptionExpression, actual.DescriptionExpression)
		assert.Equal(t, expected.DeviationThreshold, actual.DeviationThreshold)
		assert.Equal(t, expected.Enabled, actual.Enabled)
		assert.Equal(t, expected.EntitySelectors, actual.EntitySelectors)
		assert.Equal(t, expected.FloorValue, actual.FloorValue)
		assert.Equal(t, expected.GroupByFields, actual.GroupByFields)
		assert.Equal(t, expected.IsPrototype, actual.IsPrototype)
		assert.Equal(t, expected.MatchExpression, actual.MatchExpression)
		assert.Equal(t, expected.Name, actual.Name)
		assert.Equal(t, expected.NameExpression, actual.NameExpression)
		assert.Equal(t, expected.RetentionWindowSize, actual.RetentionWindowSize)
		assert.Equal(t, expected.Severity, actual.Severity)
		assert.Equal(t, expected.SummaryExpression, actual.SummaryExpression)
		assert.Equal(t, expected.Tags, actual.Tags)
		assert.Equal(t, expected.WindowSize, windowSizeField(actual.WindowSizeName))
		assert.Equal(t, expected.SuppressionWindowSize, actual.SuppressionWindowSize)

		return nil
	}
}
