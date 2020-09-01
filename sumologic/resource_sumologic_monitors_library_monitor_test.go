package sumologic

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicMonitorsLibraryMonitor_basic(t *testing.T) {
	var monitorsLibraryMonitor MonitorsLibraryMonitor
	testName := "name-JMH7N"
	// testName := acctest.RandString(16)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorsLibraryMonitorDestroy(monitorsLibraryMonitor),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMonitorsLibraryMonitor(testName),
			},
			{
				ResourceName:      "sumologic_monitor.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccMonitorsLibraryMonitor_create(t *testing.T) {
	var monitorsLibraryMonitor MonitorsLibraryMonitor
	testNameSuffix := acctest.RandString(16)

	testCreatedAt := "2020-09-01T23:15:22.735Z"
	testModifiedAt := "2020-09-01T23:15:22.735Z"
	testCreatedBy := "000000000AD5976D"
	testModifiedBy := "000000000AD5976D"
	testIsLocked := false
	testIsSystem := false
	testIsMutable := true
	testIsDisabled := false
	testName := "terraform_test_monitor_" + testNameSuffix
	testDescription := "terraform_test_monitor_description"
	testMonitorType := "Logs"
	testParentID := "0000000000000002"
	testType := "MonitorsLibraryMonitor"
	testContentType := "Monitor"
	testVersion := 0
	// testNotifications := []MonitorNotification{}
	// testTriggers := []TriggerCondition{}
	// testQueries := []MonitorQuery{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorsLibraryMonitorDestroy(monitorsLibraryMonitor),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMonitorsLibraryMonitor(testName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorsLibraryMonitorExists("sumologic_monitor.test", &monitorsLibraryMonitor, t),
					testAccCheckMonitorsLibraryMonitorAttributes("sumologic_monitor.test"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "modified_at", testModifiedAt),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "created_by", testCreatedBy),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_locked", strconv.FormatBool(testIsLocked)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "monitor_type", testMonitorType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_system", strconv.FormatBool(testIsSystem)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_disabled", strconv.FormatBool(testIsDisabled)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "parent_id", testParentID),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_mutable", strconv.FormatBool(testIsMutable)),
					// resource.TestCheckResourceAttr("sumologic_monitor.test", "notifications.0", strings.Replace(testNotifications[0], "\"", "", 2)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "type", testType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "version", strconv.Itoa(testVersion)),
					// resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0", strings.Replace(testTriggers[0], "\"", "", 2)),
					// resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0", strings.Replace(testQueries[0], "\"", "", 2)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "description", testDescription),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "modified_by", testModifiedBy),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "created_at", testCreatedAt),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "content_type", testContentType),
				),
			},
		},
	})
}

