package sumologic

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicKinesisLogSource_create(t *testing.T) {
	var kinesisLogSource KinesisLogSource
	var collector Collector
	cName, cDescription, cCategory := getRandomizedParams()
	sName, sDescription, sCategory := getRandomizedParams()
	kinesisLogResourceName := "sumologic_kinesis_log_source.kinesisLog"
	testAwsRoleArn := os.Getenv("SUMOLOGIC_TEST_ROLE_ARN")
	testAwsBucket := os.Getenv("SUMOLOGIC_TEST_BUCKET_NAME")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKinesisLogSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicKinesisLogSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsRoleArn, testAwsBucket),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisLogSourceExists(kinesisLogResourceName, &kinesisLogSource),
					testAccCheckKinesisLogSourceValues(&kinesisLogSource, sName, sDescription, sCategory),
					testAccCheckCollectorExists("sumologic_collector.test", &collector),
					testAccCheckCollectorValues(&collector, cName, cDescription, cCategory, "Etc/UTC", ""),
					resource.TestCheckResourceAttrSet(kinesisLogResourceName, "id"),
					resource.TestCheckResourceAttr(kinesisLogResourceName, "name", sName),
					resource.TestCheckResourceAttr(kinesisLogResourceName, "description", sDescription),
					resource.TestCheckResourceAttr(kinesisLogResourceName, "category", sCategory),
					resource.TestCheckResourceAttr(kinesisLogResourceName, "content_type", "KinesisLog"),
					resource.TestCheckResourceAttr(kinesisLogResourceName, "path.0.type", "KinesisLogPath"),
				),
			},
		},
	})
}
func TestAccSumologicKinesisLogSource_update(t *testing.T) {
	var kinesisLogSource KinesisLogSource
	cName, cDescription, cCategory := getRandomizedParams()
	sName, sDescription, sCategory := getRandomizedParams()
	sNameUpdated, sDescriptionUpdated, sCategoryUpdated := getRandomizedParams()
	kinesisLogResourceName := "sumologic_kinesis_Log_source.kinesisLog"
	testAwsRoleArn := os.Getenv("SUMOLOGIC_TEST_ROLE_ARN")
	testAwsBucket := os.Getenv("SUMOLOGIC_TEST_BUCKET_NAME")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHTTPSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicKinesisLogSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsRoleArn, testAwsBucket),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisLogSourceExists(kinesisLogResourceName, &kinesisLogSource),
					testAccCheckKinesisLogSourceValues(&kinesisLogSource, sName, sDescription, sCategory),
					resource.TestCheckResourceAttrSet(kinesisLogResourceName, "id"),
					resource.TestCheckResourceAttr(kinesisLogResourceName, "name", sName),
					resource.TestCheckResourceAttr(kinesisLogResourceName, "description", sDescription),
					resource.TestCheckResourceAttr(kinesisLogResourceName, "category", sCategory),
					resource.TestCheckResourceAttr(kinesisLogResourceName, "content_type", "KinesisLog"),
					resource.TestCheckResourceAttr(kinesisLogResourceName, "path.0.type", "KinesisLogPath"),
				),
			},
			{
				Config: testAccSumologicKinesisLogSourceConfig(cName, cDescription, cCategory, sNameUpdated, sDescriptionUpdated, sCategoryUpdated, testAwsRoleArn, testAwsBucket),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKinesisLogSourceExists(kinesisLogResourceName, &kinesisLogSource),
					testAccCheckKinesisLogSourceValues(&kinesisLogSource, sNameUpdated, sDescriptionUpdated, sCategoryUpdated),
					resource.TestCheckResourceAttrSet(kinesisLogResourceName, "id"),
					resource.TestCheckResourceAttr(kinesisLogResourceName, "name", sNameUpdated),
					resource.TestCheckResourceAttr(kinesisLogResourceName, "description", sDescriptionUpdated),
					resource.TestCheckResourceAttr(kinesisLogResourceName, "category", sCategoryUpdated),
					resource.TestCheckResourceAttr(kinesisLogResourceName, "content_type", "KinesisLog"),
					resource.TestCheckResourceAttr(kinesisLogResourceName, "path.0.type", "KinesisLogPath"),
				),
			},
		},
	})
}
func testAccCheckKinesisLogSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_kinesis_log_source" {
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
		s, err := client.GetKinesisLogSource(collectorID, id)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("KinesisLog Source still exists")
		}
	}
	return nil
}
func testAccCheckKinesisLogSourceExists(n string, kinesisLogSource *KinesisLogSource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("KinesisLog Source ID is not set")
		}
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("KinesisLog Source id should be int; got %s", rs.Primary.ID)
		}
		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		c := testAccProvider.Meta().(*Client)
		kinesisLogSourceResp, err := c.GetKinesisLogSource(collectorID, id)
		if err != nil {
			return err
		}
		*kinesisLogSource = *kinesisLogSourceResp
		return nil
	}
}
func testAccCheckKinesisLogSourceValues(kinesisLogSource *KinesisLogSource, name, description, category string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if kinesisLogSource.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, kinesisLogSource.Name)
		}
		if kinesisLogSource.Description != description {
			return fmt.Errorf("bad description, expected \"%s\", got: %#v", description, kinesisLogSource.Description)
		}
		if kinesisLogSource.Category != category {
			return fmt.Errorf("bad category, expected \"%s\", got: %#v", category, kinesisLogSource.Category)
		}
		return nil
	}
}
func testAccSumologicKinesisLogSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsRoleArn, testAwsBucket string) string {
	return fmt.Sprintf(`
resource "sumologic_collector" "test" {
	name = "%s"
	description = "%s"
	category = "%s"
}

resource "sumologic_kinesis_log_source" "kinesisLog" {
	name = "%s"
	description = "%s"
	category = "%s"
	content_type  = "KinesisLog"
	collector_id = "${sumologic_collector.test.id}"
	authentication {
    type = "AWSRoleBasedAuthentication"
    role_arn = "%s"
  }
	path {
		type = "KinesisLogPath"
    bucket_name     = "%s"
    path_expression = "http-endpoint-failed/*"
    scan_interval   = 30000
	}
}
`, cName, cDescription, cCategory, sName, sDescription, sCategory, testAwsRoleArn, testAwsBucket)
}

func testAccSumologicKinesisLogSourceNoAuthConfig(cName, cDescription, cCategory, sName, sDescription, sCategory string) string {
	return fmt.Sprintf(`
resource "sumologic_collector" "test" {
	name = "%s"
	description = "%s"
	category = "%s"
}

resource "sumologic_kinesis_log_source" "kinesisLog" {
	name = "%s"
	description = "%s"
	category = "%s"
	content_type  = "KinesisLog"
	collector_id = "${sumologic_collector.test.id}"
	authentication {
    type = "NoAuthentication"
  }
	path {
		type = "NoPathExpression"
	}
}
`, cName, cDescription, cCategory, sName, sDescription, sCategory)
}
