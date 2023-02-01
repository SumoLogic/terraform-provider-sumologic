---
layout: "sumologic"
page_title: "SumoLogic: sumologic_collector"
description: |-
  Provides a Sumologic (Hosted) Collector.
---

# sumologic_collector
Provides a [Sumologic (Hosted) Collector][1].

## Example Usage
```hcl
resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Just testing this"
  fields  = {
    environment = "production"
  }
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the collector. This is required, and has to be unique. Changing this will force recreation the collector.
  * `description` - (Optional) The description of the collector.
  * `category` - (Optional) The default source category for any source attached to this collector. Can be overridden in the configuration of said sources.
  * `timezone` - (Optional) The time zone to use for this collector. The value follows the [tzdata][2] naming convention.
  * `fields` - (Optional) Map containing [key/value pairs][3].

### See also
   * [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties)

## Attributes Reference
The following attributes are exported:

  * `id` - The internal ID of the collector. This can be used to attach sources to the collector.

## Import
Collectors can be imported using the collector id, e.g.:

```hcl
terraform import sumologic_collector.test 1234567890
```

Collectors can also be imported using the collector name, which is unique per Sumo Logic account, e.g.:

```hcl
terraform import sumologic_collector.test my_test_collector
```

[1]: https://help.sumologic.com/docs/send-data/hosted-collectors/
[2]: https://en.wikipedia.org/wiki/Tz_database
[3]: https://help.sumologic.com/Manage/Fields
[4]: https://www.terraform.io/docs/configuration/resources.html#prevent_destroy
