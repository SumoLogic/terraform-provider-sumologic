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
	var kinesisLogSource HTTPSource
	var httpOtlpSource HTTPSource
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
	kName := acctest.RandomWithPrefix("tf-acc-test")
	kDescription := acctest.RandomWithPrefix("tf-acc-test")
	kCategory := acctest.RandomWithPrefix("tf-acc-test")
	oName := acctest.RandomWithPrefix("tf-acc-test")
	oDescription := acctest.RandomWithPrefix("tf-acc-test")
	oCategory := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "sumologic_http_source.http"
	tracingResourceName := "sumologic_http_source.traces"
	kinesisResourceName := "sumologic_http_source.kinesisLog"
	otlpResourceName := "sumologic_http_source.otlp"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHTTPSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicHTTPSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, tName, tDescription, tCategory, kName, kDescription, kCategory, oName, oDescription, oCategory),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHTTPSourceExists(resourceName, &httpSource),
					testAccCheckHTTPSourceValues(&httpSource, sName, sDescription, sCategory),
					testAccCheckCollectorExists("sumologic_collector.test", &collector),
					testAccCheckCollectorValues(&collector, cName, cDescription, cCategory, "Etc/UTC", ""),
					testAccCheckHTTPSourceExists(tracingResourceName, &httpTraceSource),
					testAccCheckHTTPSourceValues(&httpTraceSource, tName, tDescription, tCategory),
					testAccCheckHTTPSourceExists(kinesisResourceName, &kinesisLogSource),
					testAccCheckHTTPSourceValues(&kinesisLogSource, kName, kDescription, kCategory),
					testAccCheckHTTPSourceExists(otlpResourceName, &httpOtlpSource),
					testAccCheckHTTPSourceValues(&httpOtlpSource, oName, oDescription, oCategory),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
					resource.TestCheckResourceAttr(resourceName, "name", sName),
					resource.TestCheckResourceAttr(resourceName, "description", sDescription),
					resource.TestCheckResourceAttr(resourceName, "message_per_request", "false"),
					resource.TestCheckResourceAttr(resourceName, "category", sCategory),
					resource.TestCheckResourceAttrSet(tracingResourceName, "id"),
					resource.TestCheckResourceAttrSet(tracingResourceName, "url"),
					resource.TestCheckResourceAttr(tracingResourceName, "name", tName),
					resource.TestCheckResourceAttr(tracingResourceName, "description", tDescription),
					resource.TestCheckResourceAttr(tracingResourceName, "category", tCategory),
					resource.TestCheckResourceAttr(tracingResourceName, "content_type", "Zipkin"),
					resource.TestCheckResourceAttrSet(kinesisResourceName, "id"),
					resource.TestCheckResourceAttrSet(kinesisResourceName, "url"),
					resource.TestCheckResourceAttr(kinesisResourceName, "name", kName),
					resource.TestCheckResourceAttr(kinesisResourceName, "description", kDescription),
					resource.TestCheckResourceAttr(kinesisResourceName, "category", kCategory),
					resource.TestCheckResourceAttr(kinesisResourceName, "content_type", "KinesisLog"),
					resource.TestCheckResourceAttrSet(otlpResourceName, "id"),
					resource.TestCheckResourceAttrSet(otlpResourceName, "url"),
					resource.TestCheckResourceAttr(otlpResourceName, "name", oName),
					resource.TestCheckResourceAttr(otlpResourceName, "description", oDescription),
					resource.TestCheckResourceAttr(otlpResourceName, "category", oCategory),
					resource.TestCheckResourceAttr(otlpResourceName, "content_type", "Otlp"),
				),
			},
		},
	})
}

