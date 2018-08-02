package sumologic

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccSumologicHTTPSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicHTTPSourceConfig,
			},
		}})
}

var testAccSumologicHTTPSourceConfig = `

resource "sumologic_collector" "test" {
  name = "MyCollector"
  description = "MyCollectorDesc"
  category = "Cat"
}

resource "sumologic_http_source" "http" {
  name = "test_http"
  description = "test_desc"
  message_per_request = true
  category = "source/category"
  collector_id = "${sumologic_collector.test.id}"
}
`
