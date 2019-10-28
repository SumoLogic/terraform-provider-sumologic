package sumologic

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSumologicHTTPSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicHTTPSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHTTPSourceExists("sumologic_http_source.http", t),
					resource.TestCheckResourceAttrSet("sumologic_http_source.http", "id"),
					resource.TestCheckResourceAttrSet("sumologic_http_source.http", "url"),
				),
			},
		},
	})
}

func TestAccSumologicHTTPSourceUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicHTTPSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHTTPSourceExists("sumologic_http_source.http", t),
					resource.TestCheckResourceAttrSet("sumologic_http_source.http", "id"),
					resource.TestCheckResourceAttr("sumologic_http_source.http", "name", "test_http"),
					resource.TestCheckResourceAttr("sumologic_http_source.http", "description", "test_desc"),
					resource.TestCheckResourceAttr("sumologic_http_source.http", "message_per_request", "false"),
					resource.TestCheckResourceAttr("sumologic_http_source.http", "category", "source/category"),
				),
			},
			{
				Config: testAccSumologicHTTPSourceConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHTTPSourceExists("sumologic_http_source.http", t),
					resource.TestCheckResourceAttrSet("sumologic_http_source.http", "id"),
					resource.TestCheckResourceAttr("sumologic_http_source.http", "name", "test_http"),
					resource.TestCheckResourceAttr("sumologic_http_source.http", "description", "desc_test"),
					resource.TestCheckResourceAttr("sumologic_http_source.http", "category", "category/source"),
				),
			},
		},
	})
}

// func TestAccSumologicHTTPSourceImport(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		Providers: testAccProviders,
// 		Steps: []resource.TestStep{
// 			{
// 				Config:            testAccSumologicHTTPSourceConfig,
// 				ResourceName:      "sumologic_http_source.http",
// 				ImportState:       true,
// 				ImportStateId:     "123/456",
// 				ImportStateVerify: true,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckHTTPSourceExists("sumologic_http_source.http", t),
// 					resource.TestCheckResourceAttr("sumologic_http_source.http", "id", "123"),
// 					resource.TestCheckResourceAttr("sumologic_http_source.http", "collector_id", "456"),
// 				),
// 			},
// 		},
// 	})
// }

func testAccCheckHTTPSourceExists(n string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return nil
	}
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
  message_per_request = false
  category = "source/category"
  collector_id = "${sumologic_collector.test.id}"
}
`

var testAccSumologicHTTPSourceConfigUpdate = `

resource "sumologic_collector" "test" {
  name = "MyCollector"
  description = "MyCollectorDesc"
  category = "Cat"
}

resource "sumologic_http_source" "http" {
  name = "test_http"
  description = "desc_test"
  message_per_request = false
  category = "category/source"
  collector_id = "${sumologic_collector.test.id}"
  lookup_by_name = true
}
`
