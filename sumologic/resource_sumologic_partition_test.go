package sumologic

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSumoLogicPartition(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	testName := strconv.Itoa(rand.Int())
	resource.Test(t, resource.TestCase{
		// PreCheck:     func() { TestAccPreCheck(t) },
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
						"sumologic_partition.foo", "analytics_tier", "enhanced"),
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
						"sumologic_partition.foo", "analytics_tier", "enhanced"),
					resource.TestCheckResourceAttr(
						"sumologic_partition.foo", "retention_period", "365"),
					resource.TestCheckResourceAttr(
						"sumologic_partition.foo", "is_compliant", "false"),
				),
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
    analytics_tier = "enhanced"
    is_compliant = false
}
`, testName)
}

func updatePartitionConfig(testName string) string {
	return fmt.Sprintf(`
resource "sumologic_partition" "foo" {
    name = "terraform_acctest_%s"
    routing_expression = "_sourcecategory=*/Terraform"
	analytics_tier = "enhanced"
	retention_period = 365
	is_compliant = false
}
`, testName)
}
