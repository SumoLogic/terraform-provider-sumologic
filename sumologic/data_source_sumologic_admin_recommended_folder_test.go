package sumologic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestDataSourceAdminRecommendedFolder_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSumologicAdminRecommendedFolderConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceCollectorCheck("data.sumologic_collector.by_name", "sumologic_collector.test"),
					testAccDataSourceCollectorCheck("data.sumologic_collector.by_id", "sumologic_collector.test"),
				),
			},
		},
	})
}

func testAccDataSourceAdminRecommendedFolderCheck(name, reference string) resource.TestCheckFunc {
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