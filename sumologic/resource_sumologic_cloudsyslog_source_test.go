package sumologic

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccSumologicCloudsyslogSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCloudsyslogSourceConfig,
			},
		}})
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
