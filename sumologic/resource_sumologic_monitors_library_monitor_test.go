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
	testNameSuffix := acctest.RandString(16)

	testName := "terraform_test_monitor_" + testNameSuffix

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
				ImportStateVerify: false,
			},
		},
	})
}
func TestAccSumologicMonitorsLibraryMonitor_create(t *testing.T) {
	var monitorsLibraryMonitor MonitorsLibraryMonitor
	testNameSuffix := acctest.RandString(16)

	testName := "terraform_test_monitor_" + testNameSuffix
	testDescription := "terraform_test_monitor_description"
	testType := "MonitorsLibraryMonitor"
	testContentType := "Monitor"
	testMonitorType := "Logs"
	testIsDisabled := false
	testQueries := []MonitorQuery{
		{
			RowID: "A",
			Query: "_sourceCategory=monitor-manager error",
		},
	}
	testTriggers := []TriggerCondition{
		{
			ThresholdType:   "GreaterThan",
			Threshold:       40.0,
			TimeRange:       "15m",
			OccurrenceType:  "ResultCount",
			TriggerSource:   "AllResults",
			TriggerType:     "Critical",
			DetectionMethod: "StaticCondition",
		},
		{
			ThresholdType:   "LessThanOrEqual",
			Threshold:       40.0,
			TimeRange:       "15m",
			OccurrenceType:  "ResultCount",
			TriggerSource:   "AllResults",
			TriggerType:     "ResolvedCritical",
			DetectionMethod: "StaticCondition",
		},
	}
	recipients := []string{"abc@example.com"}
	testRecipients := make([]interface{}, len(recipients))
	for i, v := range recipients {
		testRecipients[i] = v
	}
	triggerTypes := []string{"Critical", "ResolvedCritical"}
	testTriggerTypes := make([]interface{}, len(triggerTypes))
	for i, v := range triggerTypes {
		testTriggerTypes[i] = v
	}
	testNotificationAction := EmailNotification{
		ConnectionType: "Email",
		Recipients:     testRecipients,
		Subject:        "test tf monitor",
		TimeZone:       "PST",
		MessageBody:    "test",
	}
	testNotifications := []MonitorNotification{
		{
			Notification:       testNotificationAction,
			RunForTriggerTypes: testTriggerTypes,
		},
	}

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
					resource.TestCheckResourceAttr("sumologic_monitor.test", "monitor_type", testMonitorType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_disabled", strconv.FormatBool(testIsDisabled)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "type", testType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "description", testDescription),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "content_type", testContentType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0.row_id", testQueries[0].RowID),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.trigger_type", testTriggers[0].TriggerType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notifications.0.notification.0.connection_type", testNotifications[0].Notification.(EmailNotification).ConnectionType),
				),
			},
		},
	})
}

