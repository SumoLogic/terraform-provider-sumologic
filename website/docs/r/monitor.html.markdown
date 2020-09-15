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
  notifications {
    notification {
        action_type = "EmailAction"
        recipients = ["abc@example.com"]
        subject = "Triggered: tf logs monitor"
        time_zone = "PST"
        message_body = "testing123"
    }
    run_for_trigger_types = ["Critical"]
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
  parent_id = "0000000000000001"
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
      occurrence_type = "Always"
      trigger_source = "AllTimeSeries"
      trigger_type = "Critical"
      detection_method = "StaticCondition"
    }
  triggers {
    threshold_type = "LessThan"
    threshold = 30.0
    time_range = "15m"
    occurrence_type = "Always"
    trigger_source = "AllTimeSeries"
    trigger_type = "ResolvedCritical"
    detection_method = "StaticCondition"
    }
  notifications {
    notification_type = "EmailAction"
    notification {
        action_type = "EmailAction"
        recipients = ["abc@example.com"]
        subject = "Triggered: tf metrics monitor"
        time_zone = "PST"
        message_body = "testing123"
    }
    run_for_trigger_types = ["Critical","ResolvedCritical"]
  }
}
```

## Example Metrics Monitor with Webhook Notification
```hcl
resource "sumologic_connection" "tf_connection_1" {
  type        = "WebhookConnection"
  name        = "test-connection
  description = "My description"
  url         = "https://connection-endpoint.com"
  headers = {
    "X-Header" : "my-header"
  }
  custom_headers = {
    "X-custom" : "my-custom-header"
  }
  default_payload = <<JSON
{
  "client" : "Sumo Logic",
  "eventType" : "{{SearchName}}",
  "description" : "{{SearchDescription}}",
  "search_url" : "{{SearchQueryUrl}}",
  "num_records" : "{{NumRawResults}}",
  "search_results" : "{{AggregateResultsJson}}"
}
JSON
  webhook_type    = "Webhook"
}

resource "sumologic_monitor" "tf_metrics_monitor_1" {
  name = "Terraform Metrics Monitor"
  description = "tf metrics monitor"
  type = "MonitorsLibraryMonitor"
  is_disabled = false
  parent_id = "0000000000000001"
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
      occurrence_type = "Always"
      trigger_source = "AllTimeSeries"
      trigger_type = "Critical"
      detection_method = "StaticCondition"
    }
  triggers {
    threshold_type = "LessThan"
    threshold = 30.0
    time_range = "15m"
    occurrence_type = "Always"
    trigger_source = "AllTimeSeries"
    trigger_type = "ResolvedCritical"
    detection_method = "StaticCondition"
    }
  notifications {
    notification {
      action_type = "NamedConnectionAction"
      connection_id    = sumologic_connection.tf_connection_1.id
      time_zone        = "UTC"
      payload_override = <<JSON
{
  "client" : "Sumo Logic",
  "eventType" : "Different",
  "description" : "{{SearchDescription}}"
}
       JSON
    }
    run_for_trigger_types = ["Critical"]
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
- `notifications` - (Optional) The notifications the monitor will send when the respective trigger condition is met. See details below.
- `group_notifications` - (Optional) Whether or not to group notifications for individual items that meet the trigger condition. Defaults to true.

**notifications** is a child block with the following arguments:
- `run_for_trigger_types` - (Required) A notification will be sent for each trigger type defined in this list.

- `notification` - (Required) Defines the notification. See details below.

**notification** is a child block with the following arguments:
- `action_type` - (Required) The type of notification. Must be either `NamedConnectionAction` for a webhook or `EmailAction` for an email notification.
- `connection_id` - (Optional) The ID of the connection to be used to send the notification. This parameter is required if the `action_type` is `NamedConnectionAction`.
- `time_zone` - (Optional) Set the timzone for the notification.
- `payload_override` - (Optional) This field can be used to overwrite connection's default template. Inside a template several variables are supported, you can find a full list of variables [here][2]. Defaults to `null` means that the default template (provided when the connection was created) will be used.
- `recipients` - (Optional) List of emails to send the notification to. This parameter is required if the `action_type` is `EmailAction`.
- `subject` (Optional) Sets the subject of the email that is sent when the alert is triggered. This parameter is required if the `action_type` is `EmailAction`.
- `message_body` (Optional) Sets the body of the email that is sent when the alert is triggered. This parameter is required if the `action_type` is `EmailAction`.

Additional data provided in state:

- `id` - (Computed) The ID for this monitor.
- `status` - (Computed) The current status for this monitor. Values are:
  - `Critical`
  - `Warning`
  - `MissingData`
  - `Normal`
  - `Disabled`

## Import

Monitors can be imported using the monitor ID, such as:

```hcl
terraform import sumologic_monitor.test 1234567890
```

[1]: https://help.sumologic.com/Beta/Monitors
[2]: https://help.sumologic.com/Manage/Connections-and-Integrations/Webhook-Connections/Set_Up_Webhook_Connections#webhook-payload-variables
