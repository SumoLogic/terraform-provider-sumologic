package sumologic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicRumSource_Mincreate(t *testing.T) {
	var rumSource RumSource
	var collector Collector

	cName, cDescription, cCategory := getRandomizedParams()
	sName, sDescription, sCategory := getRandomizedParams()

	rumSourceResourceName := "sumologic_rum_source.testRumSource"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRumSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMinRumSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRumSourceExists(rumSourceResourceName, &rumSource),
					testAccCheckRumSourceBasicValues(&rumSource, sName, sDescription, sCategory),
					testAccCheckCollectorExists("sumologic_collector.testCollector", &collector),
					testAccCheckCollectorValues(&collector, cName, cDescription, cCategory, "Etc/UTC", ""),
					resource.TestCheckResourceAttrSet(rumSourceResourceName, "id"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "name", sName),
					resource.TestCheckResourceAttr(rumSourceResourceName, "description", sDescription),
					resource.TestCheckResourceAttr(rumSourceResourceName, "category", sCategory),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.service_name", "some_service2"),
				),
			},
		},
	})
}

func TestAccSumologicRumSource_Fullcreate(t *testing.T) {
	var rumSource RumSource
	var collector Collector

	cName, cDescription, cCategory := getRandomizedParams()
	sName, sDescription, sCategory := getRandomizedParams()

	rumSourceResourceName := "sumologic_rum_source.testRumSource"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRumSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicFullRumSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRumSourceExists(rumSourceResourceName, &rumSource),
					testAccCheckRumSourceBasicValues(&rumSource, sName, sDescription, sCategory),
					testAccCheckCollectorExists("sumologic_collector.testCollector", &collector),
					testAccCheckCollectorValues(&collector, cName, cDescription, cCategory, "Etc/UTC", ""),
					resource.TestCheckResourceAttrSet(rumSourceResourceName, "id"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "name", sName),
					resource.TestCheckResourceAttr(rumSourceResourceName, "description", sDescription),
					resource.TestCheckResourceAttr(rumSourceResourceName, "category", sCategory),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.application_name", "some_application"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.service_name", "some_service"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.deployment_environment", "some_environment"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.sampling_rate", "0.5"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.ignore_urls.#", "2"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.ignore_urls.0", "asd.com"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.ignore_urls.1", "dsa.com"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.custom_tags.%", "1"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.custom_tags.some_tag", "some_value"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.propagate_trace_header_cors_urls.#", "2"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.propagate_trace_header_cors_urls.0", "xyz.com"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.propagate_trace_header_cors_urls.1", "zyx.com"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.selected_country", "Poland"),
				),
			},
		},
	})
}

func TestAccSumologicRumSource_update(t *testing.T) {
	var rumSource RumSource

	cName, cDescription, cCategory := getRandomizedParams()
	sName, sDescription, sCategory := getRandomizedParams()
	sNameUpdated, sDescriptionUpdated, sCategoryUpdated := getRandomizedParams()

	rumSourceResourceName := "sumologic_rum_source.testRumSource"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRumSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMinRumSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRumSourceExists(rumSourceResourceName, &rumSource),
					testAccCheckRumSourceBasicValues(&rumSource, sName, sDescription, sCategory),
					resource.TestCheckResourceAttrSet(rumSourceResourceName, "id"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "name", sName),
					resource.TestCheckResourceAttr(rumSourceResourceName, "description", sDescription),
					resource.TestCheckResourceAttr(rumSourceResourceName, "category", sCategory),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.service_name", "some_service2"),
				),
			},
			{
				Config: testAccSumologicFullRumSourceConfig(cName, cDescription, cCategory, sNameUpdated, sDescriptionUpdated, sCategoryUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRumSourceExists(rumSourceResourceName, &rumSource),
					testAccCheckRumSourceBasicValues(&rumSource, sNameUpdated, sDescriptionUpdated, sCategoryUpdated),
					resource.TestCheckResourceAttrSet(rumSourceResourceName, "id"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "name", sNameUpdated),
					resource.TestCheckResourceAttr(rumSourceResourceName, "description", sDescriptionUpdated),
					resource.TestCheckResourceAttr(rumSourceResourceName, "category", sCategoryUpdated),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.application_name", "some_application"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.service_name", "some_service"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.deployment_environment", "some_environment"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.sampling_rate", "0.5"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.ignore_urls.#", "2"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.ignore_urls.0", "asd.com"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.ignore_urls.1", "dsa.com"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.custom_tags.%", "1"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.custom_tags.some_tag", "some_value"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.propagate_trace_header_cors_urls.#", "2"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.propagate_trace_header_cors_urls.0", "xyz.com"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.propagate_trace_header_cors_urls.1", "zyx.com"),
					resource.TestCheckResourceAttr(rumSourceResourceName, "path.0.selected_country", "Poland"),
				),
			},
		},
	})
}

func testAccCheckRumSourceExists(n string, rumSource *RumSource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("RumSource ID is not set")
		}
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("RumSource ID should be int; got %s", rs.Primary.ID)
		}
		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		c := testAccProvider.Meta().(*Client)
		rumSourceResp, err := c.GetRumSource(collectorID, id)
		if err != nil {
			return err
		}
		*rumSource = *rumSourceResp
		return nil
	}
}

func testAccCheckRumSourceBasicValues(rumSource *RumSource, name, description, category string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if rumSource.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, rumSource.Name)
		}
		if rumSource.Description != description {
			return fmt.Errorf("bad description, expected \"%s\", got: %#v", description, rumSource.Description)
		}
		if rumSource.Category != category {
			return fmt.Errorf("bad category, expected \"%s\", got: %#v", category, rumSource.Category)
		}
		return nil
	}
}

func testAccCheckRumSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_rum_source" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Rum Source destruction check: Rum Source ID is not set")
		}
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		s, err := client.GetRumSource(collectorID, id)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("Rum Source still exists")
		}
	}
	return nil
}

func testAccSumologicFullRumSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory string) string {
	return fmt.Sprintf(`
resource "sumologic_collector" "testCollector" {
	name = "%s"
	description = "%s"
	category = "%s"
}

resource "sumologic_rum_source" "testRumSource" {
	name = "%s"
	description = "%s"
	category = "%s"
	collector_id = "${sumologic_collector.testCollector.id}"
	path {
		application_name = "some_application"
		service_name = "some_service"
		deployment_environment = "some_environment"
		sampling_rate = 0.5
		ignore_urls = ["asd.com", "dsa.com"]
		custom_tags = { some_tag = "some_value" }
		propagate_trace_header_cors_urls = ["xyz.com", "zyx.com"]
		selected_country = "Poland"
	}
}
`, cName, cDescription, cCategory, sName, sDescription, sCategory)
}

func testAccSumologicMinRumSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory string) string {
	return fmt.Sprintf(`
resource "sumologic_collector" "testCollector" {
	name = "%s"
	description = "%s"
	category = "%s"
}

resource "sumologic_rum_source" "testRumSource" {
	name = "%s"
	description = "%s"
	category = "%s"
	collector_id = "${sumologic_collector.testCollector.id}"
	path {
		service_name = "some_service2"
	}
}
`, cName, cDescription, cCategory, sName, sDescription, sCategory)
}
