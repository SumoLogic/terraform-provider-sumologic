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

func TestAccSumologicCloudWatchSource_create(t *testing.T) {
	var cloudWatchSource PollingSource
	var collector Collector
	cName := acctest.RandomWithPrefix("tf-acc-test")
	cDescription := acctest.RandomWithPrefix("tf-acc-test")
	cCategory := acctest.RandomWithPrefix("tf-acc-test")
	cwName := acctest.RandomWithPrefix("tf-acc-test")
	cwDescription := acctest.RandomWithPrefix("tf-acc-test")
	cwCategory := acctest.RandomWithPrefix("tf-acc-test")
	cloudWatchResourceName := "sumologic_cloudwatch_source.cloudwatch"
	testAwsID := os.Getenv("SUMOLOGIC_TEST_AWS_ID")
	testAwsKey := os.Getenv("SUMOLOGIC_TEST_AWS_KEY")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudWatchSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCloudWatchSourceConfig(cName, cDescription, cCategory, cwName, cwDescription, cwCategory, testAwsID, testAwsKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudWatchSourceExists(cloudWatchResourceName, &cloudWatchSource),
					testAccCheckCloudWatchSourceValues(&cloudWatchSource, cwName, cwDescription, cwCategory),
					testAccCheckCollectorExists("sumologic_collector.test", &collector),
					testAccCheckCollectorValues(&collector, cName, cDescription, cCategory, "Etc/UTC", ""),
					resource.TestCheckResourceAttrSet(cloudWatchResourceName, "id"),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "name", cwName),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "description", cwDescription),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "category", cwCategory),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "content_type", "AwsCloudWatch"),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "path.0.type", "CloudWatchPath"),
				),
			},
		},
	})
}
func TestAccSumologicCloudWatchSource_update(t *testing.T) {
	var cloudWatchSource PollingSource
	cName := acctest.RandomWithPrefix("tf-acc-test")
	cDescription := acctest.RandomWithPrefix("tf-acc-test")
	cCategory := acctest.RandomWithPrefix("tf-acc-test")
	cwName := acctest.RandomWithPrefix("tf-acc-test")
	cwDescription := acctest.RandomWithPrefix("tf-acc-test")
	cwCategory := acctest.RandomWithPrefix("tf-acc-test")
	cwNameUpdated := acctest.RandomWithPrefix("tf-acc-test")
	cwDescriptionUpdated := acctest.RandomWithPrefix("tf-acc-test")
	cwCategoryUpdated := acctest.RandomWithPrefix("tf-acc-test")
	cloudWatchResourceName := "sumologic_cloudwatch_source.cloudwatch"
	testAwsID := os.Getenv("SUMOLOGIC_TEST_AWS_ID")
	testAwsKey := os.Getenv("SUMOLOGIC_TEST_AWS_KEY")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHTTPSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCloudWatchSourceConfig(cName, cDescription, cCategory, cwName, cwDescription, cwCategory, testAwsID, testAwsKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudWatchSourceExists(cloudWatchResourceName, &cloudWatchSource),
					testAccCheckCloudWatchSourceValues(&cloudWatchSource, cwName, cwDescription, cwCategory),
					resource.TestCheckResourceAttrSet(cloudWatchResourceName, "id"),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "name", cwName),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "description", cwDescription),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "category", cwCategory),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "content_type", "AwsCloudWatch"),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "path.0.type", "CloudWatchPath"),
				),
			},
			{
				Config: testAccSumologicCloudWatchSourceConfig(cName, cDescription, cCategory, cwNameUpdated, cwDescriptionUpdated, cwCategoryUpdated, testAwsID, testAwsKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudWatchSourceExists(cloudWatchResourceName, &cloudWatchSource),
					testAccCheckCloudWatchSourceValues(&cloudWatchSource, cwNameUpdated, cwDescriptionUpdated, cwCategoryUpdated),
					resource.TestCheckResourceAttrSet(cloudWatchResourceName, "id"),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "name", cwNameUpdated),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "description", cwDescriptionUpdated),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "category", cwCategoryUpdated),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "content_type", "AwsCloudWatch"),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "path.0.type", "CloudWatchPath"),
				),
			},
		},
	})
}
func testAccCheckCloudWatchSourceDestroy(s *terraform.State) error {
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
func testAccCheckCloudWatchSourceExists(n string, pollingSource *PollingSource) resource.TestCheckFunc {
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
func testAccCheckCloudWatchSourceValues(pollingSource *PollingSource, name, description, category string) resource.TestCheckFunc {
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
func testAccSumologicCloudWatchSourceConfig(cName, cDescription, cCategory, cwName, cwDescription, cwCategory, testAwsId, testAwsKey string) string {
	return fmt.Sprintf(`
resource "sumologic_collector" "test" {
	name = "%s"
	description = "%s"
	category = "%s"
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
		tag_filters {
			type = "TagFilters"
          	namespace = "All"
          	tags = ["k3=v3"]
		}
		tag_filters {
			type = "TagFilters"
          	namespace = "AWS/Route53"
          	tags = ["k1=v1"]
		}
	  }
}
`, cName, cDescription, cCategory, cwName, cwDescription, cwCategory, testAwsId, testAwsKey)
}
