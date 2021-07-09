---
layout: 'sumologic'
page_title: 'SumoLogic: sumologic_monitor'
description: |-
  Provides the ability to create, read, delete, and update monitors.
---

# sumologic_monitor

Provides the ability to create, read, delete, and update [Monitors][1].

## Example Logs Monitor

```hcl
resource "sumologic_monitor" "tf_logs_monitor_1" {
  name = "Terraform Logs Monitor"
  description = "tf logs monitor"
  type = "MonitorsLibraryMonitor"
  is_disabled = false
  content_type = "Monitor"
  monitor_type = "Logs"
  queries {
      row_id = "A"
      query = "_sourceCategory=event-action info"
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
      recipients = [
        "abc@example.com",
      ]
      subject = "Monitor Alert: {{TriggerType}} on {{Name}}"
      time_zone = "PST"
      message_body = "Triggered {{TriggerType}} Alert on {{Name}}: {{QueryURL}}"
    }
    run_for_trigger_types = ["Critical", "ResolvedCritical"]
  }
  notifications {
    notification {
      connection_type = "Webhook"
      connection_id = "0000000000ABC123"
    }
    run_for_trigger_types = ["Critical", "ResolvedCritical"]
  }
}
```

## Example Metrics Monitor

```hcl
resource "sumologic_monitor" "tf_metrics_monitor_1" {
  name = "Terraform Metrics Monitor"
  description = "tf metrics monitor"
  type = "MonitorsLibraryMonitor"
  is_disabled = false
  content_type = "Monitor"
  monitor_type = "Metrics"
  queries {
      row_id = "A"
      query = "metric=CPU_Idle _sourceCategory=event-action"
  }
  triggers  {
      threshold_type = "GreaterThanOrEqual"
      threshold = 40.0
      time_range = "15m"
      occurrence_type = "AtLeastOnce"
      trigger_source = "AnyTimeSeries"
      trigger_type = "Critical"
      detection_method = "StaticCondition"
    }
  triggers {
    threshold_type = "LessThan"
    threshold = 40.0
    time_range = "15m"
    occurrence_type = "Always"
    trigger_source = "AnyTimeSeries"
    trigger_type = "ResolvedCritical"
    detection_method = "StaticCondition"
    }
  notifications {
    notification {
      connection_type = "Email"
      recipients = ["abc@example.com"]
      subject = "Triggered {{TriggerType}} Alert on Monitor {{Name}}"
      time_zone = "PST"
      message_body = "Triggered {{TriggerType}} Alert on {{Name}}: {{QueryURL}}"
    }
    run_for_trigger_types = ["Critical","ResolvedCritical"]
  }
}
```

## Example Logs Monitor with Webhook Connection and Folder

```hcl
resource "sumologic_monitor_folder" "tf_monitor_folder_1" {
  name = "Terraform Managed Folder 1"
  description = "A folder for Monitors"
}

resource "sumologic_connection" "example_pagerduty_connection" {
  name = "example_pagerduty_connection"
  description = "PagerDuty connection for notifications from Monitors"
  type = "WebhookConnection"
  webhook_type = "PagerDuty"
  url = "https://events.pagerduty.com/"
  default_payload = <<JSON
{
  "service_key": "pagerduty_api_integration_key",
  "event_type": "trigger",
  "description": "PagerDuty connection for notifications",
  "client": "Sumo Logic",
  "client_url": ""
}
JSON
}

resource "sumologic_monitor" "tf_logs_monitor_2" {
  name = "Terraform Logs Monitor with Webhook Connection"
  description = "tf logs monitor with webhook"
  type = "MonitorsLibraryMonitor"
  parent_id = sumologic_monitor_folder.tf_monitor_folder_1.id
  is_disabled = false
  content_type = "Monitor"
  monitor_type = "Logs"
  queries {
      row_id = "A"
      query = "_sourceCategory=event-action info"
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
      recipients = [
        "abc@example.com",
      ]
      subject = "Monitor Alert: {{TriggerType}} on {{Name}}"
      time_zone = "PST"
      message_body = "Triggered {{TriggerType}} Alert on {{Name}}: {{QueryURL}}"
    }
    run_for_trigger_types = ["Critical", "ResolvedCritical"]
  }
  notifications {
    notification {
      connection_type = "PagerDuty"
      connection_id = sumologic_connection.example_pagerduty_connection.id
      payload_override = <<JSON
{
  "service_key": "your_pagerduty_api_integration_key",
  "event_type": "trigger",
  "description": "Alert: Triggered {{TriggerType}} for Monitor {{Name}}",
  "client": "Sumo Logic",
  "client_url": "{{QueryUrl}}"
}
JSON
    }
    run_for_trigger_types = ["Critical", "ResolvedCritical"]
  }
}
```

