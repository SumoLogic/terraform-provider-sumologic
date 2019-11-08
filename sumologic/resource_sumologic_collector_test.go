package sumologic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicCollectorMinimal(t *testing.T) {
	var collector *Collector
	resourceName := "sumologic_collector.test"
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCollectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfigMinimal,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCollectorExists(resourceName, &collector, t),
					testAccCheckCollectorAttributes(resourceName, &collector),
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

func TestAccSumologicCollectorSimple(t *testing.T) {
	var collector *Collector
	resourceName := "sumologic_collector.test"
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCollectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCollectorExists(resourceName, &collector, t),
					testAccCheckCollectorAttributes(resourceName, &collector),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "MyTerraformCollector2"),
					resource.TestCheckResourceAttr(resourceName, "description", "MyCollectorDesc"),
					resource.TestCheckResourceAttr(resourceName, "category", "Cat"),
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

// func TestAccSumologicCollectorLookupByName(t *testing.T) {
// 	var collector *Collector
// 	resourceName := "sumologic_collector.test"
// 	resource.Test(t, resource.TestCase{
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckCollectorDestroy,
// 		// TODO: if we keep lookup_by_name, we need to beef up the tests and have 2 steps
// 		// TODO: first step creates the resource
// 		// TODO: second step looks it up by name
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccSumologicCollectorConfigLookupByName,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckCollectorExists(resourceName, &collector, t),
// 					testAccCheckCollectorAttributes(resourceName, &collector),
// 					resource.TestCheckResourceAttrSet(resourceName, "id"),
// 				),
// 			},
// 			{
// 				ResourceName:            resourceName,
// 				ImportState:             true,
// 				ImportStateVerify:       true,
// 				ImportStateVerifyIgnore: []string{"lookup_by_name", "destroy"},
// 			},
// 		},
// 	})
// }

// func TestAccSumologicCollectorAllConfig(t *testing.T) {
// 	var collector *Collector
// 	resourceName := "sumologic_collector.test"
// 	resource.Test(t, resource.TestCase{
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckCollectorDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccSumologicCollectorConfigAll,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckCollectorExists(resourceName, &collector, t),
// 					testAccCheckCollectorAttributes(resourceName, &collector),
// 					resource.TestCheckResourceAttrSet(resourceName, "id"),
// 					resource.TestCheckResourceAttr(resourceName, "name", "CollectorName"),
// 					resource.TestCheckResourceAttr(resourceName, "description", "CollectorDesc"),
// 					resource.TestCheckResourceAttr(resourceName, "category", "Category"),
// 					resource.TestCheckResourceAttr(resourceName, "timezone", "Europe/Berlin"),
// 				),
// 			},
// 			{
// 				ResourceName:            resourceName,
// 				ImportState:             true,
// 				ImportStateVerify:       true,
// 				ImportStateVerifyIgnore: []string{"lookup_by_name", "destroy"},
// 			},
// 		},
// 	})
// }

func TestAccSumologicCollectorChangeConfig(t *testing.T) {
	var collector *Collector
	resourceName := "sumologic_collector.test"
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCollectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCollectorExists(resourceName, &collector, t),
					testAccCheckCollectorAttributes(resourceName, &collector),
					resource.TestCheckResourceAttr(resourceName, "name", "MyTerraformCollector2"),
					resource.TestCheckResourceAttr(resourceName, "description", "MyCollectorDesc"),
					resource.TestCheckResourceAttr(resourceName, "category", "Cat"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Etc/UTC"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lookup_by_name", "destroy"},
			},
			{
				Config: testAccSumologicCollectorConfigAll,
				Check: resource.ComposeTestCheckFunc(
					// check the id of this resource is the same as the one in the previous step
					testAccCheckCollectorId(resourceName, &collector),
					testAccCheckCollectorExists(resourceName, &collector, t),
					testAccCheckCollectorAttributes(resourceName, &collector),
					resource.TestCheckResourceAttr(resourceName, "name", "CollectorName"),
					resource.TestCheckResourceAttr(resourceName, "description", "CollectorDesc"),
					resource.TestCheckResourceAttr(resourceName, "category", "Category"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "Europe/Berlin"),
				),
			},
		},
	})
}

func TestAccSumologicCollectorManualDeletion(t *testing.T) {
	var collector *Collector

	deleteCollector := func() {
		c := testAccProvider.Meta().(*Client)
		_, err := c.GetCollector(collector.ID)
		if err != nil {
			t.Fatal(fmt.Sprintf("attempted to delete collector %d but it does not exist (%s)", collector.ID, err))
		}
		err = c.DeleteCollector(collector.ID)
		if err != nil {
			t.Fatal(fmt.Sprintf("failed to delete collector %d (%s)", collector.ID, err))
		}
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCollectorConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCollectorExists("sumologic_collector.test", &collector, t),
					resource.TestCheckResourceAttr("sumologic_collector.test", "name", "MyTerraformCollector2"),
					resource.TestCheckResourceAttr("sumologic_collector.test", "description", "MyCollectorDesc"),
					resource.TestCheckResourceAttr("sumologic_collector.test", "category", "Cat"),
					resource.TestCheckResourceAttr("sumologic_collector.test", "timezone", "Etc/UTC"),
				),
			},
			{
				PreConfig: deleteCollector, // simulate a manual deletion by deleting the collector between the 2 applies
				Config:    testAccSumologicCollectorConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCollectorExists("sumologic_collector.test", &collector, t),
					resource.TestCheckResourceAttr("sumologic_collector.test", "name", "MyTerraformCollector2"),
					resource.TestCheckResourceAttr("sumologic_collector.test", "description", "MyCollectorDesc"),
					resource.TestCheckResourceAttr("sumologic_collector.test", "category", "Cat"),
					resource.TestCheckResourceAttr("sumologic_collector.test", "timezone", "Etc/UTC"),
				),
			},
		},
	})
}

// TODO: if we keep the collector's destroy attribute we need to include a test checking if destroy=false works as expected

// Returns a function checking that the collector with the id from the state file has an expected id.
// The expected id is specified in the collector passed as parameter
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

// Returns a function checking that the collector with the id from the state exists.
// If the collecor exists, its attributes are updated in *collector
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
			return fmt.Errorf("collector %d not found", id)
		}
		return nil
	}
}

// Returns a function checking that the attributes in the state match that attributes of the actual resource created
func testAccCheckCollectorAttributes(name string, expected **Collector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(name, "name", (**expected).Name),
			resource.TestCheckResourceAttr(name, "description", (**expected).Description),
			resource.TestCheckResourceAttr(name, "category", (**expected).Category),
			resource.TestCheckResourceAttr(name, "timezone", (**expected).TimeZone),
		)
		return f(s)
	}
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
