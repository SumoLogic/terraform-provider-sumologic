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

func TestAccSumologicCSEOutlierRule_Override(t *testing.T) {
	SkipCseTest(t)

	var OutlierRule CSEOutlierRule
	descriptionExpression := "Observes granting of administrative privileges in Windows environments where the number of systems accessed by a single user exceeds standard deviation of what is expected for the user based on a historic baseline. The minumum floor of unique systems expected by default is set to 1."

	resourceName := "sumologic_cse_outlier_rule.sumo_outlier_rule_test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEOutlierRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config:                  testOverrideCSEOutlierRuleConfig(descriptionExpression),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateId:           "OUTLIER-S00007",
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"name"}, // Ignore fields that might differ
				ImportStatePersist:      true,
			},
			{
				Config: testOverrideCSEOutlierRuleConfig(fmt.Sprintf("Updated %s", descriptionExpression)),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEOutlierRuleExists(resourceName, &OutlierRule),
					testCheckOutlierRuleOverrideValues(&OutlierRule, fmt.Sprintf("Updated %s", descriptionExpression)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "id", "OUTLIER-S00007"),
				),
			},
			{
				Config: testOverrideCSEOutlierRuleConfig(descriptionExpression),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEOutlierRuleExists(resourceName, &OutlierRule),
					testCheckOutlierRuleOverrideValues(&OutlierRule, descriptionExpression),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "id", "OUTLIER-S00007"),
					removeState("sumologic_cse_outlier_rule.sumo_outlier_rule_test"),
				),
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

func testOverrideCSEOutlierRuleConfig(descriptionExpression string) string {
	return fmt.Sprintf(`
resource "sumologic_cse_outlier_rule" "sumo_outlier_rule_test" {
    baseline_window_size   = "432000000"
    description_expression = "%s"
    deviation_threshold    = 2
    enabled                = true
    floor_value            = 1
    group_by_fields        = [
        "user_username",
    ]
    is_prototype           = true
    match_expression       = <<-EOT
        metadata_vendor = 'Microsoft'
        AND metadata_product = 'Windows'
        AND metadata_deviceEventId = 'Security-4672'
        AND NOT user_username = 'system'
        AND NOT user_username RLIKE '(\$$)'
        AND NOT user_username RLIKE '(dwm\-)'
        AND NOT user_username RLIKE '(local service|network service)'
        AND NOT user_username RLIKE '(iusr)'AND NOT (
            LOWER(user_username) LIKE '%%svc%%'
        )
    EOT
    name                   = "Spike in Windows Administrative Privileges Granted for User"
    name_expression        = "Spike in Windows Administrative Privileges Granted for User: {{user_username}}"
    retention_window_size  = "7776000000"
    severity               = 2
    summary_expression     = "Outlier in distinct count of systems identified with Windows Administrative Privileges Granted for user: {{user_username}} based on hourly historic activity"
    tags                   = [
        "_mitreAttackTactic:TA0004",
        "_mitreAttackTechnique:T1078.002",
    ]
    window_size            = "T60M"

    aggregation_functions {
        arguments = [
            "device_hostname",
        ]
        function  = "count_distinct"
        name      = "current"
    }

    entity_selectors {
        entity_type = "_username"
        expression  = "user_username"
    }
}
`, descriptionExpression)
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

func testCheckOutlierRuleOverrideValues(outlierRule *CSEOutlierRule, descriptionExpression string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if outlierRule.DescriptionExpression != descriptionExpression {
			return fmt.Errorf("bad descriptionExpression, expected \"%s\", got %#v", descriptionExpression, outlierRule.DescriptionExpression)
		}
		return nil
	}
}
