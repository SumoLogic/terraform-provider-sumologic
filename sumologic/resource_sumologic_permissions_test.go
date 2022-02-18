package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPermission_create(t *testing.T) {
	var response PermissionsResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPermissionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicPermission(
					otherResource, false, "create", "role", "sumologic_role.permission_test_role"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPermissionExists("sumologic_content_permission.content_permission_test", &response, t),
					testAccCheckPermissionAttributes("sumologic_content_permission.content_permission_test"),
					resource.TestCheckResourceAttr(
						"sumologic_content_permission.content_permission_test", "notify_recipient", "false"),
					resource.TestCheckResourceAttr(
						"sumologic_content_permission.content_permission_test", "notification_message", "create"),
					// Need to upgrade to plugin SDKv2 to test TypeSet objects
				),
			},
		},
	})
}

func TestAccPermission_update(t *testing.T) {
	var response PermissionsResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPermissionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicPermission(
					otherResource, false, "create", "role", "sumologic_role.permission_test_role"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPermissionExists("sumologic_content_permission.content_permission_test", &response, t),
					testAccCheckPermissionAttributes("sumologic_content_permission.content_permission_test"),
					resource.TestCheckResourceAttr(
						"sumologic_content_permission.content_permission_test", "notify_recipient", "false"),
					resource.TestCheckResourceAttr(
						"sumologic_content_permission.content_permission_test", "notification_message", "create"),
				),
			}, {
				Config: testAccSumologicPermissionUpdate(
					otherResource, true, "update", "user", "sumologic_user.permission_test_user"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPermissionExists("sumologic_content_permission.content_permission_test", &response, t),
					testAccCheckPermissionAttributes("sumologic_content_permission.content_permission_test"),
					resource.TestCheckResourceAttr(
						"sumologic_content_permission.content_permission_test", "notify_recipient", "true"),
					resource.TestCheckResourceAttr(
						"sumologic_content_permission.content_permission_test", "notification_message", "update"),
				),
			},
		},
	})
}

func TestAccPermission_delete(t *testing.T) {
	var response PermissionsResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPermissionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicPermission(
					otherResource, false, "create", "role", "sumologic_role.permission_test_role"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPermissionExists("sumologic_content_permission.content_permission_test", &response, t),
					testAccCheckPermissionAttributes("sumologic_content_permission.content_permission_test"),
					resource.TestCheckResourceAttr(
						"sumologic_content_permission.content_permission_test", "notify_recipient", "false"),
					resource.TestCheckResourceAttr(
						"sumologic_content_permission.content_permission_test", "notification_message", "create"),
				),
			}, {
				Config: testAccSumologicPermissionDelete(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPermissionDelete("sumologic_content.permission_test_content"),
				),
			},
		},
	})
}

func testAccCheckPermissionDestroy(s *terraform.State) error {
	// keep it here at this moment
	client := testAccProvider.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_content_permission" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Content Permission ID is not set")
		}
		id := rs.Primary.ID
		u, err := client.GetPermissions(id)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if u != nil {
			return fmt.Errorf("Content and permissions still exists")
		}
	}
	return nil
}

func testAccCheckPermissionDelete(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		meta := testAccProvider.Meta()
		client := meta.(*Client)
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Content not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Content ID is not set")
		}

		contentId := rs.Primary.ID
		permissions, err := client.GetPermissions(contentId)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s == nil {
			return fmt.Errorf("All permissions have been deleted")
		}
		creatorId, _ := getCreatorId(contentId, meta)
		if creatorId == "" {
			return fmt.Errorf("Empty creator's id")
		}
		for _, permission := range permissions.ExplicitPermissions {
			if creatorId != permission.SourceId || permission.SourceType != "user" {
				return fmt.Errorf("Contains more than creator's permissions")
			}
		}
		return nil
	}
}

func testAccCheckPermissionExists(name string, response *PermissionsResponse,
	t *testing.T) resource.TestCheckFunc {
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
			resource.TestCheckResourceAttrSet(name, "notify_recipient"),
			resource.TestCheckResourceAttrSet(name, "notification_message"),
		)
		return f(s)
	}
}

func testAccSumologicPermission(resource string, notify_recipient bool, notification_message string,
	source_type string, source_id string) string {
	return fmt.Sprintf(`
		%s

		resource "sumologic_content_permission" "content_permission_test" {
			content_id = sumologic_content.permission_test_content.id
			notify_recipient = %t
			notification_message = "%s"
			permission {
				permission_name = "View"
				source_type = "%s"
				source_id = %s.id
			}
		}`,
		resource, notify_recipient, notification_message, source_type, source_id,
	)
}

func testAccSumologicPermissionUpdate(resource string, notify_recipient bool,
	notification_message string, source_type string, source_id string) string {
	return fmt.Sprintf(`
		%s

		resource "sumologic_content_permission" "content_permission_test" {
			content_id = sumologic_content.permission_test_content.id
			notify_recipient = %t
			notification_message = "%s"
			permission {
				permission_name = "View"
				source_type = "role"
				source_id = sumologic_role.permission_test_role.id
			}
			permission {
				permission_name = "View"
				source_type = "%s"
				source_id = %s.id
			}
		}`,
		resource, notify_recipient, notification_message, source_type, source_id,
	)
}

func testAccSumologicPermissionDelete() string {
	return otherResource
}

var otherResource = `
	data "sumologic_personal_folder" "personalFolder" {}

	resource "sumologic_content" "permission_test_content" {
		parent_id = data.sumologic_personal_folder.personalFolder.id
		config = jsonencode({
			"type": "FolderSyncDefinition",
			"name": "test_permission_resource_folder",
			"description": "",
			"children": []
		})
	}

	resource "sumologic_role" "permission_test_role" {
		name        = "permission_test_role"
		description = "Testing content permission resource"
	}

	resource "sumologic_user" "permission_test_user" {
		first_name   = "test"
		last_name    = "permission"
		email        = "testpermission@gmail.com"
		is_active    = true
		role_ids     = [sumologic_role.permission_test_role.id]
		transfer_to  = ""
	}`
