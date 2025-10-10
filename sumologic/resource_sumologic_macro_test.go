package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicMacro_crud(t *testing.T) {
	var macro Macro
	testNameSuffix := acctest.RandString(16)

	testName := "terraform_test_macro_" + testNameSuffix

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMacroDestroy(macro),
		Steps: []resource.TestStep{
			{
				// create
				Config: testAccSumologicMacro(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_macro.test", "name", "terraform_macro_"+testName),
					resource.TestCheckResourceAttr("sumologic_macro.test", "description", ""),
					resource.TestCheckResourceAttr("sumologic_macro.test", "definition", "_sourceCategory=stream {{arg1}} | count by _timeslice"),
					resource.TestCheckResourceAttr("sumologic_macro.test", "enabled", "true"),
					resource.TestCheckResourceAttr("sumologic_macro.test", "argument.0.name", "arg1"),
					resource.TestCheckResourceAttr("sumologic_macro.test", "argument.0.type", "String"),
					resource.TestCheckResourceAttr("sumologic_macro.test", "argument_validation.0.eval_expression", "arg1 > 3"),
					resource.TestCheckResourceAttr("sumologic_macro.test", "argument_validation.0.error_message", "This is an error"),
				),
			},

			{
				// update
				Config: testAccSumologicMacroEdit(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_macro.test", "name", "terraform_macro_"+testName),
					resource.TestCheckResourceAttr("sumologic_macro.test", "description", "edited"),
					resource.TestCheckResourceAttr("sumologic_macro.test", "definition", "_sourceCategory=stream {{arg2}} | count by _timeslice"),
					resource.TestCheckResourceAttr("sumologic_macro.test", "enabled", "false"),
					resource.TestCheckResourceAttr("sumologic_macro.test", "argument.0.name", "arg2"),
					resource.TestCheckResourceAttr("sumologic_macro.test", "argument.0.type", "Any"),
					resource.TestCheckResourceAttr("sumologic_macro.test", "argument_validation.0.eval_expression", "arg2 > 3"),
					resource.TestCheckResourceAttr("sumologic_macro.test", "argument_validation.0.error_message", "This is an updated error"),
				),
			},
			{
				ResourceName:      "sumologic_macro.test",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccSumologicMacro(testName string) string {
	return fmt.Sprintf(`
    resource "sumologic_macro" "test" {
	name = "terraform_macro_%s"
	definition = "_sourceCategory=stream {{arg1}} | count by _timeslice"
        argument {
        	name = "arg1"
        	type = "String"
        }
        argument_validation {
        	eval_expression = "arg1 > 3"
        	error_message = "This is an error"
        }
    }
    `, testName)
}

func testAccSumologicMacroEdit(testName string) string {
	return fmt.Sprintf(`
    resource "sumologic_macro" "test" {
	name = "terraform_macro_%s"
	description = "edited"
	enabled = false
	definition = "_sourceCategory=stream {{arg2}} | count by _timeslice"
        argument {
        	name = "arg2"
        	type = "Any"
        }
        argument_validation {
        	eval_expression = "arg2 > 3"
        	error_message = "This is an updated error"
        }
    }
    `, testName)
}

func testAccCheckMacroDestroy(macro Macro) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			u, err := client.GetMacro(id)
			if err != nil {
				return fmt.Errorf("Encountered an error: " + err.Error())
			}
			if u != nil {
				return fmt.Errorf("Macro %s still exists", id)
			}
		}
		return nil
	}
}
