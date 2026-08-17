package sumologic

import (
	"fmt"
	"os"
	"regexp"
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

// Exercises the S3 update contract, which accepts neither the id nor the bucket name.
// The bucket is deliberately held constant so the steps update in place instead of
// forcing a replacement.
func TestAccSumologicDataArchivingDestination_updateS3RoleBased(t *testing.T) {
	name := "terraform_test_archive_" + acctest.RandString(10)
	updatedName := name + "_updated"
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
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.enabled", "true"),
				),
			},
			{
				Config: testAccDataArchivingDestinationS3RoleBasedUpdated(updatedName, testAwsBucket, testAwsRegion, testAwsRoleArn, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataArchivingDestinationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "destination_name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.bucket_name", testAwsBucket),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.region", testAwsRegion),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "modified_at"),
				),
			},
			{
				Config: testAccDataArchivingDestinationS3RoleBasedUpdated(updatedName, testAwsBucket, testAwsRegion, testAwsRoleArn, "updated by acceptance test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataArchivingDestinationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.description", "updated by acceptance test"),
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

// bucket_name and region are validated on every plan, not only on create, so dropping
// either one from an existing destination has to fail at plan time rather than silently
// leaving a diff that never converges.
func TestAccSumologicDataArchivingDestination_s3RequiresBucketAndRegionOnUpdate(t *testing.T) {
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
				Check:  testAccCheckDataArchivingDestinationExists(resourceName),
			},
			{
				Config:      testAccDataArchivingDestinationS3RoleBased(name, testAwsBucket, "", testAwsRoleArn),
				ExpectError: regexp.MustCompile(`region is required for S3 destination`),
			},
			{
				Config:      testAccDataArchivingDestinationS3RoleBased(name, "", testAwsRegion, testAwsRoleArn),
				ExpectError: regexp.MustCompile(`bucket_name is required for S3 destination`),
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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// The API never returns the password, so an imported destination
				// cannot reproduce it.
				ImportStateVerifyIgnore: []string{"destination_config.0.password"},
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
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"destination_config.0.password"},
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

func TestAccSumologicDataArchivingDestination_updateHitachi(t *testing.T) {
	name := "terraform_test_archive_" + acctest.RandString(10)
	updatedName := name + "_updated"
	resourceName := "sumologic_data_archiving_destination.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDataArchiving(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataArchivingDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataArchivingDestinationHitachi(name),
				Check:  testAccCheckDataArchivingDestinationExists(resourceName),
			},
			{
				Config: testAccDataArchivingDestinationHitachiUpdated(updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataArchivingDestinationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "destination_name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.url", "https://hitachi-updated.example.com"),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.object_id", "archive_path_updated/logname_{day}_{hour}_{uuid}.log"),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.username", "testuser_updated"),
				),
			},
		},
	})
}

func TestAccSumologicDataArchivingDestination_updateRestAPI(t *testing.T) {
	name := "terraform_test_archive_" + acctest.RandString(10)
	updatedName := name + "_updated"
	resourceName := "sumologic_data_archiving_destination.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDataArchiving(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataArchivingDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataArchivingDestinationRestAPI(name),
				Check:  testAccCheckDataArchivingDestinationExists(resourceName),
			},
			{
				Config: testAccDataArchivingDestinationRestAPIUpdated(updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataArchivingDestinationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "destination_name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.url", "https://example.com/receiver/v1/http/token_updated"),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.object_id", "archive_test_updated"),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.username", "testuser_updated"),
				),
			},
		},
	})
}

// A RestAPI destination can be created without a username but not updated without one,
// so the second step has to fail at plan time instead of returning an API error.
func TestAccSumologicDataArchivingDestination_restAPIRequiresUsernameOnUpdate(t *testing.T) {
	name := "terraform_test_archive_" + acctest.RandString(10)
	resourceName := "sumologic_data_archiving_destination.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDataArchiving(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataArchivingDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataArchivingDestinationRestAPIWithoutUsername(name, "https://example.com/receiver/v1/http/token"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataArchivingDestinationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.username", ""),
				),
			},
			{
				Config:      testAccDataArchivingDestinationRestAPIWithoutUsername(name, "https://example.com/receiver/v1/http/token_updated"),
				ExpectError: regexp.MustCompile(`username is required to update a RestAPI destination`),
			},
		},
	})
}

