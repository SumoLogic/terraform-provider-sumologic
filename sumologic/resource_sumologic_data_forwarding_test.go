package sumologic

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

//func getRandomizedParams() (string, string, string) {
//	name := acctest.RandomWithPrefix("tf-acc-test")
//	description := acctest.RandomWithPrefix("tf-acc-test")
//	category := acctest.RandomWithPrefix("tf-acc-test")
//	return name, description, category }

func TestAccSumologicDataForwarding_create(t *testing.T) {
	//var dataForwarding DataForwarding
	testAwsAccessKey := ""
	testAwsSecretKey := ""

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicDataForwardingConfig(testAwsAccessKey, testAwsSecretKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataForwardingExists(),
					resource.TestCheckResourceAttr("sumologic_data_forwarding.test", "destination_name", "abc")),
			},
		},
	})

}

func testAccCheckDataForwardingExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			if _, err := client.GetPartition(id); err != nil {
				return fmt.Errorf("Received an error retrieving data forwarding %s", err)
			}
		}
		return nil
	}
}

func testAccSumologicDataForwardingConfig(testAwsAccessKey string, testAwsSecretKey string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test_bucket" {
  bucket = "S3-tests-data-forwarding"
}

resource "aws_iam_role" "sumo_role" {
  name = "SumoRole"
  assume_role_policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Principal": {
          "AWS": "arn:aws:iam::246946804217:root"
        },
        "Action": "sts:AssumeRole",
        "Condition": {
          "StringEquals": {
            "sts:ExternalId": "long:0000000000000005"
          }
        }
      }
    ]
  })
  path = "/"
}

resource "aws_iam_policy" "sumo_policy" {
  name   = "SumoPolicy"
  policy = jsonencode({
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:PutObject"],
      "Resource": ["arn:aws:s3:::${aws_s3_bucket.test_bucket.bucket}/*"]
    }
  ]
})
}

resource "aws_iam_role_policy_attachment" "sumo_policy_attach" {
  role       = aws_iam_role.sumo_role.name
  policy_arn = aws_iam_policy.sumo_policy.arn
}

output "sumo_role_arn" {
  description = "ARN of the created role. Copy this ARN back to Sumo to complete the data forwarding destination creation process."
  value       = aws_iam_role.sumo_role.arn
}

resource "sumologic_data_forwarding" "test" {
	destination_name = "abc"
	bucket_name = "${aws_s3_bucket.test_bucket.bucket}"
	authentication {
		type = "RoleBased"
		role_arn = "${aws_iam_role_policy_attachment.sumo_policy_attach.policy_arn}"
	}
	destination_type = "temp"	
}
`)
}
