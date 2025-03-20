---
layout: "sumologic"
page_title: "SumoLogic: sumologic_ot_collector"
description: |-
  Provides a Sumologic OT Collector
---

# sumologic_OT_Collector
Provides a [Sumologic OT Collector][1]
**NOTE**: This resource can only be imported and deleted by terraform. Creation/Updation of this resource via terrform is not supported.
Tag Edit functionality will be supported in future.

## Example Usage
```hcl
resource "sumologic_ot_collector" "example_ot_collector" {
    description = "Testing OT collector using terraform"
    time_zone = "UTC"
    category = "apache"
    is_remotely_managed = "true"
    ephemeral = "false"
    name = "test OT Collector"
}
```
## Argument reference

The following arguments are supported:

- `description` - (Optional) Description of the OT Collector.
- `time_zone` - (Optional) The time zone to use for this collector. The value follows the [tzdata][4] naming convention.
- `category` - (Optional) The default source category for any source attached to this collector. Can be overridden in the configuration of said sources.
- `is_remotely_managed` - (Optional) Management Status of the OT Collector based on if it is remotely or locally managed.
- `ephemeral` - (Optional) When true, the collector will be deleted after 12 hours of inactivity. For more information, see [Setting a Collector as Ephemeral][2].
- `name` - (Required) Name of the OT Collector.
- `tags` - (Optional) Map containing [key/value pairs][3].
 
The following attributes are exported:

- `id` - The internal ID of the OT collector

## Import
OT Collectors can be imported using the collector id, e.g.:

```hcl
terraform import sumologic_ot_collector.test 00005AF3107A4007
```

[1]: https://help.sumologic.com/docs/send-data/opentelemetry-collector
[2]: https://help.sumologic.com/03Send-Data/Installed-Collectors/05Reference-Information-for-Collector-Installation/11Set-a-Collector-as-Ephemeral
[3]: https://help.sumologic.com/Manage/Fields
[4]: https://en.wikipedia.org/wiki/Tz_database
