---
layout: "sumologic"
page_title: "SumoLogic: sumologic_http_source"
description: |-
  Provides a way to retrieve Sumo Logic HTTP Source details (id, names, etc) for an HTTP Source managed by another terraform stack.
---

# sumologic_http_source

Provides a way to retrieve Sumo Logic HTTP Source details (id, names, etc) for an HTTP Source managed by another terraform stack.


## Example Usage
```hcl
data "sumologic_http_source" "this" {
  collector_id = 121212
  name = "source_name"
}
```

A HTTP Source can be looked up by using a combination of `collector_id` & `name`.
If either `id` or `name` are not present, the data source block fails with a panic (at this point).

## Attributes reference

The following attributes are exported:

- `id` - The internal ID of the collector. This can be used to attach sources to the collector.
- `name` - The name of the collector.
- `description` - The description of the collector.
- `category` - The default source category for any source attached to this collector.
- `timezone` - The time zone to use for this collector. The value follows the [tzdata][2] naming convention.
- `multiline` - Multiline processing enabled or not.
- `url` - The HTTP endpoint to use for sending data to this source.
- `token` - The token to use for sending data to this source.
- `base_url` - The base URL for the HTTP source endpoint.


