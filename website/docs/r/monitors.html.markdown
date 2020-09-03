---
layout: 'sumologic'
page_title: 'SumoLogic: sumologic_monitor'
description: |-
  Provides the ability to create, read, delete, update monitors.
---

# sumologic_monitor

Provides the ability to create, read, delete, update [Monitors (New)][1].

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

## Argument reference

The following arguments are supported:

- `type` - (Required) Type of monitor. Valid values are `Logs` or `Metrics`
- `name` - (Required) Name of monitor. Name should be a valid alphanumeric value.
- `description` - (Optional) Description of the monitor.

Additional data provided in state

- `id` - (Computed) The Id for this monitor.
- `status` - (Computed) The current status for this monitor. Values are `Critical`, `Warning`, `MissingData`, `Normal`, or `Disabled`

## Import

Monitors can be imported using the monitor id, e.g.:

```hcl
terraform import sumologic_monitor.test 1234567890
```

[1]: https://help.sumologic.com/Beta/Monitors
