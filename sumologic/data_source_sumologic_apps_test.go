package sumologic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
					//testAccCheckSumoLogicAppsGreaterThanZero("data.sumologic_apps.test"),
					resource.TestCheckResourceAttrSet("data.sumologic_apps.test", "id"),
					//resource.TestCheckResourceAttrSet("data.sumologic_apps.test", "apps"),
					resource.TestCheckResourceAttrSet("data.sumologic_apps.test", "apps.0.uuid"),
					resource.TestCheckResourceAttrSet("data.sumologic_apps.test", "apps.0.name"),
					resource.TestCheckResourceAttrSet("data.sumologic_apps.test", "apps.0.description"),
					resource.TestCheckResourceAttrSet("data.sumologic_apps.test", "apps.0.latestVersion"),
					resource.TestCheckResourceAttrSet("data.sumologic_apps.test", "apps.0.icon"),
					resource.TestCheckResourceAttrSet("data.sumologic_apps.test", "apps.0.author"),
					resource.TestCheckResourceAttrSet("data.sumologic_apps.test", "apps.0.appType"),
					//resource.TestCheckResourceAttrSet("data.sumologic_apps.test", "apps.0.attributes"),
					resource.TestCheckResourceAttrSet("data.sumologic_apps.test", "apps.0.installable"),
					resource.TestCheckResourceAttrSet("data.sumologic_apps.test", "apps.0.showOnMarketplace"),
				),
			},
		},
	})
}

// func TestAccDataSourceSumoLogicApps_filtered(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheck(t) },
// 		ProviderFactories: providerFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccDataSourceSumoLogicAppsConfig_filtered,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckSumoLogicAppsDataSourceID("data.sumologic_apps.filtered"),
// 					resource.TestCheckResourceAttr("data.sumologic_apps.filtered", "apps.#", "1"),
// 					resource.TestCheckResourceAttr("data.sumologic_apps.filtered", "apps.0.name", "AWS CloudTrail"),
// 					resource.TestCheckResourceAttr("data.sumologic_apps.filtered", "apps.0.author", "Sumo Logic"),
// 				),
// 			},
// 		},
// 	})
// }

func testAccCheckSumoLogicAppsGreaterThanZero(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find SumoLogic Apps data source: %s", n)
		}

		appsCount, ok := rs.Primary.Attributes["apps.#"]
		if !ok {
			return fmt.Errorf("Apps count not found in state")
		}

		count, err := strconv.Atoi(appsCount)
		if err != nil {
			return fmt.Errorf("Failed to parse apps count: %v", err)
		}

		if count <= 0 {
			return fmt.Errorf("No apps returned, expected at least one")
		}

		return nil
	}
}

func testAccCheckSumoLogicAppsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find SumoLogic Apps data source: %s", n)
		}

		fmt.Printf("%v\n", s.RootModule().Resources)
		fmt.Printf("%v\n", rs.Primary)
		if rs.Primary.ID == "" {
			return fmt.Errorf("SumoLogic Apps data source ID not set")
		}
		return nil
	}
}

const testAccDataSourceSumoLogicAppsConfig_basic = `
	data "sumologic_apps" "test" {
		name = "MySQL - OpenTelemetry"
		author = "Sumo Logic"
	}
`

const testAccDataSourceSumoLogicAppsConfig_filtered = `
	data "sumologic_apps" "filtered" {
		name = "AWS CloudTrail"
		author = "Sumo Logic"
	}
`
