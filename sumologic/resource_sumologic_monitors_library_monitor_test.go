package sumologic

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestSumologicMonitorsLibraryMonitor_conversionsToFromTriggerConditionsShouldBeInverses(t *testing.T) {
	sortTriggerConditions := func(slice []TriggerCondition) {
		sort.SliceStable(slice, func(i, j int) bool {
			return slice[i].DetectionMethod < slice[j].DetectionMethod
		})
	}
	testTriggerConditions := [][]TriggerCondition{
		{
			exampleLogsStaticTriggerCondition("Critical", 100, "GreaterThan"),
			exampleLogsStaticTriggerCondition("ResolvedCritical", 90, "LessThanOrEqual"),
		},
		{
			exampleLogsStaticTriggerCondition("Warning", 90, "GreaterThan"),
			exampleLogsStaticTriggerCondition("ResolvedWarning", 80, "LessThanOrEqual"),
		},
		{
			exampleMetricsStaticTriggerCondition("Critical", 100, "GreaterThan"),
			exampleMetricsStaticTriggerCondition("ResolvedCritical", 90, "LessThanOrEqual"),
		},
		{
			exampleLogsOutlierTriggerCondition("Critical", 3),
			exampleLogsOutlierTriggerCondition("ResolvedCritical", 3),
		},
		{
			exampleMetricsOutlierTriggerCondition("Critical", 3),
			exampleMetricsOutlierTriggerCondition("ResolvedCritical", 3),
		},
		{
			exampleLogsMissingDataTriggerCondition("MissingData"),
			exampleLogsMissingDataTriggerCondition("ResolvedMissingData"),
		},
		{
			exampleMetricsMissingDataTriggerCondition("MissingData"),
			exampleMetricsMissingDataTriggerCondition("ResolvedMissingData"),
		},
		{
			exampleLogsStaticTriggerCondition("Critical", 100, "GreaterThan"),
			exampleLogsStaticTriggerCondition("ResolvedCritical", 90, "LessThanOrEqual"),
			exampleLogsStaticTriggerCondition("Warning", 90, "GreaterThan"),
			exampleLogsStaticTriggerCondition("ResolvedWarning", 80, "LessThanOrEqual"),
			exampleLogsMissingDataTriggerCondition("MissingData"),
			exampleLogsMissingDataTriggerCondition("ResolvedMissingData"),
		},
		{
			exampleMetricsOutlierTriggerCondition("Critical", 3),
			exampleMetricsOutlierTriggerCondition("ResolvedCritical", 3),
			exampleMetricsOutlierTriggerCondition("Warning", 2),
			exampleMetricsOutlierTriggerCondition("ResolvedWarning", 2),
			exampleMetricsMissingDataTriggerCondition("MissingData"),
			exampleMetricsMissingDataTriggerCondition("ResolvedMissingData"),
		},
	}
	for _, triggerConditions := range testTriggerConditions {
		triggerConditionsAfterRoundTrip := triggerConditionsBlockToJson(jsonToTriggerConditionsBlock(triggerConditions))
		sortTriggerConditions(triggerConditionsAfterRoundTrip)
		sortTriggerConditions(triggerConditions)
		if len(triggerConditionsAfterRoundTrip) != len(triggerConditions) {
			log.Fatalln("Test case:", triggerConditions, "Lengths differ: Expected", len(triggerConditions), "got", len(triggerConditionsAfterRoundTrip))
		}
		for i := range triggerConditions {
			if triggerConditionsAfterRoundTrip[i] != triggerConditions[i] {
				log.Fatalln("Test case:", triggerConditions, "Expected", triggerConditions[i], "got", triggerConditionsAfterRoundTrip[i])
			}
		}
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
	testEvaluationDelay := "5m"
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
	testAlertName := "Alert from {{Name}}"
	testGroupFields := [2]string{"groupingField1", "groupingField2"}

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
					resource.TestCheckResourceAttr("sumologic_monitor.test", "evaluation_delay", testEvaluationDelay),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "content_type", testContentType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0.row_id", testQueries[0].RowID),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.trigger_type", testTriggers[0].TriggerType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.time_range", testTriggers[0].TimeRange),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notifications.0.notification.0.connection_type", testNotifications[0].Notification.(EmailNotification).ConnectionType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "alert_name", testAlertName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notification_group_fields.0", testGroupFields[0]),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notification_group_fields.1", testGroupFields[1]),
				),
			},
		},
	})
}

