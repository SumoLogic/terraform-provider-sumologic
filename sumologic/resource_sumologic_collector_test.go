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
					resource.TestCheckResourceAttrSet("sumologic_collector.test", "id"),
				),
			},
		},
	})
}

func TestAccSumologicCollectorLookupByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfigLookupByName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCollectorExists("sumologic_collector.test", t),
					resource.TestCheckResourceAttrSet("sumologic_collector.test", "id"),
				),
			},
		},
	})
}

func TestAccSumologicCollectorAllConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfigAll,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCollectorExists("sumologic_collector.test", t),
					resource.TestCheckResourceAttrSet("sumologic_collector.test", "id"),
				),
			},
		},
	})
}

func TestAccSumologicCollectorChangeConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCollectorExists("sumologic_collector.test", t),
					resource.TestCheckResourceAttr("sumologic_collector.test", "name", "MyCollector"),
					resource.TestCheckResourceAttr("sumologic_collector.test", "description", "MyCollectorDesc"),
					resource.TestCheckResourceAttr("sumologic_collector.test", "category","Cat"),
					resource.TestCheckResourceAttr("sumologic_collector.test", "timezone", "Etc/UTC"),
				),
			},
			{
				Config: testAccSumologicCollectorConfigAll,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCollectorExists("sumologic_collector.test", t),
					resource.TestCheckResourceAttr("sumologic_collector.test", "name", "CollectorName"),
					resource.TestCheckResourceAttr("sumologic_collector.test", "description", "CollectorDesc"),
					resource.TestCheckResourceAttr("sumologic_collector.test", "category","Category"),
					resource.TestCheckResourceAttr("sumologic_collector.test", "timezone", "Europe/Berlin"),
				),
			},
		},
	})
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
  timezone = "Etc/UTC"
}

`

var testAccSumologicCollectorConfigLookupByName = `

resource "sumologic_collector" "test" {
  name = "MyOtherCollector"
  description = "MyCollectorDesc"
  category = "Cat"
  timezone = "Europe/Berlin"
  lookup_by_name=true
}

`

var testAccSumologicCollectorConfigAll = `
resource "sumologic_collector" "test" {
  name="CollectorName"
  description="CollectorDesc"
  category="Category"
  timezone="Europe/Berlin"
  lookup_by_name=true
  destroy=true
}
`
