---
layout: "sumologic"
page_title: "SumoLogic: sumologic_collector"
description: |-
  Provides a way to retrieve Sumo Logic collector details (id, names, etc) for a collector managed by another terraform stack.
---

# sumologic_collector

Provides a way to retrieve Sumo Logic collector details (id, names, etc) for a collector
managed by another terraform stack.


## Example Usage
```hcl
data "sumologic_collector" "this" {
  name = "MyCollector"
}
```

```hcl
data "sumologic_collector" "that" {
  id = "1234567890"
}
```

A collector can be looked up by either `id` or `name`. One of those attributes needs to be specified.

If both `id` and `name` have been specified, `id` takes precedence.

## Attributes reference

The following attributes are exported:

- `id` - The internal ID of the collector. This can be used to attach sources to the collector.
- `name` - The name of the collector.
- `description` - The description of the collector.
- `category` - The default source category for any source attached to this collector.
- `timezone` - The time zone to use for this collector. The value follows the [tzdata][2] naming convention.


[Back to Index][0]

[0]: ../README.md

