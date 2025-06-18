package sumologic

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

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

func TestAccSumologicLogSearch_withValidIntervalTimeType(t *testing.T) {
	name := "TF IntervalTimeType Valid"
	description := "Testing interval_time_type with valid value"
	queryString := "error | count"
	parsingMode := "Manual"
	intervalTimeType := "receiptTime"
	literalRangeName := "today"
	tfResourceName := "tf_valid_interval_time_type"
	resourceName := fmt.Sprintf("sumologic_log_search.%s", tfResourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogSearchDestroy(LogSearch{}),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "sumologic_personal_folder" "personalFolder" {}
					resource "sumologic_log_search" "%s" {
						name = "%s"
						description = "%s"
						query_string = "%s"
						parsing_mode = "%s"
						parent_id = data.sumologic_personal_folder.personalFolder.id
						interval_time_type = "%s"
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
				`, tfResourceName, name, description, queryString, parsingMode, intervalTimeType, literalRangeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "interval_time_type", intervalTimeType),
				),
			},
		},
	})
}

func TestAccSumologicLogSearch_withInvalidIntervalTimeType(t *testing.T) {
	t := &testing.T{}
	name := "TF IntervalTimeType Invalid"
	queryString := "error | count"
	intervalTimeType := "invalidTime"
	tfResourceName := "tf_invalid_interval_time_type"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "sumologic_personal_folder" "personalFolder" {}
					resource "sumologic_log_search" "%s" {
						name = "%s"
						description = "Invalid test"
						query_string = "%s"
						parent_id = data.sumologic_personal_folder.personalFolder.id
						interval_time_type = "%s"
						time_range {
							begin_bounded_time_range {
								from {
									literal_time_range {
										range_name = "today"
									}
								}
							}
						}
					}
				`, tfResourceName, name, queryString, intervalTimeType),
				ExpectError: regexp.MustCompile("should be either 'messageTime' or 'receiptTime' or 'searchableTime'"),
			},
		},
	})
}

