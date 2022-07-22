package sumologic

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicCloudWatchSource_create(t *testing.T) {
	var cloudWatchSource PollingSource
	var collector Collector
	cName, cDescription, cCategory := getRandomizedParams()
	sName, sDescription, sCategory := getRandomizedParams()
	cloudWatchResourceName := "sumologic_cloudwatch_source.cloudwatch"
	testAwsRoleArn := os.Getenv("SUMOLOGIC_TEST_ROLE_ARN")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudWatchSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCloudWatchSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsRoleArn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudWatchSourceExists(cloudWatchResourceName, &cloudWatchSource),
					testAccCheckCloudWatchSourceValues(&cloudWatchSource, sName, sDescription, sCategory),
					testAccCheckCollectorExists("sumologic_collector.test", &collector),
					testAccCheckCollectorValues(&collector, cName, cDescription, cCategory, "Etc/UTC", ""),
					resource.TestCheckResourceAttrSet(cloudWatchResourceName, "id"),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "name", sName),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "description", sDescription),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "category", sCategory),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "content_type", "AwsCloudWatch"),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "path.0.type", "CloudWatchPath"),
				),
			},
		},
	})
}
func TestAccSumologicCloudWatchSource_update(t *testing.T) {
	var cloudWatchSource PollingSource
	cName, cDescription, cCategory := getRandomizedParams()
	sName, sDescription, sCategory := getRandomizedParams()
	sNameUpdated, sDescriptionUpdated, sCategoryUpdated := getRandomizedParams()
	cloudWatchResourceName := "sumologic_cloudwatch_source.cloudwatch"
	testAwsRoleArn := os.Getenv("SUMOLOGIC_TEST_ROLE_ARN")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHTTPSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCloudWatchSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsRoleArn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudWatchSourceExists(cloudWatchResourceName, &cloudWatchSource),
					testAccCheckCloudWatchSourceValues(&cloudWatchSource, sName, sDescription, sCategory),
					resource.TestCheckResourceAttrSet(cloudWatchResourceName, "id"),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "name", sName),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "description", sDescription),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "category", sCategory),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "content_type", "AwsCloudWatch"),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "path.0.type", "CloudWatchPath"),
				),
			},
			{
				Config: testAccSumologicCloudWatchSourceConfig(cName, cDescription, cCategory, sNameUpdated, sDescriptionUpdated, sCategoryUpdated, testAwsRoleArn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudWatchSourceExists(cloudWatchResourceName, &cloudWatchSource),
					testAccCheckCloudWatchSourceValues(&cloudWatchSource, sNameUpdated, sDescriptionUpdated, sCategoryUpdated),
					resource.TestCheckResourceAttrSet(cloudWatchResourceName, "id"),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "name", sNameUpdated),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "description", sDescriptionUpdated),
					resource.TestCheckResourceAttr(cloudWatchResourceName, "category", sCategoryUpdated),
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
func testAccSumologicCloudWatchSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsRoleArn string) string {
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
		type = "AWSRoleBasedAuthentication"
		role_arn = "%s"
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
		use_versioned_api = false
	  }
}
`, cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsRoleArn)
}
