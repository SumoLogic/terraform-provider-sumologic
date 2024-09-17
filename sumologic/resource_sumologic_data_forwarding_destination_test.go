package sumologic

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"strconv"
	"strings"
	"testing"
)

func TestAccSumologicDataForwardingDestination_basic(t *testing.T) {
	var dataForwardingDestination DataForwardingDestination
	testDestinationName := "df-destination"
	testRoleArn := "roleArn"
	testAuthenticationMode := "RoleBased"
	testRegion := "us-east-1"
	testSecretAccessKey := "secretAccessKey"
	testBucketName := "df-bucket"
	testEnabled := false
	testDescription := "description-Fr4YU"
	testAccessKeyId := "accessKeyId"
	testEncrypted := false

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataForwardingDestinationDestroy(dataForwardingDestination),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSumologicDataForwardingDestinationConfigImported(testDestinationName, testRoleArn, testAuthenticationMode, testRegion, testSecretAccessKey, testBucketName, testEnabled, testDescription, testAccessKeyId, testEncrypted),
			},
			{
				ResourceName:      "sumologic_data_forwarding_destination.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDataForwardingDestination_create(t *testing.T) {
	var dataForwardingDestination DataForwardingDestination
	testDestinationName := "df-destination"
	testRoleArn := "roleArn"
	testAuthenticationMode := "RoleBased"
	testRegion := "us-east-1"
	testSecretAccessKey := "secretAccessKey"
	testBucketName := "df-bucket"
	testEnabled := false
	testDescription := "description-uhKep"
	testAccessKeyId := "accessKeyId"
	testEncrypted := false
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataForwardingDestinationDestroy(dataForwardingDestination),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicDataForwardingDestination(testDestinationName, testRoleArn, testAuthenticationMode, testRegion, testSecretAccessKey, testBucketName, testEnabled, testDescription, testAccessKeyId, testEncrypted),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataForwardingDestinationExists("sumologic_data_forwarding_destination.test", &dataForwardingDestination, t),
					testAccCheckDataForwardingDestinationAttributes("sumologic_data_forwarding_destination.test"),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "destination_name", testDestinationName),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "role_arn", testRoleArn),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "authentication_mode", testAuthenticationMode),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "region", testRegion),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "secret_access_key", testSecretAccessKey),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "bucket_name", testBucketName),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "enabled", strconv.FormatBool(testEnabled)),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "description", testDescription),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "access_key_id", testAccessKeyId),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "encrypted", strconv.FormatBool(testEncrypted)),
				),
			},
		},
	})
}

func TestAccDataForwardingDestination_update(t *testing.T) {
	var dataForwardingDestination DataForwardingDestination
	testDestinationName := "df-destination"
	testRoleArn := "roleArn"
	testAuthenticationMode := "RoleBased"
	testRegion := "us-east-1"
	testSecretAccessKey := "secretAccessKey"
	testBucketName := "df-bucket"
	testEnabled := false
	testDescription := "description-Up4Sy"
	testAccessKeyId := "accessKeyId"
	testEncrypted := false

	testUpdatedDestinationName := "df-destinationUpdate"
	testUpdatedRoleArn := "roleArnUpdate"
	testUpdatedAuthenticationMode := "RoleBasedUpdate"
	testUpdatedRegion := "us-east-1Update"
	testUpdatedSecretAccessKey := "secretAccessKeyUpdate"
	testUpdatedBucketName := "?!xn--z7z"
	testUpdatedEnabled := false
	testUpdatedDescription := "description-TOW5aUpdate"
	testUpdatedAccessKeyId := "accessKeyIdUpdate"
	testUpdatedEncrypted := false

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataForwardingDestinationDestroy(dataForwardingDestination),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicDataForwardingDestination(testDestinationName, testRoleArn, testAuthenticationMode, testRegion, testSecretAccessKey, testBucketName, testEnabled, testDescription, testAccessKeyId, testEncrypted),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataForwardingDestinationExists("sumologic_data_forwarding_destination.test", &dataForwardingDestination, t),
					testAccCheckDataForwardingDestinationAttributes("sumologic_data_forwarding_destination.test"),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "destination_name", testDestinationName),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "role_arn", testRoleArn),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "authentication_mode", testAuthenticationMode),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "region", testRegion),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "secret_access_key", testSecretAccessKey),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "bucket_name", testBucketName),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "enabled", strconv.FormatBool(testEnabled)),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "description", testDescription),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "access_key_id", testAccessKeyId),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "encrypted", strconv.FormatBool(testEncrypted)),
				),
			},
			{
				Config: testAccSumologicDataForwardingDestinationUpdate(testUpdatedDestinationName, testUpdatedRoleArn, testUpdatedAuthenticationMode, testUpdatedRegion, testUpdatedSecretAccessKey, testUpdatedBucketName, testUpdatedEnabled, testUpdatedDescription, testUpdatedAccessKeyId, testUpdatedEncrypted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "destination_name", testUpdatedDestinationName),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "role_arn", testUpdatedRoleArn),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "authentication_mode", testUpdatedAuthenticationMode),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "region", testUpdatedRegion),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "secret_access_key", testUpdatedSecretAccessKey),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "bucket_name", testUpdatedBucketName),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "enabled", strconv.FormatBool(testUpdatedEnabled)),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "description", testUpdatedDescription),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "access_key_id", testUpdatedAccessKeyId),
					resource.TestCheckResourceAttr("sumologic_data_forwarding_destination.test", "encrypted", strconv.FormatBool(testUpdatedEncrypted)),
				),
			},
		},
	})
}

