package sumologic

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicField_basic(t *testing.T) {
	var field Field
	testFieldName := "fields_provider_test"
	testDataType := "String"
	testState := "Enabled"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFieldDestroy(field),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSumologicFieldConfigImported(testFieldName, testDataType, testState),
			},
			{
				ResourceName:      "sumologic_field.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccField_create(t *testing.T) {
	var field Field
	testFieldName := "fields_provider_test"
	testDataType := "String"
	testState := "Enabled"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFieldDestroy(field),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicField(testFieldName, testDataType, testState),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFieldExists("sumologic_field.test", &field, t),
					testAccCheckFieldAttributes("sumologic_field.test"),
					resource.TestCheckResourceAttr("sumologic_field.test", "field_name", testFieldName),
					resource.TestCheckResourceAttr("sumologic_field.test", "data_type", testDataType),
					resource.TestCheckResourceAttr("sumologic_field.test", "state", testState),
				),
			},
		},
	})
}

func testAccCheckFieldDestroy(field Field) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			u, err := client.GetField(id)
			if err != nil {
				return fmt.Errorf("Encountered an error: " + err.Error())
			}
			if u != nil {
				return fmt.Errorf("Field %s still exists", id)
			}
		}
		return nil
	}
}

func testAccCheckFieldExists(name string, field *Field, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			//need this so that we don't get an unused import error for strconv in some cases
			return fmt.Errorf("Error = %s. Field not found: %s", strconv.FormatBool(ok), name)
		}

		//need this so that we don't get an unused import error for strings in some cases
		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("Field ID is not set")
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)
		newField, err := client.GetField(id)
		if err != nil {
			return fmt.Errorf("Field %s not found", id)
		}
		field = newField
		return nil
	}
}
func testAccCheckSumologicFieldConfigImported(fieldName string, dataType string, state string) string {
	return fmt.Sprintf(`
resource "sumologic_field" "foo" {
      field_name = "%s"
      data_type = "%s"
      state = "%s"
}
`, fieldName, dataType, state)
}

func testAccSumologicField(fieldName string, dataType string, state string) string {
	return fmt.Sprintf(`
resource "sumologic_field" "test" {
    field_name = "%s"
    data_type = "%s"
    state = "%s"
}
`, fieldName, dataType, state)
}

func testAccSumologicFieldUpdate(fieldName string, dataType string, state string) string {
	return fmt.Sprintf(`
resource "sumologic_field" "test" {
      field_name = "%s"
      data_type = "%s"
      state = "%s"
}
`, fieldName, dataType, state)
}

func testAccCheckFieldAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(name, "field_name"),
			resource.TestCheckResourceAttrSet(name, "field_id"),
			resource.TestCheckResourceAttrSet(name, "data_type"),
			resource.TestCheckResourceAttrSet(name, "state"),
		)
		return f(s)
	}
}
