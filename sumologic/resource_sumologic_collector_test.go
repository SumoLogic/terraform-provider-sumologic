package sumologic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicCollector_basic(t *testing.T) {
	rname := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "sumologic_collector.test"
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCollectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfigBasic(rname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rname),
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
	rname := acctest.RandomWithPrefix("tf-acc-test")
	rdescription := acctest.RandomWithPrefix("tf-acc-test")
	rcategory := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "sumologic_collector.test"
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCollectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfig(rname, rdescription, rcategory),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rname),
					resource.TestCheckResourceAttr(resourceName, "description", rdescription),
					resource.TestCheckResourceAttr(resourceName, "category", rcategory),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Etc/UTC"),
				),
			},
		},
	})
}

func TestAccSumologicCollector_update(t *testing.T) {
	rname := acctest.RandomWithPrefix("tf-acc-test")
	rdescription := acctest.RandomWithPrefix("tf-acc-test")
	rcategory := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "sumologic_collector.test"
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCollectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfig(rname, rdescription, rcategory),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rname),
					resource.TestCheckResourceAttr(resourceName, "description", rdescription),
					resource.TestCheckResourceAttr(resourceName, "category", rcategory),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Etc/UTC"),
				),
			},
			{
				Config: testAccSumologicCollectorConfigUpdate(rname, rdescription, rcategory),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rname),
					resource.TestCheckResourceAttr(resourceName, "description", rdescription),
					resource.TestCheckResourceAttr(resourceName, "category", rcategory),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Europe/Berlin"),
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

func testAccSumologicCollectorConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "sumologic_collector" "test" {
	name = "%s"
}`, name)
}

func testAccSumologicCollectorConfig(name, description, category string) string {
	return fmt.Sprintf(`
resource "sumologic_collector" "test" {
	name = "%s"
	description = "%s"
	category = "%s"
	timezone = "Etc/UTC"
}`, name, description, category)
}

func testAccSumologicCollectorConfigUpdate(name, description, category string) string {
	return fmt.Sprintf(`
resource "sumologic_collector" "test" {
	name = "%s"
	description = "%s"
	category = "%s"
	timezone = "Europe/Berlin"
}`, name, description, category)
}
