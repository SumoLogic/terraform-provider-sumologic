package sumologic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourcSumologicRole_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAccSumologicRoleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceRoleCheck("data.sumologic_role.by_name", "sumologic_role.test"),
					testAccDataSourceRoleCheck("data.sumologic_role.by_id", "sumologic_role.test"),
				),
			},
		},
	})
}

func testAccDataSourceRoleCheck(name, reference string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrPair(name, "id", reference, "id"),
		resource.TestCheckResourceAttrPair(name, "name", reference, "name"),
		resource.TestCheckResourceAttrPair(name, "description", reference, "description"),
		resource.TestCheckResourceAttrPair(name, "filter_predicate", reference, "filter_predicate"),
		resource.TestCheckResourceAttrPair(name, "capabilities", reference, "capabilities"),
	)
}

var testDataSourceAccSumologicRoleConfig = `
resource "sumologic_role" "test" {
  name = "My_SumoRole"
  description = "My_SumoRoleDesc"
  filter_predicate = "_sourceCategory=Test"
  capabilities = ["viewCollectors"]
}

data "sumologic_role" "by_name" {
  name = "${sumologic_role.test.name}"
}

data "sumologic_role" "by_id" {
  id = "${sumologic_role.test.id}"
}
`
