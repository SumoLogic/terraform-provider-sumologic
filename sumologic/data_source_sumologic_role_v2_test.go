package sumologic

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSumologicRoleV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAccSumologicRoleV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceRoleV2Check("data.sumologic_role_v2.by_name", "sumologic_role_v2.test"),
					testAccDataSourceRoleV2Check("data.sumologic_role_v2.by_id", "sumologic_role_v2.test"),
				),
			},
		},
	})
}

func testAccDataSourceRoleV2Check(name, reference string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrPair(name, "id", reference, "id"),
		resource.TestCheckResourceAttrPair(name, "name", reference, "name"),
		resource.TestCheckResourceAttrPair(name, "description", reference, "description"),
		resource.TestCheckResourceAttrPair(name, "capabilities", reference, "capabilities"),
		resource.TestCheckResourceAttrPair(name, "selection_type", reference, "selection_type"),
		resource.TestCheckResourceAttrPair(name, "audit_data_filter", reference, "audit_data_filter"),
		resource.TestCheckResourceAttrPair(name, "security_data_filter", reference, "security_data_filter"),
		resource.TestCheckResourceAttrPair(name, "selected_views", reference, "selected_views"),
	)
}

func TestAccDataSourceSumologicRoleV2_role_name_doesnt_exist(t *testing.T) {
	roleDoestExistConfig := `
  data "sumologic_role_v2" "role_name_doesnt_exist" {
    name = "someRoleNameDoesntExist87461"
  }`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      roleDoestExistConfig,
				ExpectError: regexp.MustCompile("role with name 'someRoleNameDoesntExist87461' does not exist"),
			},
		},
	})
}

func TestAccDataSourceSumologicRoleV2_role_id_doesnt_exist(t *testing.T) {
	roleDoestExistConfig := `
  data "sumologic_role_v2" "role_id_doesnt_exist" {
    id = 99999999999999
  }`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      roleDoestExistConfig,
				ExpectError: regexp.MustCompile("role with id 99999999999999 not found"),
			},
		},
	})
}

var testDataSourceAccSumologicRoleV2Config = `
resource "sumologic_role_v2" "test" {
  name = "My_SumoRole_V2"
  description = "My_SumoRoleDesc"
  capabilities = ["viewCollectors"]
  selection_type = "All"
  audit_data_filter = "info"
  security_data_filter = "error"
  log_analytics_filter = "!_sourceCategory=collector"
}

data "sumologic_role_v2" "by_name" {
  name = "${sumologic_role_v2.test.name}"
}

data "sumologic_role_v2" "by_id" {
  id = "${sumologic_role_v2.test.id}"
}
`
