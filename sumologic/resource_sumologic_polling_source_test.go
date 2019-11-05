package sumologic

// func TestAccSumologicPollingSource(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		Providers: testAccProviders,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccSumologicPollingSourceConfig,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttrSet("sumologic_polling_source.s3_audit", "content_type"),
// 				),
// 			},
// 		}})
// }

var testAccSumologicPollingSourceConfig = `
variable "aws_accessid" {}
variable "aws_accesskey" {}

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
  content_type = "AwsS3Bucket"
  scan_interval = 1000
  paused = false

  authentication {
	type	   = "S3BucketAuthentication"
    access_key = "${var.aws_accessid}"
    secret_key = "${var.aws_accesskey}"
  }

  path {
    bucket_name = "Bucket1"
    path_expression = "AWSLogs/*"
  }
}
`
