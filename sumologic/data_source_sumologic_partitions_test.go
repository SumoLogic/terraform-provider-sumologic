package sumologic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePartitions_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSumologicPartitionsConfig,
				Check: resource.ComposeTestCheckFunc(
					checkResourceAttrGreaterThanZero("data.sumologic_partitions.test", "partitions.#"),
				),
			},
		},
	})
}

var testDataSourceSumologicPartitionsConfig = `
	data "sumologic_partitions" "test" {}
`
