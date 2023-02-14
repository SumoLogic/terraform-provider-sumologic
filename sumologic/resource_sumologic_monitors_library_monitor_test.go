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

func TestAccSumologicMonitorsLibraryMonitor_schemaTriggerValidations(t *testing.T) {
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

func TestAccSumologicMonitorsLibraryMonitor_schemaTriggerConditionValidations(t *testing.T) {
	var monitorsLibraryMonitor MonitorsLibraryMonitor
	for _, monitorConfig := range allInvalidTriggerConditionMonitorResources {
		testNameSuffix := acctest.RandString(16)

		testName := "terraform_test_invalid_monitor_" + testNameSuffix

		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckMonitorsLibraryMonitorDestroy(monitorsLibraryMonitor),
			Steps: []resource.TestStep{
				{
					Config:      monitorConfig(testName),
					ExpectError: regexp.MustCompile("config is invalid"),
				},
			},
		})
	}
}

func TestAccSumologicMonitorsLibraryMonitor_triggersTimeRangeDiffSuppression(t *testing.T) {
	var monitorsLibraryMonitor MonitorsLibraryMonitor
	canonicalTimeRange := "1h"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorsLibraryMonitorDestroy(monitorsLibraryMonitor),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMonitorsLibraryMonitor("triggers_negative_expanded_hour"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.time_range", canonicalTimeRange),
				),
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
	canonicalTestEvaluationDelay := "1h"
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
			TimeRange:       "1h",
			OccurrenceType:  "ResultCount",
			TriggerSource:   "AllResults",
			TriggerType:     "Critical",
			DetectionMethod: "StaticCondition",
		},
		{
			ThresholdType:    "LessThanOrEqual",
			Threshold:        40.0,
			TimeRange:        "1h",
			OccurrenceType:   "ResultCount",
			TriggerSource:    "AllResults",
			TriggerType:      "ResolvedCritical",
			DetectionMethod:  "StaticCondition",
			ResolutionWindow: "5m",
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
					resource.TestCheckResourceAttr("sumologic_monitor.test", "evaluation_delay", canonicalTestEvaluationDelay),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "content_type", testContentType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0.row_id", testQueries[0].RowID),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.trigger_type", testTriggers[0].TriggerType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.time_range", testTriggers[0].TimeRange),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.1.resolution_window", testTriggers[1].ResolutionWindow),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notifications.0.notification.0.connection_type", testNotifications[0].Notification.(EmailNotification).ConnectionType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "alert_name", testAlertName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notification_group_fields.0", testGroupFields[0]),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notification_group_fields.1", testGroupFields[1]),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "obj_permission.#", "2"),
					testAccCheckMonitorsLibraryMonitorFGPBackend("sumologic_monitor.test", t, genExpectedPermStmtsMonitor),
				),
			},
		},
	})
}

