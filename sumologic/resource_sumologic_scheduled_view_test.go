package sumologic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccSumologicScheduledView_basic tests the creation and basic reading of a scheduled view.
func TestAccSumologicScheduledView_basic(t *testing.T) {
	resourceName := "sumologic_scheduled_view.test_scheduled_view"
	name := acctest.RandomWithPrefix("tf_test_scheduled_view_")
	nameNoTimeZone := acctest.RandomWithPrefix("tf_test_scheduled_view_no_tz")
	initialQuery := "_sourceCategory=terraform/test/scheduledview/basic | count by _source"
	initialTimeZone := "America/Los_Angeles"
	initialRetention := 30 // days

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSumologicScheduledViewDestroy,
		Steps: []resource.TestStep{
			// Test 1: Create a scheduled view with all initial fields, including time_zone.
			{
				Config: testAccSumologicScheduledViewConfig_Create(
					name,
					initialQuery,
					initialTimeZone,
					initialRetention,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "index_name", name),
					resource.TestCheckResourceAttr(resourceName, "query", initialQuery),
					resource.TestCheckResourceAttr(resourceName, "time_zone", initialTimeZone),
					resource.TestCheckResourceAttr(resourceName, "retention_period", fmt.Sprintf("%d", initialRetention)),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Test 2: Update only time_zone and retention_period.
			// Other fields (name, query, index_alias) should remain unchanged.
			{
				Config: testAccSumologicScheduledViewConfig_Update(
					name,            // Keep original name
					initialQuery,    // Keep original query
					"Europe/Berlin", // NEW time_zone
					60,              // NEW retention_period
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "index_name", name),
					resource.TestCheckResourceAttr(resourceName, "query", initialQuery), // Verify query is unchanged
					resource.TestCheckResourceAttr(resourceName, "time_zone", "Europe/Berlin"),
					resource.TestCheckResourceAttr(resourceName, "retention_period", "60"),
				),
			},
			// Test 3: Create a scheduled view WITHOUT time_zone initially (check default behavior)
			{
				Config: testAccSumologicScheduledViewConfig_CreateNoTimeZone(
					nameNoTimeZone,
					initialQuery,
					initialRetention,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "index_name", nameNoTimeZone),
					resource.TestCheckResourceAttr(resourceName, "query", initialQuery),
					resource.TestCheckResourceAttr(resourceName, "retention_period", fmt.Sprintf("%d", initialRetention)),
					// defaults to "UTC" if time_zone not provided
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC"),
				),
			},
			// Test 4: For the no-time-zone scheduled view, now set the time_zone and update retention.
			{
				Config: testAccSumologicScheduledViewConfig_Update(
					nameNoTimeZone, // Keep original name
					initialQuery,   // Keep original query
					"Asia/Kolkata", // NEW time_zone
					90,             // NEW retention_period
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "index_name", nameNoTimeZone),
					resource.TestCheckResourceAttr(resourceName, "query", initialQuery), // Verify query is unchanged
					resource.TestCheckResourceAttr(resourceName, "time_zone", "Asia/Kolkata"),
					resource.TestCheckResourceAttr(resourceName, "retention_period", "90"),
				),
			},
		},
	})
}

// testAccCheckSumologicScheduledViewDestroy verifies that the scheduled view is deleted upon destroy.
func testAccCheckSumologicScheduledViewDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_scheduled_view" {
			continue
		}
		fmt.Printf("DEBUG: Checking for id: %s", rs.Primary.ID)
		sv, err := client.GetScheduledView(rs.Primary.ID)
		if sv != nil {
			return fmt.Errorf("scheduled view (%s) still exists", rs.Primary.ID)
		}

		if err != nil && !regexp.MustCompile(`not found|404|view:scheduled_view_not_found`).MatchString(err.Error()) {
			return fmt.Errorf("error checking scheduled view %s: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

// testAccSumologicScheduledViewConfig_Create generates the HCL for initial creation.
func testAccSumologicScheduledViewConfig_Create(name, query, timeZone string, retention int) string {
	return fmt.Sprintf(`
resource "sumologic_scheduled_view" "test_scheduled_view" {
  index_name               = "%s"
  query                    = "%s"
  start_time               = "2024-01-01T00:00:00Z"
  retention_period         = %d
  time_zone                = "%s"
}
`, name, query, retention, timeZone)
}

// testAccSumologicScheduledViewConfig_CreateNoTimeZone generates HCL for creation without time_zone.
func testAccSumologicScheduledViewConfig_CreateNoTimeZone(name, query string, retention int) string {
	return fmt.Sprintf(`
resource "sumologic_scheduled_view" "test_scheduled_view" {
  index_name               = "%s"
  query                    = "%s"
  start_time               = "2024-01-01T00:00:00Z"
  retention_period         = %d
  # time_zone intentionally omitted to test default behavior
}
`, name, query, retention)
}

// testAccSumologicScheduledViewConfig_Update generates HCL for updating specific fields.
// It assumes that 'name', 'query', and 'indexAlias' are part of the resource config
// but are NOT changed in the update scenario (they are just passed through from initial state).
func testAccSumologicScheduledViewConfig_Update(name, query, newTimeZone string, newRetention int) string {
	return fmt.Sprintf(`
resource "sumologic_scheduled_view" "test_scheduled_view" {
  index_name                           = "%s"
  query                                = "%s"
  start_time                           = "2024-01-01T00:00:00Z" # Keep start_time consistent
  retention_period                     = %d # Updated retention_period
  reduce_retention_period_immediately  = false
  time_zone                            = "%s" # Updated time_zone
}
`, name, query, newRetention, newTimeZone)
}