func TestAccMonitorsLibraryMonitor_update(t *testing.T) {
	var monitorsLibraryMonitor MonitorsLibraryMonitor
	testNameSuffix := acctest.RandString(16)

	testCreatedAt := "2020-09-01T23:15:22.735Z"
	testModifiedAt := "2020-09-01T23:15:22.735Z"
	testCreatedBy := "000000000AD5976D"
	testModifiedBy := "000000000AD5976D"
	testIsLocked := false
	testIsSystem := false
	testIsMutable := true
	testIsDisabled := false
	testName := "terraform_test_monitor_" + testNameSuffix
	testDescription := "terraform_test_monitor_description"
	testMonitorType := "Logs"
	testParentID := "0000000000000002"
	testType := "MonitorsLibraryMonitor"
	testContentType := "Monitor"
	testVersion := 0
	// testNotifications := []MonitorNotification{}
	// testTriggers := []TriggerCondition{}
	// testQueries := []MonitorQuery{}

	testUpdatedCreatedAt := "2020-09-01T23:15:22.735Z"
	testUpdatedModifiedAt := "2020-09-01T23:15:22.735Z"
	testUpdatedCreatedBy := "000000000AD5976D"
	testUpdatedModifiedBy := "000000000AD5976D"
	testUpdatedIsLocked := false
	testUpdatedIsSystem := false
	testUpdatedIsMutable := true
	testUpdatedIsDisabled := true
	testUpdatedName := "terraform_test_monitor_" + testNameSuffix
	testUpdatedDescription := "terraform_test_monitor_description"
	testUpdatedMonitorType := "Logs"
	testUpdatedParentID := "0000000000000002"
	testUpdatedType := "MonitorsLibraryMonitor"
	testUpdatedContentType := "Monitor"
	testUpdatedVersion := 0
	// NOTE: updated notifications are webhook instead of email
	// testUpdatedNotifications := []MonitorNotification{}
	// testUpdatedTriggers := []TriggerCondition{}
	// testUpdatedQueries := []MonitorQuery{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorsLibraryMonitorDestroy(monitorsLibraryMonitor),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMonitorsLibraryMonitor(testNameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorsLibraryMonitorExists("sumologic_monitor.test", &monitorsLibraryMonitor, t),
					testAccCheckMonitorsLibraryMonitorAttributes("sumologic_monitor.test"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "modified_at", testModifiedAt),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "created_by", testCreatedBy),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_locked", strconv.FormatBool(testIsLocked)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "monitor_type", testMonitorType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_system", strconv.FormatBool(testIsSystem)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_disabled", strconv.FormatBool(testIsDisabled)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "parent_id", testParentID),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_mutable", strconv.FormatBool(testIsMutable)),
					// resource.TestCheckResourceAttr("sumologic_monitor.test", "notifications", strings.Replace(testNotifications[0], "\"", "", 2)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "type", testType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "version", strconv.Itoa(testVersion)),
					// resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers", strings.Replace(testTriggers[0], "\"", "", 2)),
					// resource.TestCheckResourceAttr("sumologic_monitor.test", "queries", strings.Replace(testQueries[0], "\"", "", 2)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "description", testDescription),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "modified_by", testModifiedBy),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "created_at", testCreatedAt),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "content_type", testContentType),
				),
			},
			{
				Config: testAccSumologicMonitorsLibraryMonitorUpdate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_monitor.test", "modified_at", testUpdatedModifiedAt),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "created_by", testUpdatedCreatedBy),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_locked", strconv.FormatBool(testUpdatedIsLocked)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "monitor_type", testUpdatedMonitorType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_system", strconv.FormatBool(testUpdatedIsSystem)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_disabled", strconv.FormatBool(testUpdatedIsDisabled)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "name", testUpdatedName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "parent_id", testUpdatedParentID),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_mutable", strconv.FormatBool(testUpdatedIsMutable)),
					// resource.TestCheckResourceAttr("sumologic_monitor.test", "notifications.0", strings.Replace(testUpdatedNotifications[0], "\"", "", 2)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "type", testUpdatedType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "version", strconv.Itoa(testUpdatedVersion)),
					// resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0", strings.Replace(testUpdatedTriggers[0], "\"", "", 2)),
					// resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0", strings.Replace(testUpdatedQueries[0], "\"", "", 2)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "description", testUpdatedDescription),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "modified_by", testUpdatedModifiedBy),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "created_at", testUpdatedCreatedAt),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "content_type", testUpdatedContentType),
				),
			},
		},
	})
}

func testAccCheckMonitorsLibraryMonitorDestroy(monitorsLibraryMonitor MonitorsLibraryMonitor) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			u, err := client.MonitorsRead(id)
			if err != nil {
				return fmt.Errorf("Encountered an error: " + err.Error())
			}
			if u != nil {
				return fmt.Errorf("MonitorsLibraryMonitor %s still exists", id)
			}
		}
		return nil
	}
}

