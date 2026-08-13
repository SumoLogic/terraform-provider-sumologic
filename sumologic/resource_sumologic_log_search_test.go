package sumologic

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccPreCheckMultiNotification(t *testing.T) {
	if os.Getenv("SUMOLOGIC_MULTI_NOTIFICATION_ENABLED") == "" {
		t.Skip("Skipping multi-notification test: set SUMOLOGIC_MULTI_NOTIFICATION_ENABLED=true to run")
	}
}

func TestAccSumologicLogSearch_basic(t *testing.T) {
	var logSearch LogSearch
	name := "TF Import Search Test"
	description := "TF Import Search Test Description"
	queryString := "error | timeslice {{timeslice}} | count by _timeslice"
	parsingMode := "Manual"
	literalRangeName := "today"

	queryParameter := LogSearchQueryParameter{
		Name:        "timeslice",
		Description: "timeslice query param",
		DataType:    "ANY",
		Value:       "1d",
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
	searchParameters := []ScheduleSearchParameter{
		{
			Name:  "timeslice",
			Value: "15m",
		},
	}
	notificationThreshold := SearchNotificationThreshold{
		ThresholdType: "group",
		Operator:      "gt",
		Count:         10,
	}
	schedule := LogSearchSchedule{
		CronExpression:     "0 0 6 ? * 3 *",
		ParseableTimeRange: boundedTimeRange,
		TimeZone:           "America/Los_Angeles",
		Threshold:          &notificationThreshold,
		Parameters:         searchParameters,
		MuteErrorEmails:    true,
		Notification:       emailNotification,
		ScheduleType:       "Custom",
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
					runByReceiptTime, queryParameter, literalRangeName, schedule),
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
	queryString := "error | timeslice {{timeslice}} | count by _timeslice"
	parsingMode := "Manual"
	literalRangeName := "today"

	queryParameter := LogSearchQueryParameter{
		Name:        "timeslice",
		Description: "timeslice query param",
		DataType:    "ANY",
		Value:       "1d",
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
	searchParameters := []ScheduleSearchParameter{
		{
			Name:  "timeslice",
			Value: "15m",
		},
	}
	notificationThreshold := SearchNotificationThreshold{
		ThresholdType: "group",
		Operator:      "gt",
		Count:         10,
	}
	schedule := LogSearchSchedule{
		CronExpression:     "0 0 6 ? * 3 *",
		ParseableTimeRange: boundedTimeRange,
		TimeZone:           "America/Los_Angeles",
		Threshold:          &notificationThreshold,
		Parameters:         searchParameters,
		MuteErrorEmails:    true,
		Notification:       emailNotification,
		ScheduleType:       "Custom",
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
					runByReceiptTime, queryParameter, literalRangeName, schedule),
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
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.parameter.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.parameter.0.name", schedule.Parameters[0].Name),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.parameter.0.value", schedule.Parameters[0].Value),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.threshold.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.threshold.0.count",
						strconv.Itoa(schedule.Threshold.Count)),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.threshold.0.operator",
						schedule.Threshold.Operator),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.threshold.0.threshold_type",
						schedule.Threshold.ThresholdType),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.time_zone", schedule.TimeZone),

					resource.TestCheckResourceAttr(tfSearchResource, "query_parameter.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource, "query_parameter.0.name", queryParameter.Name),
					resource.TestCheckResourceAttr(tfSearchResource, "query_parameter.0.description", queryParameter.Description),
					resource.TestCheckResourceAttr(tfSearchResource, "query_parameter.0.data_type", queryParameter.DataType),
					resource.TestCheckResourceAttr(tfSearchResource, "query_parameter.0.value", queryParameter.Value),
				),
			},
		},
	})
}

