package sumologic

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicMonitorsLibraryMonitor_schemaValidations(t *testing.T) {
	var monitorsLibraryMonitor MonitorsLibraryMonitor
	config := `
       resource "sumologic_monitor" "test" { 
         name = "test"
         type = "MonitorsLibraryMonitor"
         monitor_type = "Logs"
         triggers {
           threshold_type = "foo"
         }
       }`
	expectedError := regexp.MustCompile(".*expected triggers.0.threshold_type to be one of \\[LessThan LessThanOrEqual GreaterThan GreaterThanOrEqual\\], got foo.*")
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorsLibraryMonitorDestroy(monitorsLibraryMonitor),
		Steps: []resource.TestStep{
			{
				Config:      config,
				PlanOnly:    true,
				ExpectError: expectedError,
			},
		},
	})
}

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
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.time_range", testTriggers[0].TimeRange),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notifications.0.notification.0.connection_type", testNotifications[0].Notification.(EmailNotification).ConnectionType),
				),
			},
		},
	})
}

func TestAccSumologicMonitorsLibraryMonitor_createLogsStaticMonitors(t *testing.T) {
	var monitorsLibraryMonitor MonitorsLibraryMonitor
	testNameSuffix := acctest.RandString(16)
	testName := "terraform_test_monitor_" + testNameSuffix
	testType := "MonitorsLibraryMonitor"
	testField := "time_taken"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorsLibraryMonitorDestroy(monitorsLibraryMonitor),
		Steps: []resource.TestStep{
			{
				Config: exampleLogsStaticMonitor(testNameSuffix, testField),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorsLibraryMonitorExists("sumologic_monitor.test", &monitorsLibraryMonitor, t),
					testAccCheckMonitorsLibraryMonitorAttributes("sumologic_monitor.test"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "monitor_type", "Logs"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_disabled", strconv.FormatBool(false)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "type", testType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "content_type", "Monitor"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0.row_id", "A"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0.query", fmt.Sprintf(`_sourceCategory=monitor-manager error | parse "field=*," as %s`, testField)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.trigger_type", "Critical"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.time_range", "15m"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.threshold", "40"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.threshold_type", "GreaterThan"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.field", "time_taken"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.detection_method", "LogsStaticCondition"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notifications.0.notification.0.connection_type", "Email"),
				),
			},
		},
	})
}

func TestAccSumologicMonitorsLibraryMonitor_createMetricsStaticMonitors(t *testing.T) {
	var monitorsLibraryMonitor MonitorsLibraryMonitor
	testNameSuffix := acctest.RandString(16)
	testName := "terraform_test_monitor_" + testNameSuffix
	testType := "MonitorsLibraryMonitor"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorsLibraryMonitorDestroy(monitorsLibraryMonitor),
		Steps: []resource.TestStep{
			{
				Config: exampleMetricsStaticMonitor(testNameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorsLibraryMonitorExists("sumologic_monitor.test", &monitorsLibraryMonitor, t),
					testAccCheckMonitorsLibraryMonitorAttributes("sumologic_monitor.test"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "monitor_type", "Metrics"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_disabled", strconv.FormatBool(false)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "type", testType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "content_type", "Monitor"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0.row_id", "A"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0.query", "_sourceCategory=monitor-manager error"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.trigger_type", "Critical"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.time_range", "15m"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.threshold", "40"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.threshold_type", "GreaterThan"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.occurrence_type", "Always"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.detection_method", "MetricsStaticCondition"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notifications.0.notification.0.connection_type", "Email"),
				),
			},
		},
	})
}

func TestAccSumologicMonitorsLibraryMonitor_createLogsOutlierMonitors(t *testing.T) {
	var monitorsLibraryMonitor MonitorsLibraryMonitor
	testNameSuffix := acctest.RandString(16)
	testName := "terraform_test_monitor_" + testNameSuffix
	testType := "MonitorsLibraryMonitor"
	testField := "time_taken"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorsLibraryMonitorDestroy(monitorsLibraryMonitor),
		Steps: []resource.TestStep{
			{
				Config: exampleLogsOutlierMonitor(testNameSuffix, testField),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorsLibraryMonitorExists("sumologic_monitor.test", &monitorsLibraryMonitor, t),
					testAccCheckMonitorsLibraryMonitorAttributes("sumologic_monitor.test"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "monitor_type", "Logs"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_disabled", strconv.FormatBool(false)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "type", testType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "content_type", "Monitor"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0.row_id", "A"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0.query", fmt.Sprintf(`_sourceCategory=monitor-manager error | parse "field=*," as %s | timeslice 1m | avg(%s) as %s by _timeslice`, testField, testField, testField)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.trigger_type", "Critical"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.field", "time_taken"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.window", strconv.Itoa(5)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.consecutive", strconv.Itoa(1)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.direction", "Both"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.detection_method", "LogsOutlierCondition"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notifications.0.notification.0.connection_type", "Email"),
				),
			},
		},
	})
}

