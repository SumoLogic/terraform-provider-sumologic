package sumologic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceCSELogMappingVendorProduct_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSumologicCSELogMappingVendorProduct,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceSumologicCSELogMappingVendorProductCheck("data.sumologic_cse_log_mapping_vendor_product.web_gateway"),
				),
			},
		},
	})
}

func testAccDataSourceSumologicCSELogMappingVendorProductCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "product"),
		resource.TestCheckResourceAttr(name, "product", "Web Gateway"),
		resource.TestCheckResourceAttrSet(name, "vendor"),
		resource.TestCheckResourceAttr(name, "vendor", "McAfee"),
		resource.TestCheckResourceAttrSet(name, "guid"),
	)
}

var testDataSourceSumologicCSELogMappingVendorProduct = `
data "sumologic_cse_log_mapping_vendor_product" "web_gateway" {
  	product = "Web Gateway"
	vendor = "McAfee"
}
`
