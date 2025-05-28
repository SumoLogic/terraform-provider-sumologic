package sumologic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
		resource.TestCheckResourceAttr(name, "id", "00000000000AD88D"),
		resource.TestCheckResourceAttr(name, "name", "apache"),
		resource.TestCheckResourceAttr(name, "routing_expression", "_sourcecategory=*/Apache"),
		resource.TestCheckResourceAttr(name, "analytics_tier", "continuous"),
		resource.TestCheckResourceAttr(name, "retention_period", "365"),
		resource.TestCheckResourceAttr(name, "is_compliant", "false"),
		resource.TestCheckResourceAttr(name, "is_active", "true"),
		resource.TestCheckResourceAttr(name, "total_bytes", "0"),
		resource.TestCheckResourceAttr(name, "data_forwarding_id", ""),
		resource.TestCheckResourceAttr(name, "index_type", "Partition"),
		resource.TestCheckResourceAttr(name, "reduce_retention_period_immediately", "false"),
		resource.TestCheckResourceAttr(name, "is_included_in_default_search", "true"),
	)
}

var testDataSourceSumologicPartitionConfig = `
	data "sumologic_partition" "test" {
		id = "00000000000AD88D"
	}
`