func testAccCheckDataForwardingDestinationDestroy(dataForwardingDestination DataForwardingDestination) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			u, err := client.GetDataForwardingDestination(id)
			if err != nil {
				return fmt.Errorf("Encountered an error: " + err.Error())
			}
			if u != nil {
				return fmt.Errorf("DataForwardingDestination %s still exists", id)
			}
		}
		return nil
	}
}

func testAccCheckDataForwardingDestinationExists(name string, dataForwardingDestination *DataForwardingDestination, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			//need this so that we don't get an unused import error for strconv in some cases
			return fmt.Errorf("Error = %s. DataForwardingDestination not found: %s", strconv.FormatBool(ok), name)
		}

		//need this so that we don't get an unused import error for strings in some cases
		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("DataForwardingDestination ID is not set")
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)
		newDataForwardingDestination, err := client.GetDataForwardingDestination(id)
		if err != nil {
			return fmt.Errorf("DataForwardingDestination %s not found", id)
		}
		dataForwardingDestination = newDataForwardingDestination
		return nil
	}
}
func testAccCheckSumologicDataForwardingDestinationConfigImported(destinationName string, roleArn string, authenticationMode string, region string, secretAccessKey string, bucketName string, enabled bool, description string, accessKeyId string, encrypted bool) string {
	return fmt.Sprintf(`
resource "sumologic_data_forwarding_destination" "foo" {
      destination_name = "%s"
      role_arn = "%s"
      authentication_mode = "%s"
      region = "%s"
      secret_access_key = "%s"
      bucket_name = "%s"
      enabled = %t
      description = "%s"
      access_key_id = "%s"
      encrypted = %t
}
`, destinationName, roleArn, authenticationMode, region, secretAccessKey, bucketName, enabled, description, accessKeyId, encrypted)
}

func testAccSumologicDataForwardingDestination(destinationName string, roleArn string, authenticationMode string, region string, secretAccessKey string, bucketName string, enabled bool, description string, accessKeyId string, encrypted bool) string {
	return fmt.Sprintf(`
resource "sumologic_data_forwarding_destination" "test" {
    destination_name = "%s"
    role_arn = "%s"
    authentication_mode = "%s"
    region = "%s"
    secret_access_key = "%s"
    bucket_name = "%s"
    enabled = %t
    description = "%s"
    access_key_id = "%s"
    encrypted = %t
}
`, destinationName, roleArn, authenticationMode, region, secretAccessKey, bucketName, enabled, description, accessKeyId, encrypted)
}

func testAccSumologicDataForwardingDestinationUpdate(destinationName string, roleArn string, authenticationMode string, region string, secretAccessKey string, bucketName string, enabled bool, description string, accessKeyId string, encrypted bool) string {
	return fmt.Sprintf(`
resource "sumologic_data_forwarding_destination" "test" {
      destination_name = "%s"
      role_arn = "%s"
      authentication_mode = "%s"
      region = "%s"
      secret_access_key = "%s"
      bucket_name = "%s"
      enabled = %t
      description = "%s"
      access_key_id = "%s"
      encrypted = %t
}
`, destinationName, roleArn, authenticationMode, region, secretAccessKey, bucketName, enabled, description, accessKeyId, encrypted)
}

func testAccCheckDataForwardingDestinationAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(name, "destination_name"),
			resource.TestCheckResourceAttrSet(name, "role_arn"),
			resource.TestCheckResourceAttrSet(name, "authentication_mode"),
			resource.TestCheckResourceAttrSet(name, "region"),
			resource.TestCheckResourceAttrSet(name, "secret_access_key"),
			resource.TestCheckResourceAttrSet(name, "bucket_name"),
			resource.TestCheckResourceAttrSet(name, "enabled"),
			resource.TestCheckResourceAttrSet(name, "description"),
			resource.TestCheckResourceAttrSet(name, "access_key_id"),
			resource.TestCheckResourceAttrSet(name, "encrypted"),
		)
		return f(s)
	}
}