func TestAccSumologicMonitorsLibraryMonitor_createMetricsOutlierMonitors(t *testing.T) {
	var monitorsLibraryMonitor MonitorsLibraryMonitor
	testNameSuffix := acctest.RandString(16)
	testName := "terraform_test_monitor_" + testNameSuffix
	testType := "MonitorsLibraryMonitor"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorsLibraryMonitorDestroy(monitorsLibraryMonitor),
		Steps: []resource.TestStep{
			{
				Config: exampleMetricsOutlierMonitor(testNameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorsLibraryMonitorExists("sumologic_monitor.test", &monitorsLibraryMonitor, t),
					testAccCheckMonitorsLibraryMonitorAttributes("sumologic_monitor.test"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "monitor_type", "Metrics"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_disabled", strconv.FormatBool(false)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "type", testType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "content_type", "Monitor"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0.row_id", "A"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0.query", "_sourceCategory=monitor-manager error"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.trigger_type", "Critical"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.baseline_window", "15m"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.threshold", "3"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.direction", "Both"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.detection_method", "MetricsOutlierCondition"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notifications.0.notification.0.connection_type", "Email"),
				),
			},
		},
	})
}

func TestAccSumologicMonitorsLibraryMonitor_createLogsMissingDataMonitors(t *testing.T) {
	var monitorsLibraryMonitor MonitorsLibraryMonitor
	testNameSuffix := acctest.RandString(16)
	testName := "terraform_test_monitor_" + testNameSuffix
	testType := "MonitorsLibraryMonitor"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorsLibraryMonitorDestroy(monitorsLibraryMonitor),
		Steps: []resource.TestStep{
			{
				Config: exampleLogsMissingDataMonitor(testNameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorsLibraryMonitorExists("sumologic_monitor.test", &monitorsLibraryMonitor, t),
					testAccCheckMonitorsLibraryMonitorAttributes("sumologic_monitor.test"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "monitor_type", "Logs"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_disabled", strconv.FormatBool(false)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "type", testType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "content_type", "Monitor"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0.row_id", "A"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0.query", "_sourceCategory=monitor-manager info"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.trigger_type", "MissingData"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.time_range", "15m"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.detection_method", "LogsMissingDataCondition"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notifications.0.notification.0.connection_type", "Email"),
				),
			},
		},
	})
}

