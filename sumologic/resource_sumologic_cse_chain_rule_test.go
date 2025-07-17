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

func TestAccSumologicCSEChainRule_createAndUpdateWithCustomWindowSize(t *testing.T) {
	SkipCseTest(t)

	payload := getCSEChainRuleTestPayload()
	payload.WindowSize = "CUSTOM"
	payload.WindowSizeMilliseconds = "10800000" // 3h

	updatedPayload := payload
	updatedPayload.WindowSizeMilliseconds = "14400000" // 4h
	updatedSuppressionWindow := 15000000
	updatedPayload.SuppressionWindowSize = &updatedSuppressionWindow

	var chainRule CSEChainRule
	resourceName := "sumologic_cse_chain_rule.chain_rule"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEChainRuleDestroy,
		Steps: []resource.TestStep{
			{ // create
				Config: testCreateCSEChainRuleConfig(t, &payload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEChainRuleExists(resourceName, &chainRule),
					testCheckCSEChainRuleValues(t, &payload, &chainRule),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{ // update
				Config: testCreateCSEChainRuleConfig(t, &updatedPayload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEChainRuleExists(resourceName, &chainRule),
					testCheckCSEChainRuleValues(t, &updatedPayload, &chainRule),
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

func TestAccSumologicCSEChainRule_createAndUpdateToCustomWindowSize(t *testing.T) {
	SkipCseTest(t)

	payload := getCSEChainRuleTestPayload()
	payload.WindowSize = "T30M"
	payload.WindowSizeMilliseconds = "irrelevant"
	suppressionWindow := 2000000
	payload.SuppressionWindowSize = &suppressionWindow

	updatedPayload := payload
	updatedPayload.WindowSize = "CUSTOM"
	updatedPayload.WindowSizeMilliseconds = "14400000" // 4h
	updatedSuppressionWindow := 20000000
	updatedPayload.SuppressionWindowSize = &updatedSuppressionWindow

	var chainRule CSEChainRule
	resourceName := "sumologic_cse_chain_rule.chain_rule"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEChainRuleDestroy,
		Steps: []resource.TestStep{
			{ // create
				Config: testCreateCSEChainRuleConfig(t, &payload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEChainRuleExists(resourceName, &chainRule),
					testCheckCSEChainRuleValues(t, &payload, &chainRule),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{ // update
				Config: testCreateCSEChainRuleConfig(t, &updatedPayload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEChainRuleExists(resourceName, &chainRule),
					testCheckCSEChainRuleValues(t, &updatedPayload, &chainRule),
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

func TestAccSumologicCSEChainRule_createAndUpdate(t *testing.T) {
	SkipCseTest(t)

	payload := getCSEChainRuleTestPayload()
	payload.WindowSize = "T30M"
	payload.WindowSizeMilliseconds = "irrelevant"
	suppressionWindow := 35 * 60 * 1000
	payload.SuppressionWindowSize = &suppressionWindow

	updatedPayload := payload
	updatedPayload.Name = fmt.Sprintf("Updated Chain Rule %s", uuid.New())
	updatedPayload.WindowSize = "T12H"
	updatedSuppressionWindow := 13 * 60 * 60 * 1000
	updatedPayload.SuppressionWindowSize = &updatedSuppressionWindow

	var chainRule CSEChainRule
	resourceName := "sumologic_cse_chain_rule.chain_rule"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEChainRuleDestroy,
		Steps: []resource.TestStep{
			{ // create
				Config: testCreateCSEChainRuleConfig(t, &payload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEChainRuleExists(resourceName, &chainRule),
					testCheckCSEChainRuleValues(t, &payload, &chainRule),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{ // update
				Config: testCreateCSEChainRuleConfig(t, &updatedPayload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEChainRuleExists(resourceName, &chainRule),
					testCheckCSEChainRuleValues(t, &updatedPayload, &chainRule),
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

func testCreateCSEChainRuleConfig(t *testing.T, payload *CSEChainRule) string {
	resourceTemplate := `
		resource "sumologic_cse_chain_rule" "chain_rule" {
			description = "{{ .Description }}"
			enabled = {{ .Enabled }}
			{{ range .EntitySelectors }}
			entity_selectors {
				entity_type = "{{ .EntityType }}"
				expression = "{{ .Expression }}"
			}
			{{ end }}
			{{ range .ExpressionsAndLimits }}
			expressions_and_limits {
				expression = "{{ .Expression }}"
				limit = {{ .Limit }}
			}
			{{ end }}
			group_by_fields = {{ quoteStringArray .GroupByFields }}
			is_prototype = {{ .IsPrototype }}
			ordered = {{ .Ordered }}
			name = "{{ .Name }}"
			severity = {{ .Severity }}
			summary_expression = "{{ .SummaryExpression }}"
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

	configTemplate := template.Must(template.New("chain_rule").Funcs(template.FuncMap{
		"quoteStringArray": func(arr []string) string {
			return `["` + strings.Join(arr, `","`) + `"]`
		},
	}).Parse(resourceTemplate))

	var buffer bytes.Buffer
	if err := configTemplate.Execute(&buffer, payload); err != nil {
		t.Error(err)
	}

	return buffer.String()
}

func getCSEChainRuleTestPayload() CSEChainRule {
	return CSEChainRule{
		Description:            "Test description",
		Enabled:                true,
		EntitySelectors:        []EntitySelector{{EntityType: "_ip", Expression: "srcDevice_ip"}},
		ExpressionsAndLimits:   []ExpressionAndLimit{{Expression: "foo = bar", Limit: 5}, {Expression: "baz = qux", Limit: 1}},
		GroupByFields:          []string{"destPort"},
		IsPrototype:            false,
		Ordered:                true,
		Name:                   fmt.Sprintf("Test Chain Rule %s", uuid.New()),
		Severity:               5,
		SummaryExpression:      "Signal Summary",
		Tags:                   []string{"foo"},
		WindowSize:             windowSizeField("CUSTOM"),
		WindowSizeMilliseconds: "10800000",
		SuppressionWindowSize:  nil,
	}
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

func testCheckCSEChainRuleValues(t *testing.T, expected *CSEChainRule, actual *CSEChainRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		assert.Equal(t, expected.Description, actual.Description)
		assert.Equal(t, expected.Enabled, actual.Enabled)
		assert.Equal(t, expected.EntitySelectors, actual.EntitySelectors)
		assert.Equal(t, expected.ExpressionsAndLimits, actual.ExpressionsAndLimits)
		assert.Equal(t, expected.GroupByFields, actual.GroupByFields)
		assert.Equal(t, expected.IsPrototype, actual.IsPrototype)
		assert.Equal(t, expected.Ordered, actual.Ordered)
		assert.Equal(t, expected.Name, actual.Name)
		assert.Equal(t, expected.Severity, actual.Severity)
		assert.Equal(t, expected.SummaryExpression, actual.SummaryExpression)
		assert.Equal(t, expected.Tags, actual.Tags)
		assert.Equal(t, string(expected.WindowSize), actual.WindowSizeName)
		if strings.EqualFold(actual.WindowSizeName, "CUSTOM") {
			assert.Equal(t, expected.WindowSizeMilliseconds, string(actual.WindowSize))
		}
		assert.Equal(t, expected.SuppressionWindowSize, actual.SuppressionWindowSize)
		return nil
	}
}
