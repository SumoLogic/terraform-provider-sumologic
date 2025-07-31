package sumologic

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicFieldExtractionRule_basic(t *testing.T) {
	var fieldextractionrule FieldExtractionRule
	testName := FieldsMap["FieldExtractionRule"]["name"] + acctest.RandString(8)
	testScope := FieldsMap["FieldExtractionRule"]["scope"]
	testParseExpression := FieldsMap["FieldExtractionRule"]["parseExpression"]
	testEnabled, _ := strconv.ParseBool(FieldsMap["FieldExtractionRule"]["enabled"])

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFieldExtractionRuleDestroy(fieldextractionrule),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSumologicFieldExtractionRuleConfigImported(testName, testScope, testParseExpression, testEnabled),
			},
			{
				ResourceName:      "sumologic_field_extraction_rule.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSumologicFieldExtractionRule_create(t *testing.T) {
	var fieldextractionrule FieldExtractionRule
	testName := FieldsMap["FieldExtractionRule"]["name"] + acctest.RandString(8)
	testScope := FieldsMap["FieldExtractionRule"]["scope"]
	testParseExpression := FieldsMap["FieldExtractionRule"]["parseExpression"]
	testEnabled, _ := strconv.ParseBool(FieldsMap["FieldExtractionRule"]["enabled"])
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFieldExtractionRuleDestroy(fieldextractionrule),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicFieldExtractionRule(testName, testScope, testParseExpression, testEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFieldExtractionRuleExists("sumologic_field_extraction_rule.test", &fieldextractionrule, t),
					testAccCheckFieldExtractionRuleAttributes("sumologic_field_extraction_rule.test"),
					resource.TestCheckResourceAttr("sumologic_field_extraction_rule.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_field_extraction_rule.test", "scope", testScope),
					resource.TestCheckResourceAttr("sumologic_field_extraction_rule.test", "parse_expression", testParseExpression),
					resource.TestCheckResourceAttr("sumologic_field_extraction_rule.test", "enabled", strconv.FormatBool(testEnabled)),
				),
			},
		},
	})
}

func testAccCheckFieldExtractionRuleDestroy(fieldextractionrule FieldExtractionRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			u, err := client.GetFieldExtractionRule(id)
			if err != nil {
				return fmt.Errorf("Encountered an error: " + err.Error())
			}
			if u != nil {
				return fmt.Errorf("FieldExtractionRule still exists")
			}
		}
		return nil
	}
}

func testAccCheckFieldExtractionRuleExists(name string, fieldextractionrule *FieldExtractionRule, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			//need this so that we don't get an unused import error for strconv in some cases
			return fmt.Errorf("Error = %s. FieldExtractionRule not found: %s", strconv.FormatBool(ok), name)
		}

		//need this so that we don't get an unused import error for strings in some cases
		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("FieldExtractionRule ID is not set")
		}

		id := rs.Primary.ID
		c := testAccProvider.Meta().(*Client)
		newFieldExtractionRule, err := c.GetFieldExtractionRule(id)
		if err != nil {
			return fmt.Errorf("FieldExtractionRule %s not found", id)
		}
		fieldextractionrule = newFieldExtractionRule
		return nil
	}
}

func TestAccSumologicFieldExtractionRule_update(t *testing.T) {
	var fieldextractionrule FieldExtractionRule
	randomSuffix := acctest.RandString(8)
	testName := FieldsMap["FieldExtractionRule"]["name"] + randomSuffix
	testScope := FieldsMap["FieldExtractionRule"]["scope"]
	testParseExpression := FieldsMap["FieldExtractionRule"]["parseExpression"]
	testEnabled, _ := strconv.ParseBool(FieldsMap["FieldExtractionRule"]["enabled"])

	testUpdatedName := FieldsMap["FieldExtractionRule"]["updatedName"] + randomSuffix
	testUpdatedScope := FieldsMap["FieldExtractionRule"]["updatedScope"]
	testUpdatedParseExpression := FieldsMap["FieldExtractionRule"]["updatedParseExpression"]
	testUpdatedEnabled, _ := strconv.ParseBool(FieldsMap["FieldExtractionRule"]["updatedEnabled"])

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFieldExtractionRuleDestroy(fieldextractionrule),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicFieldExtractionRule(testName, testScope, testParseExpression, testEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFieldExtractionRuleExists("sumologic_field_extraction_rule.test", &fieldextractionrule, t),
					testAccCheckFieldExtractionRuleAttributes("sumologic_field_extraction_rule.test"),
					resource.TestCheckResourceAttr("sumologic_field_extraction_rule.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_field_extraction_rule.test", "scope", testScope),
					resource.TestCheckResourceAttr("sumologic_field_extraction_rule.test", "parse_expression", testParseExpression),
					resource.TestCheckResourceAttr("sumologic_field_extraction_rule.test", "enabled", strconv.FormatBool(testEnabled)),
				),
			},
			{
				Config: testAccSumologicFieldExtractionRuleUpdate(testUpdatedName, testUpdatedScope, testUpdatedParseExpression, testUpdatedEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFieldExtractionRuleExists("sumologic_field_extraction_rule.test", &fieldextractionrule, t),
					testAccCheckFieldExtractionRuleAttributes("sumologic_field_extraction_rule.test"),
					resource.TestCheckResourceAttr("sumologic_field_extraction_rule.test", "name", testUpdatedName),
					resource.TestCheckResourceAttr("sumologic_field_extraction_rule.test", "scope", testUpdatedScope),
					resource.TestCheckResourceAttr("sumologic_field_extraction_rule.test", "parse_expression", testUpdatedParseExpression),
					resource.TestCheckResourceAttr("sumologic_field_extraction_rule.test", "enabled", strconv.FormatBool(testUpdatedEnabled)),
				),
			},
		},
	})
}

func testAccCheckSumologicFieldExtractionRuleConfigImported(name string, scope string, parseExpression string, enabled bool) string {
	return fmt.Sprintf(`
resource "sumologic_field_extraction_rule" "foo" {
      name = "%s"
      scope = "%s"
      parse_expression = "%s"
      enabled = %t
}
`, name, scope, parseExpression, enabled)
}

func testAccSumologicFieldExtractionRule(name string, scope string, parseExpression string, enabled bool) string {
	return fmt.Sprintf(`
resource "sumologic_field_extraction_rule" "test" {
    name = "%s"
    scope = "%s"
    parse_expression = "%s"
    enabled = %t
}
`, name, scope, parseExpression, enabled)
}

func testAccSumologicFieldExtractionRuleUpdate(name string, scope string, parseExpression string, enabled bool) string {
	return fmt.Sprintf(`
resource "sumologic_field_extraction_rule" "test" {
      name = "%s"
      scope = "%s"
      parse_expression = "%s"
      enabled = %t
}
`, name, scope, parseExpression, enabled)
}

func testAccCheckFieldExtractionRuleAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(name, "name"),
			resource.TestCheckResourceAttrSet(name, "scope"),
			resource.TestCheckResourceAttrSet(name, "parse_expression"),
			resource.TestCheckResourceAttrSet(name, "enabled"),
		)
		return f(s)
	}
}