// destination_type is ForceNew, so changing it has to replace the destination rather
// than attempt an update the API cannot perform.
func TestAccSumologicDataArchivingDestination_forceNewOnDestinationTypeChange(t *testing.T) {
	name := "terraform_test_archive_" + acctest.RandString(10)
	resourceName := "sumologic_data_archiving_destination.test"
	var firstID string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDataArchiving(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataArchivingDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataArchivingDestinationSyslog(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataArchivingDestinationExists(resourceName),
					testAccCaptureDataArchivingDestinationID(resourceName, &firstID),
				),
			},
			{
				Config: testAccDataArchivingDestinationHitachi(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataArchivingDestinationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "destination_config.0.destination_type", "Hitachi"),
					testAccCheckDataArchivingDestinationIDChanged(resourceName, &firstID),
				),
			},
		},
	})
}

// Every plan-time validation branch, exercised without touching the API: the plan fails
// before any request is made.
func TestAccSumologicDataArchivingDestination_invalidConfig(t *testing.T) {
	const name = "terraform_test_archive_invalid"

	cases := []struct {
		testName    string
		config      string
		expectError *regexp.Regexp
	}{
		{
			testName: "descriptionOnNonS3",
			config: testAccDataArchivingDestinationConfig(name, `
    destination_type = "Syslog"
    description      = "not supported outside S3"
    protocol         = "udp"
    host             = "10.20.30.40"
    port             = 514`),
			expectError: regexp.MustCompile(`description is only supported for S3 destinations`),
		},
		{
			testName: "s3MissingBucketName",
			config: testAccDataArchivingDestinationConfig(name, `
    destination_type = "S3"
    region           = "us-east-1"
    encrypted        = true
    enabled          = true`+testAccArchivingRoleBasedAuthBlock),
			expectError: regexp.MustCompile(`bucket_name is required for S3 destination`),
		},
		{
			testName: "s3MissingRegion",
			config: testAccDataArchivingDestinationConfig(name, `
    destination_type = "S3"
    bucket_name      = "terraform-test-bucket"
    encrypted        = true
    enabled          = true`+testAccArchivingRoleBasedAuthBlock),
			expectError: regexp.MustCompile(`region is required for S3 destination`),
		},
		{
			testName: "s3MissingEncrypted",
			config: testAccDataArchivingDestinationConfig(name, `
    destination_type = "S3"
    bucket_name      = "terraform-test-bucket"
    region           = "us-east-1"
    enabled          = true`+testAccArchivingRoleBasedAuthBlock),
			expectError: regexp.MustCompile(`encrypted is required for S3 destination`),
		},
		{
			testName: "s3MissingEnabled",
			config: testAccDataArchivingDestinationConfig(name, `
    destination_type = "S3"
    bucket_name      = "terraform-test-bucket"
    region           = "us-east-1"
    encrypted        = true`+testAccArchivingRoleBasedAuthBlock),
			expectError: regexp.MustCompile(`enabled is required for S3 destination`),
		},
		{
			testName: "s3MissingAuthConfig",
			config: testAccDataArchivingDestinationConfig(name, `
    destination_type = "S3"
    bucket_name      = "terraform-test-bucket"
    region           = "us-east-1"
    encrypted        = true
    enabled          = true`),
			expectError: regexp.MustCompile(`auth_config is required for S3 destination`),
		},
		{
			testName: "s3AccessKeyMissingKeyID",
			config: testAccDataArchivingDestinationConfig(name, `
    destination_type = "S3"
    bucket_name      = "terraform-test-bucket"
    region           = "us-east-1"
    encrypted        = true
    enabled          = true

    auth_config {
      authentication_mode = "AccessKey"
      access_key_secret   = "test-secret"
    }`),
			expectError: regexp.MustCompile(`access_key_id is required when authentication_mode is AccessKey`),
		},
		{
			testName: "s3AccessKeyMissingSecret",
			config: testAccDataArchivingDestinationConfig(name, `
    destination_type = "S3"
    bucket_name      = "terraform-test-bucket"
    region           = "us-east-1"
    encrypted        = true
    enabled          = true

    auth_config {
      authentication_mode = "AccessKey"
      access_key_id       = "AKIAIOSFODNN7EXAMPLE"
    }`),
			expectError: regexp.MustCompile(`access_key_secret is required when authentication_mode is AccessKey`),
		},
		{
			testName: "s3RoleBasedMissingRoleArn",
			config: testAccDataArchivingDestinationConfig(name, `
    destination_type = "S3"
    bucket_name      = "terraform-test-bucket"
    region           = "us-east-1"
    encrypted        = true
    enabled          = true

    auth_config {
      authentication_mode = "RoleBased"
    }`),
			expectError: regexp.MustCompile(`role_arn is required when authentication_mode is RoleBased`),
		},
		{
			testName: "syslogMissingProtocol",
			config: testAccDataArchivingDestinationConfig(name, `
    destination_type = "Syslog"
    host             = "10.20.30.40"
    port             = 514`),
			expectError: regexp.MustCompile(`protocol is required for Syslog destination`),
		},
		{
			testName: "syslogMissingHost",
			config: testAccDataArchivingDestinationConfig(name, `
    destination_type = "Syslog"
    protocol         = "udp"
    port             = 514`),
			expectError: regexp.MustCompile(`host is required for Syslog destination`),
		},
		{
			testName: "syslogMissingPort",
			config: testAccDataArchivingDestinationConfig(name, `
    destination_type = "Syslog"
    protocol         = "udp"
    host             = "10.20.30.40"`),
			expectError: regexp.MustCompile(`port is required for Syslog destination`),
		},
		{
			testName: "hitachiMissingURL",
			config: testAccDataArchivingDestinationConfig(name, `
    destination_type = "Hitachi"
    object_id        = "archive_path/logname.log"
    username         = "testuser"
    password         = "testpassword"`),
			expectError: regexp.MustCompile(`url is required for Hitachi destination`),
		},
		{
			testName: "hitachiMissingUsername",
			config: testAccDataArchivingDestinationConfig(name, `
    destination_type = "Hitachi"
    url              = "https://hitachi.example.com"
    object_id        = "archive_path/logname.log"
    password         = "testpassword"`),
			expectError: regexp.MustCompile(`username is required for Hitachi destination`),
		},
		{
			testName: "hitachiMissingObjectID",
			config: testAccDataArchivingDestinationConfig(name, `
    destination_type = "Hitachi"
    url              = "https://hitachi.example.com"
    username         = "testuser"
    password         = "testpassword"`),
			expectError: regexp.MustCompile(`object_id is required for Hitachi destination on create`),
		},
		{
			testName: "hitachiMissingPassword",
			config: testAccDataArchivingDestinationConfig(name, `
    destination_type = "Hitachi"
    url              = "https://hitachi.example.com"
    object_id        = "archive_path/logname.log"
    username         = "testuser"`),
			expectError: regexp.MustCompile(`password is required for Hitachi destination on create`),
		},
		{
			testName: "restAPIMissingURL",
			config: testAccDataArchivingDestinationConfig(name, `
    destination_type = "RestAPI"
    object_id        = "archive_test"
    username         = "testuser"
    password         = "testpassword"`),
			expectError: regexp.MustCompile(`url is required for RestAPI destination`),
		},
		{
			testName: "unsupportedDestinationType",
			config: testAccDataArchivingDestinationConfig(name, `
    destination_type = "GCS"`),
			expectError: regexp.MustCompile(`to be one of.*got GCS`),
		},
		{
			testName: "unsupportedSyslogProtocol",
			config: testAccDataArchivingDestinationConfig(name, `
    destination_type = "Syslog"
    protocol         = "sctp"
    host             = "10.20.30.40"
    port             = 514`),
			expectError: regexp.MustCompile(`to be one of.*got sctp`),
		},
		{
			testName: "syslogPortOutOfRange",
			config: testAccDataArchivingDestinationConfig(name, `
    destination_type = "Syslog"
    protocol         = "udp"
    host             = "10.20.30.40"
    port             = 70000`),
			expectError: regexp.MustCompile(`to be in the range.*got 70000`),
		},
		{
			testName: "emptyDestinationName",
			config: testAccDataArchivingDestinationConfig("", `
    destination_type = "Syslog"
    protocol         = "udp"
    host             = "10.20.30.40"
    port             = 514`),
			expectError: regexp.MustCompile(`expected length of destination_name to be in the range`),
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:     func() { testAccPreCheckDataArchiving(t) },
				Providers:    testAccProviders,
				CheckDestroy: testAccCheckDataArchivingDestinationDestroy,
				Steps: []resource.TestStep{
					{
						Config:      c.config,
						PlanOnly:    true,
						ExpectError: c.expectError,
					},
				},
			})
		})
	}
}

