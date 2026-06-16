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
- `pattern` - The regular expression pattern used to match sensitive data.
- `pii_type` - The type of PII the rule targets.
- `replacement` - The replacement string used when masking.
- `scope` - The scope of the rule (`org`, `child_org`, or `all_orgs`).
- `scope_target_org_ids` - The list of child org IDs the rule applies to.
- `enabled` - Whether the rule is enabled.
- `description` - The description of the rule.
- `is_active` - Whether the rule is currently active.

[1]: https://help.sumologic.com/docs/manage/data-masking/
