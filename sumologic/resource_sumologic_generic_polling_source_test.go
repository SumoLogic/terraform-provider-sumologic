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

func TestAccSumologicGenericPollingSource_create(t *testing.T) {
	var s3Source PollingSource
	var cloudWatchSource PollingSource
	var collector Collector
	cName := acctest.RandomWithPrefix("tf-acc-test")
	cDescription := acctest.RandomWithPrefix("tf-acc-test")
	cCategory := acctest.RandomWithPrefix("tf-acc-test")
	sName := acctest.RandomWithPrefix("tf-acc-test")
	sDescription := acctest.RandomWithPrefix("tf-acc-test")
	sCategory := acctest.RandomWithPrefix("tf-acc-test")
	cwName := acctest.RandomWithPrefix("tf-acc-test")
	cwDescription := acctest.RandomWithPrefix("tf-acc-test")
	cwCategory := acctest.RandomWithPrefix("tf-acc-test")
	cloudWatchResourceName := "sumologic_cloudwatch_source.cloudwatch"
	s3ResourceName := "sumologic_s3_source.s3"
	testAwsID := os.Getenv("SUMOLOGIC_TEST_AWS_ID")
	testAwsKey := os.Getenv("SUMOLOGIC_TEST_AWS_KEY")
	testAwsBucket := os.Getenv("SUMOLOGIC_TEST_BUCKET_NAME")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGenericPollingSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicGenericPollingSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsID, testAwsKey, testAwsBucket, cwName, cwDescription, cwCategory, testAwsID, testAwsKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGenericPollingSourceExists(s3ResourceName, &s3Source),
					testAccCheckGenericPollingSourceValues(&s3Source, sName, sDescription, sCategory),
					testAccCheckCollectorExists("sumologic_collector.test", &collector),
					testAccCheckCollectorValues(&collector, cName, cDescription, cCategory, "Etc/UTC", ""),
					testAccCheckGenericPollingSourceExists(cloudWatchResourceName, &cloudWatchSource),
					testAccCheckGenericPollingSourceValues(&cloudWatchSource, cwName, cwDescription, cwCategory),
					resource.TestCheckResourceAttrSet(s3ResourceName, "id"),
					resource.TestCheckResourceAttr(s3ResourceName, "name", sName),
					resource.TestCheckResourceAttr(s3ResourceName, "description", sDescription),
					resource.TestCheckResourceAttr(s3ResourceName, "category", sCategory),
					resource.TestCheckResourceAttr(s3ResourceName, "content_type", "AwsS3Bucket"),
					resource.TestCheckResourceAttrSet(cloudWatchResourceName, "id"),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "name", cwName),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "description", cwDescription),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "category", cwCategory),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "content_type", "AwsCloudWatch"),
				),
			},
		},
	})
}
func TestAccSumologicGenericPollingSource_update(t *testing.T) {
	var s3Source PollingSource
	var cloudWatchSource PollingSource
	cName := acctest.RandomWithPrefix("tf-acc-test")
	cDescription := acctest.RandomWithPrefix("tf-acc-test")
	cCategory := acctest.RandomWithPrefix("tf-acc-test")
	sName := acctest.RandomWithPrefix("tf-acc-test")
	sDescription := acctest.RandomWithPrefix("tf-acc-test")
	sCategory := acctest.RandomWithPrefix("tf-acc-test")
	cwName := acctest.RandomWithPrefix("tf-acc-test")
	cwDescription := acctest.RandomWithPrefix("tf-acc-test")
	cwCategory := acctest.RandomWithPrefix("tf-acc-test")
	sNameUpdated := acctest.RandomWithPrefix("tf-acc-test")
	sDescriptionUpdated := acctest.RandomWithPrefix("tf-acc-test")
	sCategoryUpdated := acctest.RandomWithPrefix("tf-acc-test")
	s3ResourceName := "sumologic_s3_source.s3"
	cloudWatchResourceName := "sumologic_cloudwatch_source.cloudwatch"
	testAwsID := os.Getenv("SUMOLOGIC_TEST_AWS_ID")
	testAwsKey := os.Getenv("SUMOLOGIC_TEST_AWS_KEY")
	testAwsBucket := os.Getenv("SUMOLOGIC_TEST_BUCKET_NAME")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHTTPSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicGenericPollingSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsID, testAwsKey, testAwsBucket, cwName, cwDescription, cwCategory, testAwsID, testAwsKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGenericPollingSourceExists(s3ResourceName, &s3Source),
					testAccCheckGenericPollingSourceValues(&s3Source, sName, sDescription, sCategory),
					testAccCheckGenericPollingSourceExists(cloudWatchResourceName, &cloudWatchSource),
					testAccCheckGenericPollingSourceValues(&cloudWatchSource, cwName, cwDescription, cwCategory),
					resource.TestCheckResourceAttrSet(s3ResourceName, "id"),
					resource.TestCheckResourceAttr(s3ResourceName, "name", sName),
					resource.TestCheckResourceAttr(s3ResourceName, "description", sDescription),
					resource.TestCheckResourceAttr(s3ResourceName, "category", sCategory),
					resource.TestCheckResourceAttr(s3ResourceName, "content_type", "AwsS3Bucket"),
					resource.TestCheckResourceAttrSet(cloudWatchResourceName, "id"),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "name", cwName),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "description", cwDescription),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "category", cwCategory),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "content_type", "AwsCloudWatch"),
				),
			},
			{
				Config: testAccSumologicGenericPollingSourceConfig(cName, cDescription, cCategory, sNameUpdated, sDescriptionUpdated, sCategoryUpdated, testAwsID, testAwsKey, testAwsBucket, cwName, cwDescription, cwCategory, testAwsID, testAwsKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGenericPollingSourceExists(s3ResourceName, &s3Source),
					testAccCheckGenericPollingSourceValues(&s3Source, sNameUpdated, sDescriptionUpdated, sCategoryUpdated),
					resource.TestCheckResourceAttrSet(s3ResourceName, "id"),
					resource.TestCheckResourceAttr(s3ResourceName, "name", sNameUpdated),
					resource.TestCheckResourceAttr(s3ResourceName, "description", sDescriptionUpdated),
					resource.TestCheckResourceAttr(s3ResourceName, "category", sCategoryUpdated),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "content_type", "AwsCloudWatch"),
				),
			},
		},
	})
}
func testAccCheckGenericPollingSourceDestroy(s *terraform.State) error {
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
func testAccCheckGenericPollingSourceExists(n string, pollingSource *PollingSource) resource.TestCheckFunc {
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
func testAccCheckGenericPollingSourceValues(pollingSource *PollingSource, name, description, category string) resource.TestCheckFunc {
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
func testAccSumologicGenericPollingSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsID, testAwsKey, testAwsBucket, cwName, cwDescription, cwCategory, testAwsId2, testAwsKey2 string) string {
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
		type = "S3BucketAuthentication"
		access_key = "%s"
		secret_key = "%s"
	  }
	  path {
		type = "S3BucketPathExpression"
		bucket_name     = "%s"
		path_expression = "*"
	  }
	}

resource "sumologic_cloudwatch_source" "cloudwatch" {
	name = "%s"
	description = "%s"
	category = "%s"
	content_type  = "AwsCloudWatch"
 	scan_interval = 300000
  	paused        = false
	collector_id = "${sumologic_collector.test.id}"
	authentication {
		type = "S3BucketAuthentication"
		access_key = "%s"
		secret_key = "%s"
	  }
	path {
		type = "CloudWatchPath"
		limit_to_regions = ["us-west-2"]
		limit_to_namespaces = ["AWS/Route53","AWS/S3","customNamespace"]
	  }
}
`, cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsID, testAwsKey, testAwsBucket, cwName, cwDescription, cwCategory, testAwsId2, testAwsKey2)
}
