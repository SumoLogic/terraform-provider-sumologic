package sumologic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSumologicCollector(t *testing.T) {
	var collector *Collector
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCollectorExists("sumologic_collector.test", &collector, t),
					resource.TestCheckResourceAttrSet("sumologic_collector.test", "id"),
				),
			},
		},
	})
}

func TestAccSumologicCollectorLookupByName(t *testing.T) {
	var collector *Collector
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfigLookupByName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCollectorExists("sumologic_collector.test", &collector, t),
					resource.TestCheckResourceAttrSet("sumologic_collector.test", "id"),
				),
			},
		},
	})
}

func TestAccSumologicCollectorAllConfig(t *testing.T) {
	var collector *Collector
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfigAll,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCollectorExists("sumologic_collector.test", &collector, t),
					resource.TestCheckResourceAttrSet("sumologic_collector.test", "id"),
				),
			},
		},
	})
}

func TestAccSumologicCollectorChangeConfig(t *testing.T) {
	var collector *Collector
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCollectorExists("sumologic_collector.test", &collector, t),
					resource.TestCheckResourceAttr("sumologic_collector.test", "name", "MyCollector"),
					resource.TestCheckResourceAttr("sumologic_collector.test", "description", "MyCollectorDesc"),
					resource.TestCheckResourceAttr("sumologic_collector.test", "category", "Cat"),
					resource.TestCheckResourceAttr("sumologic_collector.test", "timezone", "Etc/UTC"),
				),
			},
			{
				Config: testAccSumologicCollectorConfigAll,
				Check: resource.ComposeTestCheckFunc(
					// check the id of this resource is the same as the one in the previous step
					testAccCheckCollectorId("sumologic_collector.test", &collector),
					testAccCheckCollectorExists("sumologic_collector.test", &collector, t),
					resource.TestCheckResourceAttr("sumologic_collector.test", "name", "CollectorName"),
					resource.TestCheckResourceAttr("sumologic_collector.test", "description", "CollectorDesc"),
					resource.TestCheckResourceAttr("sumologic_collector.test", "category", "Category"),
					resource.TestCheckResourceAttr("sumologic_collector.test", "timezone", "Europe/Berlin"),
				),
			},
		},
	})
}

func testAccCheckCollectorId(name string, collector **Collector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("collector ID is not set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("collector id should be int; got %s", rs.Primary.ID)
		}

		expectedId := (**collector).ID
		if id != expectedId {
			return fmt.Errorf("incorrect collector id: got %d; expected %d", id, expectedId)
		}
		return nil
	}
}

func testAccCheckCollectorExists(name string, collector **Collector, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("collector ID is not set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("collector id should be int; got %s", rs.Primary.ID)
		}
		c := testAccProvider.Meta().(*Client)
		*collector, err = c.GetCollector(id)
		if err != nil {
			return fmt.Errorf("collector %d not foung", id)
		}
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
