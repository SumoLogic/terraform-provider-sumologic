---
layout: "sumologic"
page_title: "SumoLogic: sumologic_log_search"
description: |-
  Provides a Sumologic Log Search
---

# sumologic_log_search
Provides a Sumologic Log Search.

## Example Usage
```hcl
data "sumologic_personal_folder" "personalFolder" {}

resource "sumologic_log_search" "example_log_search" {
    name = "Demo Search"
    description = "Demo search description"
    parent_id = data.sumologic_personal_folder.personalFolder.id
    query_string = "_sourceCategory=api | parse \"auth=*:*:*:*:*,\" as a,email,cid,uid,x | parse \"[logger=*,\" as loggerx | where loggerx matches {{logger}} | where cid matches {{customerID}} | count by _sourceHost"
    parsing_mode =  "AutoParse"
    run_by_receipt_time = true

    time_range {
        begin_bounded_time_range {
            from {
                relative_time_range {
                    relative_time = "-30m"
                }
            }
        }
    }
    
    query_parameter {
        name          = "logger"
        description   = "The logger for which the result will be returned"
        data_type     = "STRING"
        value = "*"
    }
    query_parameter {
        name          = "customerID"
        description   = "The customer id for which the result will be returned"
        data_type     = "STRING"
        value = "*"
    }

    schedule {
        cron_expression = "0 0 * * * ? *"
        mute_error_emails = false
        notification {
            email_search_notification {
                include_csv_attachment = false
                include_histogram = false
                include_query = true
                include_result_set = true
                subject_template = "Search Alert: {{TriggerCondition}} found for {{SearchName}}"
                to_list = [
                    "will@acme.com",
                ]
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
        schedule_type = "1Week"
        threshold {
            count = 10
            operator = "gt"
            threshold_type = "group"
        }
        time_zone = "America/Los_Angeles"
    }
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) Name of the search.
- `description` - (Optional) Description of the search.
- `parent_id` - (Required) The identifier of the folder to create the log search in.
- `query_string` - (Required) Log query to perform.
- `query_parameter` - (Block List, Optional) Upto 10 query_parameter blocks can be added one for each parameter in the query string. 
    See [query parameter schema](#schema-for-query_parameter).
- `parsing_mode` - (Optional) Define the parsing mode to scan the JSON format log messages. Possible values are:
    `AutoParse` and  `Manual`. Default value is `Manual`.

    In `AutoParse` mode, the system automatically figures out fields to parse based on the search query. While in
    the `Manual` mode, no fields are parsed out automatically. For more information see
    [Dynamic Parsing](https://help.sumologic.com/?cid=0011).
- `time_range` - (Block List, Max: 1, Required) Time range of the log search. See [time range schema](#schema-for-time_range)
- `schedule` - (Block List, Max: 1, Optional) Schedule of the log search. See [schedule schema](#schema-for-schedule)
- `run_by_receipt_time` - (Optional) This has the value `true` if the search is to be run by receipt time and
    `false` if it is to be run by message time. Default value is `false`.

### Schema for `query_parameter`
- `name` - (Required) The name of the parameter.
- `description` - (Optional) A description of the parameter.
- `data_type` - (Required) The data type of the parameter. Supported values are:
  1. `NUMBER`
  2. `STRING`
  3. `ANY`
  4. `KEYWORD`
- `value` - (Required) The default value for the parameter. Should be compatible with the type set in dataType field.


### Schema for `schedule`
- `cron_expression` - (Optional) Cron-like expression specifying the search's schedule. `schedule_type` must be set
    to "Custom", otherwise, `schedule_type` takes precedence over `cron_expression`.
- `schedule_type` - (Required) Run schedule of the scheduled search. Set to "Custom" to specify the schedule with
    a CRON expression. Possible schedule types are: `RealTime`, `15Minutes`, `1Hour`, `2Hours`, `4Hours`, `6Hours`,
    `8Hours`, `12Hours`, `1Day`, `1Week`, `Custom`.

    -> With `Custom`, `1Day` and `1Week` schedule types you need to provide the corresponding cron expression
    to determine when to actually run the search. E.g. valid cron for `1Day` is `0 0 16 ? * 2-6 *`.
- `parseable_time_range` - (Block List, Max: 1, Required) Time range of the scheduled log search. See
    [time range schema](#schema-for-time_range)
- `time_zone` - (Required) Time zone for the scheduled log search. Either an abbreviation such as "PST",
    a full name such as "America/Los_Angeles", or a custom ID such as "GMT-8:00". Note that the support of
    abbreviations is for JDK 1.1.x compatibility only and full names should be used.
- `threshold` - (Block List, Max: 1, Optional) Threshold for when to send notification. See
    [threshold schema](#schema-threshold)
- `notification` - (Block List, Max: 1, Required) Notification of the log search. See
    [notification schema](#schema-for-notification)
- `mute_error_emails` - (Optional) If enabled, emails are not sent out in case of errors with the search.
- `parameters` - (Block List, Optional) A list of scheduled search parameters. See
    [parameter schema](#schema-for-parameter)


### Schema for `time_range`
- `complete_literal_time_range` - (Block List, Max: 1, Optional) Literal time range. See
[complete_literal_time_range schema](#schema-for-complete_literal_time_range) for details.
- `begin_bounded_time_range` - (Block List, Max: 1, Optional) Bounded time range. See
[begin_bounded_time_range schema](#schema-for-begin_bounded_time_range) schema for details.

### Schema for `complete_literal_time_range`
- `range_name` - (Required) Name of complete literal time range. One of `today`, `yesterday`, `previous_week`, and
    `previous_month`.

### Schema for `begin_bounded_time_range`
- `from` - (Block List, Max: 1, Required) Start boundary of bounded time range. See
[time_range_boundary schema](#schema-for-time_range_boundary) for details.
- `to` - (Block List, Max: 1, Optional) End boundary of bounded time range. See
[time_range_boundary schema](#schema-for-time_range_boundary) for details.

### Schema for `time_range_boundary`
- `epoch_time_range` - (Block List, Optional) Time since the epoch.
    - `epoch_millis` - (Required) Time as a number of milliseconds since the epoch.

- `iso8601_time_range` - (Block List, Optional) Time in ISO 8601 format.
    - `iso8601_time` - (Required) Time as a string in ISO 8601 format.

- `relative_time_range` - (Block List, Optional) Time in relative format.
    - `relative_time` - (Required) Relative time as a string consisting of following elements:
      1. `-` (optional): minus sign indicates time in the past,
      2. `<number>`: number of time units,
      3. `<time_unit>`: time unit; possible values are: `w` (week), `d` (day), `h` (hour), `m` (minute), `s` (second).

      Multiple pairs of `<number><time_unit>` may be provided, and they may be in any order. For example,
      `-2w5d3h` points to the moment in time 2 weeks, 5 days and 3 hours ago.

- `literal_time_range` - (Block List, Optional) Time in literal format.
    - `range_name` - (Required) One of `now`, `second`, `minute`, `hour`, `day`, `today`, `week`, `month`, `year`.


### Schema for `notification`
- `alert_search_notification` - (Block List, Max: 1, Optional) Run an script action. See
[alert_search_notification schema](#schema-for-alert_search_notification) for details.
- `cse_signal_notification` - (Block List, Max: 1, Optional) Create a CSE signal with a scheduled search.
See [cse_signal_notification schema](#schema-for-cse_signal_notification) schema for details.
- `email_search_notification` - (Block List, Max: 1, Optional) Send an alert via email. See
[email_search_notification schema](#schema-for-email_search_notification) schema for details.
- `save_to_lookup_notification` - (Block List, Max: 1, Optional) Save results to a Lookup Table. See
[save_to_lookup_notification schema](#schema-for-save_to_lookup_notification) schema for details.
- `save_to_view_notification` - (Block List, Max: 1, Optional) Save results to an index. See
[save_to_view_notification schema](#schema-for-save_to_view_notification) schema for details.
- `service_now_search_notification` - (Block List, Max: 1, Optional) Send results to Service Now. See
[service_now_search_notification schema](#schema-for-service_now_search_notification) schema for details.
- `webhook_search_notification` - (Block List, Max: 1, Optional) Send an alert via Webhook. See
[webhook_search_notification schema](#schema-for-webhook_search_notification) schema for details.

### Schema for `alert_search_notification`
- `source_id` - (Required) Identifier of the collector's source.

### Schema for `cse_signal_notification`
- `record_type` - (Required) Name of the Cloud SIEM Enterprise Record to be created.

### Schema for `email_search_notification`
- `subject_template` - (Optional) Subject of the email. If the notification is scheduled with a threshold,
    the default subject template will be `Search Alert: {{AlertCondition}} results found for {{SearchName}}`.
    For email notifications without a threshold, the default subject template is `Search Results: {{SearchName}}`.
- `to_list` - (Block List, Required) A list of email recipients.
- `include_query` - (Optional) If the search query should be included in the notification email.
- `include_result_set` - (Optional) If the search result set should be included in the notification email.
- `include_histogram` - (Optional) If the search result histogram should be included in the notification email.
- `include_csv_attachment` - (Optional) If the search results should be included in the notification email
    as a CSV attachment.

### Schema for `save_to_lookup_notification`
- `lookup_file_path` - (Required) Path of the lookup table to save the results to.
- `is_lookup_merge_operation` - (Required) Whether to merge the file contents with existing data in the lookup table.

### Schema for `save_to_view_notification`
- `view_name` - (Required) Name of the View(Index) to save the results to.

### Schema for `service_now_search_notification`
- `external_id` - (Required) Service Now Identifier.
- `fields` - (Block List, Optional) Service Now fields.
    - `event_type` - (Optional) The category that the event source uses to identify the event.
    - `severity` - (Optional) An integer value representing the severity of the alert. Supported values are:
        * 0 for Clear
        * 1 for Critical
        * 2 for Major
        * 3 for Minor
        * 4 for Warning
    - `resource` - (Optional) The component on the node to which the event applies.
    - `node` - (Optional) The physical or virtual device on which the event occurred.


### Schema for `webhook_search_notification`
- `webhook_id` - (Required) Identifier of the webhook connection.
- `payload` - (Optional) A JSON object in the format required by the target WebHook URL.
- `itemize_alerts` - (Optional) If set to true, one webhook per result will be sent when the trigger conditions are met.
- `max_itemized_alerts` - (Optional) The maximum number of results for which we send separate alerts.


### Schema for `threshold`
- `threshold_type` - (Required) Threshold type for the scheduled log search. Possible values are: `message` and `group`.
    Use `group` as threshold type if the search query is of aggregate type. For non-aggregate queries, set it
    to `message`.
- `operator` - (Required) Criterion to be applied when comparing actual result count with expected count. Possible
    values are: `eq`, `gt`, `ge`, `lt`, and `le`.
- `count` - (Required) Expected result count.


### Schema for `parameter`
- `name` - (Required) Name of scheduled search parameter.
- `value` - (Required) Value of scheduled search parameter.


## Attributes reference
In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the log search.


## Import
A log search can be imported using it's identifier, e.g.:
```hcl
terraform import sumologic_log_search.example_search 0000000007FFD79D
```