func TestAccSumologicLogSearch_update(t *testing.T) {
	var logSearch LogSearch
	name := "TF Update Search Test"
	description := "TF Update Search Test Description"
	queryString := "error | timeslice {{timeslice}} | count by _timeslice"
	parsingMode := "Manual"
	literalRangeName := "today"

	queryParameter := LogSearchQueryParameter{
		Name:        "timeslice",
		Description: "timeslice query param",
		DataType:    "ANY",
		Value:       "1d",
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
			Name:  "timeslice",
			Value: "15m",
		},
	}
	notificationThreshold := SearchNotificationThreshold{
		ThresholdType: "group",
		Operator:      "gt",
		Count:         10,
	}
	schedule := LogSearchSchedule{
		CronExpression:     "0 0 6 ? * 3 *",
		ParseableTimeRange: boundedTimeRange,
		TimeZone:           "America/Los_Angeles",
		Threshold:          &notificationThreshold,
		Parameters:         searchParameters,
		MuteErrorEmails:    true,
		Notification:       emailNotification,
		ScheduleType:       "Custom",
	}
	runByReceiptTime := false

	// updated values
	newName := "TF Update Search Test New"
	newDescription := "TF Update Search Test New Description"
	newQueryString := "_sourceCategory={{source}} error | timeslice {{timeslice}} | count by _timeslice"
	newParsingMode := "AutoParse"
	newLiteralRangeName := "hour"

	newQueryParameters := []LogSearchQueryParameter{
		{
			Name:        "timeslice",
			Description: "timeslice query param",
			DataType:    "ANY",
			Value:       "1d",
		},
		{
			Name:        "source",
			Description: "source query param",
			DataType:    "STRING",
			Value:       "api",
		},
	}

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
	newSearchParameters := []ScheduleSearchParameter{
		{
			Name:  "timeslice",
			Value: "15m",
		},
		{
			Name:  "source",
			Value: "api",
		},
	}
	newSchedule.Parameters = newSearchParameters

	tfResourceName := "tf_update_search_test"
	tfSearchResource := fmt.Sprintf("sumologic_log_search.%s", tfResourceName)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogSearchDestroy(logSearch),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicLogSearch(tfResourceName, name, description, queryString, parsingMode,
					runByReceiptTime, queryParameter, literalRangeName, schedule),
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
					// query_parameters
					resource.TestCheckResourceAttr(tfSearchResource, "query_parameter.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource, "query_parameter.0.name", queryParameter.Name),
					// schedule
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.schedule_type", schedule.ScheduleType),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.mute_error_emails",
						strconv.FormatBool(schedule.MuteErrorEmails)),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.notification.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.notification.0.email_search_notification.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.notification.0.email_search_notification.0.subject_template", emailNotification.SubjectTemplate),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.parameter.#", "1"),
				),
			},
			{
				Config: testAccSumologicUpdatedLogSearch(tfResourceName, newName, newDescription, newQueryString, newParsingMode,
					runByReceiptTime, newQueryParameters, newLiteralRangeName, newSchedule),
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
					// query_parameters
					resource.TestCheckResourceAttr(tfSearchResource, "query_parameter.#", "2"),
					// schedule
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.schedule_type", newSchedule.ScheduleType),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.mute_error_emails",
						strconv.FormatBool(newSchedule.MuteErrorEmails)),
					// schedule notification
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.notification.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.notification.0.email_search_notification.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.notification.0.email_search_notification.0.include_histogram",
						strconv.FormatBool(newEmailNotification.IncludeHistogram)),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.notification.0.email_search_notification.0.include_query",
						strconv.FormatBool(newEmailNotification.IncludeQuery)),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.notification.0.email_search_notification.0.subject_template", newEmailNotification.SubjectTemplate),
					// schedule search parameters
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.parameter.#", "2"),
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
				return fmt.Errorf("Encountered an error: %w", err)
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
	parsingMode string, runByReceiptTime bool, queryParameter LogSearchQueryParameter, literalRangeName string,
	schedule LogSearchSchedule) string {

	emailNotification := schedule.Notification.(EmailSearchNotification)
	relativeTimeRange := schedule.ParseableTimeRange.(BeginBoundedTimeRange).From.(RelativeTimeRangeBoundary)
	tfSchedule := fmt.Sprintf(`
		schedule {
			cron_expression = "%s"
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

			parameter {
			  name = "%s"
			  value = "%s"
			}

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
		`, schedule.CronExpression, schedule.MuteErrorEmails,
		emailNotification.IncludeCsvAttachment, emailNotification.IncludeHistogram, emailNotification.IncludeQuery,
		emailNotification.IncludeResultSet, emailNotification.SubjectTemplate, emailNotification.ToList[0],
		schedule.Parameters[0].Name, schedule.Parameters[0].Value,
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
		query_parameter {
			name = "%s"
			description = "%s"
			data_type = "%s"
			value = "%s"
		}
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
	`, tfResourceName, name, description, queryString, parsingMode, runByReceiptTime,
		queryParameter.Name, queryParameter.Description, queryParameter.DataType, queryParameter.Value,
		literalRangeName, tfSchedule)
}

func testAccSumologicUpdatedLogSearch(tfResourceName string, name string, description string, queryString string,
	parsingMode string, runByReceiptTime bool, queryParameters []LogSearchQueryParameter, literalRangeName string,
	schedule LogSearchSchedule) string {

	emailNotification := schedule.Notification.(EmailSearchNotification)
	relativeTimeRange := schedule.ParseableTimeRange.(BeginBoundedTimeRange).From.(RelativeTimeRangeBoundary)
	tfSchedule := fmt.Sprintf(`
		schedule {
			cron_expression = "%s"
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

			parameter {
			  name = "%s"
			  value = "%s"
			}

			parameter {
			  name = "%s"
			  value = "%s"
			}

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
		`, schedule.CronExpression, schedule.MuteErrorEmails,
		emailNotification.IncludeCsvAttachment, emailNotification.IncludeHistogram, emailNotification.IncludeQuery,
		emailNotification.IncludeResultSet, emailNotification.SubjectTemplate, emailNotification.ToList[0],
		schedule.Parameters[0].Name, schedule.Parameters[0].Value, schedule.Parameters[1].Name, schedule.Parameters[1].Value,
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
		query_parameter {
			name = "%s"
			description = "%s"
			data_type = "%s"
			value = "%s"
		}
		query_parameter {
			name = "%s"
			description = "%s"
			data_type = "%s"
			value = "%s"
		}
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
	`, tfResourceName, name, description, queryString, parsingMode, runByReceiptTime,
		queryParameters[0].Name, queryParameters[0].Description, queryParameters[0].DataType, queryParameters[0].Value,
		queryParameters[1].Name, queryParameters[1].Description, queryParameters[1].DataType, queryParameters[1].Value,
		literalRangeName, tfSchedule)
}

func TestAccSumologicLogSearch_multi_notification(t *testing.T) {
	testAccPreCheckMultiNotification(t)
	var logSearch LogSearch
	name := "TF Multi Notification Search Test"
	description := "TF Multi Notification Search Test Description"
	queryString := "error | timeslice {{timeslice}} | count by _timeslice"
	parsingMode := "Manual"
	literalRangeName := "today"
	runByReceiptTime := false

	queryParameter := LogSearchQueryParameter{
		Name:        "timeslice",
		Description: "timeslice query param",
		DataType:    "ANY",
		Value:       "1d",
	}

	tfResourceName := "tf_multi_notif_test"
	tfSearchResource := fmt.Sprintf("sumologic_log_search.%s", tfResourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogSearchDestroy(logSearch),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicLogSearchMultiNotification(tfResourceName, name, description,
					queryString, parsingMode, runByReceiptTime, queryParameter, literalRangeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogSearchExists(tfSearchResource, &logSearch, t),
					resource.TestCheckResourceAttr(tfSearchResource, "name", name),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.notifications.#", "2"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.notifications.0.email_search_notification.0.to_list.0",
						"tf_multi_notif_1@sumologic.com"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.notifications.1.email_search_notification.0.to_list.0",
						"tf_multi_notif_2@sumologic.com"),
				),
			},
			{
				ResourceName:      tfSearchResource,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSumologicLogSearchMultiNotification(tfResourceName string, name string, description string,
	queryString string, parsingMode string, runByReceiptTime bool, queryParameter LogSearchQueryParameter,
	literalRangeName string) string {

	return fmt.Sprintf(`
	data "sumologic_personal_folder" "personalFolder" {}

	resource "sumologic_log_search" "%s" {
		name = "%s"
		description = "%s"
		query_string = "%s"
		parsing_mode = "%s"
		parent_id = data.sumologic_personal_folder.personalFolder.id
		run_by_receipt_time = %t
		query_parameter {
			name = "%s"
			description = "%s"
			data_type = "%s"
			value = "%s"
		}
		time_range {
			begin_bounded_time_range {
				from {
					literal_time_range {
						range_name = "%s"
					}
				}
			}
		}
		schedule {
			cron_expression = "0 0 6 ? * 3 *"
			mute_error_emails = false
			notifications {
				email_search_notification {
					include_csv_attachment = false
					include_histogram = true
					include_query = true
					include_result_set = true
					subject_template = "Alert 1: {{TriggerCondition}} for {{SearchName}}"
					to_list = [
						"tf_multi_notif_1@sumologic.com",
					]
				}
			}
			notifications {
				email_search_notification {
					include_csv_attachment = false
					include_histogram = false
					include_query = false
					include_result_set = true
					subject_template = "Alert 2: {{TriggerCondition}} for {{SearchName}}"
					to_list = [
						"tf_multi_notif_2@sumologic.com",
					]
				}
			}
			parameter {
				name = "timeslice"
				value = "15m"
			}
			parseable_time_range {
				begin_bounded_time_range {
					from {
						relative_time_range {
							relative_time = "-15m"
						}
					}
				}
			}
			schedule_type = "Custom"
			threshold {
				count = 10
				operator = "gt"
				threshold_type = "group"
			}
			time_zone = "America/Los_Angeles"
		}
	}
	`, tfResourceName, name, description, queryString, parsingMode, runByReceiptTime,
		queryParameter.Name, queryParameter.Description, queryParameter.DataType, queryParameter.Value,
		literalRangeName)
}

func TestAccSumologicLogSearch_multi_notification_email_and_webhook(t *testing.T) {
	testAccPreCheckMultiNotification(t)
	var logSearch LogSearch
	name := "TF Multi Notif Email+Webhook Test"
	description := "TF Multi Notification with Email and Webhook"
	queryString := "error | timeslice {{timeslice}} | count by _timeslice"
	parsingMode := "Manual"
	literalRangeName := "today"
	runByReceiptTime := false

	queryParameter := LogSearchQueryParameter{
		Name:        "timeslice",
		Description: "timeslice query param",
		DataType:    "ANY",
		Value:       "1d",
	}

	tfResourceName := "tf_multi_notif_email_webhook_test"
	tfSearchResource := fmt.Sprintf("sumologic_log_search.%s", tfResourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogSearchDestroy(logSearch),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicLogSearchMultiNotificationEmailAndWebhook(tfResourceName, name,
					description, queryString, parsingMode, runByReceiptTime, queryParameter, literalRangeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogSearchExists(tfSearchResource, &logSearch, t),
					testAccCheckLogSearchMultiNotifEmailAndWebhook(tfSearchResource, &logSearch, t),
					resource.TestCheckResourceAttr(tfSearchResource, "name", name),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.notifications.#", "2"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.notifications.0.email_search_notification.0.to_list.0",
						"tf_multi_notif_email_webhook@sumologic.com"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.notifications.1.webhook_search_notification.#", "1"),
				),
			},
			{
				ResourceName:      tfSearchResource,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLogSearchMultiNotifEmailAndWebhook(name string, logSearch *LogSearch, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("LogSearch not found: %s", name)
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)
		fetchedLogSearch, err := client.GetLogSearch(id)
		if err != nil {
			return fmt.Errorf("Error fetching LogSearch %s via API: %v", id, err)
		}

		if fetchedLogSearch.Schedule == nil {
			return fmt.Errorf("LogSearch %s has no schedule", id)
		}

		notifications := fetchedLogSearch.Schedule.Notifications
		if notifications == nil || len(notifications) == 0 {
			return fmt.Errorf("API response missing 'notifications' field for LogSearch %s", id)
		}

		if len(notifications) != 2 {
			return fmt.Errorf("Expected 2 notifications in API response, got %d", len(notifications))
		}

		foundEmail := false
		foundWebhook := false
		for _, n := range notifications {
			notifMap, ok := n.(map[string]interface{})
			if !ok {
				continue
			}
			taskType, _ := notifMap["taskType"].(string)
			if taskType == "EmailSearchNotificationSyncDefinition" {
				foundEmail = true
			}
			if taskType == "WebhookSearchNotificationSyncDefinition" {
				foundWebhook = true
			}
		}

		if !foundEmail {
			return fmt.Errorf("API response notifications missing Email notification for LogSearch %s", id)
		}
		if !foundWebhook {
			return fmt.Errorf("API response notifications missing Webhook notification for LogSearch %s", id)
		}

		t.Logf("API verification passed: LogSearch %s has 2 notifications (Email + Webhook)", id)
		return nil
	}
}

func testAccSumologicLogSearchMultiNotificationEmailAndWebhook(tfResourceName string, name string,
	description string, queryString string, parsingMode string, runByReceiptTime bool,
	queryParameter LogSearchQueryParameter, literalRangeName string) string {

	return fmt.Sprintf(`
	data "sumologic_personal_folder" "personalFolder" {}

	resource "sumologic_connection" "tf_test_webhook" {
		name        = "TF Test Webhook for Multi Notif"
		type        = "WebhookConnection"
		description = "Webhook connection for multi-notification test"
		url         = "https://example.com"
		webhook_type = "Webhook"
		default_payload = <<JSON
{"eventType": "{{Name}}"}
JSON
	}

	resource "sumologic_log_search" "%s" {
		name = "%s"
		description = "%s"
		query_string = "%s"
		parsing_mode = "%s"
		parent_id = data.sumologic_personal_folder.personalFolder.id
		run_by_receipt_time = %t
		query_parameter {
			name = "%s"
			description = "%s"
			data_type = "%s"
			value = "%s"
		}
		time_range {
			begin_bounded_time_range {
				from {
					literal_time_range {
						range_name = "%s"
					}
				}
			}
		}
		schedule {
			cron_expression = "0 0 6 ? * 3 *"
			mute_error_emails = false
			notifications {
				email_search_notification {
					include_csv_attachment = false
					include_histogram = true
					include_query = true
					include_result_set = true
					subject_template = "Alert: {{TriggerCondition}} for {{SearchName}}"
					to_list = [
						"tf_multi_notif_email_webhook@sumologic.com",
					]
				}
			}
			notifications {
				webhook_search_notification {
					webhook_id = sumologic_connection.tf_test_webhook.id
					payload    = "{\"alertName\": \"{{SearchName}}\"}"
				}
			}
			parameter {
				name = "timeslice"
				value = "15m"
			}
			parseable_time_range {
				begin_bounded_time_range {
					from {
						relative_time_range {
							relative_time = "-15m"
						}
					}
				}
			}
			schedule_type = "Custom"
			threshold {
				count = 10
				operator = "gt"
				threshold_type = "group"
			}
			time_zone = "America/Los_Angeles"
		}
	}
	`, tfResourceName, name, description, queryString, parsingMode, runByReceiptTime,
		queryParameter.Name, queryParameter.Description, queryParameter.DataType, queryParameter.Value,
		literalRangeName)
}

func TestAccSumologicLogSearch_singular_notification_backward_compat(t *testing.T) {
	var logSearch LogSearch
	name := "TF Singular Notification Backward Compat"
	description := "Verifies singular notification still works"
	queryString := "error | timeslice {{timeslice}} | count by _timeslice"
	parsingMode := "Manual"
	literalRangeName := "today"
	runByReceiptTime := false

	queryParameter := LogSearchQueryParameter{
		Name:        "timeslice",
		Description: "timeslice query param",
		DataType:    "ANY",
		Value:       "1d",
	}

	tfResourceName := "tf_singular_notif_compat_test"
	tfSearchResource := fmt.Sprintf("sumologic_log_search.%s", tfResourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogSearchDestroy(logSearch),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicLogSearchSingularNotification(tfResourceName, name, description,
					queryString, parsingMode, runByReceiptTime, queryParameter, literalRangeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogSearchExists(tfSearchResource, &logSearch, t),
					resource.TestCheckResourceAttr(tfSearchResource, "name", name),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.notification.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.notification.0.email_search_notification.0.to_list.0",
						"tf_singular_compat@sumologic.com"),
				),
			},
			{
				ResourceName:      tfSearchResource,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSumologicLogSearchSingularNotification(tfResourceName string, name string, description string,
	queryString string, parsingMode string, runByReceiptTime bool, queryParameter LogSearchQueryParameter,
	literalRangeName string) string {

	return fmt.Sprintf(`
	data "sumologic_personal_folder" "personalFolder" {}

	resource "sumologic_log_search" "%s" {
		name = "%s"
		description = "%s"
		query_string = "%s"
		parsing_mode = "%s"
		parent_id = data.sumologic_personal_folder.personalFolder.id
		run_by_receipt_time = %t
		query_parameter {
			name = "%s"
			description = "%s"
			data_type = "%s"
			value = "%s"
		}
		time_range {
			begin_bounded_time_range {
				from {
					literal_time_range {
						range_name = "%s"
					}
				}
			}
		}
		schedule {
			cron_expression = "0 0 6 ? * 3 *"
			mute_error_emails = false
			notification {
				email_search_notification {
					include_csv_attachment = false
					include_histogram = true
					include_query = true
					include_result_set = true
					subject_template = "Alert: {{TriggerCondition}} for {{SearchName}}"
					to_list = [
						"tf_singular_compat@sumologic.com",
					]
				}
			}
			parameter {
				name = "timeslice"
				value = "15m"
			}
			parseable_time_range {
				begin_bounded_time_range {
					from {
						relative_time_range {
							relative_time = "-15m"
						}
					}
				}
			}
			schedule_type = "Custom"
			threshold {
				count = 10
				operator = "gt"
				threshold_type = "group"
			}
			time_zone = "America/Los_Angeles"
		}
	}
	`, tfResourceName, name, description, queryString, parsingMode, runByReceiptTime,
		queryParameter.Name, queryParameter.Description, queryParameter.DataType, queryParameter.Value,
		literalRangeName)
}

func TestAccSumologicLogSearch_both_notification_fields_errors(t *testing.T) {
	var logSearch LogSearch
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogSearchDestroy(logSearch),
		Steps: []resource.TestStep{
			{
				Config: `
				data "sumologic_personal_folder" "personalFolder" {}

				resource "sumologic_log_search" "tf_both_notif_test" {
					name = "TF Both Notif Error Test"
					description = "Should fail validation"
					query_string = "error"
					parsing_mode = "Manual"
					parent_id = data.sumologic_personal_folder.personalFolder.id
					run_by_receipt_time = false
					time_range {
						begin_bounded_time_range {
							from {
								relative_time_range {
									relative_time = "-15m"
								}
							}
						}
					}
					schedule {
						cron_expression = "0 0 6 ? * 3 *"
						mute_error_emails = false
						notification {
							email_search_notification {
								include_csv_attachment = false
								include_histogram = true
								include_query = true
								include_result_set = true
								subject_template = "Alert"
								to_list = ["a@sumologic.com"]
							}
						}
						notifications {
							email_search_notification {
								include_csv_attachment = false
								include_histogram = true
								include_query = true
								include_result_set = true
								subject_template = "Alert"
								to_list = ["b@sumologic.com"]
							}
						}
						parseable_time_range {
							begin_bounded_time_range {
								from {
									relative_time_range {
										relative_time = "-15m"
									}
								}
							}
						}
						schedule_type = "Custom"
						time_zone = "America/Los_Angeles"
					}
				}
				`,
				ExpectError: regexp.MustCompile(`"schedule.0.notification": only one of`),
			},
		},
	})
}

func TestAccSumologicLogSearch_single_notification_in_notifications_array(t *testing.T) {
	testAccPreCheckMultiNotification(t)
	var logSearch LogSearch
	name := "TF Single Notif In Array"
	description := "Single notification using the notifications (plural) field"
	queryString := "error | timeslice {{timeslice}} | count by _timeslice"
	parsingMode := "Manual"
	literalRangeName := "today"
	runByReceiptTime := false

	queryParameter := LogSearchQueryParameter{
		Name:        "timeslice",
		Description: "timeslice query param",
		DataType:    "ANY",
		Value:       "1d",
	}

	tfResourceName := "tf_single_in_array_test"
	tfSearchResource := fmt.Sprintf("sumologic_log_search.%s", tfResourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogSearchDestroy(logSearch),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicLogSearchSingleNotifInArray(tfResourceName, name, description,
					queryString, parsingMode, runByReceiptTime, queryParameter, literalRangeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogSearchExists(tfSearchResource, &logSearch, t),
					resource.TestCheckResourceAttr(tfSearchResource, "name", name),
					resource.TestCheckResourceAttr(tfSearchResource, "schedule.0.notifications.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"schedule.0.notifications.0.email_search_notification.0.to_list.0",
						"tf_single_in_array@sumologic.com"),
				),
			},
			{
				ResourceName:      tfSearchResource,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSumologicLogSearchSingleNotifInArray(tfResourceName string, name string, description string,
	queryString string, parsingMode string, runByReceiptTime bool, queryParameter LogSearchQueryParameter,
	literalRangeName string) string {

	return fmt.Sprintf(`
	data "sumologic_personal_folder" "personalFolder" {}

	resource "sumologic_log_search" "%s" {
		name = "%s"
		description = "%s"
		query_string = "%s"
		parsing_mode = "%s"
		parent_id = data.sumologic_personal_folder.personalFolder.id
		run_by_receipt_time = %t
		query_parameter {
			name = "%s"
			description = "%s"
			data_type = "%s"
			value = "%s"
		}
		time_range {
			begin_bounded_time_range {
				from {
					literal_time_range {
						range_name = "%s"
					}
				}
			}
		}
		schedule {
			cron_expression = "0 0 6 ? * 3 *"
			mute_error_emails = false
			notifications {
				email_search_notification {
					include_csv_attachment = false
					include_histogram = true
					include_query = true
					include_result_set = true
					subject_template = "Alert: {{TriggerCondition}} for {{SearchName}}"
					to_list = [
						"tf_single_in_array@sumologic.com",
					]
				}
			}
			parameter {
				name = "timeslice"
				value = "15m"
			}
			parseable_time_range {
				begin_bounded_time_range {
					from {
						relative_time_range {
							relative_time = "-15m"
						}
					}
				}
			}
			schedule_type = "Custom"
			threshold {
				count = 10
				operator = "gt"
				threshold_type = "group"
			}
			time_zone = "America/Los_Angeles"
		}
	}
	`, tfResourceName, name, description, queryString, parsingMode, runByReceiptTime,
		queryParameter.Name, queryParameter.Description, queryParameter.DataType, queryParameter.Value,
		literalRangeName)
}

func TestAccSumologicLogSearch_intervalTimeType(t *testing.T) {
	var logSearch LogSearch
	name := "TF IntervalTimeType Test"
	description := "Verifies interval_time_type is correctly persisted"
	queryString := "error | timeslice {{timeslice}} | count by _timeslice"
	parsingMode := "Manual"
	literalRangeName := "today"
	runByReceiptTime := false
	intervalTimeType := "messageTime"

	queryParameter := LogSearchQueryParameter{
		Name:        "timeslice",
		Description: "Time slicing param",
		DataType:    "ANY",
		Value:       "1h",
	}

	tfResourceName := "tf_interval_time_type_test"
	tfSearchResource := fmt.Sprintf("sumologic_log_search.%s", tfResourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogSearchDestroy(logSearch),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "sumologic_personal_folder" "personalFolder" {}

					resource "sumologic_log_search" "%s" {
						name                = "%s"
						description         = "%s"
						query_string        = "%s"
						parsing_mode        = "%s"
						parent_id           = data.sumologic_personal_folder.personalFolder.id
						run_by_receipt_time = %t
						interval_time_type  = "%s"

						query_parameter {
							name        = "%s"
							description = "%s"
							data_type   = "%s"
							value       = "%s"
						}

						time_range {
							begin_bounded_time_range {
								from {
									literal_time_range {
										range_name = "%s"
									}
								}
							}
						}
					}
				`,
					tfResourceName, name, description, queryString, parsingMode,
					runByReceiptTime, intervalTimeType,
					queryParameter.Name, queryParameter.Description, queryParameter.DataType, queryParameter.Value,
					literalRangeName,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogSearchExists(tfSearchResource, &logSearch, t),
					resource.TestCheckResourceAttr(tfSearchResource, "interval_time_type", intervalTimeType),
				),
			},
		},
	})
}