func TestAccSumologicMonitorsLibraryMonitor_create_with_no_resolution_window(t *testing.T) {
	var monitorsLibraryMonitor MonitorsLibraryMonitor
	testNameSuffix := acctest.RandString(16)

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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorsLibraryMonitorDestroy(monitorsLibraryMonitor),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMonitorsLibraryMonitorWithNoResolutionWindow(testNameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorsLibraryMonitorExists("sumologic_monitor.test", &monitorsLibraryMonitor, t),
					testAccCheckMonitorsLibraryMonitorAttributes("sumologic_monitor.test"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.trigger_type", testTriggers[0].TriggerType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.time_range", testTriggers[0].TimeRange),
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
	canonicalTestEvaluationDelay := "1h"
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
			TimeRange:       "1h",
			OccurrenceType:  "ResultCount",
			TriggerSource:   "AllResults",
			TriggerType:     "Critical",
			DetectionMethod: "StaticCondition",
		},
		{
			ThresholdType:    "LessThanOrEqual",
			Threshold:        40.0,
			TimeRange:        "1h",
			OccurrenceType:   "ResultCount",
			TriggerSource:    "AllResults",
			TriggerType:      "ResolvedCritical",
			DetectionMethod:  "StaticCondition",
			ResolutionWindow: "5m",
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
			ThresholdType:    "LessThanOrEqual",
			Threshold:        40.0,
			TimeRange:        "30m",
			OccurrenceType:   "ResultCount",
			TriggerSource:    "AllResults",
			TriggerType:      "ResolvedCritical",
			DetectionMethod:  "StaticCondition",
			ResolutionWindow: "15m",
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
					resource.TestCheckResourceAttr("sumologic_monitor.test", "evaluation_delay", canonicalTestEvaluationDelay),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "content_type", testContentType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "queries.0.row_id", testQueries[0].RowID),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.trigger_type", testTriggers[0].TriggerType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.0.time_range", testTriggers[0].TimeRange),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.1.resolution_window", testTriggers[1].ResolutionWindow),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notifications.0.notification.0.connection_type", testNotifications[0].Notification.(EmailNotification).ConnectionType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "playbook", testPlaybook),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "alert_name", testAlertName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notification_group_fields.0", testGroupFields[0]),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notification_group_fields.1", testGroupFields[1]),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "obj_permission.#", "2"),
					testAccCheckMonitorsLibraryMonitorFGPBackend("sumologic_monitor.test", t, genExpectedPermStmtsMonitor),
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
					resource.TestCheckResourceAttr("sumologic_monitor.test", "triggers.1.resolution_window", testUpdatedTriggers[1].ResolutionWindow),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notifications.0.notification.0.connection_type", testUpdatedNotifications[0].Notification.(EmailNotification).ConnectionType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "playbook", testUpdatedPlaybook),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "alert_name", testUpdatedAlertName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notification_group_fields.0", testUpdatedGroupFields[0]),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "notification_group_fields.1", testUpdatedGroupFields[1]),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "obj_permission.#", "1"),
					// 1, instead of 2
					testAccCheckMonitorsLibraryMonitorFGPBackend("sumologic_monitor.test", t, genExpectedPermStmtsForMonitorUpdate),
				),
			},
		},
	})
}

func TestAccSumologicMonitorsLibraryMonitor_driftingCorrectionFGP(t *testing.T) {
	var monitorsLibraryMonitor MonitorsLibraryMonitor
	testNameSuffix := acctest.RandString(16)
	tfResourceKey := "sumologic_monitor.test"
	testName := "terraform_test_monitor_" + testNameSuffix

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorsLibraryMonitorDestroy(monitorsLibraryMonitor),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMonitorsLibraryMonitor(testNameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorsLibraryMonitorExists(tfResourceKey, &monitorsLibraryMonitor, t),
					testAccCheckMonitorsLibraryMonitorAttributes(tfResourceKey),

					resource.TestCheckResourceAttr("sumologic_monitor.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "description",
						"terraform_test_monitor_description"),

					resource.TestCheckResourceAttr("sumologic_monitor.test",
						"obj_permission.#", "2"),
					testAccCheckMonitorsLibraryMonitorFGPBackend(tfResourceKey, t, genExpectedPermStmtsMonitor),
					// Emulating Drifting at the Backend
					testAccEmulateFGPDriftingMonitor(t),
				),
				// "After applying this step and refreshing, the plan was not empty"
				// Non-Empty Plan would occur, after the above step that emulates FGP drifting
				ExpectNonEmptyPlan: true,
			},
			// the following Test Step emulates running "terraform apply" again.
			// This step would detect and correct Drifting
			{
				Config: testAccSumologicMonitorsLibraryMonitor(testNameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorsLibraryMonitorExists(tfResourceKey, &monitorsLibraryMonitor, t),
					testAccCheckMonitorsLibraryMonitorAttributes(tfResourceKey),

					resource.TestCheckResourceAttr("sumologic_monitor.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "description",
						"terraform_test_monitor_description"),

					resource.TestCheckResourceAttr("sumologic_monitor.test",
						"obj_permission.#", "2"),
					testAccCheckMonitorsLibraryFolderFGPBackend(tfResourceKey, t, genExpectedPermStmtsMonitor),
				),
			},
		},
	})
}

