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
  - `regex_pattern` - The regular expression pattern used to match sensitive data.
  - `mask_string` - The replacement string used when masking.
  - `enabled` - Whether the rule is enabled.
  - `description` - The description of the rule.

[1]: https://help.sumologic.com/docs/manage/data-masking/
