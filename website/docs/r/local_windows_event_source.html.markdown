---
layout: "sumologic"
page_title: "SumoLogic: sumologic_local_windows_event_log_source"
description: |-
  Provides a Sumologic Local Windows Event Log Source.
---

# sumologic_local_windows_event_log_source
Provides a [Sumologic Local Windows Event Log Source][1].

Note: installed collector sources are a special case for terraform based management. It is not possible to install a local collector via the API, only to add sources to it. The installed collector must be done installed locally on the instance and be set to cloud managed mode to allow for API based configuration.

## Example Usage

Each example declares the installed collector it attaches the source to, so it can be copied on its own. The installed collector must already exist and be registered in Sumo Logic; the `sumologic_collector` data source looks it up by name.

### System and Application logs in JSON format
Configures JSON format and picks up the System and Application logs, with `/os/windows/events` as the source category.

```hcl
data "sumologic_collector" "installed_collector" {
  name = "terraform_source_testing"
}

resource "sumologic_local_windows_event_log_source" "local" {
  name         = "windows_logs"
  description  = "windows system and application logs in json format"
  category     = "/os/windows/events"
  collector_id = data.sumologic_collector.installed_collector.id
  log_names    = ["System", "Application"]
  event_format = 1
}
```

### Custom logs and a deny list
Collects two additional channels and drops events with the listed event IDs.

```hcl
data "sumologic_collector" "installed_collector" {
  name = "terraform_source_testing"
}

resource "sumologic_local_windows_event_log_source" "local_custom" {
  name         = "windows_logs_custom"
  description  = "windows logs in json format"
  category     = "/os/windows/events"
  collector_id = data.sumologic_collector.installed_collector.id
  log_names = [
    "System",
    "Application",
    "Microsoft-Windows-PowerShell/Operational",
    "Microsoft-Windows-TaskScheduler/Operational",
  ]
  deny_list    = "9999,7890"
  event_format = 1
}
```

### Legacy (XML) format with an allow list
Collects only the listed event IDs from the Security log, in the legacy format with the complete event message.

```hcl
data "sumologic_collector" "installed_collector" {
  name = "terraform_source_testing"
}

resource "sumologic_local_windows_event_log_source" "local_security" {
  name            = "windows_security"
  description     = "windows security events in legacy format"
  category        = "/os/windows/security"
  collector_id    = data.sumologic_collector.installed_collector.id
  log_names       = ["Security"]
  event_format    = 0
  event_message   = 0
  render_messages = true
  allow_list      = "4624,4625"
}
```

## Argument Reference

In addition to the [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties), the following arguments are supported:

* `log_names` - (Required) List of Windows event log channels to collect, for example `System`, `Application`, `Security`, or a custom channel such as `Microsoft-Windows-PowerShell/Operational`.
* `event_format` - (Optional) Output format of collected events. `0` for legacy format (XML), `1` for JSON format. Allowed values: `0`, `1`. Defaults to `0`.
* `event_message` - (Optional) How much of the event message to collect (with JSON format the UI calls this Event Collection Level: Complete, Concise, Metadata Only). `0` for the complete message, `1` for the message title only, `2` for metadata only. Allowed values: `0`, `1`, `2`. The schema description states this is required when `event_format` is `0`; the provider does not enforce it, so set it explicitly when using the legacy format. When omitted, no value is sent.
* `render_messages` - (Optional) When using the legacy format (`event_format = 0`), whether full event messages are collected. Defaults to `true`.
* `allow_list` - (Optional) Comma-separated list of event IDs to collect (the help page recommends using only one of `allow_list` and `deny_list` at a time, and notes event ID filters need a recent installed collector version), for example `"4624,4625"`. When omitted, no allow list is sent.
* `deny_list` - (Optional) Comma-separated list of event IDs to exclude, for example `"9999,7890"`. When omitted, no deny list is sent.

### See also

* [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties)
* [Windows Event Source properties](https://help.sumologic.com/docs/send-data/use-json-configure-sources/json-parameters-installed-sources/#local-windows-event-logsource)

## Attributes Reference
The following attributes are exported:

  * `id` - The internal ID of the source.

## Import
Local Windows Event Log Sources can be imported using the collector and source IDs (`collector/source`), e.g.:

```hcl
terraform import sumologic_local_windows_event_log_source.test 123/456
```

Local Windows Event Log Sources can also be imported using the collector name and source name (`collectorName/sourceName`), e.g.:

```hcl
terraform import sumologic_local_windows_event_log_source.test my-test-collector/my-test-source
```

[1]: https://www.sumologic.com/help/docs/send-data/installed-collectors/sources/local-windows-event-log-source/