func testAccCaptureDataArchivingDestinationID(name string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}
		*id = rs.Primary.ID
		return nil
	}
}

func testAccCheckDataArchivingDestinationIDChanged(name string, previousID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}
		if rs.Primary.ID == *previousID {
			return fmt.Errorf("expected destination to be replaced, but id is still %s", *previousID)
		}
		return nil
	}
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

// description is only valid for S3, so it is omitted entirely when empty rather than
// rendered as an empty string.
func testAccDataArchivingDestinationS3RoleBasedUpdated(name, bucket, region, roleArn, description string) string {
	descriptionLine := ""
	if description != "" {
		descriptionLine = fmt.Sprintf("\n    description      = %q", description)
	}

	return fmt.Sprintf(`
resource "sumologic_data_archiving_destination" "test" {
  destination_name = "%s"

  destination_config {
    destination_type = "S3"%s
    bucket_name      = "%s"
    region           = "%s"
    encrypted        = true
    enabled          = false

    auth_config {
      authentication_mode = "RoleBased"
      role_arn            = "%s"
    }
  }
}
`, name, descriptionLine, bucket, region, roleArn)
}

func testAccDataArchivingDestinationHitachiUpdated(name string) string {
	return fmt.Sprintf(`
resource "sumologic_data_archiving_destination" "test" {
  destination_name = "%s"

  destination_config {
    destination_type = "Hitachi"
    url              = "https://hitachi-updated.example.com"
    object_id        = "archive_path_updated/logname_{day}_{hour}_{uuid}.log"
    username         = "testuser_updated"
    password         = "testpassword_updated"
  }
}
`, name)
}

func testAccDataArchivingDestinationRestAPIUpdated(name string) string {
	return fmt.Sprintf(`
resource "sumologic_data_archiving_destination" "test" {
  destination_name = "%s"

  destination_config {
    destination_type = "RestAPI"
    url              = "https://example.com/receiver/v1/http/token_updated"
    object_id        = "archive_test_updated"
    username         = "testuser_updated"
    password         = "testpassword_updated"
  }
}
`, name)
}

func testAccDataArchivingDestinationRestAPIWithoutUsername(name, url string) string {
	return fmt.Sprintf(`
resource "sumologic_data_archiving_destination" "test" {
  destination_name = "%s"

  destination_config {
    destination_type = "RestAPI"
    url              = "%s"
    object_id        = "archive_test"
  }
}
`, name, url)
}

const testAccArchivingRoleBasedAuthBlock = `

    auth_config {
      authentication_mode = "RoleBased"
      role_arn            = "arn:aws:iam::123456789012:role/terraform-test"
    }`

func testAccDataArchivingDestinationConfig(name, destinationConfigBody string) string {
	return fmt.Sprintf(`
resource "sumologic_data_archiving_destination" "test" {
  destination_name = "%s"

  destination_config {%s
  }
}
`, name, destinationConfigBody)
}
