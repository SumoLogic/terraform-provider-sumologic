package sumologic

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicLogSearch_basic(t *testing.T) {
	var logSearch LogSearch
	name := "TF Import Search Test"
	description := "TF Import Search Test Description"
	queryString := "error | count by _sourceCategory"
	parsingMode := "Manual"
	literalRangeName := "today"

	queryParameters := []LogSearchQueryParameter{
		{
			Name:        "",
			Description: "",
			DataType:    "",
			Value:       "",
		},
	}

	boundedTimeRange := BeginBoundedTimeRange{
		From: RelativeTimeRangeBoundary{
			RelativeTime: "-15m",
		},
	}
	emailNotification := EmailSearchNotification{
		ToList:               []string{"tf_import_search_test@sumologic.com"},
		SubjectTemplate:      "Search Alert: {{TriggerCondition}} found for {{SearchName}}",
		IncludeQuery:         false,
		IncludeResultSet:     true,
		IncludeHistogram:     true,
		IncludeCsvAttachment: false,
	}
	// TODO test parameters for scheduled search
	searchParameters := []ScheduleSearchParameter{
		{
			Name:  "TODO",
			Value: "TODO",
		},
	}
	notificationThreshold := SearchNotificationThreshold{
		ThresholdType: "group",
		Operator:      "gt",
		Count:         10,
	}
	schedule := LogSearchSchedule{
		CronExpression:       "0 0 6 ? * 3 *",
		DisplayableTimeRange: "-15m",
		ParseableTimeRange:   boundedTimeRange,
		TimeZone:             "America/Los_Angeles",
		Threshold:            &notificationThreshold,
		Parameters:           searchParameters,
		MuteErrorEmails:      true,
		Notification:         emailNotification,
		ScheduleType:         "Custom",
	}
	runByReceiptTime := false

	tfResourceName := "tf_import_search_test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogSearchDestroy(logSearch),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicLogSearch(tfResourceName, name, description, queryString, parsingMode,
					runByReceiptTime, queryParameters, literalRangeName, schedule),
			},
			{
				ResourceName:      fmt.Sprintf("sumologic_log_search.%s", tfResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSumologicLogSearch_create(t *testing.T) {
	var logSearch LogSearch
	name := "TF Create Search Test"
	description := "TF Create Search Test Description"
	queryString := "error | count by _sourceCategory"
	parsingMode := "Manual"
	literalRangeName := "today"

	// TODO test query parameters
	queryParameters := []LogSearchQueryParameter{
		{
			Name:        "",
			Description: "",
			DataType:    "",
			Value:       "",
		},
	}

	boundedTimeRange := BeginBoundedTimeRange{
		From: RelativeTimeRangeBoundary{
			RelativeTime: "-15m",
		},
	}
	emailNotification := EmailSearchNotification{
		ToList:               []string{"tf_create_search_test@sumologic.com"},
		SubjectTemplate:      "Search Alert: {{TriggerCondition}} found for {{SearchName}}",
		IncludeQuery:         false,
		IncludeResultSet:     true,
		IncludeHistogram:     true,
		IncludeCsvAttachment: false,
	}
	// TODO test parameters for scheduled search
	searchParameters := []ScheduleSearchParameter{
		{
			Name:  "TODO",
			Value: "TODO",
		},
	}
	notificationThreshold := SearchNotificationThreshold{
		ThresholdType: "group",
		Operator:      "gt",
		Count:         10,
	}
	schedule := LogSearchSchedule{
		CronExpression:       "0 0 6 ? * 3 *",
		DisplayableTimeRange: "-15m",
		ParseableTimeRange:   boundedTimeRange,
		TimeZone:             "America/Los_Angeles",
		Threshold:            &notificationThreshold,
		Parameters:           searchParameters,
		MuteErrorEmails:      true,
		Notification:         emailNotification,
		ScheduleType:         "Custom",
	}
	runByReceiptTime := false

	tfResourceName := "tf_create_search_test"
	tfSearchResource := fmt.Sprintf("sumologic_log_search.%s", tfResourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogSearchDestroy(logSearch),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicLogSearch(tfResourceName, name, description, queryString, parsingMode,
					runByReceiptTime, queryParameters, literalRangeName, schedule),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogSearchExists(tfSearchResource, &logSearch, t),

					resource.TestCheckResourceAttr(tfSearchResource, "name", name),
					resource.TestCheckResourceAttr(tfSearchResource, "description", description),
					resource.TestCheckResourceAttr(tfSearchResource, "query_string", queryString),
					resource.TestCheckResourceAttr(tfSearchResource, "parsing_mode", parsingMode),
					resource.TestCheckResourceAttr(tfSearchResource, "run_by_receipt_time", strconv.FormatBool(runByReceiptTime)),

					// timerange
					resource.TestCheckResourceAttr(tfSearchResource, "time_range.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"time_range.0.begin_bounded_time_range.0.from.0.literal_time_range.0.range_name",
						literalRangeName),

					// schedule
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.cron_expression", schedule.CronExpression),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.displayable_time_range",
						schedule.DisplayableTimeRange),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.mute_error_emails",
						strconv.FormatBool(schedule.MuteErrorEmails)),
					// email notification
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.notification.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.notification.0.email_search_notification.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.notification.0.email_search_notification.0.include_csv_attachment",
						strconv.FormatBool(emailNotification.IncludeCsvAttachment)),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.notification.0.email_search_notification.0.include_histogram",
						strconv.FormatBool(emailNotification.IncludeHistogram)),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.notification.0.email_search_notification.0.include_query",
						strconv.FormatBool(emailNotification.IncludeQuery)),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.notification.0.email_search_notification.0.include_result_set",
						strconv.FormatBool(emailNotification.IncludeResultSet)),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.notification.0.email_search_notification.0.subject_template", emailNotification.SubjectTemplate),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.notification.0.email_search_notification.0.to_list.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.notification.0.email_search_notification.0.to_list.0", emailNotification.ToList[0]),

					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.parseable_time_range.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.parseable_time_range.0.begin_bounded_time_range.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.parseable_time_range.0.begin_bounded_time_range.0.from.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.parseable_time_range.0.begin_bounded_time_range.0.from.0.relative_time_range.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.parseable_time_range.0.begin_bounded_time_range.0.from.0.relative_time_range.0.relative_time",
						schedule.ParseableTimeRange.(BeginBoundedTimeRange).From.(RelativeTimeRangeBoundary).RelativeTime),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.schedule_type", schedule.ScheduleType),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.threshold.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.threshold.0.count",
						strconv.Itoa(schedule.Threshold.Count)),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.threshold.0.operator",
						schedule.Threshold.Operator),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.threshold.0.threshold_type",
						schedule.Threshold.ThresholdType),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.time_zone", schedule.TimeZone),

					// TODO query parameters
					// resource.TestCheckResourceAttr(tfSearchResource, "query_parameter.0", strings.Replace(testQueryParameters[0], "\"", "", 2)),
				),
			},
		},
	})
}

