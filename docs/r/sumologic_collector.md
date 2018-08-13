# sumologic_collector
Provides a [Sumologic (Hosted) Collector][1].

## Example Usage
```hcl
resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Just testing this"
}
```

## Argument reference
The following arguments are supported:
- `name` - (Required) The name of the collector. This is required, and has to be unique. Changing this will force recreation the collector.
- `description` - (Optional) The description of the collector.
- `category` - (Optional) The default source category for any source attached to this collector. Can be overridden in the configuration of said sources.
- `timezone` - (Optional) The time zone to use for this collector. The value follows the [tzdata][2] naming convention.
- `lookup_by_name` - (Optional) Configures an existent collector using the same 'name' or creates a new one if non existent. Defaults to false.
- `destroy` - (Optional) Whether or not to delete the collector in Sumo when it is removed from Terraform.  Defaults to true.

## Attributes reference
The following attributes are exported:
- `id` - The internal ID of the collector. This can be used to attach sources to the collector.

[Back to Index][0]

[0]: ../README.md
[1]: https://help.sumologic.com/Send_Data/Hosted_Collectors
[2]: https://en.wikipedia.org/wiki/Tz_database
