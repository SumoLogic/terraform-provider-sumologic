package sumologic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicHTTPSource_create(t *testing.T) {
	var httpSource HTTPSource
	var httpTraceSource HTTPSource
	var collector Collector
	cName := acctest.RandomWithPrefix("tf-acc-test")
	cDescription := acctest.RandomWithPrefix("tf-acc-test")
	cCategory := acctest.RandomWithPrefix("tf-acc-test")
	sName := acctest.RandomWithPrefix("tf-acc-test")
	sDescription := acctest.RandomWithPrefix("tf-acc-test")
	sCategory := acctest.RandomWithPrefix("tf-acc-test")
	tName := acctest.RandomWithPrefix("tf-acc-test")
	tDescription := acctest.RandomWithPrefix("tf-acc-test")
	tCategory := acctest.RandomWithPrefix("tf-acc-test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHTTPSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicHTTPSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, tName, tDescription, tCategory),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHTTPSourceExists("sumologic_http_source.http", &httpSource),
					testAccCheckHTTPSourceValues(&httpSource, sName, sDescription, sCategory),
					testAccCheckCollectorExists("sumologic_collector.test", &collector),
					testAccCheckCollectorValues(&collector, cName, cDescription, cCategory, "Etc/UTC", ""),
					testAccCheckHTTPSourceExists("sumologic_http_source.traces", &httpTraceSource),
					testAccCheckHTTPSourceValues(&httpTraceSource, tName, tDescription, tCategory),
					resource.TestCheckResourceAttrSet("sumologic_http_source.http", "id"),
					resource.TestCheckResourceAttrSet("sumologic_http_source.http", "url"),
					resource.TestCheckResourceAttr("sumologic_http_source.http", "name", sName),
					resource.TestCheckResourceAttr("sumologic_http_source.http", "description", sDescription),
					resource.TestCheckResourceAttr("sumologic_http_source.http", "message_per_request", "false"),
					resource.TestCheckResourceAttr("sumologic_http_source.http", "category", sCategory),
					resource.TestCheckResourceAttrSet("sumologic_http_source.traces", "id"),
					resource.TestCheckResourceAttrSet("sumologic_http_source.traces", "url"),
					resource.TestCheckResourceAttr("sumologic_http_source.traces", "name", tName),
					resource.TestCheckResourceAttr("sumologic_http_source.traces", "description", tDescription),
					resource.TestCheckResourceAttr("sumologic_http_source.traces", "category", tCategory),
					resource.TestCheckResourceAttr("sumologic_http_source.traces", "content_type", "Zipkin"),
				),
			},
		},
	})
}

