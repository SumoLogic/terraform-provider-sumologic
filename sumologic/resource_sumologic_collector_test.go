package sumologic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicCollector_basic(t *testing.T) {
	resourceName := "sumologic_collector.test"
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCollectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfigMinimal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "MyTerraformCollector1"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "category", ""),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Etc/UTC"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lookup_by_name", "destroy"},
			},
		},
	})
}

func TestAccSumologicCollector_create(t *testing.T) {
	resourceName := "sumologic_collector.test"
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCollectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "MyTerraformCollector2"),
					resource.TestCheckResourceAttr(resourceName, "description", "MyCollectorDesc"),
					resource.TestCheckResourceAttr(resourceName, "category", "Cat"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Etc/UTC"),
				),
			},
		},
	})
}

func TestAccSumologicCollector_update(t *testing.T) {
	resourceName := "sumologic_collector.test"
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCollectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "MyTerraformCollector2"),
					resource.TestCheckResourceAttr(resourceName, "description", "MyCollectorDesc"),
					resource.TestCheckResourceAttr(resourceName, "category", "Cat"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Etc/UTC"),
				),
			},
			{
				Config: testAccSumologicCollectorConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "MyTerraformCollector2Updated"),
					resource.TestCheckResourceAttr(resourceName, "description", "MyCollectorDescUpdated"),
					resource.TestCheckResourceAttr(resourceName, "category", "Cat"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Etc/UTC"),
				),
			},
		},
	})
}

func testAccCheckCollectorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_collector" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("collector destruction check: collector ID is not set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("collector destruction check: collector id should be int; got %s", rs.Primary.ID)
		}
		c, err := client.GetCollector(id)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if c != nil {
			return fmt.Errorf("Collector still exists")
		}
	}
	return nil
}

var testAccSumologicCollectorConfigMinimal = `
resource "sumologic_collector" "test" {
  name = "MyTerraformCollector1"
}
`

var testAccSumologicCollectorConfig = `
resource "sumologic_collector" "test" {
  name = "MyTerraformCollector2"
  description = "MyCollectorDesc"
  category = "Cat"
  timezone = "Etc/UTC"
}
`

var testAccSumologicCollectorConfigUpdate = `
resource "sumologic_collector" "test" {
  name = "MyTerraformCollector2Updated"
  description = "MyCollectorDescUpdated"
  category = "Cat"
  timezone = "Etc/UTC"
}
`
