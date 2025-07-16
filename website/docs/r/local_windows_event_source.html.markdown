---
layout: "sumologic"
page_title: "SumoLogic: sumologic_local_windows_event_log_source"
description: |-
  Provides a Sumologic Local Windows Event Log Source.
---

# sumologic_local_windows_event_source
Provides a [Sumologic Local Windows Event Log Source][1].

Note that installed collector sources must be treated as a special case as the user must have a pipeline to install them outside of terraform as it is not possible to install a local collector via the API, that must be done locally on the instance. Make sure the collector is in cloud managed not local json file mode to allow for API based configuration.

Use the installed collector data source to map to installed collector instances by name or id.

## Example Usage

Example: 1 This will configure JSON format with "concise" setting and pick up System and Application logs with /os/windows/events as the source category.

```hcl
data "sumologic_collector" "installed_collector" {
  name = "terraform_source_testing"
}

resource "sumologic_local_windows_event_log_source" "local" {
  name             = "windows_logs"
  description      = "windows system and application logs in json format"
  category         = "/os/windows/events"
  collector_id     = "${data.sumologic_collector.installed_collector.id}"
  log_names  = ["System","Application"]
  event_format = 1 // 0 = XML, 1 = JSON
}
```

Example 2: Using custom logs and a deny list
```hcl
resource "sumologic_local_windows_event_log_source" "local" {
  name             = "windows_logs"
  description      = "windows logs in json format"
  category         = "/os/windows/events"
  collector_id     = "${data.sumologic_collector.installed_collector.id}"
  log_names  = ["System","Application","Microsoft-Windows-PowerShell/Operational", "Microsoft-Windows-TaskScheduler/Operational"]
  deny_list = "9999,7890"
  event_format = 1 // 0 = XML, 1 = JSON
}
```



## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the local file source. This is required, and has to be unique. Changing this will force recreation the source.
  * `description` - (Optional) The description of the source.
  * `log_names` - List of Windows log types to collect (e.g., Security, Application, System)
  * `render_messages` - When using legacy format, indicates if full event messages are collected
  * `event_format` - 0 for legacy format (XML), 1 for JSON format. Default 0.
  * `event_message` - 0 for complete message, 1 for message title, 2 for metadata only. Required if event_format is 0
  * `deny_list` - Comma-separated list of event IDs to deny
  * `category` - (Optional) The default source category for the source.
  * `fields` - (Optional) Map containing [key/value pairs][2].
  * `denylist` - (Optional) Comma-separated list of event IDs to deny. This is used to exclude specific events from being collected.

### See also
  * [Common Source Properties](https://github.com/terraform-providers/terraform-provider-sumologic/tree/master/website#common-source-properties)

## Attributes Reference
The following attributes are exported:

  * `id` - The internal ID of the local file source.

## Import
Local Windows Event Log Sources can be imported using the collector and source IDs, e.g.:

```hcl
terraform import sumologic_local_windows_event_source.test 123/456

```hcl
terraform import sumologic_local_windows_event_source.test my-test-collector/my-test-source
```

[1]: https://help.sumologic.com/docs/send-data/installed-collectors/sources/local-windows-event-log-source/
[2]: https://help.sumologic.com/Manage/Fields
[3]: https://help.sumologic.com/docs/send-data/use-json-configure-sources/json-parameters-installed-sources/#local-windows-event-logsource