func TestAccSumologicMonitorsLibraryMonitor_create_all_monitor_types(t *testing.T) {
	var monitorsLibraryMonitor MonitorsLibraryMonitor
	for _, monitorConfig := range allExampleMonitors {
		testNameSuffix := acctest.RandString(16)

		testName := "terraform_test_monitor_" + testNameSuffix

		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckMonitorsLibraryMonitorDestroy(monitorsLibraryMonitor),
			Steps: []resource.TestStep{
				{
					Config: monitorConfig(testName),
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
	testPlaybook := "This is a test playbook"
	testIsDisabled := false
	testEvaluationDelay := "5m"
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
	testAlertName := "Alert from {{Name}}"
	testGroupFields := [2]string{"groupingField1", "groupingField2"}

	// updated fields
	testUpdatedName := "terraform_test_monitor_" + testNameSuffix
	testUpdatedDescription := "terraform_test_monitor_description"
	testUpdatedType := "MonitorsLibraryMonitor"
	testUpdatedContentType := "Monitor"
	testUpdatedMonitorType := "Logs"
	testUpdatedPlaybook := "This is an updated test playbook"
	testUpdatedIsDisabled := true
	testUpdatedEvaluationDelay := "8m"
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
	testUpdatedAlertName := "Updated Alert from {{Name}}"
	testUpdatedGroupFields := [2]string{"groupingField3", "groupingField4"}

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
					resource.TestCheckResourceAttr("sumologic_monitor.test", "evaluation_delay", testEvaluationDelay),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "content_type", testContentType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0.row_id", testQueries[0].RowID),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.trigger_type", testTriggers[0].TriggerType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.time_range", testTriggers[0].TimeRange),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notifications.0.notification.0.connection_type", testNotifications[0].Notification.(EmailNotification).ConnectionType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "playbook", testPlaybook),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "alert_name", testAlertName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notification_group_fields.0", testGroupFields[0]),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notification_group_fields.1", testGroupFields[1]),
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
					resource.TestCheckResourceAttr("sumologic_monitor.test", "evaluation_delay", testUpdatedEvaluationDelay),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "content_type", testUpdatedContentType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0.row_id", testUpdatedQueries[0].RowID),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.trigger_type", testUpdatedTriggers[0].TriggerType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.time_range", testUpdatedTriggers[0].TimeRange),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notifications.0.notification.0.connection_type", testUpdatedNotifications[0].Notification.(EmailNotification).ConnectionType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "playbook", testUpdatedPlaybook),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "alert_name", testUpdatedAlertName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notification_group_fields.0", testUpdatedGroupFields[0]),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notification_group_fields.1", testUpdatedGroupFields[1]),
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
			resource.TestCheckResourceAttrSet(name, "evaluation_delay"),
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
	evaluation_delay = "5m"
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
	playbook = "This is a test playbook"  
	alert_name =  "Alert from {{Name}}"
	notification_group_fields = ["groupingField1", "groupingField2"]
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
	evaluation_delay = "8m"
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
	playbook = "This is an updated test playbook"
	alert_name = "Updated Alert from {{Name}}"
	notification_group_fields = ["groupingField3", "groupingField4"]
}`, testName)
}

func exampleMonitorWithTriggerCondition(
	testName string,
	monitorType string,
	query string,
	trigger string,
	triggerTys []string) string {
	triggerTysStr := `"` + strings.Join(triggerTys, `","`) + `"`
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
    trigger_conditions {
      %s
    }
	notifications {
		notification {
			connection_type = "Email"
			recipients = ["abc@example.com"]
			subject = "test tf monitor"
			time_zone = "PST"
			message_body = "test"
		  }
		run_for_trigger_types = [%s]
	  }
	playbook = "This is a test playbook"
}`, testName, monitorType, query, trigger, triggerTysStr)
}

var exampleLogsStaticTriggerConditionBlock = `
   logs_static_condition {
     critical {
       time_range = "30m"
       alert {
         threshold = 100.0
         threshold_type = "GreaterThan"
       }
       resolution {
         threshold = 90
         threshold_type = "LessThanOrEqual"
       }
     }
     field = "field"
   }`

var exampleMetricsStaticTriggerConditionBlock1 = `
   metrics_static_condition {
     critical {
       time_range = "30m"
       occurrence_type = "AtLeastOnce"
       alert {
         threshold = 100.0
         threshold_type = "GreaterThan"
       }
       resolution {
         threshold = 90
         threshold_type = "LessThanOrEqual"
       }
     }
   }`

var exampleMetricsStaticTriggerConditionBlock2 = `
   metrics_static_condition {
     critical {
       time_range = "30m"
       occurrence_type = "Always"
       alert {
         threshold = 100.0
         threshold_type = "GreaterThan"
       }
       resolution {
         threshold = 90
         threshold_type = "LessThanOrEqual"
       }
     }
   }`

var exampleLogsOutlierTriggerConditionBlock = `
   logs_outlier_condition {
     critical {
       window = 5
       consecutive = 1
       threshold = 3.0
     }
     field = "field"
     direction = "Both"
   }`

var exampleMetricsOutlierTriggerConditionBlock = `
   metrics_outlier_condition {
     critical {
       baseline_window = "15m"
       threshold = 3.0
     }
     direction = "Both"
   }`