func TestAccSumologicLogSearch_update(t *testing.T) {
	var logSearch LogSearch
	name := "TF Update Search Test"
	description := "TF Update Search Test Description"
	queryString := "error | count by _sourceCategory"
	parsingMode := "Manual"
	literalRangeName := "today"

	// TODO test query parameters
	queryParameters := []LogSearchQueryParameter{
		{
			Name:        "",
			Description: "",
			DataType:    "",
			Value:       "",
		},
	}

	boundedTimeRange := BeginBoundedTimeRange{
		From: RelativeTimeRangeBoundary{
			RelativeTime: "-15m",
		},
	}
	emailNotification := EmailSearchNotification{
		ToList:               []string{"tf_update_search_test@sumologic.com"},
		SubjectTemplate:      "Search Alert: {{TriggerCondition}} found for {{SearchName}}",
		IncludeQuery:         false,
		IncludeResultSet:     true,
		IncludeHistogram:     true,
		IncludeCsvAttachment: false,
	}
	searchParameters := []ScheduleSearchParameter{
		{
			Name:  "TODO",
			Value: "TODO",
		},
	}
	notificationThreshold := SearchNotificationThreshold{
		ThresholdType: "group",
		Operator:      "gt",
		Count:         10,
	}
	schedule := LogSearchSchedule{
		CronExpression:       "0 0 6 ? * 3 *",
		DisplayableTimeRange: "-15m",
		ParseableTimeRange:   boundedTimeRange,
		TimeZone:             "America/Los_Angeles",
		Threshold:            &notificationThreshold,
		Parameters:           searchParameters,
		MuteErrorEmails:      true,
		Notification:         emailNotification,
		ScheduleType:         "Custom",
	}
	runByReceiptTime := false

	// updated values
	newName := "TF Update Search Test New"
	newDescription := "TF Update Search Test New Description"
	newQueryString := "warn | count by _sourceCategory"
	newParsingMode := "AutoParse"
	newLiteralRangeName := "hour"

	newEmailNotification := emailNotification
	newEmailNotification.ToList = []string{
		"tf_update_search_test@sumologic.com", "tf_update_new_search_test@sumologic.com",
	}
	newEmailNotification.IncludeHistogram = false
	newEmailNotification.IncludeQuery = true
	newEmailNotification.SubjectTemplate = "{{TriggerCondition}} found for {{SearchName}}"
	newSchedule := schedule
	newSchedule.ScheduleType = "1Day"
	newSchedule.MuteErrorEmails = false
	newSchedule.Notification = newEmailNotification

	tfResourceName := "tf_update_search_test"
	tfSearchResource := fmt.Sprintf("sumologic_log_search.%s", tfResourceName)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogSearchDestroy(logSearch),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicLogSearch(tfResourceName, name, description, queryString, parsingMode,
					runByReceiptTime, queryParameters, literalRangeName, schedule),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogSearchExists(tfSearchResource, &logSearch, t),

					resource.TestCheckResourceAttr(tfSearchResource, "name", name),
					resource.TestCheckResourceAttr(tfSearchResource, "description", description),
					resource.TestCheckResourceAttr(tfSearchResource, "query_string", queryString),
					resource.TestCheckResourceAttr(tfSearchResource, "parsing_mode", parsingMode),
					resource.TestCheckResourceAttr(tfSearchResource, "run_by_receipt_time", strconv.FormatBool(runByReceiptTime)),

					// timerange
					resource.TestCheckResourceAttr(tfSearchResource, "time_range.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"time_range.0.begin_bounded_time_range.0.from.0.literal_time_range.0.range_name",
						literalRangeName),
				),
			},
			{
				Config: testAccSumologicLogSearch(tfResourceName, newName, newDescription, newQueryString, newParsingMode,
					runByReceiptTime, queryParameters, newLiteralRangeName, newSchedule),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfSearchResource, "name", newName),
					resource.TestCheckResourceAttr(tfSearchResource, "description", newDescription),
					resource.TestCheckResourceAttr(tfSearchResource, "query_string", newQueryString),
					resource.TestCheckResourceAttr(tfSearchResource, "parsing_mode", newParsingMode),
					resource.TestCheckResourceAttr(tfSearchResource, "run_by_receipt_time", strconv.FormatBool(runByReceiptTime)),

					// timerange
					resource.TestCheckResourceAttr(tfSearchResource, "time_range.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"time_range.0.begin_bounded_time_range.0.from.0.literal_time_range.0.range_name",
						newLiteralRangeName),
				),
			},
		},
	})
}

