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
  * `lookup_by_name` - DEPRECATED - Use the data source or import the collector - (Optional) Configures an existent collector using the same 'name' or creates a new one if non existent. Defaults to false.
  * `destroy` - DEPRECATED - Use lifecycle management attribute [`prevent_destroy`][4] - (Optional) Whether or not to delete the collector in Sumo when it is removed from Terraform.  Defaults to true.

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

[Back to Index][0]

[0]: ../README.md
[1]: https://help.sumologic.com/Send_Data/Hosted_Collectors
[2]: https://en.wikipedia.org/wiki/Tz_database
[3]: https://help.sumologic.com/Manage/Fields
[4]: https://www.terraform.io/docs/configuration/resources.html#prevent_destroy