var exampleLogsMissingDataTriggerConditionBlock = `
   logs_missing_data_condition {
     time_range = "30m"
   }`

var exampleMetricsMissingDataTriggerConditionBlock = `
   metrics_missing_data_condition {
     time_range = "30m"
     trigger_source = "AnyTimeSeries"
   }`

func exampleLogsStaticMonitor(testName string) string {
	query := "error | timeslice 1m | count as field by _timeslice"
	return exampleMonitorWithTriggerCondition(testName, "Logs", query,
		exampleLogsStaticTriggerConditionBlock, []string{"Critical", "ResolvedCritical"})
}

func exampleMetricsStaticMonitor1(testName string) string {
	query := "error _sourceCategory=category"
	return exampleMonitorWithTriggerCondition(testName, "Metrics", query,
		exampleMetricsStaticTriggerConditionBlock1, []string{"Critical", "ResolvedCritical"})
}

func exampleMetricsStaticMonitor2(testName string) string {
	query := "error _sourceCategory=category"
	return exampleMonitorWithTriggerCondition(testName, "Metrics", query,
		exampleMetricsStaticTriggerConditionBlock2, []string{"Critical", "ResolvedCritical"})
}

func exampleLogsOutlierMonitor(testName string) string {
	query := "error | timeslice 1m | count as field by _timeslice"
	return exampleMonitorWithTriggerCondition(testName, "Logs", query,
		exampleLogsOutlierTriggerConditionBlock, []string{"Critical", "ResolvedCritical"})
}

func exampleMetricsOutlierMonitor(testName string) string {
	query := "error _sourceCategory=category"
	return exampleMonitorWithTriggerCondition(testName, "Metrics", query,
		exampleMetricsOutlierTriggerConditionBlock, []string{"Critical", "ResolvedCritical"})
}

func exampleLogsMissingDataMonitor(testName string) string {
	query := "error | timeslice 1m | count as field by _timeslice"
	return exampleMonitorWithTriggerCondition(testName, "Logs", query,
		exampleLogsMissingDataTriggerConditionBlock, []string{"MissingData", "ResolvedMissingData"})
}

func exampleMetricsMissingDataMonitor(testName string) string {
	query := "error _sourceCategory=category"
	return exampleMonitorWithTriggerCondition(testName, "Metrics", query,
		exampleMetricsMissingDataTriggerConditionBlock, []string{"MissingData", "ResolvedMissingData"})
}

var allExampleMonitors = []func(testName string) string{
	exampleLogsStaticMonitor,
	exampleMetricsStaticMonitor1,
	exampleMetricsStaticMonitor2,
	exampleLogsOutlierMonitor,
	exampleMetricsOutlierMonitor,
	exampleLogsMissingDataMonitor,
	exampleMetricsMissingDataMonitor,
}

func exampleLogsStaticTriggerCondition(triggerType string, threshold float64, thresholdType string) TriggerCondition {
	return TriggerCondition{
		TimeRange:       "30m",
		TriggerType:     triggerType,
		Threshold:       threshold,
		ThresholdType:   thresholdType,
		Field:           "field",
		DetectionMethod: "LogsStaticCondition",
	}
}

func exampleMetricsStaticTriggerCondition(triggerType string, threshold float64, thresholdType string) TriggerCondition {
	return TriggerCondition{
		TimeRange:       "30m",
		TriggerType:     triggerType,
		Threshold:       threshold,
		ThresholdType:   thresholdType,
		OccurrenceType:  "Always",
		DetectionMethod: "MetricsStaticCondition",
	}
}

func exampleLogsOutlierTriggerCondition(triggerType string, threshold float64) TriggerCondition {
	return TriggerCondition{
		TriggerType:     triggerType,
		Window:          5,
		Consecutive:     1,
		Direction:       "Both",
		Threshold:       threshold,
		Field:           "field",
		DetectionMethod: "LogsOutlierCondition",
	}
}

func exampleMetricsOutlierTriggerCondition(triggerType string, threshold float64) TriggerCondition {
	return TriggerCondition{
		TriggerType:     triggerType,
		Threshold:       threshold,
		BaselineWindow:  "30m",
		Direction:       "Both",
		DetectionMethod: "MetricsOutlierCondition",
	}
}

func exampleLogsMissingDataTriggerCondition(triggerType string) TriggerCondition {
	return TriggerCondition{
		TimeRange:       "30m",
		TriggerType:     triggerType,
		DetectionMethod: "LogsMissingDataCondition",
	}
}

func exampleMetricsMissingDataTriggerCondition(triggerType string) TriggerCondition {
	return TriggerCondition{
		TimeRange:       "30m",
		TriggerType:     triggerType,
		TriggerSource:   "AllTimeSeries",
		DetectionMethod: "MetricsMissingDataCondition",
	}
}
