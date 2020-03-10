package sumologic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicCloudsyslogSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCloudsyslogSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudSyslogSourceExists("sumologic_cloudsyslog_source.cloudsyslog", t),
					resource.TestCheckResourceAttrSet("sumologic_cloudsyslog_source.cloudsyslog", "id"),
				),
			},
		}})
}

func testAccCheckCloudSyslogSourceExists(n string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
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