func TestAccSumologicHTTPSource_update(t *testing.T) {
	var httpSource HTTPSource
	var httpTraceSource HTTPSource
	var kinesisLogSource HTTPSource
	var httpOtlpSource HTTPSource
	cName := acctest.RandomWithPrefix("tf-acc-test")
	cDescription := acctest.RandomWithPrefix("tf-acc-test")
	cCategory := acctest.RandomWithPrefix("tf-acc-test")
	sName := acctest.RandomWithPrefix("tf-acc-test")
	sDescription := acctest.RandomWithPrefix("tf-acc-test")
	sCategory := acctest.RandomWithPrefix("tf-acc-test")
	tName := acctest.RandomWithPrefix("tf-acc-test")
	tDescription := acctest.RandomWithPrefix("tf-acc-test")
	tCategory := acctest.RandomWithPrefix("tf-acc-test")
	kName := acctest.RandomWithPrefix("tf-acc-test")
	kDescription := acctest.RandomWithPrefix("tf-acc-test")
	kCategory := acctest.RandomWithPrefix("tf-acc-test")
	oName := acctest.RandomWithPrefix("tf-acc-test")
	oDescription := acctest.RandomWithPrefix("tf-acc-test")
	oCategory := acctest.RandomWithPrefix("tf-acc-test")
	sNameUpdated := acctest.RandomWithPrefix("tf-acc-test")
	sDescriptionUpdated := acctest.RandomWithPrefix("tf-acc-test")
	sCategoryUpdated := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "sumologic_http_source.http"
	tracingResourceName := "sumologic_http_source.traces"
	kinesisResourceName := "sumologic_http_source.kinesisLog"
	otlpResourceName := "sumologic_http_source.otlp"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHTTPSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicHTTPSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, tName, tDescription, tCategory, kName, kDescription, kCategory, oName, oDescription, oCategory),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHTTPSourceExists(resourceName, &httpSource),
					testAccCheckHTTPSourceValues(&httpSource, sName, sDescription, sCategory),
					testAccCheckHTTPSourceExists(tracingResourceName, &httpTraceSource),
					testAccCheckHTTPSourceValues(&httpTraceSource, tName, tDescription, tCategory),
					testAccCheckHTTPSourceExists(kinesisResourceName, &kinesisLogSource),
					testAccCheckHTTPSourceValues(&kinesisLogSource, kName, kDescription, kCategory),
					testAccCheckHTTPSourceExists(otlpResourceName, &httpOtlpSource),
					testAccCheckHTTPSourceValues(&httpOtlpSource, oName, oDescription, oCategory),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
					resource.TestCheckResourceAttr(resourceName, "name", sName),
					resource.TestCheckResourceAttr(resourceName, "description", sDescription),
					resource.TestCheckResourceAttr(resourceName, "message_per_request", "false"),
					resource.TestCheckResourceAttr(resourceName, "category", sCategory),
					resource.TestCheckResourceAttrSet(tracingResourceName, "id"),
					resource.TestCheckResourceAttrSet(tracingResourceName, "url"),
					resource.TestCheckResourceAttr(tracingResourceName, "name", tName),
					resource.TestCheckResourceAttr(tracingResourceName, "description", tDescription),
					resource.TestCheckResourceAttr(tracingResourceName, "category", tCategory),
					resource.TestCheckResourceAttr(tracingResourceName, "content_type", "Zipkin"),
					resource.TestCheckResourceAttrSet(kinesisResourceName, "id"),
					resource.TestCheckResourceAttrSet(kinesisResourceName, "url"),
					resource.TestCheckResourceAttr(kinesisResourceName, "name", kName),
					resource.TestCheckResourceAttr(kinesisResourceName, "description", kDescription),
					resource.TestCheckResourceAttr(kinesisResourceName, "category", kCategory),
					resource.TestCheckResourceAttr(kinesisResourceName, "content_type", "KinesisLog"),
					resource.TestCheckResourceAttrSet(otlpResourceName, "id"),
					resource.TestCheckResourceAttrSet(otlpResourceName, "url"),
					resource.TestCheckResourceAttr(otlpResourceName, "name", oName),
					resource.TestCheckResourceAttr(otlpResourceName, "description", oDescription),
					resource.TestCheckResourceAttr(otlpResourceName, "category", oCategory),
					resource.TestCheckResourceAttr(otlpResourceName, "content_type", "Otlp"),
				),
			},
			{
				Config: testAccSumologicHTTPSourceConfig(cName, cDescription, cCategory, sNameUpdated, sDescriptionUpdated, sCategoryUpdated, tName, tDescription, tCategory, kName, kDescription, kCategory, oName, oDescription, oCategory),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHTTPSourceExists(resourceName, &httpSource),
					testAccCheckHTTPSourceValues(&httpSource, sNameUpdated, sDescriptionUpdated, sCategoryUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
					resource.TestCheckResourceAttr(resourceName, "name", sNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", sDescriptionUpdated),
					resource.TestCheckResourceAttr(resourceName, "category", sCategoryUpdated),
					resource.TestCheckResourceAttr(tracingResourceName, "content_type", "Zipkin"),
					resource.TestCheckResourceAttr(kinesisResourceName, "content_type", "KinesisLog"),
					resource.TestCheckResourceAttr(otlpResourceName, "content_type", "Otlp"),
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

func testAccSumologicHTTPSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, tName, tDescription, tCategory, kName, kDescription, kCategory, oName, oDescription, oCategory string) string {
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

resource "sumologic_http_source" "kinesisLog" {
	name = "%s"
	description = "%s"
	category = "%s"
	content_type = "KinesisLog"
	collector_id = "${sumologic_collector.test.id}"
}

resource "sumologic_http_source" "otlp" {
	name = "%s"
	description = "%s"
	category = "%s"
	content_type = "Otlp"
	collector_id = "${sumologic_collector.test.id}"
}
`, cName, cDescription, cCategory, sName, sDescription, sCategory, tName, tDescription, tCategory, kName, kDescription, kCategory, oName, oDescription, oCategory)
}
