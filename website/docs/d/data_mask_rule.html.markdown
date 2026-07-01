---
layout: "sumologic"
page_title: "SumoLogic: sumologic_data_mask_rule"
description: |-
  Provides a way to retrieve a Sumo Logic data mask rule.
---

# sumologic_data_mask_rule

Provides a way to retrieve a [Sumologic Data Mask Rule][1] by its ID.

## Example Usage

```hcl
data "sumologic_data_mask_rule" "example" {
  id = "000000000ABC1234"
}
```

## Argument Reference

The following arguments are supported:

- `id` - (Required) The ID of the data mask rule to retrieve.

## Attributes Reference

The following attributes are exported:

- `name` - The name of the data mask rule.
- `regex_pattern` - The regular expression pattern used to match sensitive data.
- `mask_string` - The replacement string used when masking.
- `enabled` - Whether the rule is enabled.
- `description` - The description of the rule.

[1]: https://help.sumologic.com/docs/manage/data-masking/
