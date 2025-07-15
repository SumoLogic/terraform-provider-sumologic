package sumologic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMonitorFolder_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: monitorFolderConfig("/Monitor/Terraform Test/Subfolder"),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceMonitorFolderCheck("data.sumologic_monitor_folder.test"),
				),
			},
		},
	})
}

func TestAccDataSourceMonitorFolder_folder_does_not_exist(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: monitorFolderConfig("/Monitor/Terraform Test/Subfolder/DoesNotExist"),
				ExpectError: regexp.MustCompile(
					`folder with path '/Monitor/Terraform Test/Subfolder/DoesNotExist' does not exist`),
			},
		},
	})
}

func testAccDataSourceMonitorFolderCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

func monitorFolderConfig(path string) string {
	return fmt.Sprintf(`
		data "sumologic_monitor_folder" "test" {
			path = "%s"
		}
	`, path)
}
