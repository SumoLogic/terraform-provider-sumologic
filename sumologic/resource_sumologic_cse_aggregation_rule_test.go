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

func TestAccSumologicCSEAggregationRule_createAndUpdateWithCustomWindowSize(t *testing.T) {
	SkipCseTest(t)

	payload := getCSEAggregationRuleTestPayload()
	payload.WindowSize = "CUSTOM"
	payload.WindowSizeMilliseconds = "10800000" // 3h

	updatedPayload := payload
	updatedPayload.WindowSizeMilliseconds = "14400000" // 4h
	updatedSuppressionWindow := 5 * 60 * 60 * 1000
	updatedPayload.SuppressionWindowSize = &updatedSuppressionWindow

	var aggregationRule CSEAggregationRule
	resourceName := "sumologic_cse_aggregation_rule.aggregation_rule"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEAggregationRuleDestroy,
		Steps: []resource.TestStep{
			{ // create
				Config: testCreateCSEAggregationRuleConfig(t, &payload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEAggregationRuleExists(resourceName, &aggregationRule),
					testCheckCSEAggregationRuleValues(t, &payload, &aggregationRule),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{ // update
				Config: testCreateCSEAggregationRuleConfig(t, &updatedPayload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEAggregationRuleExists(resourceName, &aggregationRule),
					testCheckCSEAggregationRuleValues(t, &updatedPayload, &aggregationRule),
				),
			},
			{ // import
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSumologicCSEAggregationRule_createAndUpdateToCustomWindowSize(t *testing.T) {
	SkipCseTest(t)

	payload := getCSEAggregationRuleTestPayload()
	payload.WindowSize = "T30M"
	payload.WindowSizeMilliseconds = "irrelevant"
	suppressionWindow := 35 * 60 * 1000
	payload.SuppressionWindowSize = &suppressionWindow

	updatedPayload := payload
	updatedPayload.WindowSize = "CUSTOM"
	updatedPayload.WindowSizeMilliseconds = "14400000" // 4h
	updatedSuppressionWindow := 5 * 60 * 60 * 1000
	updatedPayload.SuppressionWindowSize = &updatedSuppressionWindow

	var aggregationRule CSEAggregationRule
	resourceName := "sumologic_cse_aggregation_rule.aggregation_rule"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEAggregationRuleDestroy,
		Steps: []resource.TestStep{
			{ // create
				Config: testCreateCSEAggregationRuleConfig(t, &payload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEAggregationRuleExists(resourceName, &aggregationRule),
					testCheckCSEAggregationRuleValues(t, &payload, &aggregationRule),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{ // update
				Config: testCreateCSEAggregationRuleConfig(t, &updatedPayload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEAggregationRuleExists(resourceName, &aggregationRule),
					testCheckCSEAggregationRuleValues(t, &updatedPayload, &aggregationRule),
				),
			},
			{ // import
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSumologicCSEAggregationRule_createAndUpdate(t *testing.T) {
	SkipCseTest(t)

	payload := getCSEAggregationRuleTestPayload()
	payload.WindowSize = "T30M"
	payload.WindowSizeMilliseconds = "irrelevant"
	suppressionWindow := 35 * 60 * 1000
	payload.SuppressionWindowSize = &suppressionWindow

	updatedPayload := payload
	updatedPayload.Name = fmt.Sprintf("Updated Aggregation Rule %s", uuid.New())
	updatedPayload.WindowSize = "T12H"
	updatedSuppressionWindow := 13 * 60 * 60 * 1000
	updatedPayload.SuppressionWindowSize = &updatedSuppressionWindow

	var aggregationRule CSEAggregationRule
	resourceName := "sumologic_cse_aggregation_rule.aggregation_rule"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEAggregationRuleDestroy,
		Steps: []resource.TestStep{
			{ // create
				Config: testCreateCSEAggregationRuleConfig(t, &payload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEAggregationRuleExists(resourceName, &aggregationRule),
					testCheckCSEAggregationRuleValues(t, &payload, &aggregationRule),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{ // update
				Config: testCreateCSEAggregationRuleConfig(t, &updatedPayload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEAggregationRuleExists(resourceName, &aggregationRule),
					testCheckCSEAggregationRuleValues(t, &updatedPayload, &aggregationRule),
				),
			},
			{ // import
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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

func testCreateCSEAggregationRuleConfig(t *testing.T, payload *CSEAggregationRule) string {
	resourceTemplate := `
		resource "sumologic_cse_aggregation_rule" "aggregation_rule" {
			{{ range .AggregationFunctions }}
			aggregation_functions {
					name = "{{ .Name }}"
					function = "{{ .Function }}"
					arguments = {{ quoteStringArray .Arguments }}
			}
			{{ end }}
			description_expression = "{{ .DescriptionExpression }}"
			enabled = {{ .Enabled }}
			{{ range .EntitySelectors }}
			entity_selectors {
				entity_type = "{{ .EntityType }}"
				expression = "{{ .Expression }}"
			}
			{{ end }}
			group_by_entity = {{ .GroupByEntity }}
			group_by_fields = {{ quoteStringArray .GroupByFields }}
			is_prototype = {{ .IsPrototype }}
			match_expression = "{{ js .MatchExpression }}"
			name = "{{ .Name }}"
			name_expression = "{{ .NameExpression }}"
			severity_mapping {
				type = "{{ .SeverityMapping.Type }}"
				default = {{ .SeverityMapping.Default }}
			}
			summary_expression = "{{ .SummaryExpression }}"
			trigger_expression = "{{ js .TriggerExpression }}"
			tags = {{ quoteStringArray .Tags }}
			window_size = "{{ .WindowSize }}"
			{{ if eq .WindowSize "CUSTOM" }}
			window_size_millis = "{{ .WindowSizeMilliseconds }}"
			{{ end }}
			{{ if .SuppressionWindowSize }}
			suppression_window_size = {{ .SuppressionWindowSize }}
			{{ end }}
		}
	`

	configTemplate := template.Must(template.New("aggregation_rule").Funcs(template.FuncMap{
		"quoteStringArray": func(arr []string) string {
			return `["` + strings.Join(arr, `","`) + `"]`
		},
		"js": func(in string) string {
			escaped := strings.Replace(in, `"`, `\"`, -1)
			escaped = strings.Replace(escaped, `$`, `$$`, -1) // Escape Terraform interpolation
			return escaped
		},
	}).Parse(resourceTemplate))

	var buffer bytes.Buffer
	if err := configTemplate.Execute(&buffer, payload); err != nil {
		t.Error(err)
	}

	return buffer.String()
}

func getCSEAggregationRuleTestPayload() CSEAggregationRule {
	return CSEAggregationRule{
		AggregationFunctions:  []AggregationFunction{{Name: "distinct_eventid_count", Function: "count_distinct", Arguments: []string{"metadata_deviceEventId"}}},
		DescriptionExpression: "Test description",
		Enabled:               true,
		EntitySelectors:       []EntitySelector{{EntityType: "_ip", Expression: "srcDevice_ip"}},
		GroupByEntity:         true,
		GroupByFields:         []string{"dstDevice_hostname"},
		IsPrototype:           false,
		MatchExpression:       "foo = bar",
		Name:                  fmt.Sprintf("Test Aggregation Rule %s", uuid.New()),
		NameExpression:        "Signal Name",
		SeverityMapping: SeverityMapping{
			Type:    "constant",
			Default: 5,
		},
		SummaryExpression:      "Signal Summary",
		TriggerExpression:      "foo = bar",
		Tags:                   []string{"foo"},
		WindowSize:             windowSizeField("CUSTOM"),
		WindowSizeMilliseconds: "10800000",
		SuppressionWindowSize:  nil,
	}
}

func testCheckCSEAggregationRuleExists(n string, AggregationRule *CSEAggregationRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("aggregation rule ID is not set")
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

func testCheckCSEAggregationRuleValues(t *testing.T, expected *CSEAggregationRule, actual *CSEAggregationRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		assert.Equal(t, expected.AggregationFunctions, actual.AggregationFunctions)
		assert.Equal(t, expected.DescriptionExpression, actual.DescriptionExpression)
		assert.Equal(t, expected.Enabled, actual.Enabled)
		assert.Equal(t, expected.EntitySelectors, actual.EntitySelectors)
		assert.Equal(t, expected.GroupByEntity, actual.GroupByEntity)
		assert.Equal(t, expected.GroupByFields, actual.GroupByFields)
		assert.Equal(t, expected.IsPrototype, actual.IsPrototype)
		assert.Equal(t, expected.MatchExpression, actual.MatchExpression)
		assert.Equal(t, expected.Name, actual.Name)
		assert.Equal(t, expected.NameExpression, actual.NameExpression)
		assert.Equal(t, expected.SeverityMapping, actual.SeverityMapping)
		assert.Equal(t, expected.SummaryExpression, actual.SummaryExpression)
		assert.Equal(t, expected.TriggerExpression, actual.TriggerExpression)
		assert.Equal(t, expected.Tags, actual.Tags)
		assert.Equal(t, string(expected.WindowSize), actual.WindowSizeName)
		if strings.EqualFold(actual.WindowSizeName, "CUSTOM") {
			assert.Equal(t, expected.WindowSizeMilliseconds, string(actual.WindowSize))
		}
		assert.Equal(t, expected.SuppressionWindowSize, actual.SuppressionWindowSize)
		return nil
	}
}
