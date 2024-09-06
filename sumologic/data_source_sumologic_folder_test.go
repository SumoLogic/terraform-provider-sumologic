package sumologic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceFolder_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: folderConfig("/Library/Installed Apps"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSumoLogicfDataSourceID("data.sumologic_folder.personal_folder"),
					testAccDataSourceFolderCheck("data.sumologic_folder.personal_folder"),
				),
			},
		},
	})
}

func testAccCheckSumoLogicfDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find SumoLogic Apps data source: %s", n)
		}

		fmt.Printf("%v\n", s.RootModule().Resources)
		fmt.Printf("%v\n", rs.Primary)
		if rs.Primary.ID == "" {
			return fmt.Errorf("SumoLogic Apps data source ID not set")
		}
		return nil
	}
}

func TestAccDataSourceFolder_folder_does_not_exist(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: folderConfig("/Library/Users/dgould+terraform@sumologic.com/doesNotExist"),
				ExpectError: regexp.MustCompile(
					`folder with path '/Library/Users/dgould\+terraform@sumologic.com/doesNotExist' does not exist`),
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
