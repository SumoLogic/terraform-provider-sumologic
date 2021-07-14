package sumologic

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestSumologicMonitorsLibraryMonitor_conversionsToFromTriggerConditionShouldBeInverses(t *testing.T) {
	for _, trigger := range allExampleTriggerConditions {
		triggerAfterRoundTrip := triggersBlockToTriggerCondition(trigger.toTriggersBlock())

		if triggerAfterRoundTrip != trigger {
			log.Fatalln("Expected", trigger, "got", triggerAfterRoundTrip)
		}
	}
}

func TestSumologicMonitorsLibraryMonitor_triggerConditionConvertersShouldWorkWithLegacyVersion(t *testing.T) {
	legacyTriggerCondition := exampleStaticTriggerCondition
	legacyTriggerCondition.Field = "" // Important, since legacy trigger conditions do not support 'field' argument.
	triggerAfterRoundTrip := triggersBlockToTriggerCondition(legacyTriggerCondition.toLegacyTriggersBlock())

	if triggerAfterRoundTrip != legacyTriggerCondition {
		log.Fatalln("Expected", legacyTriggerCondition, "got", triggerAfterRoundTrip)
	}
}

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

func TestAccSumologicMonitorsLibraryMonitor_create_all_monitor_types(t *testing.T) {
	var monitorsLibraryMonitor MonitorsLibraryMonitor
	testNameSuffix := acctest.RandString(16)

	testName := "terraform_test_monitor_" + testNameSuffix
	for _, monitorConfig := range allExampleMonitors(testName) {
		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckMonitorsLibraryMonitorDestroy(monitorsLibraryMonitor),
			Steps: []resource.TestStep{
				{
					Config: monitorConfig,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckMonitorsLibraryMonitorExists("sumologic_monitor.test", &monitorsLibraryMonitor, t),
						testAccCheckMonitorsLibraryMonitorAttributes("sumologic_monitor.test"),
						resource.TestCheckResourceAttr("sumologic_monitor.test", "is_disabled", strconv.FormatBool(false)),
						resource.TestCheckResourceAttr("sumologic_monitor.test", "name", testName),
					),
				},
			},
		})
	}
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

func exampleMonitorWithTriggers(testName string,
	monitorType string,
	query string,
	trigger1 string, triggerTy1 string,
	trigger2 string, triggerTy2 string) string {
	return fmt.Sprintf(`
resource "sumologic_monitor" "test" {
	name = "%s"
	description = "terraform_test_monitor_description"
	type = "MonitorsLibraryMonitor"
	is_disabled = false
	content_type = "Monitor"
	monitor_type = "%s"
	queries {
		row_id = "A"
		query = "%s"
	  }
   %s
   %s
	notifications {
		notification {
			connection_type = "Email"
			recipients = ["abc@example.com"]
			subject = "test tf monitor"
			time_zone = "PST"
			message_body = "test"
		  }
		run_for_trigger_types = ["%s", "%s"]
	  }
}`, testName, monitorType, query, trigger1, trigger2, triggerTy1, triggerTy2)
}

var exampleLegacyAlertTrigger = `
triggers {
		threshold_type = "GreaterThan"
		threshold = 100.0
		time_range = "30m"
		occurrence_type = "ResultCount"
		trigger_source = "AllResults"
		trigger_type = "Critical"
		detection_method = "StaticCondition"
}`

var exampleLegacyResolutionTrigger = `
triggers {
		threshold_type = "LessThanOrEqual"
		threshold = 100.0
		time_range = "30m"
		occurrence_type = "ResultCount"
		trigger_source = "AllResults"
		trigger_type = "ResolvedCritical"
		detection_method = "StaticCondition"
}`

var exampleStaticConditionAlertTrigger = `
triggers {
   static_condition {
		trigger_type = "Critical"
        threshold = 100.0
        threshold_type = "GreaterThan"
        field = "some field"
		time_range = "30m"
		trigger_source = "AllResults"
		occurrence_type = "ResultCount"
   }
}`

var exampleStaticConditionResolutionTrigger = `
triggers {
   static_condition {
		trigger_type = "ResolvedCritical"
        threshold = 100.0
        threshold_type = "LessThanOrEqual"
        field = "field"
		time_range = "30m"
		trigger_source = "AllResults"
		occurrence_type = "ResultCount"
   }
}`

