package sumologic

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicCloudsyslogSource_basic(t *testing.T) {
	var cloudsyslogSource CloudSyslogSource
	resourceName := "sumologic_cloudsyslog_source.cloudsyslog"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudsyslogSourceDestroy(resourceName, cloudsyslogSource),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCloudsyslogSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudSyslogSourceExists(resourceName, &cloudsyslogSource, t),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "test_cloudsyslog"),
					resource.TestCheckResourceAttr(resourceName, "description", "test_desc"),
					resource.TestCheckResourceAttr(resourceName, "category", "source/category"),
				),
			},
			{
				Config: testAccSumologicCloudsyslogSourceUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudSyslogSourceExists(resourceName, &cloudsyslogSource, t),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "test_cloudsyslog_update"),
					resource.TestCheckResourceAttr(resourceName, "description", "test_desc_update"),
					resource.TestCheckResourceAttr(resourceName, "category", "source/category"),
				),
			},
		}})
}

func testAccCheckCloudsyslogSourceDestroy(name string, cloudsyslogSource CloudSyslogSource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			//need this so that we don't get an unused import error for strconv in some cases
			return fmt.Errorf("Error = %s. CloudsyslogSource not found: %s", strconv.FormatBool(ok), name)
		}
		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		u, err := client.GetCloudSyslogSource(collectorID, cloudsyslogSource.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if u != nil {
			return fmt.Errorf("FieldExtractionRule still exists")
		}
		return nil
	}
}

func testAccCheckCloudSyslogSourceExists(name string, cloudsyslogSource *CloudSyslogSource, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			//need this so that we don't get an unused import error for strconv in some cases
			return fmt.Errorf("Error = %s. CloudsyslogSource not found: %s", strconv.FormatBool(ok), name)
		}

		//need this so that we don't get an unused import error for strings in some cases
		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("CloudsyslogSource, ID is not set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		c := testAccProvider.Meta().(*Client)
		newCloudsyslogSource, err := c.GetCloudSyslogSource(collectorID, id)
		if err != nil {
			return fmt.Errorf("CloudsyslogSource, %v not found", id)
		}
		cloudsyslogSource = newCloudsyslogSource
		return nil
	}
}

var testAccSumologicCloudsyslogSourceConfig = `

resource "sumologic_collector" "test" {
  name = "MyCollector"
  description = "MyCollectorDesc"
  category = "Cat"
}

resource "sumologic_cloudsyslog_source" "cloudsyslog" {
  name = "test_cloudsyslog"
  description = "test_desc"
  category = "source/category"
  collector_id = "${sumologic_collector.test.id}"
}
`

var testAccSumologicCloudsyslogSourceUpdateConfig = `

resource "sumologic_collector" "test" {
  name = "MyCollector"
  description = "MyCollectorDesc"
  category = "Cat"
}

resource "sumologic_cloudsyslog_source" "cloudsyslog" {
  name = "test_cloudsyslog_update"
  description = "test_desc_update"
  category = "source/category"
  collector_id = "${sumologic_collector.test.id}"
}
`
