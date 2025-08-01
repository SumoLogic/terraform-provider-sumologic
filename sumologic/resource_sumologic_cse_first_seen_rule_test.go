package sumologic

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
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
		FilterExpression:      `objectType="Network"`,
		GroupByFields:         []string{"user_username"},
		IsPrototype:           false,
		Name:                  acctest.RandomWithPrefix("FirstSeenRuleTerraformTest"),
		NameExpression:        "FirstSeenRuleTerraformTest - {{ user_username }}",
		RetentionWindowSize:   "86400000",
		Severity:              1,
		ValueFields:           []string{"dstDevice_hostname"},
		Version:               1,
		SuppressionWindowSize: nil,
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

func TestAccSumologicCSEFirstSeenRule_Override(t *testing.T) {
	SkipCseTest(t)

	var FirstSeenRule CSEFirstSeenRule
	descriptionExpression := "Observes for a user performing a RDP logon for the first time."

	resourceName := "sumologic_cse_first_seen_rule.sumo_first_seen_rule_test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEFirstSeenRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config:                  testOverrideCSEFirstSeenRuleConfig(descriptionExpression),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateId:           "FIRST-S00009",
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"name"}, // Ignore fields that might differ
				ImportStatePersist:      true,
			},
			{
				Config: testOverrideCSEFirstSeenRuleConfig(fmt.Sprintf("Updated %s", descriptionExpression)),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEFirstSeenRuleExists(resourceName, &FirstSeenRule),
					testCheckFirstSeenRuleOverrideValues(&FirstSeenRule, fmt.Sprintf("Updated %s", descriptionExpression)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "id", "FIRST-S00009"),
				),
			},
			{
				Config: testOverrideCSEFirstSeenRuleConfig(fmt.Sprintf(descriptionExpression)),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEFirstSeenRuleExists(resourceName, &FirstSeenRule),
					testCheckFirstSeenRuleOverrideValues(&FirstSeenRule, fmt.Sprintf(descriptionExpression)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "id", "FIRST-S00009"),
				),
			},
			{
				Config: getFirstSeenRuleRemovedBlock(),
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

func testOverrideCSEFirstSeenRuleConfig(descriptionExpression string) string {
	return fmt.Sprintf(`
resource "sumologic_cse_first_seen_rule" "sumo_first_seen_rule_test" {
    baseline_type          = "GLOBAL"
    baseline_window_size   = "1814400000"
    description_expression = "%s"
    enabled                = true
    filter_expression      = <<-EOT
        metadata_vendor = 'Microsoft'
        AND metadata_product = 'Windows'
        AND metadata_deviceEventId = 'Security-4624'
        AND fields["EventData.LogonType"] = "10"
        AND user_username != "local service"
    EOT
    group_by_fields        = []
    is_prototype           = true
    name                   = "First Seen RDP Logon From User"
    name_expression        = "First Seen RDP Logon From User"
    retention_window_size  = "7776000000"
    severity               = 2
    summary_expression     = "First Seen RDP logon by the user: {{user_username}}"
    tags                   = [
        "_mitreAttackTactic:TA0008",
        "_mitreAttackTechnique:T1021",
        "_mitreAttackTechnique:T1021.001",
    ]
    value_fields           = []

    entity_selectors {
        entity_type = "_username"
        expression  = "user_username"
    }
}
`, descriptionExpression)
}

func getFirstSeenRuleRemovedBlock() string {
	return fmt.Sprintf(`
removed {
  from = sumologic_cse_first_seen_rule.sumo_first_seen_rule_test

  lifecycle {
	destroy = false
  }
}
`)
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
		assert.Equal(t, expected.SuppressionWindowSize, actual.SuppressionWindowSize)

		return nil
	}
}

func testCheckFirstSeenRuleOverrideValues(firstSeenRule *CSEFirstSeenRule, descriptionExpression string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if firstSeenRule.DescriptionExpression != descriptionExpression {
			return fmt.Errorf("bad descriptionExpression, expected \"%s\", got %#v", descriptionExpression, firstSeenRule.DescriptionExpression)
		}
		return nil
	}
}
