package sumologic

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSumologicCollector(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCollectorExists("sumologic_collector.test", t),
				),
			},
		}})
}

func testAccCheckCollectorExists(n string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return nil
	}
}

var testAccSumologicCollectorConfig = `

resource "sumologic_collector" "test" {
  name = "MyCollector"
  description = "MyCollectorDesc"
  category = "Cat"
  timezone = "Europe/Berlin"
}

`
