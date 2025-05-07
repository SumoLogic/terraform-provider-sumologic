package sumologic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourcePartition_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSumologicPartitionConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourcePartitionCheck("data.sumologic_partition.test"),
				),
			},
		},
	})
}

func testAccDataSourcePartitionCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttr(name, "name", "sumologic_default"),
		resource.TestCheckResourceAttrSet(name, "analytics_tier"),
		resource.TestCheckResourceAttrSet(name, "retention_period"),
		resource.TestCheckResourceAttrSet(name, "is_compliant"),
		resource.TestCheckResourceAttrSet(name, "is_active"),
		resource.TestCheckResourceAttrSet(name, "total_bytes"),
		resource.TestCheckResourceAttrSet(name, "index_type"),
		resource.TestCheckResourceAttrSet(name, "is_included_in_default_search"),
	)
}

var testDataSourceSumologicPartitionConfig = `
	data "sumologic_partition" "test" {
		id = "000000000005A4A3"
	}
`
