package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccFolder_create(t *testing.T) {
	var folder Folder
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFolderDestroy(folder),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicFolder(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFolderExists("sumologic_folder.test", &folder, t),
					testAccCheckFolderAttributes("sumologic_folder.test"),
					resource.TestCheckResourceAttr("sumologic_folder.test", "description", "test"),
					resource.TestCheckResourceAttr("sumologic_folder.test", "name", rName),
				),
			},
		},
	})
}

func TestAccFolder_update(t *testing.T) {
	var folder Folder
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFolderDestroy(folder),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicFolder(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFolderExists("sumologic_folder.test", &folder, t),
					testAccCheckFolderAttributes("sumologic_folder.test"),
					resource.TestCheckResourceAttr("sumologic_folder.test", "description", "test"),
					resource.TestCheckResourceAttr("sumologic_folder.test", "name", rName),
				),
			},
			{
				Config: testAccSumologicFolderUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFolderExists("sumologic_folder.test", &folder, t),
					testAccCheckFolderAttributes("sumologic_folder.test"),
					resource.TestCheckResourceAttr("sumologic_folder.test", "description", "Update test"),
					resource.TestCheckResourceAttr("sumologic_folder.test", "name", rName),
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

func testAccSumologicFolder(name string) string {
	return fmt.Sprintf(`
data "sumologic_personal_folder" "personalFolder" {}
resource "sumologic_folder" "test" {
  name = "%s"
  parent_id = "${data.sumologic_personal_folder.personalFolder.id}"
  description = "test"
}
`, name)
}

func testAccSumologicFolderUpdate(name string) string {
	return fmt.Sprintf(`
data "sumologic_personal_folder" "personalFolder" {}
resource "sumologic_folder" "test" {
  name = "%s"
  parent_id = "${data.sumologic_personal_folder.personalFolder.id}"
  description = "Update test"
}
`, name)
}
