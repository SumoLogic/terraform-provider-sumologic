package sumologic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumoLogicPartition_basic(t *testing.T) {
	testName := acctest.RandString(16)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPartitionDestroy,
		Steps: []resource.TestStep{
			// Create a Partition
			{
				Config: newPartitionConfig(testName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPartitionExists("sumologic_partition.foo"),
					resource.TestCheckResourceAttr(
						"sumologic_partition.foo", "name", "terraform_acctest_"+testName),
					resource.TestCheckResourceAttr(
						"sumologic_partition.foo", "routing_expression", "_sourcecategory=*/Terraform"),
					resource.TestCheckResourceAttr(
						"sumologic_partition.foo", "analytics_tier", "continuous"),
					resource.TestCheckResourceAttr(
						"sumologic_partition.foo", "retention_period", "30"),
					resource.TestCheckResourceAttr(
						"sumologic_partition.foo", "is_compliant", "false"),
				),
			},
			// Update a Partition
			{
				Config: updatePartitionConfig(testName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPartitionExists("sumologic_partition.foo"),
					resource.TestCheckResourceAttr(
						"sumologic_partition.foo", "name", "terraform_acctest_"+testName),
					resource.TestCheckResourceAttr(
						"sumologic_partition.foo", "routing_expression", "_sourcecategory=*/Terraform"),
					resource.TestCheckResourceAttr(
						"sumologic_partition.foo", "analytics_tier", "continuous"),
					resource.TestCheckResourceAttr(
						"sumologic_partition.foo", "retention_period", "365"),
					resource.TestCheckResourceAttr(
						"sumologic_partition.foo", "is_compliant", "false"),
				),
			},
			// allow change in casing of analytics_tier
			{
				Config: updatePartitionAnalyticsTierCase(testName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPartitionExists("sumologic_partition.foo"),
					resource.TestCheckResourceAttr(
						"sumologic_partition.foo", "name", "terraform_acctest_"+testName),
					resource.TestCheckResourceAttr(
						"sumologic_partition.foo", "routing_expression", "_sourcecategory=*/Terraform"),
					resource.TestCheckResourceAttr(
						"sumologic_partition.foo", "analytics_tier", "continuous"),
					resource.TestCheckResourceAttr(
						"sumologic_partition.foo", "retention_period", "366"),
					resource.TestCheckResourceAttr(
						"sumologic_partition.foo", "is_compliant", "false"),
				),
			},
			// Update analytics tier to a different value and assert error
			{
				Config:      updatePartitionAnalyticsTierConfig(testName),
				ExpectError: regexp.MustCompile(`(?i)analytics_tier of a partition can only be updated post creation if partition has been moved to flex tier`),
			},
		},
	})
}

func testAccCheckPartitionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, r := range s.RootModule().Resources {
		id := r.Primary.ID
		p, err := client.GetPartition(id)

		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}

		if p != nil {
			return fmt.Errorf("Partition still exists")
		}
	}
	return nil
}

func testAccCheckPartitionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			if _, err := client.GetPartition(id); err != nil {
				return fmt.Errorf("Received an error retrieving partition %s", err)
			}
		}
		return nil
	}
}

func newPartitionConfig(testName string) string {
	return fmt.Sprintf(`
resource "sumologic_partition" "foo" {
    name = "terraform_acctest_%s"
    routing_expression = "_sourcecategory=*/Terraform"
    is_compliant = false
    retention_period = 30
    analytics_tier = "continuous"
}
`, testName)
}

func updatePartitionConfig(testName string) string {
	return fmt.Sprintf(`
resource "sumologic_partition" "foo" {
    name = "terraform_acctest_%s"
    routing_expression = "_sourcecategory=*/Terraform"
	retention_period = 365
	is_compliant = false
	analytics_tier = "continuous"
}
`, testName)
}

func updatePartitionAnalyticsTierConfig(testName string) string {
	return fmt.Sprintf(`
resource "sumologic_partition" "foo" {
    name = "terraform_acctest_%s"
    routing_expression = "_sourcecategory=*/Terraform"
	retention_period = 365
	is_compliant = false
	analytics_tier = "infrequent" // Changed from "continuous" to trigger error
}
`, testName)
}

func updatePartitionAnalyticsTierCase(testName string) string {
	return fmt.Sprintf(`
	resource "sumologic_partition" "foo" {
		name = "terraform_acctest_%s"
		routing_expression = "_sourcecategory=*/Terraform"
		retention_period = 366
		is_compliant = false
		analytics_tier = "ContinuouS" // Changed case for "continuous"
	}
	`, testName)
}
