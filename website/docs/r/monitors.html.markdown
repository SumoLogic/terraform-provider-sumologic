---
layout: 'sumologic'
page_title: 'SumoLogic: sumologic_monitor'
description: |-
  Provides the ability to create, read, delete, update monitors.
---

# sumologic_monitor

Provides the ability to create, read, delete, update monitors.

## Example Usage

```hcl
resource "sumologic_monitor" "tf_logs_monitor_1" {
  name = "Terraform Monitor"
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
    notification_type = "EmailAction"
    notification {
        action_type = "EmailAction"
        recipients = ["rohit@sumologic.com"]
        subject = "Triggered: tf logs monitor"
        time_zone = "PST"
        message_body = "testing123"
    }
    run_for_trigger_types = ["Critical"]
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
