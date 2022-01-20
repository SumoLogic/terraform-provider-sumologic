package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPermission_create(t *testing.T) {
	var response PermissionsResponse
	var permission Permission

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPermissionDestroy(permission),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicPermission("View", "role"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPermissionExists("sumologic_permission.test_permission", &response, t),
					testAccCheckPermissionAttributes("sumologic_permission.test_permission"),
					resource.TestCheckResourceAttr("sumologic_permission.test_permission", "permission.0.permission_name", "View"),
					resource.TestCheckResourceAttr("sumologic_permission.test_permission", "permission.0.source_type", "role"),
				),
			},
		},
	})
}

func testAccCheckPermissionDestroy(permission Permission) resource.TestCheckFunc {
	// ??
	return func(s *terraform.State) error {
		return nil
	}
}

func testAccCheckPermissionExists(name string, response *PermissionsResponse, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Content permissions not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Content ID is not set")
		}

		id := rs.Primary.ID
		c := testAccProvider.Meta().(*Client)
		newResponse, err := c.GetPermissions(id)
		if err != nil {
			return fmt.Errorf("Content %s not found", id)
		}
		response = newResponse
		return nil
	}
}

func testAccCheckPermissionAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(name, "content_id"),
			resource.TestCheckResourceAttrSet(name, "permission"),
		)
		return f(s)
	}
}

func testAccSumologicPermission(permission_name string, source_type string) string {
	return fmt.Sprintf(`
data "sumologic_personal_folder" "personalFolder" {}

data "sumologic_role" "role" { 
	name = "test-role"
}	

resource "sumologic_content" "test_content" {
	parent_id = data.sumologic_personal_folder.personalFolder.id
	config = jsonencode({})
}

resource "sumologic_permissions" "test_permission" {
    content_id = sumologic_content.test_content.id
	permission {
		permission_name = "%s"
		source_type = "%s"
		sourceId = data.sumologic_role.role.id
	}
}
`, permission_name, source_type)
}
