---
layout: "sumologic"
page_title: "SumoLogic: sumologic_data_mask_rules"
description: |-
  Provides a way to retrieve all Sumo Logic data mask rules.
---

# sumologic_data_mask_rules

Provides a way to retrieve all [Sumologic Data Mask Rules][1].

## Example Usage

```hcl
data "sumologic_data_mask_rules" "all" {}
```

## Attributes Reference

The following attributes are exported:

- `rules` - A list of data mask rules. Each rule has the following attributes:
  - `id` - The unique identifier of the rule.
  - `name` - The name of the rule.
  - `pattern` - The regular expression pattern used to match sensitive data.
  - `pii_type` - The type of PII the rule targets.
  - `replacement` - The replacement string used when masking.
  - `scope` - The scope of the rule (`org`, `child_org`, or `all_orgs`).
  - `scope_target_org_ids` - The list of child org IDs the rule applies to.
  - `enabled` - Whether the rule is enabled.
  - `description` - The description of the rule.
  - `is_active` - Whether the rule is currently active.

[1]: https://help.sumologic.com/docs/manage/data-masking/
