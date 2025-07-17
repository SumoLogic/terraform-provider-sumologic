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
  field_name = "mystring"
}

resource "sumologic_field" "field" {
  field_name = "my_other_string"
  state      = "Disabled"
}
```

## Argument reference

The following arguments are supported:

- `field_name` - (Required)  Name of the field.
- `state` - (Optional) State of the field. Possible values are `Enabled` or `Disabled` (default: `Enabled`).

## Import
Fields can be imported using the field id, e.g.:

```hcl
terraform import sumologic_field.field 000000000ABC1234
```

[1]: https://help.sumologic.com/Manage/Fields

