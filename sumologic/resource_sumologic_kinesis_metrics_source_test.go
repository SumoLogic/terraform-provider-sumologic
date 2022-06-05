package sumologic

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicKinesisMetricsSource_create(t *testing.T) {
	var kinesisMetricsSource KinesisMetricsSource
	var collector Collector
	cName, cDescription, cCategory := getRandomizedParams()
	sName, sDescription, sCategory := getRandomizedParams()
	kinesisMetricsResourceName := "sumologic_kinesis_metrics_source.kinesisMetrics"
	testAwsRoleArn := os.Getenv("SUMOLOGIC_TEST_ROLE_ARN")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKinesisMetricsSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicKinesisMetricsSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsRoleArn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisMetricsSourceExists(kinesisMetricsResourceName, &kinesisMetricsSource),
					testAccCheckKinesisMetricsSourceValues(&kinesisMetricsSource, sName, sDescription, sCategory),
					testAccCheckCollectorExists("sumologic_collector.test", &collector),
					testAccCheckCollectorValues(&collector, cName, cDescription, cCategory, "Etc/UTC", ""),
					resource.TestCheckResourceAttrSet(kinesisMetricsResourceName, "id"),
					resource.TestCheckResourceAttr(kinesisMetricsResourceName, "name", sName),
					resource.TestCheckResourceAttr(kinesisMetricsResourceName, "description", sDescription),
					resource.TestCheckResourceAttr(kinesisMetricsResourceName, "category", sCategory),
					resource.TestCheckResourceAttr(kinesisMetricsResourceName, "content_type", "KinesisMetric"),
					resource.TestCheckResourceAttr(kinesisMetricsResourceName, "path.0.type", "KinesisMetricPath"),
				),
			},
		},
	})
}
func TestAccSumologicKinesisMetricsSource_update(t *testing.T) {
	var kinesisMetricsSource KinesisMetricsSource
	cName, cDescription, cCategory := getRandomizedParams()
	sName, sDescription, sCategory := getRandomizedParams()
	sNameUpdated, sDescriptionUpdated, sCategoryUpdated := getRandomizedParams()
	kinesisMetricsResourceName := "sumologic_kinesis_metrics_source.kinesisMetrics"
	testAwsRoleArn := os.Getenv("SUMOLOGIC_TEST_ROLE_ARN")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHTTPSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicKinesisMetricsSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsRoleArn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisMetricsSourceExists(kinesisMetricsResourceName, &kinesisMetricsSource),
					testAccCheckKinesisMetricsSourceValues(&kinesisMetricsSource, sName, sDescription, sCategory),
					resource.TestCheckResourceAttrSet(kinesisMetricsResourceName, "id"),
					resource.TestCheckResourceAttr(kinesisMetricsResourceName, "name", sName),
					resource.TestCheckResourceAttr(kinesisMetricsResourceName, "description", sDescription),
					resource.TestCheckResourceAttr(kinesisMetricsResourceName, "category", sCategory),
					resource.TestCheckResourceAttr(kinesisMetricsResourceName, "content_type", "KinesisMetric"),
					resource.TestCheckResourceAttr(kinesisMetricsResourceName, "path.0.type", "KinesisMetricPath"),
				),
			},
			{
				Config: testAccSumologicKinesisMetricsSourceConfig(cName, cDescription, cCategory, sNameUpdated, sDescriptionUpdated, sCategoryUpdated, testAwsRoleArn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisMetricsSourceExists(kinesisMetricsResourceName, &kinesisMetricsSource),
					testAccCheckKinesisMetricsSourceValues(&kinesisMetricsSource, sNameUpdated, sDescriptionUpdated, sCategoryUpdated),
					resource.TestCheckResourceAttrSet(kinesisMetricsResourceName, "id"),
					resource.TestCheckResourceAttr(kinesisMetricsResourceName, "name", sNameUpdated),
					resource.TestCheckResourceAttr(kinesisMetricsResourceName, "description", sDescriptionUpdated),
					resource.TestCheckResourceAttr(kinesisMetricsResourceName, "category", sCategoryUpdated),
					resource.TestCheckResourceAttr(kinesisMetricsResourceName, "content_type", "KinesisMetric"),
					resource.TestCheckResourceAttr(kinesisMetricsResourceName, "path.0.type", "KinesisMetricPath"),
				),
			},
		},
	})
}
func testAccCheckKinesisMetricsSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_kinesis_metrics_source" {
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
		s, err := client.GetKinesisMetricsSource(collectorID, id)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("KinesisMetric Source still exists")
		}
	}
	return nil
}
func testAccCheckKinesisMetricsSourceExists(n string, kinesisMetricsSource *KinesisMetricsSource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("KinesisMetrics Source ID is not set")
		}
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("KinesisMetrics Source id should be int; got %s", rs.Primary.ID)
		}
		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		c := testAccProvider.Meta().(*Client)
		kinesisMetricsSourceResp, err := c.GetKinesisMetricsSource(collectorID, id)
		if err != nil {
			return err
		}
		*kinesisMetricsSource = *kinesisMetricsSourceResp
		return nil
	}
}
func testAccCheckKinesisMetricsSourceValues(kinesisMetricsSource *KinesisMetricsSource, name, description, category string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if kinesisMetricsSource.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, kinesisMetricsSource.Name)
		}
		if kinesisMetricsSource.Description != description {
			return fmt.Errorf("bad description, expected \"%s\", got: %#v", description, kinesisMetricsSource.Description)
		}
		if kinesisMetricsSource.Category != category {
			return fmt.Errorf("bad category, expected \"%s\", got: %#v", category, kinesisMetricsSource.Category)
		}
		return nil
	}
}
func testAccSumologicKinesisMetricsSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsRoleArn string) string {
	return fmt.Sprintf(`
resource "sumologic_collector" "test" {
	name = "%s"
	description = "%s"
	category = "%s"
}

resource "sumologic_kinesis_metrics_source" "kinesisMetrics" {
	name = "%s"
	description = "%s"
	category = "%s"
	content_type  = "KinesisMetric"
	collector_id = "${sumologic_collector.test.id}"
	authentication {
		type = "AWSRoleBasedAuthentication"
		role_arn = "%s"
	  }
	path {
		type = "KinesisMetricPath"
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
`, cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsRoleArn)
}
