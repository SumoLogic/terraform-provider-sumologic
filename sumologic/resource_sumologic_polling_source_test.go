package sumologic

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccSumologicPollingSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicPollingSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("sumologic_polling_source.s3_audit", "content_type"),
				),
			},
		}})
}

var testAccSumologicPollingSourceConfig = `
resource "sumologic_collector" "AWS" {
  name = "AWS2"
  description = "AWS logs"
  category = "category"
}

resource "sumologic_polling_source" "s3_audit" {
  collector_id = "${sumologic_collector.AWS.id}"
  name = "Amazon S3 Audit"
  description = "test_desc"
  category = "some/category"
  content_type = "AwsS3AuditBucket"
  scan_interval = 1
  paused = false

  authentication {
    access_key = "AKIAIOSFODNN7EXAMPLE"
    secret_key = "******"
  }

  path {
    bucket_name = "Bucket1"
    path_expression = "*"
  }
}
`
