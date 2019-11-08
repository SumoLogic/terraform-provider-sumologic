package sumologic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourcSumologicCollector(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAccSumologicCollectorConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceCollectorCheck("data.sumologic_collector.by_name", "sumologic_collector.test"),
					testAccDataSourceCollectorCheck("data.sumologic_collector.by_id", "sumologic_collector.test"),
				),
			},
		},
	})
}

func testAccDataSourceCollectorCheck(name, reference string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrPair(name, "id", reference, "id"),
		resource.TestCheckResourceAttrPair(name, "name", reference, "name"),
		resource.TestCheckResourceAttrPair(name, "description", reference, "description"),
		resource.TestCheckResourceAttrPair(name, "category", reference, "category"),
		resource.TestCheckResourceAttrPair(name, "timezone", reference, "timezone"),
	)
}

var testDataSourceAccSumologicCollectorConfig = `
resource "sumologic_collector" "test" {
  name = "MyCollector"
  description = "MyCollectorDesc"
  category = "Cat"
  timezone = "Europe/Berlin"
}

data "sumologic_collector" "by_name" {
  name = "${sumologic_collector.test.name}"
}

data "sumologic_collector" "by_id" {
  id = "${sumologic_collector.test.id}"
}
`
