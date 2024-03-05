package sumologic

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccSumologicCSEThresholdRule_createAndUpdateWithCustomWindowSize(t *testing.T) {
	SkipCseTest(t)

	payload := getCSEThresholdRuleTestPayload()
	payload.WindowSize = "CUSTOM"
	payload.WindowSizeMilliseconds = "10800000" //3h

	updatedPayload := payload
	updatedPayload.WindowSizeMilliseconds = "14400000" //4h

	var thresholdRule CSEThresholdRule
	resourceName := "sumologic_cse_threshold_rule.threshold_rule"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEThresholdRuleDestroy,
		Steps: []resource.TestStep{
			{ // create
				Config: testCreateCSEThresholdRuleConfigWithCustomWindowSize(t, &payload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEThresholdRuleExists(resourceName, &thresholdRule),
					testCheckCSEThresholdRuleValues(t, &payload, &thresholdRule),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{ // update
				Config: testCreateCSEThresholdRuleConfigWithCustomWindowSize(t, &updatedPayload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEThresholdRuleExists(resourceName, &thresholdRule),
					testCheckCSEThresholdRuleValues(t, &updatedPayload, &thresholdRule),
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
func TestAccSumologicCSEThresholdRule_createAndUpdateToCustomWindowSize(t *testing.T) {
	SkipCseTest(t)

	payload := getCSEThresholdRuleTestPayload()
	payload.WindowSize = "T30M"
	payload.WindowSizeMilliseconds = "irrelevant"

	updatedPayload := payload
	updatedPayload.WindowSize = "CUSTOM"
	updatedPayload.WindowSizeMilliseconds = "14400000" //4h

	var thresholdRule CSEThresholdRule
	resourceName := "sumologic_cse_threshold_rule.threshold_rule"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEThresholdRuleDestroy,
		Steps: []resource.TestStep{
			{ // create
				Config: testCreateCSEThresholdRuleConfigWithCustomWindowSize(t, &payload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEThresholdRuleExists(resourceName, &thresholdRule),
					testCheckCSEThresholdRuleValues(t, &payload, &thresholdRule),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{ // update
				Config: testCreateCSEThresholdRuleConfigWithCustomWindowSize(t, &updatedPayload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEThresholdRuleExists(resourceName, &thresholdRule),
					testCheckCSEThresholdRuleValues(t, &updatedPayload, &thresholdRule),
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

func TestAccSumologicCSEThresholdRule_createAndUpdate(t *testing.T) {
	SkipCseTest(t)

	payload := getCSEThresholdRuleTestPayload()
	payload.WindowSizeName = "T30M"
	payload.WindowSizeMilliseconds = "irrelevant"

	updatedPayload := payload
	updatedPayload.Name = fmt.Sprintf("Updated Threshold Rule %s", uuid.New())
	updatedPayload.WindowSizeName = "T12H"

	var thresholdRule CSEThresholdRule
	resourceName := "sumologic_cse_threshold_rule.threshold_rule"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEThresholdRuleDestroy,
		Steps: []resource.TestStep{
			{ // create
				Config: testCreateCSEThresholdRuleConfigWithCustomWindowSize(t, &payload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEThresholdRuleExists(resourceName, &thresholdRule),
					testCheckCSEThresholdRuleValues(t, &payload, &thresholdRule),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{ // update
				Config: testCreateCSEThresholdRuleConfigWithCustomWindowSize(t, &updatedPayload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEThresholdRuleExists(resourceName, &thresholdRule),
					testCheckCSEThresholdRuleValues(t, &updatedPayload, &thresholdRule),
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

func testCreateCSEThresholdRuleConfigWithCustomWindowSize(t *testing.T, payload *CSEThresholdRule) string {
	resourceTemlpate := `
		resource "sumologic_cse_threshold_rule" "threshold_rule" {
			count_distinct = {{ .CountDistinct }}
			count_field = "{{ .CountField }}"
			description = "{{ .Description }}"
			enabled = {{ .Enabled }}
			{{ range .EntitySelectors }}
			entity_selectors {
				entity_type = "{{ .EntityType }}"
				expression = "{{ .Expression }}"
			}
			{{ end }}
			expression = "{{ js .Expression }}"
			group_by_fields = {{ quoteStringArray .GroupByFields }}
			is_prototype = {{ .IsPrototype }}
			limit = {{ .Limit }}
			name = "{{ .Name }}"
			severity = {{ .Severity }}
			summary_expression = "{{ .SummaryExpression }}"
			tags = {{ quoteStringArray .Tags }}
			window_size = "{{ .WindowSize }}"
			{{ if eq .WindowSize "CUSTOM" }}
			window_size_millis = "{{ .WindowSizeMilliseconds }}"
			{{ end }}
		}
		`

	configTemplate := template.Must(template.New("threshold_rule").Funcs(template.FuncMap{
		"quoteStringArray": func(arr []string) string {
			return `["` + strings.Join(arr, `","`) + `"]`
		},
	}).Parse(resourceTemlpate))

	var buffer bytes.Buffer
	if err := configTemplate.Execute(&buffer, payload); err != nil {
		t.Error(err)
	}

	return buffer.String()
}

func getCSEThresholdRuleTestPayload() CSEThresholdRule {
	return CSEThresholdRule{
		CountDistinct:          true,
		CountField:             "dstDevice_hostname",
		Description:            "Test description",
		Enabled:                true,
		EntitySelectors:        []EntitySelector{{EntityType: "_ip", Expression: "srcDevice_ip"}},
		Expression:             "foo = bar",
		GroupByFields:          []string{"destPort"},
		IsPrototype:            false,
		Limit:                  20,
		Name:                   fmt.Sprintf("Test Threshold Rule With Custom WindowSize %s", uuid.New()),
		Severity:               5,
		SummaryExpression:      "Signal Summary",
		Tags:                   []string{"foo"},
		WindowSize:             windowSizeField("CUSTOM"),
		WindowSizeMilliseconds: "10800000",
	}
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

func testCheckCSEThresholdRuleValues(t *testing.T, expected *CSEThresholdRule, actual *CSEThresholdRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		assert.Equal(t, expected.CountDistinct, actual.CountDistinct)
		assert.Equal(t, expected.CountField, actual.CountField)
		assert.Equal(t, expected.Description, actual.Description)
		assert.Equal(t, expected.Enabled, actual.Enabled)
		assert.Equal(t, expected.EntitySelectors, actual.EntitySelectors)
		assert.Equal(t, expected.Expression, actual.Expression)
		assert.Equal(t, expected.GroupByFields, actual.GroupByFields)
		assert.Equal(t, expected.IsPrototype, actual.IsPrototype)
		assert.Equal(t, expected.Limit, actual.Limit)
		assert.Equal(t, expected.Name, actual.Name)
		assert.Equal(t, expected.Severity, actual.Severity)
		assert.Equal(t, expected.SummaryExpression, actual.SummaryExpression)
		assert.Equal(t, expected.Tags, actual.Tags)
		assert.Equal(t, string(expected.WindowSize), actual.WindowSizeName)
		if strings.EqualFold(actual.WindowSizeName, "CUSTOM") {
			assert.Equal(t, expected.WindowSizeMilliseconds, string(actual.WindowSize))
		}
		return nil
	}
}