var exampleLogsStaticConditionAlertTrigger = `
triggers {
   logs_static_condition {
		trigger_type = "Critical"
        threshold = 100.0
        threshold_type = "GreaterThan"
        field = "field"
		time_range = "30m"
   }
}`

var exampleLogsStaticConditionResolutionTrigger = `
triggers {
   logs_static_condition {
		trigger_type = "ResolvedCritical"
        threshold = 100.0
        threshold_type = "LessThanOrEqual"
        field = "field"
		time_range = "30m"
   }
}`

var exampleMetricsStaticConditionAlertTrigger = `
triggers {
   metrics_static_condition {
		trigger_type = "Critical"
        threshold = 100.0
        threshold_type = "GreaterThan"
		time_range = "30m"
        occurrence_type = "Always"
   }
}`

var exampleMetricsStaticConditionResolutionTrigger = `
triggers {
   metrics_static_condition {
		trigger_type = "ResolvedCritical"
        threshold = 100.0
        threshold_type = "LessThanOrEqual"
		time_range = "30m"
        occurrence_type = "Always"
   }
}`

var exampleLogsOutlierConditionAlertTrigger = `
triggers {
   logs_outlier_condition {
		trigger_type = "Critical"
        threshold = 3.0
        field = "field"
        window = 5
        consecutive = 1
        direction = "Both"
   }
}`

var exampleLogsOutlierConditionResolutionTrigger = `
triggers {
   logs_outlier_condition {
		trigger_type = "ResolvedCritical"
        threshold = 3.0
        field = "field"
        window = 5
        consecutive = 1
        direction = "Both"
   }
}`

var exampleMetricsOutlierConditionAlertTrigger = `
triggers {
   metrics_outlier_condition {
		trigger_type = "Critical"
        threshold = 3.0
        baseline_window = "15m"
        direction = "Both"
   }
}`

var exampleMetricsOutlierConditionResolutionTrigger = `
triggers {
   metrics_outlier_condition {
		trigger_type = "ResolvedCritical"
        threshold = 3.0
        baseline_window = "15m"
        direction = "Both"
   }
}`

var exampleLogsMissingDataConditionAlertTrigger = `
triggers {
   logs_missing_data_condition {
		trigger_type = "MissingData"
        time_range = "30m"
   }
}`

var exampleLogsMissingDataConditionResolutionTrigger = `
triggers {
   logs_missing_data_condition {
		trigger_type = "ResolvedMissingData"
        time_range = "30m"
   }
}`

var exampleMetricsMissingDataConditionAlertTrigger = `
triggers {
   metrics_missing_data_condition {
		trigger_type = "MissingData"
        time_range = "30m"
        trigger_source = "AllTimeSeries"
   }
}`

var exampleMetricsMissingDataConditionResolutionTrigger = `
triggers {
   metrics_missing_data_condition {
		trigger_type = "ResolvedMissingData"
        time_range = "30m"
        trigger_source = "AllTimeSeries"
   }
}`

func exampleLegacyTriggerMonitor(testName string) string {
	query := "error | timeslice 1m | count as field by _timeslice"
	return exampleMonitorWithTriggers(testName, "Logs", query,
		exampleLegacyAlertTrigger, "Critical",
		exampleLegacyResolutionTrigger, "ResolvedCritical")
}

func exampleStaticMonitor(testName string) string {
	query := "error | timeslice 1m | count as field by _timeslice"
	return exampleMonitorWithTriggers(testName, "Logs", query,
		exampleStaticConditionAlertTrigger, "Critical",
		exampleStaticConditionResolutionTrigger, "ResolvedCritical")
}

func exampleLogsStaticMonitor(testName string) string {
	query := "error | timeslice 1m | count as field by _timeslice"
	return exampleMonitorWithTriggers(testName, "Logs", query,
		exampleLogsStaticConditionAlertTrigger, "Critical",
		exampleLogsStaticConditionResolutionTrigger, "ResolvedCritical")
}

func exampleMetricsStaticMonitor(testName string) string {
	query := "error _sourceCategory=category"
	return exampleMonitorWithTriggers(testName, "Metrics", query,
		exampleMetricsStaticConditionAlertTrigger, "Critical",
		exampleMetricsStaticConditionResolutionTrigger, "ResolvedCritical")
}

