package sumologic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceFolder_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: folderConfig("/Library/Users/dgould+terraform@sumologic.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceFolderCheck("data.sumologic_folder.personal_folder"),
				),
			},
		},
	})
}

func TestAccDataSourceFolder_folder_does_not_exist(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: folderConfig("/Library/Users/dgould+terraform@sumologic.com/doesNotExist"),
				ExpectError: regexp.MustCompile(
					"folder with path '/Library/Users/dgould\\+terraform@sumologic.com/doesNotExist' does not exist"),
			},
		},
	})
}

func testAccDataSourceFolderCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

func folderConfig(path string) string {
	return fmt.Sprintf(`
		data "sumologic_folder" "personal_folder" {
			path = "%s"
		}
	`, path)
}
