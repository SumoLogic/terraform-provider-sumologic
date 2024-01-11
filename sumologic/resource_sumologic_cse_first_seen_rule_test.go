package sumologic

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"text/template"
)

func TestAccSumologicCSEFirstSeenRule_createAndUpdate(t *testing.T) {
	SkipCseTest(t)

	var payload = CSEFirstSeenRule{
		BaselineType:          "PER_ENTITY",
		BaselineWindowSize:    "35000",
		DescriptionExpression: "FirstSeenRuleTerraformTest - {{ user_username }}",
		Enabled:               true,
		EntitySelectors: []EntitySelector{
			{EntityType: "_username", Expression: "user_username"},
			{EntityType: "_hostname", Expression: "dstDevice_hostname"},
		},
		FilterExpression:    `objectType="Network"`,
		GroupByFields:       []string{"user_username"},
		IsPrototype:         false,
		Name:                "FirstSeenRuleTerraformTest",
		NameExpression:      "FirstSeenRuleTerraformTest - {{ user_username }}",
		RetentionWindowSize: "86400000",
		Severity:            1,
		ValueFields:         []string{"dstDevice_hostname"},
		Version:             1,
	}
	updatedPayload := payload
	updatedPayload.Enabled = false

	var result CSEFirstSeenRule

	resourceName := "sumologic_cse_first_seen_rule.first_seen_rule"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEFirstSeenRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSEFirstSeenRuleConfig(t, &payload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEFirstSeenRuleExists(resourceName, &result),
					testCheckFirstSeenRuleValues(t, &payload, &result),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testCreateCSEFirstSeenRuleConfig(t, &updatedPayload),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEFirstSeenRuleExists(resourceName, &result),
					testCheckFirstSeenRuleValues(t, &updatedPayload, &result),
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

func testAccCSEFirstSeenRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_first_seen_rule" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE FirstSeen Rule destruction check: CSE FirstSeen Rule ID is not set")
		}

		s, err := client.GetCSEFirstSeenRule(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("FirstSeen rule still exists")
		}
	}
	return nil
}

func testCreateCSEFirstSeenRuleConfig(t *testing.T, payload *CSEFirstSeenRule) string {
	configTemplate := template.Must(template.New("first_seen_rule").Funcs(template.FuncMap{
		"quoteStringArray": func(arr []string) string {
			return `["` + strings.Join(arr, `","`) + `"]`
		}}).Parse(`
resource "sumologic_cse_first_seen_rule" "first_seen_rule" {
  baseline_type          = "{{ .BaselineType }}"
  baseline_window_size   = "{{ .BaselineWindowSize }}"
  description_expression = "{{ .DescriptionExpression }}"
  enabled                = {{ .Enabled }}
  {{ range .EntitySelectors }}
  entity_selectors {
        entity_type = "{{ .EntityType }}"
        expression = "{{ .Expression }}"
  }
  {{ end }}
  filter_expression     = "{{ js .FilterExpression }}"
  group_by_fields       = {{ quoteStringArray .GroupByFields }}
  is_prototype          = {{ .IsPrototype }}
  name                  = "{{ .Name }}"
  name_expression       = "{{ .NameExpression }}"
  retention_window_size = "{{ .RetentionWindowSize }}"
  severity              = {{ .Severity }}
  value_fields          = {{ quoteStringArray .ValueFields }}
}
`))
	var buffer bytes.Buffer
	if err := configTemplate.Execute(&buffer, payload); err != nil {
		t.Error(err)
	}
	return buffer.String()
}

func testCheckCSEFirstSeenRuleExists(n string, firstSeenRule *CSEFirstSeenRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("firstSeen rule ID is not set")
		}

		c := testAccProvider.Meta().(*Client)
		firstSeenRuleResp, err := c.GetCSEFirstSeenRule(rs.Primary.ID)
		if err != nil {
			return err
		}

		*firstSeenRule = *firstSeenRuleResp

		return nil
	}
}

func testCheckFirstSeenRuleValues(t *testing.T, expected *CSEFirstSeenRule, actual *CSEFirstSeenRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		assert.Equal(t, expected.BaselineType, actual.BaselineType)
		assert.Equal(t, expected.BaselineWindowSize, actual.BaselineWindowSize)
		assert.Equal(t, expected.DescriptionExpression, actual.DescriptionExpression)
		assert.Equal(t, expected.Enabled, actual.Enabled)
		assert.Equal(t, expected.EntitySelectors, actual.EntitySelectors)
		assert.Equal(t, expected.FilterExpression, actual.FilterExpression)
		assert.Equal(t, expected.GroupByFields, actual.GroupByFields)
		assert.Equal(t, expected.IsPrototype, actual.IsPrototype)
		assert.Equal(t, expected.Name, actual.Name)
		assert.Equal(t, expected.NameExpression, actual.NameExpression)
		assert.Equal(t, expected.RetentionWindowSize, actual.RetentionWindowSize)
		assert.Equal(t, expected.Severity, actual.Severity)
		assert.Equal(t, expected.ValueFields, actual.ValueFields)

		return nil
	}
}
