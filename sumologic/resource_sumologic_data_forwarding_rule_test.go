package sumologic

// Todo: (ngoyal, 09-19-2024) - Uncomment the test cases once the issue with PUT permissions on SUMOLOGIC_TEST_BUCKET_NAME is fixed

/*
import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"strings"
	"testing"
)

var indexId string = "00000000024C6155"
var destinationId string = "00000000000732AA"
var indexName = "tf_acctest_partition_data_forwarding_rules_" + acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
var destinationName = "tf_acctest_destination_data_forwarding_rules"

func TestAccSumologicDataForwardingRule_basic(t *testing.T) {
	testAwsRegion, testAwsRoleArn, testAwsBucket := getS3TestParams()

	dataForwardingRuleResourceName := "sumologic_data_forwarding_rule.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithAWS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataForwardingRuleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicDataForwardingDestinationCreateConfig(testAwsBucket, testAwsRoleArn, testAwsRegion),
				Check:  resource.ComposeTestCheckFunc(
					setDestinationId(),
				),
			},
			{
				Config: testPartitionConfig(),
				Check:  resource.ComposeTestCheckFunc(
					setIndexId(),
				),
			},
			{
				// creating a rule
				Config: testAccSumologicDataForwardingRuleCreateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataForwardingRuleExists(),
					resource.TestCheckResourceAttr(dataForwardingRuleResourceName, "index_id", indexId),
					resource.TestCheckResourceAttr(dataForwardingRuleResourceName, "destination_id", destinationId),
					resource.TestCheckResourceAttr(dataForwardingRuleResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataForwardingRuleResourceName, "file_prefix", "test/{index}"),
					resource.TestCheckResourceAttr(dataForwardingRuleResourceName, "payload_schema", "allFields"),
					resource.TestCheckResourceAttr(dataForwardingRuleResourceName, "format", "csv"),
				),
			},
			{
				// updating the rule
				Config: testAccSumologicDataForwardingRuleUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataForwardingRuleExists(),
					resource.TestCheckResourceAttr(dataForwardingRuleResourceName, "index_id", indexId),
					resource.TestCheckResourceAttr(dataForwardingRuleResourceName, "destination_id", destinationId),
					resource.TestCheckResourceAttr(dataForwardingRuleResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataForwardingRuleResourceName, "file_prefix", "test/"),
					resource.TestCheckResourceAttr(dataForwardingRuleResourceName, "payload_schema", "builtInFields"),
					resource.TestCheckResourceAttr(dataForwardingRuleResourceName, "format", "json"),
				),
			},
			{
				// reading the rule
				ResourceName: dataForwardingRuleResourceName,
				ImportState:  true,
			},
			{
				// deleting the rule
				Config: testAccSumologicDataForwardingRuleDeleteConfig(),
				Check:  resource.ComposeTestCheckFunc(testAccCheckDataForwardingRuleDestroy()),
			},
		},
	})
}

func getS3TestParams() (string, string, string) {
	testAwsRoleArn := os.Getenv("SUMOLOGIC_TEST_ROLE_ARN")
	testAwsBucket := os.Getenv("SUMOLOGIC_TEST_BUCKET_NAME")
	testAwsRegion := os.Getenv("SUMOLOGIC_TEST_AWS_REGION")
	return testAwsRegion, testAwsRoleArn, testAwsBucket
}

func setDestinationId() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			if destination, _ := client.getDataForwardingDestination(id); destination != nil {
				destinationId = id
				println("destination id: ", destinationId)
			}
		}
		return nil
	}
}

func setIndexId() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			if Partition, _ := client.GetPartition(id); Partition != nil {
				indexId = id
				println("index id: ", indexId)
			}
		}
		return nil
	}
}

func testAccCheckDataForwardingRuleExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			rule, err := client.getDataForwardingRule(id)
			if err != nil {
				return fmt.Errorf("Received an error retrieving data forwarding rule %s", err)
			}
			if rule != nil {
				println("rule is configured on partition ", id)
			}
		}
		return nil
	}
}

func testAccCheckDataForwardingRuleDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			println("id to be deleted : ", id)

			p1, err1 := client.getDataForwardingDestination(id)
			p2, err2 := client.getDataForwardingRule(id)
			p3, err3 := client.GetPartition(id)

			if err1 != nil {
				if strings.Contains(err1.Error(), "Data forwarding Destination doesn't exists") {
					continue
				}

				return fmt.Errorf("Encountered an error: " + err1.Error())
			}

			if p1 != nil {
				return fmt.Errorf("Data Forwarding Destination still exists!")
			}

			if err2 != nil {
				if strings.Contains(err1.Error(), "Data forwarding Rule doesn't exists") {
					continue
				}

				return fmt.Errorf("Encountered an error: " + err1.Error())
			}

			if p2 != nil {
				return fmt.Errorf("Data Forwarding Rule still exists!")
			}

			if err3 != nil {
				if strings.Contains(err1.Error(), "Partition doesn't exists") {
					continue
				}

				return fmt.Errorf("Encountered an error: " + err1.Error())
			}

			if p3 != nil {
				return fmt.Errorf("Partition still exists!")
			}
		}
		return nil
	}
}

func testAccSumologicDataForwardingDestinationCreateConfig(testAwsBucket string, testAwsRoleArn string, testAwsRegion string) string {
	return fmt.Sprintf(`
resource "sumologic_data_forwarding_destination" "test" {
    destination_name = "%s"
    bucket_name      = "%s"
    authentication {
        type     = "RoleBased"
        role_arn = "%s"
    }
	s3_region = "%s"
	s3_server_side_encryption = false
}
`, destinationName, testAwsBucket, testAwsRoleArn, testAwsRegion)
}

func testPartitionConfig() string {
	return fmt.Sprintf(`
resource "sumologic_partition" "foo" {
    name = "%s"
    routing_expression = "_sourcecategory=abc/Terraform"
    is_compliant = false
    retention_period = 30
	analytics_tier = "flex"
}
`, indexName)
}

func testAccSumologicDataForwardingRuleCreateConfig() string {
	return fmt.Sprintf(`
		resource "sumologic_data_forwarding_rule" "test" {
			index_id = "%s"
			destination_id = "%s"
			enabled = true
			file_prefix = "test/{index}"
			payload_schema = "allFields"
			format = "csv"
		}`, indexId, destinationId)
}

func testAccSumologicDataForwardingRuleUpdateConfig() string {
	return fmt.Sprintf(`
		resource "sumologic_data_forwarding_rule" "test" {
			index_id = "%s"
			destination_id = "%s"
			enabled = true
			file_prefix = "test/"
			payload_schema = "builtInFields"
			format = "json"
		}`, indexId, destinationId)
}

func testAccSumologicDataForwardingRuleDeleteConfig() string {
	return fmt.Sprintf(` `)
}
*/
