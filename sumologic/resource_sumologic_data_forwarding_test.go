package sumologic

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"strings"
	"testing"
)

func getTestParams() (string, string, string, string, string, string) {
	dataForwardingResourceName := "sumologic_s3_data_forwarding_destination.test"
	destinationName, description := getRandomizedDataForwardingParams()
	testAwsRoleArn := os.Getenv("SUMOLOGIC_TEST_ROLE_ARN")
	testAwsBucket := os.Getenv("SUMOLOGIC_TEST_BUCKET_NAME")
	testAwsRegion := os.Getenv("SUMOLOGIC_TEST_AWS_REGION")
	println("AWS Test Bucket: ", testAwsBucket)
	println("AWS Test ARN: ", testAwsRoleArn)
	println("AWS Test Region: ", testAwsRegion)
	return dataForwardingResourceName, destinationName, description, testAwsRegion, testAwsRoleArn, testAwsBucket
}

func getRandomizedDataForwardingParams() (string, string) {
	destinationName := acctest.RandomWithPrefix("tf-acc-test")
	description := acctest.RandomWithPrefix("tf-acc-test")
	return destinationName, description
}

func TestAccSumologicDataForwarding_create(t *testing.T) {
	dataForwardingResourceName, destinationName, description, testAwsRegion, testAwsRoleArn, testAwsBucket := getTestParams()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataForwardingDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicDataForwardingCreateConfig(destinationName, description, testAwsBucket, testAwsRoleArn, testAwsRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataForwardingExists(),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "destination_name", destinationName),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "description", description),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "s3_region", testAwsRegion),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "s3_server_side_encryption", "false"),
				),
			},
		},
	})
}

func TestAccSumologicDataForwarding_read(t *testing.T) {
	dataForwardingResourceName, destinationName, description, testAwsRegion, testAwsRoleArn, testAwsBucket := getTestParams()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataForwardingDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicDataForwardingCreateConfig(destinationName, description, testAwsBucket, testAwsRoleArn, testAwsRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataForwardingExists(),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "destination_name", destinationName),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "description", description),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "bucket_name", testAwsBucket),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "s3_region", testAwsRegion),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "s3_server_side_encryption", "false"),
				),
			},
			{
				ResourceName: dataForwardingResourceName,
				ImportState:  true,
			},
		},
	})
}

func TestAccSumologicDataForwarding_update(t *testing.T) {
	dataForwardingResourceName, destinationName, description, testAwsRegion, testAwsRoleArn, testAwsBucket := getTestParams()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataForwardingDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicDataForwardingCreateConfig(destinationName, description, testAwsBucket, testAwsRoleArn, testAwsRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataForwardingExists(),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "destination_name", destinationName),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "description", description),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "s3_region", testAwsRegion),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "s3_server_side_encryption", "false"),
				),
			}, {
				Config: testAccSumologicDataForwardingUpdateConfig(destinationName, description, testAwsBucket, testAwsRoleArn, testAwsRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataForwardingExists(),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "destination_name", destinationName),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "description", description),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "s3_region", testAwsRegion),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "s3_server_side_encryption", "true"),
				),
			},
		},
	})

}

func TestAccSumologicDataForwarding_delete(t *testing.T) {
	dataForwardingResourceName, destinationName, description, testAwsRegion, testAwsRoleArn, testAwsBucket := getTestParams()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataForwardingDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicDataForwardingCreateConfig(destinationName, description, testAwsBucket, testAwsRoleArn, testAwsRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataForwardingExists(),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "destination_name", destinationName),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "description", description),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "s3_region", testAwsRegion),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "s3_server_side_encryption", "false"),
				),
			}, {
				Config: testAccSumologicDataForwardingDeleteConfig(),
				Check:  resource.ComposeTestCheckFunc(testAccCheckDataForwardingDestroy()),
			},
		},
	})
}

func testAccCheckDataForwardingExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			if _, err := client.getDataForwarding(id); err != nil {
				return fmt.Errorf("Received an error retrieving data forwarding %s", err)
			}
		}
		return nil
	}
}

func testAccCheckDataForwardingDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			p, err := client.getDataForwarding(id)

			if err != nil {
				if strings.Contains(err.Error(), "Data forwarding Destination doesn't exists") {
					continue
				}

				return fmt.Errorf("Encountered an error: " + err.Error())
			}

			if p != nil {
				return fmt.Errorf("Data Forwarding still exists!")
			}
		}
		return nil
	}
}

func testAccSumologicDataForwardingCreateConfig(destinationName string, description string, testAwsBucket string, testAwsRoleArn string, testAwsRegion string) string {
	return fmt.Sprintf(`
resource "sumologic_s3_data_forwarding_destination" "test" {
    destination_name = "%s"
	description = "%s"
    bucket_name      = "%s"
    authentication {
        type     = "RoleBased"
        role_arn = "%s"
    }
	s3_region = "%s"
	s3_server_side_encryption = false
}
`, destinationName, description, testAwsBucket, testAwsRoleArn, testAwsRegion)
}

func testAccSumologicDataForwardingUpdateConfig(destinationName string, description string, testAwsBucket string, testAwsRoleArn string, testAwsRegion string) string {
	return fmt.Sprintf(`
resource "sumologic_s3_data_forwarding_destination" "test" {
    destination_name = "%s"
	description = "%s"
    bucket_name      = "%s"
    authentication {
        type     = "RoleBased"
        role_arn = "%s"
    }
	s3_region = "%s"
	s3_server_side_encryption = true
}
`, destinationName, description, testAwsBucket, testAwsRoleArn, testAwsRegion)
}

func testAccSumologicDataForwardingDeleteConfig() string {
	return fmt.Sprintf(` `)
}
