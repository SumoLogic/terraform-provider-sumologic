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
resource "sumologic_monitor" "monitor" {
  type        = "Monitor"
  name        = "test-monitor
  description = "My description"
  content_type = "MonitorsLibraryMonitor"
  monitor_type = "Logs"
  queries = []
  triggers = []
  notifications = []
  is_disabled = false
  group_notifications = true
}
```

## Argument reference

The following arguments are supported:

- `type` - (Required) Type of monitor. Valid values are `Logs` or `Metrics`
- `name` - (Required) Name of monitor. Name should be a valid alphanumeric value.
- `description` - (Optional) Description of the monitor.

Additional data provided in state

- `id` - (Computed) The Id for this monitor.

## Import

Monitors can be imported using the monitor id, e.g.:

```hcl
terraform import sumologic_monitor.test 1234567890
```