func TestAccSumologicMonitorsLibraryMonitor_folder_update(t *testing.T) {
	var monitorsLibraryMonitor MonitorsLibraryMonitor
	testNameSuffix := acctest.RandString(16)

	testName := "terraform_test_monitor_" + testNameSuffix
	testType := "MonitorsLibraryMonitor"
	testMonitorType := "Logs"
	//testAlertName := "Alert from {{Name}}"
	folder1tfResourceKey := "sumologic_monitor_folder.tf_folder_01"
	folder2tfResourceKey := "sumologic_monitor_folder.tf_folder_02"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorsLibraryMonitorDestroy(monitorsLibraryMonitor),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMonitorsLibraryMonitorFolderUpdate(testNameSuffix, folder1tfResourceKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorsLibraryMonitorExists("sumologic_monitor.test", &monitorsLibraryMonitor, t),
					testAccCheckMonitorsLibraryMonitorAttributes("sumologic_monitor.test"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "monitor_type", testMonitorType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "type", testType),
					testAccCheckMonitorsLibraryMonitorFolderMatch("sumologic_monitor.test", folder1tfResourceKey, t),
				),
			},
			{
				Config: testAccSumologicMonitorsLibraryMonitorFolderUpdate(testNameSuffix, folder2tfResourceKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorsLibraryMonitorExists("sumologic_monitor.test", &monitorsLibraryMonitor, t),
					testAccCheckMonitorsLibraryMonitorAttributes("sumologic_monitor.test"),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "monitor_type", testMonitorType),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_monitor.test", "type", testType),
					testAccCheckMonitorsLibraryMonitorFolderMatch("sumologic_monitor.test", folder2tfResourceKey, t),
				),
			},
		},
	})
}

func TestAccSumologicMonitorsLibraryMonitorConnection_override_payload(t *testing.T) {
	var monitorsLibraryMonitor MonitorsLibraryMonitor
	testNameSuffix := acctest.RandString(16)

	testName := "terraform_test_monitor_connection_" + testNameSuffix
	testType := "MonitorsLibraryMonitor"
	testMonitorType := "Logs"
	defaultPayload := "{\"eventType\" : \"{{Name}}\"}"
	resolutionPayload := "{\"eventType\" : \"{{Name}}\"}"

	overrideDefaultPayload := "{\"eventType\" : \"{{Name}}-update\"}"
	overrideResolutionPayload := "{\"eventType\" : \"{{Name}}-update\"}"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorsLibraryMonitorDestroy(monitorsLibraryMonitor),
		//destroy conection too
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMonitorsLibraryMonitorUpdateConection(testNameSuffix, defaultPayload, resolutionPayload,
					overrideDefaultPayload, overrideResolutionPayload),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorsLibraryMonitorExists("sumologic_monitor.test_monitor_connection", &monitorsLibraryMonitor, t),
					testAccCheckMonitorsLibraryMonitorAttributes("sumologic_monitor.test_monitor_connection"),
					resource.TestCheckResourceAttr("sumologic_monitor.test_monitor_connection", "monitor_type", testMonitorType),
					resource.TestCheckResourceAttr("sumologic_monitor.test_monitor_connection", "name", testName),
					resource.TestCheckResourceAttr("sumologic_monitor.test_monitor_connection", "type", testType),
					resource.TestCheckResourceAttr("sumologic_monitor.test_monitor_connection", "notifications.0.notification.0.payload_override", overrideDefaultPayload+"\n"),
					resource.TestCheckResourceAttr("sumologic_monitor.test_monitor_connection", "notifications.0.notification.0.resolution_payload_override", overrideResolutionPayload+"\n"),
				),
			},
		},
	})
}