func testAccCheckLogSearchDestroy(logSearch LogSearch) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			search, err := client.GetLogSearch(id)
			if err != nil {
				return fmt.Errorf("Encountered an error: " + err.Error())
			}
			if search != nil {
				return fmt.Errorf("LogSearch %s still exists", id)
			}
		}
		return nil
	}
}

func testAccCheckLogSearchExists(name string, logSearch *LogSearch, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			//need this so that we don't get an unused import error for strconv in some cases
			return fmt.Errorf("Error = %s. LogSearch not found: %s", strconv.FormatBool(ok), name)
		}

		//need this so that we don't get an unused import error for strings in some cases
		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("LogSearch ID is not set")
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)
		newLogSearch, err := client.GetLogSearch(id)
		if err != nil {
			return fmt.Errorf("LogSearch %s not found", id)
		}
		logSearch = newLogSearch
		return nil
	}
}

func testAccSumologicLogSearch(tfResourceName string, name string, description string, queryString string,
	parsingMode string, runByReceiptTime bool, queryParameters []LogSearchQueryParameter, literalRangeName string,
	schedule LogSearchSchedule) string {

	emailNotification := schedule.Notification.(EmailSearchNotification)
	relativeTimeRange := schedule.ParseableTimeRange.(BeginBoundedTimeRange).From.(RelativeTimeRangeBoundary)
	tfSchedule := fmt.Sprintf(`
		schedule {
			cron_expression = "%s"
			displayable_time_range = "%s"
			mute_error_emails = %t
			notification {
				email_search_notification {
					include_csv_attachment = %t
					include_histogram = %t
					include_query = %t
					include_result_set = %t
					subject_template = "%s"
					to_list = [
						"%s",
					]
				}
			}

			#parameter {
			#  name = "key"
			#  value = "value"
			#}

			parseable_time_range {
				begin_bounded_time_range {
					from {
						relative_time_range {
							relative_time = "%s"
						}
					}
				}
			}
			schedule_type = "%s"
			threshold {
				count = %d
				operator = "%s"
				threshold_type = "%s"
			}
			time_zone = "%s"
		}
		`, schedule.CronExpression, schedule.DisplayableTimeRange, schedule.MuteErrorEmails,
		emailNotification.IncludeCsvAttachment, emailNotification.IncludeHistogram, emailNotification.IncludeQuery,
		emailNotification.IncludeResultSet, emailNotification.SubjectTemplate, emailNotification.ToList[0],
		relativeTimeRange.RelativeTime, schedule.ScheduleType, schedule.Threshold.Count, schedule.Threshold.Operator,
		schedule.Threshold.ThresholdType, schedule.TimeZone)

	return fmt.Sprintf(`
	data "sumologic_personal_folder" "personalFolder" {}

	resource "sumologic_log_search" "%s" {
		name = "%s"
		description = "%s"
		query_string = "%s"
		parsing_mode = "%s"
		parent_id = data.sumologic_personal_folder.personalFolder.id
		run_by_receipt_time = %t
		time_range {
			begin_bounded_time_range {
				from {
					literal_time_range {
						range_name = "%s"
					}
				}
			}
		}
		#schedule
		%s
	}
	`, tfResourceName, name, description, queryString, parsingMode, runByReceiptTime, literalRangeName, tfSchedule)
}