func testAccCheckMonitorsLibraryMonitorExists(name string, monitorsLibraryMonitor *MonitorsLibraryMonitor, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			//need this so that we don't get an unused import error for strconv in some cases
			return fmt.Errorf("Error = %s. MonitorsLibraryMonitor not found: %s", strconv.FormatBool(ok), name)
		}

		//need this so that we don't get an unused import error for strings in some cases
		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("MonitorsLibraryMonitor ID is not set")
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)
		newMonitorsLibraryMonitor, err := client.MonitorsRead(id)
		if err != nil {
			return fmt.Errorf("MonitorsLibraryMonitor %s not found", id)
		}
		monitorsLibraryMonitor = newMonitorsLibraryMonitor
		return nil
	}
}

func testAccCheckMonitorsLibraryMonitorAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(name, "modified_at"),
			resource.TestCheckResourceAttrSet(name, "created_by"),
			resource.TestCheckResourceAttrSet(name, "is_locked"),
			resource.TestCheckResourceAttrSet(name, "monitor_type"),
			resource.TestCheckResourceAttrSet(name, "is_system"),
			resource.TestCheckResourceAttrSet(name, "name"),
			resource.TestCheckResourceAttrSet(name, "parent_id"),
			resource.TestCheckResourceAttrSet(name, "is_mutable"),
			resource.TestCheckResourceAttrSet(name, "notifications"),
			resource.TestCheckResourceAttrSet(name, "type"),
			resource.TestCheckResourceAttrSet(name, "version"),
			resource.TestCheckResourceAttrSet(name, "triggers"),
			resource.TestCheckResourceAttrSet(name, "queries"),
			resource.TestCheckResourceAttrSet(name, "description"),
			resource.TestCheckResourceAttrSet(name, "modified_by"),
			resource.TestCheckResourceAttrSet(name, "created_at"),
			resource.TestCheckResourceAttrSet(name, "content_type"),
		)
		return f(s)
	}
}

func testAccSumologicMonitorsLibraryMonitor(testName string) string {
	return fmt.Sprintf(`
resource "sumologic_monitor" "test" {
	name = "terraform_test_monitor_%s"
	description = "terraform_test_monitor_description"
	type = "MonitorsLibraryMonitor"
	parent_id = "0000000000000002"
	is_disabled = false
	content_type = "Monitor"
	monitor_type = "Logs"
	queries {
		row_id = "A"
		query = "_sourceCategory=monitor-manager error"
	}
	triggers  {
		threshold_type = "GreaterThan"
		threshold = 40.0
		time_range = "15m"
		occurrence_type = "ResultCount"
		trigger_source = "AllResults"
		trigger_type = "Critical"
		detection_method = "StaticCondition"
	}
	notifications {
		notification {
			action_type = "NamedConnectionAction"
			connection_id = "0000000000006A7E"
			payload_override = "{}"
		}
		run_for_trigger_types = ["Critical"]
	}
}`, testName)
}

func testAccSumologicMonitorsLibraryMonitorUpdate(testName string) string {
	return fmt.Sprintf(`
resource "sumologic_monitor" "test" {
	name = "terraform_test_monitor_%s"
	description = "terraform_test_monitor_description"
	type = "MonitorsLibraryMonitor"
	parent_id = "0000000000000002"
	is_disabled = true
	content_type = "Monitor"
	monitor_type = "Logs"
	queries {
		row_id = "A"
		query = "_sourceCategory=monitor-manager error"
	}
	triggers  {
		threshold_type = "GreaterThan"
		threshold = 40.0
		time_range = "15m"
		occurrence_type = "ResultCount"
		trigger_source = "AllResults"
		trigger_type = "Critical"
		detection_method = "StaticCondition"
	}
	notifications {
		notification {
			action_type = "EmailAction"
			recipients = ["abc@example.com"]
			subject = "test tf metrics monitor"
			time_zone = "PST"
			message_body = "metrics_test"
		}
		run_for_trigger_types = ["Critical","ResolvedCritical"]
	}
}`, testName)
}
