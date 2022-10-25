---
layout: "sumologic"
page_title: "SumoLogic: sumologic_installed_collector"
description: |-
  Provides a Sumologic (Installed) Collector.
---

# sumologic_collector
Provides a [Sumologic (Installed) Collector][1].
**NOTE**: This resource can only be imported and managed by terraform. Creation of this resource via terrform is not supported.

## Example Usage
```hcl
resource "sumologic_installed_collector" "installed_collector" {
  name        = "test-mac"
  category = "macos/test"
  ephemeral = true
  fields =  {
      key = "value"
    }
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the collector. This is required, and has to be unique. Changing this will force recreation the collector.
  * `ephemeral` - (Required) When true, the collector will be deleted after 12 hours of inactivity. For more information, see [Setting a Collector as Ephemeral][5].
  * `description` - (Optional) The description of the collector.
  * `category` - (Optional) The default source category for any source attached to this collector. Can be overridden in the configuration of said sources.
  * `timezone` - (Optional) The time zone to use for this collector. The value follows the [tzdata][2] naming convention.
  * `fields` - (Optional) Map containing [key/value pairs][3].
  * `cut_off_timestamp` - (Optional) Only collect data from files with a modified date more recent than this timestamp, specified as milliseconds since epoch.
  * `host_name` - (Optional) Host name of the Collector. The hostname can be a maximum of 128 characters.
  * `source_sync_mode` - (Optional) For installed Collectors, whether the Collector is using local source configuration management (using a JSON file), or cloud management (using the UI)
  * `target_cpu` - When CPU utilization exceeds this threshold, the Collector will slow down its rate of ingestion to lower its CPU utilization. Currently only Local and Remote File Sources are supported.

### See also
  * [Common Collector Properties](https://help.sumologic.com/docs/api/collectors/#response-fields)
  * [Common Source Properties](https://github.com/terraform-providers/terraform-provider-sumologic/tree/master/website#common-source-properties)

## Attributes Reference
The following attributes are exported:

  * `id` - The internal ID of the collector. This can be used to attach sources to the collector.

## Import
Collectors can be imported using the collector id, e.g.:

```hcl
terraform import sumologic_installed_collector.test 1234567890
```

Collectors can also be imported using the collector name, which is unique per Sumo Logic account, e.g.:

```hcl
terraform import sumologic_installed_collector.test my_test_collector
```

[1]: https://help.sumologic.com/03Send-Data/Installed-Collectors/01About-Installed-Collectors
[2]: https://en.wikipedia.org/wiki/Tz_database
[3]: https://help.sumologic.com/Manage/Fields
[4]: https://www.terraform.io/docs/configuration/resources.html#prevent_destroy
[5]:https://help.sumologic.com/03Send-Data/Installed-Collectors/05Reference-Information-for-Collector-Installation/11Set-a-Collector-as-Ephemeral
