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
- `trigger_conditions` - (Required if not using `triggers`) Defines the conditions of when to send notifications. NOTE: `trigger_conditions` supplants the `triggers` argument. 
- `triggers` - (Deprecated) Defines the conditions of when to send notifications.
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

## The `trigger_conditions` block
A `trigger_conditions` block configures conditions for sending notifications.
### Example
```hcl
trigger_conditions {
  static_condition {
    field           = "_count"
    time_range      = "15m"
    trigger_source  = "AllResults"
    occurrence_type = "ResultCount"

    critical {
      alert {
        threshold      = 100
        threshold_type = "GreaterThan"
      }

      resolution {
        threshold      = 90
        threshold_type = "LessThanOrEqual"
      }
    }

    warning {
      alert {
        threshold      = 80
        threshold_type = "GreaterThan"
      }

      resolution {
        threshold      = 75
        threshold_type = "LessThanOrEqual"
      }
    }
  }
   
  logs_missing_data_condition {
    time_range = "30m"
  }
}
```
### Arguments
Here is a summary of the various condition types that are supported, and the arguments each of them takes:
- `static_condition`:
  - `field`
  - `time_range` (Required)
  - `trigger_source` (Required)
  - `occurrence_type` (Required)
  - `critical`
    - `alert` (Required)
       - `threshold`
       - `threshold_type`
    - `resolution` (Required)
      - `threshold`
      - `threshold_type`
  - `warning`
    - `alert` (Required)
      - `threshold`
      - `threshold_type`
    - `resolution` (Required)
      - `threshold`
      - `threshold_type`
- `logs_static_condition`:
  - `field`
  - `time_range` (Required)
  - `critical` (See `static_condition.critical` for schema)
  - `warning`  (See `static_condition.warning` for schema)
- `metrics_static_condition`:
  - `time_range` (Required)
  - `occurrence_type` (Required)
  - `critical` (See `static_condition.critical` for schema)
  - `warning`  (See `static_condition.warning` for schema)
- `logs_outlier_condition`:
  - `field`
  - `window`
  - `consecutive`
  - `direction`
  - `critical`
     - `threshold`
  - `warning`
     - `threshold`
- `metrics_outlier_condition`:
  - `baseline_window`
  - `direction`
  - `threshold`
  - `critical`
    - `threshold`
  - `warning`
    - `threshold`
- `logs_missing_data_condition`:
  - `time_range` (Required)
- `metrics_missing_data_condition`:
  - `time_range` (Required)
  - `trigger_source` (Required)

A `trigger_conditions` block can contain at most 1 data condition:
 - `static_condition`
 - `logs_static_condition`
 - `metrics_static_condition`
 - `logs_outlier_condition`
 - `metrics_outlier_condition`
 
and at most 1 missing-data condition:
  - `logs_missing_data_condition`
  - `metrics_missing_data_condition`

## The `triggers` block
The `triggers` block is deprecated. Please use `trigger_conditions` to specify notification conditions.

## Import

Monitors can be imported using the monitor ID, such as:

```hcl
terraform import sumologic_monitor.test 1234567890
```

[1]: https://help.sumologic.com/?cid=10020
