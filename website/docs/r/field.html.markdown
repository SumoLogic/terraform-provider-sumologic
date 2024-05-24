---
layout: "sumologic"
page_title: "SumoLogic: sumologic_field"
description: |-
  Provides a Sumologic Field
---

# sumologic_field
Provides a [Sumologic Field][1].

## Example Usage
```hcl
resource "sumologic_field" "field" {
  field_name = "string_field_1"
}

resource "sumologic_field" "field" {
  field_name = "int_field_1"
}
```

## Argument reference

The following arguments are supported:

- `field_name` - (Required)  Name of the field.
- `field_id` - (Required) Field identifier.
- `data_type` - (Optional) Field type - Deprecated. Possible values are `String`, `Long`, `Int`, `Double`, and `Boolean`.
- `state` - (Optional) State of the field (either `Enabled` or `Disabled`).

## Import
Fields can be imported using the field id, e.g.:

```hcl
terraform import sumologic_field.field 000000000ABC1234
```

[1]: https://help.sumologic.com/Manage/Fields

