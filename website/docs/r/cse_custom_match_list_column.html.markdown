---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_custom_match_list_column"
description: |-
  Provides a Sumologic CSE Custom Match List Column
---

# custom_match_list_column
Provides a Sumologic CSE Custom Match List Column.

## Example Usage
```hcl
resource "sumologic_cse_custom_match_list_column" "custom_match_list_column" {
  name = "Custom Match List Column name"
  fields = ["srcDevice_ip"]
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) Custom Match List Column name.
- `fields` - (Required) Custom Match List Column fields. 

The following attributes are exported:

- `id` - The internal ID of the Custom Match List Column.

## Import

Custom Match List Column can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_custom_match_list_column.custom_match_list_column id
```

