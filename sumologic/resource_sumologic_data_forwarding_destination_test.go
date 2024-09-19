package sumologic

// Todo: (ngoyal, 09-19-2024) - Uncomment the test cases once the issue with PUT permissions on SUMOLOGIC_TEST_BUCKET_NAME is fixed


import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"strings"
	"testing"

	b64 "encoding/base64"
)

func getTestParams() (string, string, string, string, string, string) {
	dataForwardingDestinationResourceName := "sumologic_data_forwarding_destination.test"
	destinationName, description := getRandomizedDataForwardingDestinationParams()
	testAwsRoleArn := os.Getenv("SUMOLOGIC_TEST_ROLE_ARN")
	testAwsBucket := os.Getenv("SUMOLOGIC_TEST_BUCKET_NAME")
	testAwsRegion := os.Getenv("SUMOLOGIC_TEST_AWS_REGION")
	println("AWS Test Bucket b64: ", b64.StdEncoding.EncodeToString([]byte("foo."+testAwsBucket+".")))
	println("AWS Test ARN b64: ", b64.StdEncoding.EncodeToString([]byte("phoo."+testAwsRoleArn+".")))
	println("AWS Test Region b64: ", b64.StdEncoding.EncodeToString([]byte("fu."+testAwsRegion+".")))
	return dataForwardingDestinationResourceName, destinationName, description, testAwsRegion, testAwsRoleArn, testAwsBucket
}

func getRandomizedDataForwardingDestinationParams() (string, string) {
	destinationName := acctest.RandomWithPrefix("tf-acc-test")
	description := acctest.RandomWithPrefix("tf-acc-test")
	return destinationName, description
}

func TestAccSumologicDataForwarding_create(t *testing.T) {
	dataForwardingDestinationResourceName, destinationName, description, testAwsRegion, testAwsRoleArn, testAwsBucket := getTestParams()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataForwardingDestinationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicDataForwardingCreateConfig(destinationName, description, testAwsBucket, testAwsRoleArn, testAwsRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataForwardingDestinationExists(),
					resource.TestCheckResourceAttr(dataForwardingDestinationResourceName, "destination_name", destinationName),
					resource.TestCheckResourceAttr(dataForwardingDestinationResourceName, "description", description),
					resource.TestCheckResourceAttr(dataForwardingDestinationResourceName, "s3_region", testAwsRegion),
					resource.TestCheckResourceAttr(dataForwardingDestinationResourceName, "s3_server_side_encryption", "false"),
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
		CheckDestroy: testAccCheckDataForwardingDestinationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicDataForwardingCreateConfig(destinationName, description, testAwsBucket, testAwsRoleArn, testAwsRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataForwardingDestinationExists(),
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
		CheckDestroy: testAccCheckDataForwardingDestinationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicDataForwardingCreateConfig(destinationName, description, testAwsBucket, testAwsRoleArn, testAwsRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataForwardingDestinationExists(),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "destination_name", destinationName),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "description", description),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "s3_region", testAwsRegion),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "s3_server_side_encryption", "false"),
				),
			}, {
				Config: testAccSumologicDataForwardingUpdateConfig(destinationName, description, testAwsBucket, testAwsRoleArn, testAwsRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataForwardingDestinationExists(),
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
		CheckDestroy: testAccCheckDataForwardingDestinationDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicDataForwardingCreateConfig(destinationName, description, testAwsBucket, testAwsRoleArn, testAwsRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataForwardingDestinationExists(),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "destination_name", destinationName),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "description", description),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "s3_region", testAwsRegion),
					resource.TestCheckResourceAttr(dataForwardingResourceName, "s3_server_side_encryption", "false"),
				),
			}, {
				Config: testAccSumologicDataForwardingDeleteConfig(),
				Check:  resource.ComposeTestCheckFunc(testAccCheckDataForwardingDestinationDestroy()),
			},
		},
	})
}

func testAccCheckDataForwardingDestinationExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			if _, err := client.getDataForwardingDestination(id); err != nil {
				return fmt.Errorf("Received an error retrieving data forwarding destination %s", err)
			}
		}
		return nil
	}
}

func testAccCheckDataForwardingDestinationDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			p, err := client.getDataForwardingDestination(id)

			if err != nil {
				if strings.Contains(err.Error(), "Data forwarding Destination doesn't exists") {
					continue
				}

				return fmt.Errorf("Encountered an error: " + err.Error())
			}

			if p != nil {
				return fmt.Errorf("Data Forwarding Destination still exists!")
			}
		}
		return nil
	}
}

func testAccSumologicDataForwardingCreateConfig(destinationName string, description string, testAwsBucket string, testAwsRoleArn string, testAwsRegion string) string {
	return fmt.Sprintf(`
resource "sumologic_data_forwarding_destination" "test" {
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
resource "sumologic_data_forwarding_destination" "test" {
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
