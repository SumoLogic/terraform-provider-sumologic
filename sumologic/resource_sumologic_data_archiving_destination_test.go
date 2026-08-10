package sumologic

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func skipDataArchivingTest(t *testing.T) {
	if strings.ToLower(os.Getenv("SKIP_DATA_ARCHIVING_TESTS")) == "true" {
		t.Skip("Skipping Data Archiving Test")
	}
}

func testAccPreCheckDataArchivingWithAWS(t *testing.T) {
	testAccPreCheck(t)
	skipDataArchivingTest(t)
	if v := os.Getenv("SUMOLOGIC_DATA_FORWARDING_ROLE_ARN"); v == "" {
		t.Fatal("SUMOLOGIC_DATA_FORWARDING_ROLE_ARN must be set for data archiving S3 acceptance tests")
	}
	if v := os.Getenv("SUMOLOGIC_DATA_FORWARDING_BUCKET"); v == "" {
		t.Fatal("SUMOLOGIC_DATA_FORWARDING_BUCKET must be set for data archiving S3 acceptance tests")
	}
	if v := os.Getenv("SUMOLOGIC_DATA_FORWARDING_AWS_REGION"); v == "" {
		t.Fatal("SUMOLOGIC_DATA_FORWARDING_AWS_REGION must be set for data archiving S3 acceptance tests")
	}
}

func testAccPreCheckDataArchiving(t *testing.T) {
	testAccPreCheck(t)
	skipDataArchivingTest(t)
}

func TestAccSumologicDataArchivingDestination_createS3RoleBased(t *testing.T) {
	name := "terraform_test_archive_" + acctest.RandString(10)
	resourceName := "sumologic_data_archiving_destination.test"
	testAwsRoleArn := os.Getenv("SUMOLOGIC_DATA_FORWARDING_ROLE_ARN")
	testAwsBucket := os.Getenv("SUMOLOGIC_DATA_FORWARDING_BUCKET")
	testAwsRegion := os.Getenv("SUMOLOGIC_DATA_FORWARDING_AWS_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDataArchivingWithAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataArchivingDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataArchivingDestinationS3RoleBased(name, testAwsBucket, testAwsRegion, testAwsRoleArn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataArchivingDestinationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "destination_name", name),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.destination_type", "S3"),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.bucket_name", testAwsBucket),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.region", testAwsRegion),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.encrypted", "true"),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.auth_config.0.authentication_mode", "RoleBased"),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.auth_config.0.role_arn", testAwsRoleArn),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSumologicDataArchivingDestination_createSyslog(t *testing.T) {
	name := "terraform_test_archive_" + acctest.RandString(10)
	resourceName := "sumologic_data_archiving_destination.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDataArchiving(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataArchivingDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataArchivingDestinationSyslog(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataArchivingDestinationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "destination_name", name),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.destination_type", "Syslog"),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.protocol", "udp"),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.host", "10.20.30.40"),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.port", "514"),
				),
			},
		},
	})
}

func TestAccSumologicDataArchivingDestination_createHitachi(t *testing.T) {
	name := "terraform_test_archive_" + acctest.RandString(10)
	resourceName := "sumologic_data_archiving_destination.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDataArchiving(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataArchivingDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataArchivingDestinationHitachi(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataArchivingDestinationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "destination_name", name),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.destination_type", "Hitachi"),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.url", "https://hitachi.example.com"),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.object_id", "archive_path/logname_{day}_{hour}_{minute}_{second}_{uuid}.log"),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.username", "testuser"),
				),
			},
		},
	})
}

func TestAccSumologicDataArchivingDestination_createRestAPI(t *testing.T) {
	name := "terraform_test_archive_" + acctest.RandString(10)
	resourceName := "sumologic_data_archiving_destination.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDataArchiving(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataArchivingDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataArchivingDestinationRestAPI(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataArchivingDestinationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "destination_name", name),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.destination_type", "RestAPI"),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.url", "https://example.com/receiver/v1/http/token"),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.object_id", "archive_test"),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.username", "testuser"),
				),
			},
		},
	})
}

func TestAccSumologicDataArchivingDestination_update(t *testing.T) {
	name := "terraform_test_archive_" + acctest.RandString(10)
	updatedName := name + "_updated"
	resourceName := "sumologic_data_archiving_destination.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDataArchiving(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataArchivingDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataArchivingDestinationSyslog(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataArchivingDestinationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "destination_name", name),
				),
			},
			{
				Config: testAccDataArchivingDestinationSyslogUpdated(updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataArchivingDestinationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "destination_name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.host", "192.168.1.1"),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.port", "1514"),
				),
			},
		},
	})
}

func testAccCheckDataArchivingDestinationExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID not set")
		}

		c := testAccProvider.Meta().(*Client)
		dest, err := c.GetDataArchivingDestination(rs.Primary.ID)
		if err != nil {
			return err
		}
		if dest == nil {
			return fmt.Errorf("data archiving destination %s not found", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckDataArchivingDestinationDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_data_archiving_destination" {
			continue
		}
		dest, err := c.GetDataArchivingDestination(rs.Primary.ID)
		if err != nil {
			return err
		}
		if dest != nil {
			return fmt.Errorf("data archiving destination %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccDataArchivingDestinationS3RoleBased(name, bucket, region, roleArn string) string {
	return fmt.Sprintf(`
resource "sumologic_data_archiving_destination" "test" {
  destination_name = "%s"

  destination_config {
    destination_type = "S3"
    bucket_name      = "%s"
    region           = "%s"
    encrypted        = true
    enabled          = true

    auth_config {
      authentication_mode = "RoleBased"
      role_arn            = "%s"
    }
  }
}
`, name, bucket, region, roleArn)
}

func testAccDataArchivingDestinationSyslog(name string) string {
	return fmt.Sprintf(`
resource "sumologic_data_archiving_destination" "test" {
  destination_name = "%s"

  destination_config {
    destination_type = "Syslog"
    protocol         = "udp"
    host             = "10.20.30.40"
    port             = 514
  }
}
`, name)
}

func testAccDataArchivingDestinationHitachi(name string) string {
	return fmt.Sprintf(`
resource "sumologic_data_archiving_destination" "test" {
  destination_name = "%s"

  destination_config {
    destination_type = "Hitachi"
    url              = "https://hitachi.example.com"
    object_id        = "archive_path/logname_{day}_{hour}_{minute}_{second}_{uuid}.log"
    username         = "testuser"
    password         = "testpassword"
  }
}
`, name)
}

func testAccDataArchivingDestinationRestAPI(name string) string {
	return fmt.Sprintf(`
resource "sumologic_data_archiving_destination" "test" {
  destination_name = "%s"

  destination_config {
    destination_type = "RestAPI"
    url              = "https://example.com/receiver/v1/http/token"
    object_id        = "archive_test"
    username         = "testuser"
    password         = "testpassword"
  }
}
`, name)
}

func testAccDataArchivingDestinationSyslogUpdated(name string) string {
	return fmt.Sprintf(`
resource "sumologic_data_archiving_destination" "test" {
  destination_name = "%s"

  destination_config {
    destination_type = "Syslog"
    protocol         = "tcp"
    host             = "192.168.1.1"
    port             = 1514
  }
}
`, name)
}
