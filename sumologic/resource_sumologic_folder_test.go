package sumologic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// func TestAccFolder(t *testing.T) {
// 	var folder Folder
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckFolderDestroy(folder),
// 		Steps: []resource.TestStep{
// 			{
// 				Config: TestAccSumologicFolder,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckFolderExists("sumologic_folder.test", &folder, t),
// 					testAccCheckFolderAttributes("sumologic_folder.test"),
// 					resource.TestCheckResourceAttr("sumologic_folder.test", "description", ""),
// 				),
// 			},
// 			{
// 				Config: TestAccSumologicFolderUpdate,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckFolderExists("sumologic_folder.test", &folder, t),
// 					testAccCheckFolderAttributes("sumologic_folder.test"),
// 					resource.TestCheckResourceAttr("sumologic_folder.test", "description", "Update test"),
// 				),
// 			},
// 		},
// 	})
// }

func testAccCheckFolderExists(name string, folder *Folder, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Folder not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Folder ID is not set")
		}

		id := rs.Primary.ID
		c := testAccProvider.Meta().(*Client)
		newFolder, err := c.GetFolder(id)
		if err != nil {
			return fmt.Errorf("Folder %s not found", id)
		}
		*folder = newFolder
		return nil
	}
}

func testAccCheckFolderAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(name, "name"),
			resource.TestCheckResourceAttrSet(name, "parent_id"),
			resource.TestCheckResourceAttrSet(name, "created_at"),
			resource.TestCheckResourceAttrSet(name, "created_by"),
			resource.TestCheckResourceAttrSet(name, "modified_at"),
			resource.TestCheckResourceAttrSet(name, "modified_by"),
			resource.TestCheckResourceAttr(name, "item_type", "Folder"),
			resource.TestCheckResourceAttr(name, "permissions.#", "6"),
			resource.TestCheckNoResourceAttr(name, "children.#"),
		)
		return f(s)
	}
}

func testAccCheckFolderDestroy(folder Folder) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		_, err := client.GetFolder(folder.ID)
		if err == nil {
			return fmt.Errorf("Folder still exists")
		}
		return nil
	}
}

func testAccPreCheck(t *testing.T) {
	if strings.Contains(TestAccSumologicFolder, "PERSONAL_FOLDER_HEX_ID") ||
		strings.Contains(TestAccSumologicFolderUpdate, "PERSONAL_FOLDER_HEX_ID") {
		t.Fatal("The parent_id must be set for the TestAccSumologicFolder and" +
			" TestAccSumologicFolderUpdate variables")
	}
}

var TestAccSumologicFolder = `
resource "sumologic_folder" "test" {
  name = "MyFolder"
	parent_id = "<PERSONAL_FOLDER_HEX_ID>"
}
`

var TestAccSumologicFolderUpdate = `
resource "sumologic_folder" "test" {
  name = "MyFolder"
	parent_id = "<PERSONAL_FOLDER_HEX_ID>"
	description = "Update test"
}
`