func TestAccSumologicHTTPSource_update(t *testing.T) {
	var httpSource HTTPSource
	var httpTraceSource HTTPSource
	cName := acctest.RandomWithPrefix("tf-acc-test")
	cDescription := acctest.RandomWithPrefix("tf-acc-test")
	cCategory := acctest.RandomWithPrefix("tf-acc-test")
	sName := acctest.RandomWithPrefix("tf-acc-test")
	sDescription := acctest.RandomWithPrefix("tf-acc-test")
	sCategory := acctest.RandomWithPrefix("tf-acc-test")
	tName := acctest.RandomWithPrefix("tf-acc-test")
	tDescription := acctest.RandomWithPrefix("tf-acc-test")
	tCategory := acctest.RandomWithPrefix("tf-acc-test")
	sNameUpdated := acctest.RandomWithPrefix("tf-acc-test")
	sDescriptionUpdated := acctest.RandomWithPrefix("tf-acc-test")
	sCategoryUpdated := acctest.RandomWithPrefix("tf-acc-test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHTTPSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicHTTPSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, tName, tDescription, tCategory),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHTTPSourceExists("sumologic_http_source.http", &httpSource),
					testAccCheckHTTPSourceValues(&httpSource, sName, sDescription, sCategory),
					testAccCheckHTTPSourceExists("sumologic_http_source.traces", &httpTraceSource),
					testAccCheckHTTPSourceValues(&httpTraceSource, tName, tDescription, tCategory),
					resource.TestCheckResourceAttrSet("sumologic_http_source.http", "id"),
					resource.TestCheckResourceAttrSet("sumologic_http_source.http", "url"),
					resource.TestCheckResourceAttr("sumologic_http_source.http", "name", sName),
					resource.TestCheckResourceAttr("sumologic_http_source.http", "description", sDescription),
					resource.TestCheckResourceAttr("sumologic_http_source.http", "message_per_request", "false"),
					resource.TestCheckResourceAttr("sumologic_http_source.http", "category", sCategory),
					resource.TestCheckResourceAttrSet("sumologic_http_source.traces", "id"),
					resource.TestCheckResourceAttrSet("sumologic_http_source.traces", "url"),
					resource.TestCheckResourceAttr("sumologic_http_source.traces", "name", tName),
					resource.TestCheckResourceAttr("sumologic_http_source.traces", "description", tDescription),
					resource.TestCheckResourceAttr("sumologic_http_source.traces", "category", tCategory),
					resource.TestCheckResourceAttr("sumologic_http_source.traces", "content_type", "Zipkin"),
				),
			},
			{
				Config: testAccSumologicHTTPSourceConfig(cName, cDescription, cCategory, sNameUpdated, sDescriptionUpdated, sCategoryUpdated, tName, tDescription, tCategory),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHTTPSourceExists("sumologic_http_source.http", &httpSource),
					testAccCheckHTTPSourceValues(&httpSource, sNameUpdated, sDescriptionUpdated, sCategoryUpdated),
					resource.TestCheckResourceAttrSet("sumologic_http_source.http", "id"),
					resource.TestCheckResourceAttrSet("sumologic_http_source.http", "url"),
					resource.TestCheckResourceAttr("sumologic_http_source.http", "name", sNameUpdated),
					resource.TestCheckResourceAttr("sumologic_http_source.http", "description", sDescriptionUpdated),
					resource.TestCheckResourceAttr("sumologic_http_source.http", "category", sCategoryUpdated),
					resource.TestCheckResourceAttr("sumologic_http_source.traces", "content_type", "Zipkin"),
				),
			},
		},
	})
}

func testAccCheckHTTPSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_http_source" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("HTTP Source destruction check: HTTP Source ID is not set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}

		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}

		s, err := client.GetHTTPSource(collectorID, id)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("HTTP Source still exists")
		}
	}
	return nil
}

func testAccCheckHTTPSourceExists(n string, httpSource *HTTPSource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("HTTP Source ID is not set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("HTTP Source id should be int; got %s", rs.Primary.ID)
		}

		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}

		c := testAccProvider.Meta().(*Client)
		httpSourceResp, err := c.GetHTTPSource(collectorID, id)
		if err != nil {
			return err
		}

		*httpSource = *httpSourceResp

		return nil
	}
}

func testAccCheckHTTPSourceValues(httpSource *HTTPSource, name, description, category string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if httpSource.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, httpSource.Name)
		}
		if httpSource.Description != description {
			return fmt.Errorf("bad description, expected \"%s\", got: %#v", description, httpSource.Description)
		}
		if httpSource.Category != category {
			return fmt.Errorf("bad category, expected \"%s\", got: %#v", category, httpSource.Category)
		}
		return nil
	}
}

func testAccSumologicHTTPSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, tName, tDescription, tCategory string) string {
	return fmt.Sprintf(`
resource "sumologic_collector" "test" {
	name = "%s"
	description = "%s"
	category = "%s"
}
	
resource "sumologic_http_source" "http" {
	name = "%s"
	description = "%s"
	message_per_request = false
	category = "%s"
	collector_id = "${sumologic_collector.test.id}"
}

resource "sumologic_http_source" "traces" {
	name = "%s"
	description = "%s"
	category = "%s"
	content_type = "Zipkin"
	collector_id = "${sumologic_collector.test.id}"
}
`, cName, cDescription, cCategory, sName, sDescription, sCategory, tName, tDescription, tCategory)
}