func TestAccSumologicMonitorsLibraryMonitor_update(t *testing.T) {
	var monitorsLibraryMonitor MonitorsLibraryMonitor
	testNameSuffix := acctest.RandString(16)

	testName := "terraform_test_monitor_" + testNameSuffix
	testDescription := "terraform_test_monitor_description"
	testType := "MonitorsLibraryMonitor"
	testContentType := "Monitor"
	testMonitorType := "Logs"
	testIsDisabled := false
	testQueries := []MonitorQuery{
		{
			RowID: "A",
			Query: "_sourceCategory=monitor-manager error",
		},
	}
	testTriggers := []TriggerCondition{
		{
			ThresholdType:   "GreaterThan",
			Threshold:       40.0,
			TimeRange:       "15m",
			OccurrenceType:  "ResultCount",
			TriggerSource:   "AllResults",
			TriggerType:     "Critical",
			DetectionMethod: "StaticCondition",
		},
		{
			ThresholdType:   "LessThanOrEqual",
			Threshold:       40.0,
			TimeRange:       "15m",
			OccurrenceType:  "ResultCount",
			TriggerSource:   "AllResults",
			TriggerType:     "ResolvedCritical",
			DetectionMethod: "StaticCondition",
		},
	}
	recipients := []string{"abc@example.com"}
	testRecipients := make([]interface{}, len(recipients))
	for i, v := range recipients {
		testRecipients[i] = v
	}
	triggerTypes := []string{"Critical"}
	testTriggerTypes := make([]interface{}, len(triggerTypes))
	for i, v := range triggerTypes {
		testTriggerTypes[i] = v
	}
	testNotificationAction := EmailNotification{
		ConnectionType: "Email",
		Recipients:     testRecipients,
		Subject:        "test tf monitor",
		TimeZone:       "PST",
		MessageBody:    "test",
	}
	testNotifications := []MonitorNotification{
		{
			Notification:       testNotificationAction,
			RunForTriggerTypes: testTriggerTypes,
		},
	}

	// updated fields
	testUpdatedName := "terraform_test_monitor_" + testNameSuffix
	testUpdatedDescription := "terraform_test_monitor_description"
	testUpdatedType := "MonitorsLibraryMonitor"
	testUpdatedContentType := "Monitor"
	testUpdatedMonitorType := "Logs"
	testUpdatedIsDisabled := true
	testUpdatedQueries := []MonitorQuery{
		{
			RowID: "A",
			Query: "_sourceCategory=monitor-manager info",
		},
	}
	testUpdatedTriggers := []TriggerCondition{
		{
			ThresholdType:   "GreaterThan",
			Threshold:       40.0,
			TimeRange:       "15m",
			OccurrenceType:  "ResultCount",
			TriggerSource:   "AllResults",
			TriggerType:     "Critical",
			DetectionMethod: "StaticCondition",
		},
		{
			ThresholdType:   "LessThanOrEqual",
			Threshold:       40.0,
			TimeRange:       "15m",
			OccurrenceType:  "ResultCount",
			TriggerSource:   "AllResults",
			TriggerType:     "ResolvedCritical",
			DetectionMethod: "StaticCondition",
		},
	}
	updatedRecipients := []string{"abc@example.com"}
	testUpdatedRecipients := make([]interface{}, len(updatedRecipients))
	for i, v := range updatedRecipients {
		testUpdatedRecipients[i] = v
	}
	updatedTriggerTypes := []string{"Critical", "ResolvedCritical"}
	testUpdatedTriggerTypes := make([]interface{}, len(updatedTriggerTypes))
	for i, v := range updatedTriggerTypes {
		testUpdatedTriggerTypes[i] = v
	}
	testUpdatedNotificationAction := EmailNotification{
		ConnectionType: "Email",
		Recipients:     testUpdatedRecipients,
		Subject:        "test tf monitor",
		TimeZone:       "PST",
		MessageBody:    "test",
	}
	testUpdatedNotifications := []MonitorNotification{
		{
			Notification:       testUpdatedNotificationAction,
			RunForTriggerTypes: testUpdatedTriggerTypes,
		},
	}

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
					resource.TestCheckResourceAttr("sumologic_monitor.test", "monitor_type", testMonitorType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_disabled", strconv.FormatBool(testIsDisabled)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "type", testType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "description", testDescription),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "content_type", testContentType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0.row_id", testQueries[0].RowID),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.trigger_type", testTriggers[0].TriggerType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notifications.0.notification.0.connection_type", testNotifications[0].Notification.(EmailNotification).ConnectionType),
				),
			},
			{
				Config: testAccSumologicMonitorsLibraryMonitorUpdate(testNameSuffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_monitor.test", "monitor_type", testUpdatedMonitorType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_disabled", strconv.FormatBool(testUpdatedIsDisabled)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "name", testUpdatedName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "type", testUpdatedType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "description", testUpdatedDescription),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "content_type", testUpdatedContentType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0.row_id", testUpdatedQueries[0].RowID),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.trigger_type", testUpdatedTriggers[0].TriggerType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notifications.0.notification.0.connection_type", testUpdatedNotifications[0].Notification.(EmailNotification).ConnectionType),
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
			resource.TestCheckResourceAttrSet(name, "is_disabled"),
			resource.TestCheckResourceAttrSet(name, "name"),
			resource.TestCheckResourceAttrSet(name, "parent_id"),
			resource.TestCheckResourceAttrSet(name, "is_mutable"),
			resource.TestCheckResourceAttrSet(name, "type"),
			resource.TestCheckResourceAttrSet(name, "version"),
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
	triggers  {
		threshold_type = "LessThanOrEqual"
		threshold = 40.0
		time_range = "15m"
		occurrence_type = "ResultCount"
		trigger_source = "AllResults"
		trigger_type = "ResolvedCritical"
		detection_method = "StaticCondition"
	  }
	notifications {
		notification {
			connection_type = "Email"
			recipients = ["abc@example.com"]
			subject = "test tf monitor"
			time_zone = "PST"
			message_body = "test"
		  }
		run_for_trigger_types = ["Critical", "ResolvedCritical"]
	  }
}`, testName)
}

func testAccSumologicMonitorsLibraryMonitorUpdate(testName string) string {
	return fmt.Sprintf(`
resource "sumologic_monitor" "test" {
	name = "terraform_test_monitor_%s"
	description = "terraform_test_monitor_description"
	type = "MonitorsLibraryMonitor"
	is_disabled = true
	content_type = "Monitor"
	monitor_type = "Logs"
	queries {
		row_id = "A"
		query = "_sourceCategory=monitor-manager info"
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
	  triggers  {
		threshold_type = "LessThanOrEqual"
		threshold = 40.0
		time_range = "15m"
		occurrence_type = "ResultCount"
		trigger_source = "AllResults"
		trigger_type = "ResolvedCritical"
		detection_method = "StaticCondition"
	  }
	notifications {
		notification {
			connection_type = "Email"
			recipients = ["abc@example.com"]
			subject = "test tf monitor"
			time_zone = "PST"
			message_body = "test"
		}
		run_for_trigger_types = ["Critical", "ResolvedCritical"]
	  }
}`, testName)
}
