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
	var collector Collector
	rname := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "sumologic_collector.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCollectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfigBasic(rname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCollectorExists(resourceName, &collector),
					testAccCheckCollectorValues(&collector, rname, "", "", "Etc/UTC", ""),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rname),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "category", ""),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Etc/UTC"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSumologicCollector_create(t *testing.T) {
	var collector Collector
	rname := acctest.RandomWithPrefix("tf-acc-test")
	rdescription := acctest.RandomWithPrefix("tf-acc-test")
	rcategory := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "sumologic_collector.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCollectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfig(rname, rdescription, rcategory),
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the collector
					testAccCheckCollectorExists(resourceName, &collector),
					// verify remote values
					testAccCheckCollectorValues(&collector, rname, rdescription, rcategory, "Etc/UTC", ""),
					// verify local values
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
	var collector Collector
	rname := acctest.RandomWithPrefix("tf-acc-test")
	rdescription := acctest.RandomWithPrefix("tf-acc-test")
	rcategory := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "sumologic_collector.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCollectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfig(rname, rdescription, rcategory),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCollectorExists(resourceName, &collector),
					testAccCheckCollectorValues(&collector, rname, rdescription, rcategory, "Etc/UTC", ""),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rname),
					resource.TestCheckResourceAttr(resourceName, "description", rdescription),
					resource.TestCheckResourceAttr(resourceName, "category", rcategory),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Etc/UTC"),
				),
			},
			{
				Config: testAccSumologicCollectorConfigUpdate(rname, rdescription, rcategory),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCollectorExists(resourceName, &collector),
					testAccCheckCollectorValues(&collector, rname, rdescription, rcategory, "Europe/Berlin", ""),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
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

func testAccCheckCollectorExists(n string, collector *Collector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("collector ID is not set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("collector id should be int; got %s", rs.Primary.ID)
		}
		c := testAccProvider.Meta().(*Client)
		collectorResp, err := c.GetCollector(id)
		if err != nil {
			return err
		}

		*collector = *collectorResp

		return nil
	}
}

func testAccCheckCollectorValues(collector *Collector, name, description, category, timezone, budgetValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if collector.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, collector.Name)
		}
		if collector.Description != description {
			return fmt.Errorf("bad description, expected \"%s\", got: %#v", description, collector.Description)
		}
		if collector.Category != category {
			return fmt.Errorf("bad category, expected \"%s\", got: %#v", category, collector.Category)
		}
		if collector.TimeZone != timezone {
			return fmt.Errorf("bad timezone, expected \"%s\", got: %#v", timezone, collector.TimeZone)
		}
		if value, ok := collector.Fields["_budget"]; ok {
			if value != budgetValue {
				return fmt.Errorf("bad budgetValue, expected \"%s\", got: %#v", budgetValue, collector.Fields["_budget"])
			}
		}
		return nil
	}
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