## Example Monitor Folder

NOTE: Monitor folders are considered a different resource from Library content folders.

```hcl
resource "sumologic_monitor_folder" "tf_monitor_folder_1" {
  name = "test terraform folder"
  description = "a folder for monitors"
}
```

## Argument reference

The following arguments are supported:

- `type` - (Optional) The type of object model. Valid value:
  - `MonitorsLibraryMonitor`
- `name` - (Required) The name of the monitor. The name must be alphanumeric.
- `description` - (Required) The description of the monitor.
- `is_disabled` - (Optional) Whether or not the monitor is disabled. Disabled monitors will not run and will not generate or send notifications.
- `parent_id` - (Optional) The ID of the Monitor Folder that contains this monitor. Defaults to the root folder.
- `content_type` - (Optional) The type of the content object. Valid value:
  - `Monitor`
- `monitor_type` - (Required) The type of monitor. Valid values:
  - `Logs`: A logs query monitor.
  - `Metrics`: A metrics query monitor.
- `queries` - (Required) All queries from the monitor.
- `triggers` - (Required) Defines the conditions of when to send notifications.
- `notifications` - (Optional) The notifications the monitor will send when the respective trigger condition is met.
- `group_notifications` - (Optional) Whether or not to group notifications for individual items that meet the trigger condition. Defaults to true.

Additional data provided in state:

- `id` - (Computed) The ID for this monitor.
- `status` - (Computed) The current status for this monitor. Values are:
  - `Critical`
  - `Warning`
  - `MissingData`
  - `Normal`
  - `Disabled`

### Triggers arguments
The following arguments are supported for `triggers`:
- `detection_method` - Detection method of trigger conditions. Valid values:
  - `StaticCondition`: A condition that triggers based off of a static threshold.
  - `LogsStaticCondition`: A logs condition that triggers based off of a static threshold. Currently `LogsStaticCondition` is available in closed beta (Notify your Sumo Logic representative in order to get the early access).
  - `MetricsStaticCondition`: A metrics condition that triggers based off of a static threshold. Currently `MetricsStaticCondition` is available in closed beta (Notify your Sumo Logic representative in order to get the early access).
  - `LogsOutlierCondition`: A logs condition that triggers based off of a dynamic outlier threshold. Currently `LogsOutlierCondition` is available in closed beta (Notify your Sumo Logic representative in order to get the early access).
  - `MetricsOutlierCondition`: A metrics condition that triggers based off of a dynamic outlier threshold. Currently `MetricsOutlierCondition` is available in closed beta (Notify your Sumo Logic representative in order to get the early access).
  - `LogsMissingDataCondition`: A logs missing data condition that triggers based off of no data available. Currently `LogsMissingDataCondition` is available in closed beta (Notify your Sumo Logic representative in order to get the early access).
  - `MetricsMissingDataCondition`: A metrics missing data condition that triggers based off of no data available. Currently `MetricsMissingDataCondition` is available in closed beta (Notify your Sumo Logic representative in order to get the early access).
- `trigger_type`: The type of trigger condition. Valid values:
  - `Critical`: A critical condition to trigger on.
  - `Warning`: A warning condition to trigger on.
  - `MissingData`: A condition that indicates data is missing.
  - `ResolvedCritical`: A condition to resolve a Critical trigger on.
  - `ResolvedWarning`: A condition to resolve a Warning trigger on.
  - `ResolvedMissingData`: A condition to resolve a MissingData trigger.