func TestAccSumologicMonitorsLibraryMonitor_createMetricsMissingDataMonitors(t *testing.T) {
	var monitorsLibraryMonitor MonitorsLibraryMonitor
	testNameSuffix := acctest.RandString(16)
	testName := "terraform_test_monitor_" + testNameSuffix
	testType := "MonitorsLibraryMonitor"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorsLibraryMonitorDestroy(monitorsLibraryMonitor),
		Steps: []resource.TestStep{
			{
				Config: exampleMetricsMissingDataMonitor(testNameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorsLibraryMonitorExists("sumologic_monitor.test", &monitorsLibraryMonitor, t),
					testAccCheckMonitorsLibraryMonitorAttributes("sumologic_monitor.test"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "monitor_type", "Metrics"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "is_disabled", strconv.FormatBool(false)),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "type", testType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "content_type", "Monitor"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0.row_id", "A"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0.query", "_sourceCategory=monitor-manager"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.trigger_type", "MissingData"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.time_range", "15m"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.trigger_source", "AllTimeSeries"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.detection_method", "MetricsMissingDataCondition"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notifications.0.notification.0.connection_type", "Email"),
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
			TimeRange:       "30m",
			OccurrenceType:  "ResultCount",
			TriggerSource:   "AllResults",
			TriggerType:     "Critical",
			DetectionMethod: "StaticCondition",
		},
		{
			ThresholdType:   "LessThanOrEqual",
			Threshold:       40.0,
			TimeRange:       "30m",
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
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.time_range", testTriggers[0].TimeRange),
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
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.time_range", testUpdatedTriggers[0].TimeRange),
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

func exampleLogsStaticMonitor(testName string, fieldName string) string {
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
		query = "_sourceCategory=monitor-manager error | parse \"field=*,\" as %s"
	  }
	triggers  {
		threshold_type = "GreaterThan"
		threshold = 40.0
        field = "%s"
		time_range = "15m"
		trigger_type = "Critical"
		detection_method = "LogsStaticCondition"
	  }
	triggers  {
		threshold_type = "LessThanOrEqual"
		threshold = 40.0
        field = "%s"
		time_range = "15m"
		trigger_type = "ResolvedCritical"
		detection_method = "LogsStaticCondition"
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
}`, testName, fieldName, fieldName, fieldName)
}

func exampleMetricsStaticMonitor(testName string) string {
	return fmt.Sprintf(`
resource "sumologic_monitor" "test" {
	name = "terraform_test_monitor_%s"
	description = "terraform_test_monitor_description"
	type = "MonitorsLibraryMonitor"
	is_disabled = false
	content_type = "Monitor"
	monitor_type = "Metrics"
	queries {
		row_id = "A"
		query = "_sourceCategory=monitor-manager error"
	  }
	triggers  {
		threshold_type = "GreaterThan"
		threshold = 40.0
		time_range = "15m"
		trigger_type = "Critical"
        occurrence_type = "Always"
		detection_method = "MetricsStaticCondition"
	  }
	triggers  {
		threshold_type = "LessThanOrEqual"
		threshold = 40.0
		time_range = "15m"
		trigger_type = "ResolvedCritical"
        occurrence_type = "Always"
		detection_method = "MetricsStaticCondition"
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

func exampleLogsOutlierMonitor(testName string, fieldName string) string {
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
		query = "_sourceCategory=monitor-manager error | parse \"field=*,\" as %s | timeslice 1m | avg(%s) as %s by _timeslice"
	  }
	triggers  {
		threshold = 3.0
        field = "%s"
        window = 5
        consecutive = 1
        direction = "Both"
		trigger_type = "Critical"
		detection_method = "LogsOutlierCondition"
	  }
	triggers  {
		threshold = 3.0
        field = "%s"
        window = 5
        consecutive = 1
        direction = "Both"
		trigger_type = "ResolvedCritical"
		detection_method = "LogsOutlierCondition"
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
}`, testName, fieldName, fieldName, fieldName, fieldName, fieldName)
}

func exampleMetricsOutlierMonitor(testName string) string {
	return fmt.Sprintf(`
resource "sumologic_monitor" "test" {
	name = "terraform_test_monitor_%s"
	description = "terraform_test_monitor_description"
	type = "MonitorsLibraryMonitor"
	is_disabled = false
	content_type = "Monitor"
	monitor_type = "Metrics"
	queries {
		row_id = "A"
		query = "_sourceCategory=monitor-manager error"
	  }
	triggers  {
		threshold = 3.0
        baseline_window = "15m"
        direction = "Both"
		trigger_type = "Critical"
		detection_method = "MetricsOutlierCondition"
	  }
	triggers  {
		threshold = 3.0
        baseline_window = "15m"
        direction = "Both"
		trigger_type = "ResolvedCritical"
		detection_method = "MetricsOutlierCondition"
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

func exampleLogsMissingDataMonitor(testName string) string {
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
		query = "_sourceCategory=monitor-manager info"
	  }
	triggers  {
		time_range = "15m"
		trigger_type = "MissingData"
		detection_method = "LogsMissingDataCondition"
	  }
	triggers  {
		time_range = "15m"
		trigger_type = "ResolvedMissingData"
		detection_method = "LogsMissingDataCondition"
	  }
	notifications {
		notification {
			connection_type = "Email"
			recipients = ["abc@example.com"]
			subject = "test tf monitor"
			time_zone = "PST"
			message_body = "test"
		  }
		run_for_trigger_types = ["MissingData", "ResolvedMissingData"]
	  }
}`, testName)
}

func exampleMetricsMissingDataMonitor(testName string) string {
	return fmt.Sprintf(`
resource "sumologic_monitor" "test" {
	name = "terraform_test_monitor_%s"
	description = "terraform_test_monitor_description"
	type = "MonitorsLibraryMonitor"
	is_disabled = false
	content_type = "Monitor"
	monitor_type = "Metrics"
	queries {
		row_id = "A"
		query = "_sourceCategory=monitor-manager"
	  }
	triggers  {
		time_range = "15m"
		trigger_type = "MissingData"
        trigger_source = "AllTimeSeries"
		detection_method = "MetricsMissingDataCondition"
	  }
	triggers  {
		time_range = "15m"
		trigger_type = "ResolvedMissingData"
        trigger_source = "AllTimeSeries"
		detection_method = "MetricsMissingDataCondition"
	  }
	notifications {
		notification {
			connection_type = "Email"
			recipients = ["abc@example.com"]
			subject = "test tf monitor"
			time_zone = "PST"
			message_body = "test"
		  }
		run_for_trigger_types = ["MissingData", "ResolvedMissingData"]
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
		time_range = "30m"
		occurrence_type = "ResultCount"
		trigger_source = "AllResults"
		trigger_type = "Critical"
		detection_method = "StaticCondition"
	  }
	  triggers  {
		threshold_type = "LessThanOrEqual"
		threshold = 40.0
		time_range = "30m"
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
