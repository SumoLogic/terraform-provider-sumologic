package sumologic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAdminRecommendedFolder_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSumologicAdminRecommendedFolderConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAdminRecommendedFolderCheck("data.sumologic_admin_recommended_folder.test"),
				),
			},
		},
	})
}

func testAccDataSourceAdminRecommendedFolderCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttr(name, "name", "Admin Recommended"),
		resource.TestCheckResourceAttrSet(name, "description"),
	)
}

var testDataSourceSumologicAdminRecommendedFolderConfig = `
data "sumologic_admin_recommended_folder" "test" {
}
`