func testAccCheckMonitorsLibraryMonitorFolderMatch(monitorName string, folderName string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		//fetching monitor information
		monitorResource, ok := s.RootModule().Resources[monitorName]
		if !ok {
			//need this so that we don't get an unused import error for strconv in some cases
			return fmt.Errorf("Error = %s. MonitorsLibraryMonitor not found: %s", strconv.FormatBool(ok), monitorName)
		}

		//need this so that we don't get an unused import error for strings in some cases
		if strings.EqualFold(monitorResource.Primary.ID, "") {
			return fmt.Errorf("MonitorsLibraryMonitor ID is not set")
		}

		monitorResourceId := monitorResource.Primary.ID

		client := testAccProvider.Meta().(*Client)
		monitorsLibraryMonitor, err := client.MonitorsRead(monitorResourceId)

		if err != nil {
			return fmt.Errorf("MonitorsLibraryMonitor %s not found", monitorResourceId)
		}

		//fetching monitor folder information
		folderResource, ok := s.RootModule().Resources[folderName]
		if !ok {
			//need this so that we don't get an unused import error for strconv in some cases
			return fmt.Errorf("Error = %s. MonitorsLibraryFolder not found: %s", strconv.FormatBool(ok), folderName)
		}

		//need this so that we don't get an unused import error for strings in some cases
		if strings.EqualFold(folderResource.Primary.ID, "") {
			return fmt.Errorf("MonitorsLibraryFolder ID is not set")
		}

		folderResourceId := folderResource.Primary.ID
		monitorsLibraryFolder, err := client.MonitorsRead(folderResourceId)
		if err != nil {
			return fmt.Errorf("MonitorsLibraryFolder%s not found", monitorResourceId)
		}

		//checkig if the monitor parent id matches to the correct folder id
		if monitorsLibraryMonitor.ParentID != monitorsLibraryFolder.ID {
			return fmt.Errorf("Parent Id should be %s but %s", monitorsLibraryFolder.ID, monitorsLibraryMonitor.ParentID)
		}

		return nil
	}
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
	evaluation_delay = "60m"
	queries {
		row_id = "A"
		query = "_sourceCategory=monitor-manager error"
	  }
	triggers  {
		threshold_type = "GreaterThan"
		threshold = 40.0
		time_range = "-60m"
		occurrence_type = "ResultCount"
		trigger_source = "AllResults"
		trigger_type = "Critical"
		detection_method = "StaticCondition"
	  }
	triggers  {
		threshold_type = "LessThanOrEqual"
		threshold = 40.0
		time_range = "-60m"
		occurrence_type = "ResultCount"
		trigger_source = "AllResults"
		trigger_type = "ResolvedCritical"
		detection_method = "StaticCondition"
		resolution_window = "5m"
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
	obj_permission {
        subject_type = "role"
        subject_id = sumologic_role.tf_test_role_01.id
        permissions = ["Read","Update","Delete"]
    }
    obj_permission {
        subject_type = "role"
        subject_id = sumologic_role.tf_test_role_02.id
        permissions = ["Read"]
    }
}
resource "sumologic_role" "tf_test_role_01" {
	name        = "tf_test_role_01_%s"
	description = "Testing resource sumologic_role"
	capabilities = [
		"viewAlerts",
		"viewMonitorsV2",
		"manageMonitorsV2"
	]
}
resource "sumologic_role" "tf_test_role_02" {
	name        = "tf_test_role_02_%s"
	description = "Testing resource sumologic_role"
	capabilities = [
		"viewAlerts",
		"viewMonitorsV2",
		"manageMonitorsV2"
	]
}
resource "sumologic_role" "tf_test_role_03" {
	name        = "tf_test_role_03_%s"
	description = "Testing resource sumologic_role"
	capabilities = [
		"viewAlerts",
		"viewMonitorsV2",
		"manageMonitorsV2"
	]
}`, testName, testName, testName, testName)
}

func testAccSumologicMonitorsLibraryMonitorWithNoResolutionWindow(testName string) string {
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
			playbook = "This is a test playbook"  
			alert_name =  "Alert from {{Name}}"
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
		resolution_window = "15m"
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
		obj_permission {
          subject_type = "role"
          subject_id = sumologic_role.tf_test_role_01.id
          permissions = ["Read","Update"]
        }
}
resource "sumologic_role" "tf_test_role_01" {
	name        = "tf_test_role_01_%s"
	description = "Testing resource sumologic_role"
	capabilities = [
		"viewAlerts",
		"viewMonitorsV2",
		"manageMonitorsV2"
	]
}`, testName, testName)
}