- `threshold` - The data value for the condition. This defines the threshold for when to trigger. Threshold value is not applicable for MissingData and ResolvedMissingData triggerTypes and will be ignored if specified.
- `threshold_type` - The comparison type for the threshold evaluation. This defines how you want the data value compared. Valid values:
  - `LessThan`: Less than the configured threshold.
  - `GreaterThan`: Greater than the configured threshold.
  - `LessThanOrEqual`: Less than or equal to the configured threshold.
  - `GreaterThanOrEqual`: Greater than or equal to the configured threshold. ThresholdType value is not applicable for `MissingData` and `ResolvedMissingData` `trigger_type`s and will be ignored if specified.
- `field` - The name of the field that the trigger condition will alert on. The trigger could compare the value of specified field with the threshold. If `field` is not specified, monitor would default to result count instead.
- `time_range` - The relative time range of the monitor.
- `trigger_source` - Determines which time series from queries to use for Metrics MissingData and ResolvedMissingData triggers Valid values:
  - `AllTimeSeries`: Evaluate the condition against all time series. (NOTE: This option is only valid if `monitor_type` is `Metrics`)
  - `AnyTimeSeries`: Evaluate the condition against any time series. (NOTE: This option is only valid if `monitor_type` is `Metrics`)
  - `AllResults`: Evaluate the condition against results from all queries. (NOTE: This option is only valid if `monitor_type` is `Logs`)
- `occurrence_type` - The criteria to evaluate the `threshold` and `threshold_type` in the given time range. Valid values:
  - `AtLeastOnce`: Trigger if the threshold is met at least once. (NOTE: This option is valid only if `monitor_type` is `Metrics`.)
  - `Always`: Trigger if the threshold is met continuously. (NOTE: This option is valid only if `monitor_type` is `Metrics`.)
  - `ResultCount`: Trigger if the threshold is met against the count of results. (NOTE: This option is valid only if `monitor_type` is `Logs`.)
  - `MissingData`: Trigger if the data is missing. (NOTE: This is valid for both `Logs` and `Metrics` `monitor_type`s.)
- `window` - Sets the trailing number of data points to calculate mean and sigma. (NOTE: This option is only valid if `detection_method` is `LogsOutlierCondition`.)
- `consecutive` - Sets the required number of consecutive indicator data points (outliers) to trigger a violation.
- `direction` - Specifies which direction should trigger violations. Valid values:
  - `Both`
  - `Up`
  - `Down`
- `baselineWindow`: The time range used to compute the baseline. (NOTE: This option is only valid if `detection_method` is `MetricsOutlierCondition`.)

The following summarizes the various arguments supported by `triggers` for each `detection_method` (arguments not marked as `Required` are optional):
- `StaticCondition`:
  - `trigger_type` (Required)
  - `time_range` (Required)
  - `threshold`
  - `threshold_type`
  - `field`
  - `occurrence_type` (Required)
  - `trigger_source` (Required)
- `LogsStaticCondition`:
  - `trigger_type` (Required)
  - `time_range` (Required)
  - `threshold` (Required)
  - `threshold_type` (Required)
  - `field`
- `MetricsStaticCondition`:
  - `trigger_type` (Required)
  - `time_range` (Required)
  - `threshold` (Required)
  - `threshold_type` (Required)
  - `occurrence_type` (Required)
- `LogsOutlierCondition`:
  - `trigger_type` (Required)
  - `window`
  - `consecutive`
  - `direction`
  - `threshold`
  - `field`
- `MetricsOutlierCondition`:
  - `trigger_type` (Required)
  - `baseline_window`
  - `direction`
  - `threshold`
- `LogsMissingDataCondition`:
  - `trigger_type` (Required)
  - `time_range` (Required)
- `MetricsMissingDataCondition`:
  - `trigger_type` (Required)
  - `time_range` (Required)
  - `trigger_source` (Required)
## Import

Monitors can be imported using the monitor ID, such as:

```hcl
terraform import sumologic_monitor.test 1234567890
```

[1]: https://help.sumologic.com/?cid=10020
