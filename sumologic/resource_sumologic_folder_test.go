package sumologic

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccFolderCreate(t *testing.T) {
	var folder Folder
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	personalFolderId := os.Getenv("SUMOLOGIC_PF")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFolderDestroy(folder),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicFolder(rName, personalFolderId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFolderExists("sumologic_folder.test", &folder, t),
					testAccCheckFolderAttributes("sumologic_folder.test"),
					resource.TestCheckResourceAttr("sumologic_folder.test", "description", "test"),
					resource.TestCheckResourceAttr("sumologic_folder.test", "name", rName),
					resource.TestCheckResourceAttr("sumologic_folder.test", "parent_id", personalFolderId),
				),
			},
		},
	})
}

func TestAccFolderUpdate(t *testing.T) {
	var folder Folder
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	personalFolderId := os.Getenv("SUMOLOGIC_PF")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFolderDestroy(folder),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicFolder(rName, personalFolderId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFolderExists("sumologic_folder.test", &folder, t),
					testAccCheckFolderAttributes("sumologic_folder.test"),
					resource.TestCheckResourceAttr("sumologic_folder.test", "description", "test"),
					resource.TestCheckResourceAttr("sumologic_folder.test", "name", rName),
					resource.TestCheckResourceAttr("sumologic_folder.test", "parent_id", personalFolderId),
				),
			},
			{
				Config: testAccSumologicFolderUpdate(rName, personalFolderId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFolderExists("sumologic_folder.test", &folder, t),
					testAccCheckFolderAttributes("sumologic_folder.test"),
					resource.TestCheckResourceAttr("sumologic_folder.test", "description", "Update test"),
					resource.TestCheckResourceAttr("sumologic_folder.test", "name", rName),
					resource.TestCheckResourceAttr("sumologic_folder.test", "parent_id", personalFolderId),
				),
			},
		},
	})
}

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
		folder = newFolder
		return nil
	}
}

func testAccCheckFolderAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(name, "name"),
			resource.TestCheckResourceAttrSet(name, "description"),
			resource.TestCheckResourceAttrSet(name, "parent_id"),
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

func testAccSumologicFolder(name string, parentId string) string {
	return fmt.Sprintf(`
resource "sumologic_folder" "test" {
  name = "%s"
  parent_id = "%s"
  description = "test"
}
`, name, parentId)
}

func testAccSumologicFolderUpdate(name string, parentId string) string {
	return fmt.Sprintf(`
resource "sumologic_folder" "test" {
  name = "%s"
  parent_id = "%s"
  description = "Update test"
}
`, name, parentId)
}
