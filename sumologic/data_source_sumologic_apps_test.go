package sumologic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceSumoLogicApps_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSumoLogicAppsConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSumoLogicAppsDataSourceID("data.sumologic_apps.test"),
					resource.TestCheckResourceAttrSet("data.sumologic_apps.test", "apps.#"),
					checkResourceAttrGreaterThanZero("data.sumologic_apps.test", "apps.#"),
				),
			},
		},
	})
}

func TestAccDataSourceSumoLogicApps_filtered(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSumoLogicAppsConfig_filtered,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSumoLogicAppsDataSourceID("data.sumologic_apps.filtered"),
					resource.TestCheckResourceAttr("data.sumologic_apps.filtered", "apps.#", "1"),
					resource.TestCheckResourceAttr("data.sumologic_apps.filtered", "apps.0.name", "MySQL - OpenTelemetry"),
					resource.TestCheckResourceAttr("data.sumologic_apps.filtered", "apps.0.author", "Sumo Logic"),
				),
			},
		},
	})
}

func testAccCheckSumoLogicAppsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find SumoLogic Apps data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("SumoLogic Apps data source ID not set")
		}
		return nil
	}
}

func checkResourceAttrGreaterThanZero(resourceName, attributeName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		attrValue, ok := rs.Primary.Attributes[attributeName]
		if !ok {
			return fmt.Errorf("Attribute not found: %s", attributeName)
		}

		// Convert the attribute value to an integer
		count, err := strconv.Atoi(attrValue)
		if err != nil {
			return fmt.Errorf("Error converting attribute value to integer: %s", err)
		}

		// Check if count is greater than 0
		if count <= 0 {
			return fmt.Errorf("Expected %s to be greater than 0, got %d", attributeName, count)
		}

		return nil
	}
}

const testAccDataSourceSumoLogicAppsConfig_basic = `
	data "sumologic_apps" "test" {}
`

const testAccDataSourceSumoLogicAppsConfig_filtered = `
	data "sumologic_apps" "filtered" {
		name = "MySQL - OpenTelemetry"
		author = "Sumo Logic"
	}
`
