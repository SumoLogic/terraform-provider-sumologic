package sumologic

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicMetadataSource_create(t *testing.T) {
	var metadataSource MetadataSource
	var collector Collector
	cName, cDescription, cCategory := getRandomizedParams()
	sName, sDescription, sCategory := getRandomizedParams()
	metadataResourceName := "sumologic_metadata_source.metadata"
	testAwsRoleArn := os.Getenv("SUMOLOGIC_TEST_ROLE_ARN")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetadataSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMetadataSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsRoleArn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetadataSourceExists(metadataResourceName, &metadataSource),
					testAccCheckMetadataSourceValues(&metadataSource, sName, sDescription, sCategory),
					testAccCheckCollectorExists("sumologic_collector.test", &collector),
					testAccCheckCollectorValues(&collector, cName, cDescription, cCategory, "Etc/UTC", ""),
					resource.TestCheckResourceAttrSet(metadataResourceName, "id"),
					resource.TestCheckResourceAttr(metadataResourceName, "name", sName),
					resource.TestCheckResourceAttr(metadataResourceName, "description", sDescription),
					resource.TestCheckResourceAttr(metadataResourceName, "category", sCategory),
					resource.TestCheckResourceAttr(metadataResourceName, "content_type", "AwsMetadata"),
					resource.TestCheckResourceAttr(metadataResourceName, "path.0.type", "AwsMetadataPath"),
				),
			},
		},
	})
}
func TestAccSumologicMetadataSource_update(t *testing.T) {
	var metadataSource MetadataSource
	cName, cDescription, cCategory := getRandomizedParams()
	sName, sDescription, sCategory := getRandomizedParams()
	sNameUpdated, sDescriptionUpdated, sCategoryUpdated := getRandomizedParams()
	metadataResourceName := "sumologic_metadata_source.metadata"
	testAwsRoleArn := os.Getenv("SUMOLOGIC_TEST_ROLE_ARN")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHTTPSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMetadataSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsRoleArn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetadataSourceExists(metadataResourceName, &metadataSource),
					testAccCheckMetadataSourceValues(&metadataSource, sName, sDescription, sCategory),
					resource.TestCheckResourceAttrSet(metadataResourceName, "id"),
					resource.TestCheckResourceAttr(metadataResourceName, "name", sName),
					resource.TestCheckResourceAttr(metadataResourceName, "description", sDescription),
					resource.TestCheckResourceAttr(metadataResourceName, "category", sCategory),
					resource.TestCheckResourceAttr(metadataResourceName, "content_type", "AwsMetadata"),
					resource.TestCheckResourceAttr(metadataResourceName, "path.0.type", "AwsMetadataPath"),
				),
			},
			{
				Config: testAccSumologicMetadataSourceConfig(cName, cDescription, cCategory, sNameUpdated, sDescriptionUpdated, sCategoryUpdated, testAwsRoleArn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetadataSourceExists(metadataResourceName, &metadataSource),
					testAccCheckMetadataSourceValues(&metadataSource, sNameUpdated, sDescriptionUpdated, sCategoryUpdated),
					resource.TestCheckResourceAttrSet(metadataResourceName, "id"),
					resource.TestCheckResourceAttr(metadataResourceName, "name", sNameUpdated),
					resource.TestCheckResourceAttr(metadataResourceName, "description", sDescriptionUpdated),
					resource.TestCheckResourceAttr(metadataResourceName, "category", sCategoryUpdated),
					resource.TestCheckResourceAttr(metadataResourceName, "content_type", "AwsMetadata"),
					resource.TestCheckResourceAttr(metadataResourceName, "path.0.type", "AwsMetadataPath"),
				),
			},
		},
	})
}
func testAccCheckMetadataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_metadata_source" {
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
		s, err := client.GetMetadataSource(collectorID, id)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("Polling Source still exists")
		}
	}
	return nil
}
func testAccCheckMetadataSourceExists(n string, pollingSource *MetadataSource) resource.TestCheckFunc {
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
		pollingSourceResp, err := c.GetMetadataSource(collectorID, id)
		if err != nil {
			return err
		}
		*pollingSource = *pollingSourceResp
		return nil
	}
}
func testAccCheckMetadataSourceValues(pollingSource *MetadataSource, name, description, category string) resource.TestCheckFunc {
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
func testAccSumologicMetadataSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsRoleArn string) string {
	return fmt.Sprintf(`
resource "sumologic_collector" "test" {
	name = "%s"
	description = "%s"
	category = "%s"
}

resource "sumologic_metadata_source" "metadata" {
	name = "%s"
	description = "%s"
	category = "%s"
	content_type  = "AwsMetadata"
 	scan_interval = 300000
  	paused        = false
	collector_id = "${sumologic_collector.test.id}"
	authentication {
		type = "AWSRoleBasedAuthentication"
		role_arn = "%s"
	  }
	path {
		type = "AwsMetadataPath"
		limit_to_regions = ["us-west-2"]
		limit_to_namespaces = ["AWS/EC2"]
		tag_filters = ["Deploy*,", "!DeployStatus,", "Cluster"]
	  }
}
`, cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsRoleArn)
}