func exampleLogsOutlierMonitor(testName string) string {
	query := "error | timeslice 1m | count as field by _timeslice"
	return exampleMonitorWithTriggers(testName, "Logs", query,
		exampleLogsOutlierConditionAlertTrigger, "Critical",
		exampleLogsOutlierConditionResolutionTrigger, "ResolvedCritical")
}

func exampleMetricsOutlierMonitor(testName string) string {
	query := "error _sourceCategory=category"
	return exampleMonitorWithTriggers(testName, "Metrics", query,
		exampleMetricsOutlierConditionAlertTrigger, "Critical",
		exampleMetricsOutlierConditionResolutionTrigger, "ResolvedCritical")
}

func exampleLogsMissingDataMonitor(testName string) string {
	query := "error | timeslice 1m | count as field by _timeslice"
	return exampleMonitorWithTriggers(testName, "Logs", query,
		exampleLogsMissingDataConditionAlertTrigger, "MissingData",
		exampleLogsMissingDataConditionResolutionTrigger, "ResolvedMissingData")
}

func exampleMetricsMissingDataMonitor(testName string) string {
	query := "error _sourceCategory=category"
	return exampleMonitorWithTriggers(testName, "Metrics", query,
		exampleMetricsMissingDataConditionAlertTrigger, "MissingData",
		exampleMetricsMissingDataConditionResolutionTrigger, "ResolvedMissingData")
}

func allExampleMonitors(testName string) []string {
	return []string{
		exampleLegacyTriggerMonitor(testName),
		exampleStaticMonitor(testName),
		exampleLogsStaticMonitor(testName),
		exampleMetricsStaticMonitor(testName),
		exampleLogsOutlierMonitor(testName),
		exampleMetricsOutlierMonitor(testName),
		exampleLogsMissingDataMonitor(testName),
		exampleMetricsMissingDataMonitor(testName),
	}
}

var exampleStaticTriggerCondition = TriggerCondition{
	TimeRange:       "30m",
	TriggerType:     "Critical",
	Threshold:       100,
	ThresholdType:   "LessThan",
	Field:           "field",
	DetectionMethod: "StaticCondition",
}

var exampleLogsStaticTriggerCondition = TriggerCondition{
	TimeRange:       "30m",
	TriggerType:     "Critical",
	Threshold:       100,
	ThresholdType:   "LessThan",
	Field:           "field",
	DetectionMethod: "LogsStaticCondition",
}

var exampleMetricsStaticTriggerCondition = TriggerCondition{
	TimeRange:       "30m",
	TriggerType:     "Critical",
	Threshold:       100,
	ThresholdType:   "LessThan",
	OccurrenceType:  "Always",
	DetectionMethod: "MetricsStaticCondition",
}

var exampleLogsOutlierTriggerCondition = TriggerCondition{
	TriggerType:     "Critical",
	Window:          5,
	Consecutive:     1,
	Direction:       "Both",
	Threshold:       3,
	Field:           "field",
	DetectionMethod: "LogsOutlierCondition",
}

var exampleMetricsOutlierTriggerCondition = TriggerCondition{
	TriggerType:     "Critical",
	Threshold:       3,
	BaselineWindow:  "30m",
	Direction:       "Both",
	DetectionMethod: "MetricsOutlierCondition",
}

var exampleLogsMissingDataTriggerCondition = TriggerCondition{
	TimeRange:       "30m",
	TriggerType:     "Critical",
	DetectionMethod: "LogsMissingDataCondition",
}

var exampleMetricsMissingDataTriggerCondition = TriggerCondition{
	TimeRange:       "30m",
	TriggerType:     "Critical",
	TriggerSource:   "AllTimeSeries",
	DetectionMethod: "MetricsMissingDataCondition",
}

var allExampleTriggerConditions = []TriggerCondition{
	exampleStaticTriggerCondition,
	exampleLogsStaticTriggerCondition,
	exampleMetricsStaticTriggerCondition,
	exampleLogsOutlierTriggerCondition,
	exampleMetricsOutlierTriggerCondition,
	exampleLogsMissingDataTriggerCondition,
	exampleMetricsMissingDataTriggerCondition,
}