func testAccSumologicMonitorsLibraryMonitorUpdateConection(testName string, defaultPayload string, resolutionPayload string,
	overrideDefaultPayload string, overrideResolutionPayload string) string {
	return fmt.Sprintf(`
resource "sumologic_connection" "connection_01" {
		name = "%s"
		type = "WebhookConnection"
		description = "WebhookConnection"
		url = "https://example.com"
		webhook_type = "Webhook"
		default_payload = <<JSON
%s
	JSON
		resolution_payload = <<JSON
%s
	JSON
}

resource "sumologic_monitor" "test_monitor_connection" {
	name = "terraform_test_monitor_connection_%s"
	description = "terraform_test_monitor_connection"
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
		resolution_window = "15m"
	  }
	notifications {
		notification {
			connection_type = "Webhook"			
			connection_id = sumologic_connection.connection_01.id
			payload_override = <<JSON
%s
			JSON
			resolution_payload_override = <<JSON
%s
			JSON
		}
		run_for_trigger_types = ["Critical", "ResolvedCritical"]
	}
	playbook = "This is an updated test playbook"
	alert_name = "Updated Alert from {{Name}}"
	
	obj_permission {
		subject_type = "role"
		subject_id = sumologic_role.tf_test_role_01.id
		permissions = ["Read","Update"]
	}
}
resource "sumologic_role" "tf_test_role_01" {
	name        = "tf_test_role_01_%s"
	description = "Testing resource sumologic_role"
	capabilities = [
		"viewAlerts",
		"viewMonitorsV2",
		"manageMonitorsV2"
	]
}`, testName, defaultPayload, resolutionPayload, testName, overrideDefaultPayload, overrideResolutionPayload, testName)
}

func testAccEmulateFGPDriftingMonitor(
	t *testing.T,
	// expectedFGPFunc func(*terraform.State, string) ([]CmfFgpPermStatement, error),
) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		monitorTargetId, resIdErr := getResourceID(s, "sumologic_monitor.test")
		if resIdErr != nil {
			return resIdErr
		}
		role01Id, resIdErr := getResourceID(s, "sumologic_role.tf_test_role_01")
		if resIdErr != nil {
			return resIdErr
		}
		role02Id, resIdErr := getResourceID(s, "sumologic_role.tf_test_role_02")
		if resIdErr != nil {
			return resIdErr
		}
		role03Id, resIdErr := getResourceID(s, "sumologic_role.tf_test_role_03")
		if resIdErr != nil {
			return resIdErr
		}

		client := testAccProvider.Meta().(*Client)
		expectedReadPermStmts := []CmfFgpPermStatement{
			{SubjectType: "role", SubjectId: role01Id, TargetId: monitorTargetId,
				Permissions: []string{"Read", "Update"}},
			{SubjectType: "role", SubjectId: role03Id, TargetId: monitorTargetId,
				Permissions: []string{"Read"}},
		}
		// using an empty Permissions array to achieve the effect of FGP Revocation
		setFGPPermStmts := append(expectedReadPermStmts,
			CmfFgpPermStatement{SubjectType: "role", SubjectId: role02Id, TargetId: monitorTargetId,
				Permissions: []string{}})

		_, setFgpErr := client.SetCmfFgp("monitors", CmfFgpRequest{
			PermissionStatements: setFGPPermStmts})
		if setFgpErr != nil {
			return setFgpErr
		}

		readfgpResult, readFgpErr := client.GetCmfFgp("monitors", monitorTargetId)
		if readFgpErr != nil {
			return readFgpErr
		}

		var expectedPermStmts []CmfFgpPermStatement
		expectedPermStmts = append(expectedPermStmts,
			CmfFgpPermStatement{
				SubjectId:   role01Id,
				SubjectType: "role",
				TargetId:    monitorTargetId,
				Permissions: []string{"Read", "Update"},
			},
			CmfFgpPermStatement{
				SubjectId:   role03Id,
				SubjectType: "role",
				TargetId:    monitorTargetId,
				Permissions: []string{"Read"},
			},
		)

		if !CmfFgpPermStmtSetEqual(readfgpResult.PermissionStatements, expectedPermStmts) {
			return fmt.Errorf("Permission Statements are different:\n  %+v\n  %+v\n",
				readfgpResult.PermissionStatements, expectedPermStmts)
		}
		return nil
	}
}

