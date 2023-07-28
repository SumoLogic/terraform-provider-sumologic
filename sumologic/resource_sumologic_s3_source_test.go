package sumologic

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func getRandomizedParams() (string, string, string) {
	name := acctest.RandomWithPrefix("tf-acc-test")
	description := acctest.RandomWithPrefix("tf-acc-test")
	category := acctest.RandomWithPrefix("tf-acc-test")
	return name, description, category
}
func TestAccSumologicS3Source_create(t *testing.T) {
	var s3Source PollingSource
	var collector Collector
	cName, cDescription, cCategory := getRandomizedParams()
	sName, sDescription, sCategory := getRandomizedParams()
	s3ResourceName := "sumologic_s3_source.s3"
	testAwsRoleArn := os.Getenv("SUMOLOGIC_TEST_ROLE_ARN")
	testAwsBucket := os.Getenv("SUMOLOGIC_TEST_BUCKET_NAME")
	useVersionedApi := true
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckS3SourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicS3SourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsRoleArn, testAwsBucket, useVersionedApi),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3SourceExists(s3ResourceName, &s3Source),
					testAccCheckS3SourceValues(&s3Source, sName, sDescription, sCategory),
					testAccCheckCollectorExists("sumologic_collector.test", &collector),
					testAccCheckCollectorValues(&collector, cName, cDescription, cCategory, "Etc/UTC", ""),
					resource.TestCheckResourceAttrSet(s3ResourceName, "id"),
					resource.TestCheckResourceAttr(s3ResourceName, "name", sName),
					resource.TestCheckResourceAttr(s3ResourceName, "description", sDescription),
					resource.TestCheckResourceAttr(s3ResourceName, "category", sCategory),
					resource.TestCheckResourceAttr(s3ResourceName, "content_type", "AwsS3Bucket"),
					resource.TestCheckResourceAttr(s3ResourceName, "path.0.type", "S3BucketPathExpression"),
					resource.TestCheckResourceAttr(s3ResourceName, "path.0.use_versioned_api", "true"),
				),
			},
		},
	})
}
func TestAccSumologicS3Source_update(t *testing.T) {
	var s3Source PollingSource
	cName, cDescription, cCategory := getRandomizedParams()
	sName, sDescription, sCategory := getRandomizedParams()
	sNameUpdated, sDescriptionUpdated, sCategoryUpdated := getRandomizedParams()
	s3ResourceName := "sumologic_s3_source.s3"
	testAwsRoleArn := os.Getenv("SUMOLOGIC_TEST_ROLE_ARN")
	testAwsBucket := os.Getenv("SUMOLOGIC_TEST_BUCKET_NAME")
	useVersionedApi := true
	useVersionedApiUpdated := false
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHTTPSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicS3SourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsRoleArn, testAwsBucket, useVersionedApi),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3SourceExists(s3ResourceName, &s3Source),
					testAccCheckS3SourceValues(&s3Source, sName, sDescription, sCategory),
					resource.TestCheckResourceAttrSet(s3ResourceName, "id"),
					resource.TestCheckResourceAttr(s3ResourceName, "name", sName),
					resource.TestCheckResourceAttr(s3ResourceName, "description", sDescription),
					resource.TestCheckResourceAttr(s3ResourceName, "category", sCategory),
					resource.TestCheckResourceAttr(s3ResourceName, "content_type", "AwsS3Bucket"),
					resource.TestCheckResourceAttr(s3ResourceName, "path.0.type", "S3BucketPathExpression"),
					resource.TestCheckResourceAttr(s3ResourceName, "path.0.use_versioned_api", "true"),
				),
			},
			{
				Config: testAccSumologicS3SourceConfig(cName, cDescription, cCategory, sNameUpdated, sDescriptionUpdated, sCategoryUpdated, testAwsRoleArn, testAwsBucket, useVersionedApiUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3SourceExists(s3ResourceName, &s3Source),
					testAccCheckS3SourceValues(&s3Source, sNameUpdated, sDescriptionUpdated, sCategoryUpdated),
					resource.TestCheckResourceAttrSet(s3ResourceName, "id"),
					resource.TestCheckResourceAttr(s3ResourceName, "name", sNameUpdated),
					resource.TestCheckResourceAttr(s3ResourceName, "description", sDescriptionUpdated),
					resource.TestCheckResourceAttr(s3ResourceName, "category", sCategoryUpdated),
					resource.TestCheckResourceAttr(s3ResourceName, "content_type", "AwsS3Bucket"),
					resource.TestCheckResourceAttr(s3ResourceName, "path.0.type", "S3BucketPathExpression"),
					resource.TestCheckResourceAttr(s3ResourceName, "path.0.use_versioned_api", "false"),
				),
			},
		},
	})
}
func testAccCheckS3SourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_s3_source" && rs.Type != "sumologic_cloudwatch_source" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("HTTP Source destruction check: HTTP Source ID is not set")
		}
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		s, err := client.GetPollingSource(collectorID, id)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("Polling Source still exists")
		}
	}
	return nil
}
func testAccCheckS3SourceExists(n string, pollingSource *PollingSource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Polling Source ID is not set")
		}
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Polling Source id should be int; got %s", rs.Primary.ID)
		}
		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		c := testAccProvider.Meta().(*Client)
		pollingSourceResp, err := c.GetPollingSource(collectorID, id)
		if err != nil {
			return err
		}
		*pollingSource = *pollingSourceResp
		return nil
	}
}
func testAccCheckS3SourceValues(pollingSource *PollingSource, name, description, category string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if pollingSource.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, pollingSource.Name)
		}
		if pollingSource.Description != description {
			return fmt.Errorf("bad description, expected \"%s\", got: %#v", description, pollingSource.Description)
		}
		if pollingSource.Category != category {
			return fmt.Errorf("bad category, expected \"%s\", got: %#v", category, pollingSource.Category)
		}
		return nil
	}
}
func testAccSumologicS3SourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsRoleArn, testAwsBucket string, useVersionedApi bool) string {
	return fmt.Sprintf(`
resource "sumologic_collector" "test" {
	name = "%s"
	description = "%s"
	category = "%s"
}
resource "sumologic_s3_source" "s3" {
	name = "%s"
	description = "%s"
	category = "%s"
	content_type  = "AwsS3Bucket"
  	scan_interval = 300000
  	paused        = false
	collector_id = "${sumologic_collector.test.id}"
	authentication {
		type = "AWSRoleBasedAuthentication"
		role_arn = "%s"
	  }
	  path {
		type = "S3BucketPathExpression"
		bucket_name     = "%s"
		path_expression = "*"
		use_versioned_api = "%v"
	  }
	}

`, cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsRoleArn, testAwsBucket, useVersionedApi)
}