func testAccSumologicMonitorsLibraryMonitorFolderUpdate(testName string, parentIdTFString string) string {
	return fmt.Sprintf(`
resource "sumologic_monitor_folder" "tf_folder_01" {
	name = "tf_test_folder_01_%s"
	description = "1st folder"
}	
resource "sumologic_monitor_folder" "tf_folder_02" {
	name = "tf_test_folder_02_%s"
	description = "1st folder"
}
resource "sumologic_monitor" "test" {
	name = "terraform_test_monitor_%s"
	description = "terraform_test_monitor_description"
	type = "MonitorsLibraryMonitor"
	is_disabled = false
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
	parent_id = %s.id
}`, testName, testName, testName, parentIdTFString)
}

func exampleMonitorWithTriggerCondition(
	testName string,
	monitorType string,
	query string,
	trigger string,
	triggerTys []string) string {
	triggerTysStr := `"` + strings.Join(triggerTys, `","`) + `"`
	var resourceText = fmt.Sprintf(`
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
	return resourceText
}

var exampleLogsStaticTriggerConditionBlock = `
   logs_static_condition {
     critical {
       time_range = "60m"
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

var exampleLogsStaticTriggerConditionBlockWithResolutionWindow = `
   logs_static_condition {
     critical {
       time_range = "1h"
       alert {
         threshold = 100.0
         threshold_type = "GreaterThan"
       }
       resolution {
         threshold = 90
         threshold_type = "LessThanOrEqual"
		 resolution_window = "60m"
       }
     }
     field = "field"
   }`

var exampleMetricsStaticTriggerConditionBlock1 = `
   metrics_static_condition {
     critical {
       time_range = "30m"
       occurrence_type = "Always"
       alert {
         threshold = 100.0
         threshold_type = "GreaterThan"
         min_data_points = 4
       }
       resolution {
         threshold = 90
         threshold_type = "LessThanOrEqual"
         min_data_points = 7
       }
     }
   }`

var exampleMetricsStaticTriggerConditionBlock2 = `
   metrics_static_condition {
     critical {
       time_range = "60m"
       occurrence_type = "Always"
       alert {
         threshold = 100.0
         threshold_type = "GreaterThan"
         min_data_points = 6
       }
       resolution {
         threshold = 90
         threshold_type = "LessThanOrEqual"
         occurrence_type = "AtLeastOnce"
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

func exampleLogsStaticMonitorWithResolutionWindow(testName string) string {
	query := "error | timeslice 1m | count as field by _timeslice"
	return exampleMonitorWithTriggerCondition(testName, "Logs", query,
		exampleLogsStaticTriggerConditionBlockWithResolutionWindow, []string{"Critical", "ResolvedCritical"})
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
	exampleLogsStaticMonitorWithResolutionWindow,
	exampleMetricsStaticMonitor1,
	exampleMetricsStaticMonitor2,
	exampleLogsOutlierMonitor,
	exampleMetricsOutlierMonitor,
	exampleLogsMissingDataMonitor,
	exampleMetricsMissingDataMonitor,
}

func testAccSumologicMonitorsLibraryMonitorWithInvalidTriggerCondition(testName string, triggerCondition string) string {
	return fmt.Sprintf(`
resource "sumologic_monitor" "test" {
	name = "terraform_test_monitor_%s"
	description = "terraform_test_monitor_description"
	type = "MonitorsLibraryMonitor"
	is_disabled = false
	content_type = "Monitor"
	monitor_type = "Logs"
	evaluation_delay = "60m"
	queries {
		row_id = "A"
		query = "_sourceCategory=monitor-manager error"
	  }
	trigger_conditions  {
		%s
	}
}`, testName, triggerCondition)
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
		MinDataPoints:   4,
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

func genExpectedPermStmtsMonitor(s *terraform.State, targetId string) ([]CmfFgpPermStatement, error) {
	role01Id, resIdErr := getResourceID(s, "sumologic_role.tf_test_role_01")
	if resIdErr != nil {
		return nil, resIdErr
	}
	role02Id, resIdErr := getResourceID(s, "sumologic_role.tf_test_role_02")
	if resIdErr != nil {
		return nil, resIdErr
	}

	var expectedPermStmts []CmfFgpPermStatement
	expectedPermStmts = append(expectedPermStmts,
		CmfFgpPermStatement{
			SubjectId:   role01Id,
			SubjectType: "role",
			TargetId:    targetId,
			Permissions: []string{"Read", "Update", "Delete"},
		},
		CmfFgpPermStatement{
			SubjectId:   role02Id,
			SubjectType: "role",
			TargetId:    targetId,
			Permissions: []string{"Read"},
		},
	)
	return expectedPermStmts, nil
}

func genExpectedPermStmtsForMonitorUpdate(s *terraform.State, targetId string) ([]CmfFgpPermStatement, error) {
	role01Id, resIdErr := getResourceID(s, "sumologic_role.tf_test_role_01")
	if resIdErr != nil {
		return nil, resIdErr
	}

	var expectedPermStmts []CmfFgpPermStatement
	expectedPermStmts = append(expectedPermStmts,
		CmfFgpPermStatement{
			SubjectId:   role01Id,
			SubjectType: "role",
			TargetId:    targetId,
			Permissions: []string{"Read", "Update"},
		},
	)
	return expectedPermStmts, nil
}

func testAccCheckMonitorsLibraryMonitorFGPBackend(
	name string,
	t *testing.T,
	expectedFGPFunc func(*terraform.State, string) ([]CmfFgpPermStatement, error),
) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		targetId, resIdErr := getResourceID(s, name)
		if resIdErr != nil {
			return resIdErr
		}

		expectedPermStmts, resIdErr := expectedFGPFunc(s, targetId)
		if resIdErr != nil {
			return resIdErr
		}

		client := testAccProvider.Meta().(*Client)

		fgpResult, fgpErr := client.GetCmfFgp("monitors", targetId)
		if fgpErr != nil {
			return fgpErr
		}

		if !CmfFgpPermStmtSetEqual(fgpResult.PermissionStatements, expectedPermStmts) {
			return fmt.Errorf("Permission Statements are different:\n  %+v\n  %+v\n",
				fgpResult.PermissionStatements, expectedPermStmts)
		}
		return nil
	}
}

var allInvalidTriggerConditionMonitorResources = []func(testName string) string{
	invalidExampleWithNoTriggerCondition,
	invalidExampleWithEmptyLogStaticTriggerCondition,
	invalidExampleWithEmptyMetricsStaticTriggerCondition,
}

func invalidExampleWithNoTriggerCondition(testName string) string {
	query := "error | timeslice 1m | count as field by _timeslice"
	return exampleMonitorWithTriggerCondition(testName, "Logs", query,
		` `, []string{"Critical", "ResolvedCritical"})
}
func invalidExampleWithEmptyLogStaticTriggerCondition(testName string) string {
	query := "error | timeslice 1m | count as field by _timeslice"
	return exampleMonitorWithTriggerCondition(testName, "Logs", query,
		`logs_static_condition {}`, []string{"Critical", "ResolvedCritical"})
}
func invalidExampleWithEmptyMetricsStaticTriggerCondition(testName string) string {
	query := "error | timeslice 1m | count as field by _timeslice"
	return exampleMonitorWithTriggerCondition(testName, "Logs", query,
		`metrics_static_condition {}`, []string{"Critical", "ResolvedCritical"})
}
