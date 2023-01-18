---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_first_seen_rule"
description: |-
  Provides a CSE First Seen Rule
---

# sumologic_cse_first_seen_rule
Provides a Sumo Logic CSE [First Seen Rule](https://help.sumologic.com/docs/cse/rules/write-first-seen-rule/).

## Example Usage
```hcl
resource "sumologic_cse_first_seen_rule" "first_seen_rule" {
  baseline_type          = "PER_ENTITY"
  baseline_window_size   = "35000"
  description_expression = "First User Login - {{ user_username }}"
  enabled                = true

  entity_selectors {
        entity_type = "_username"
        expression = "user_username"
  }

  entity_selectors {
        entity_type = "_hostname"
        expression = "dstDevice_hostname"
  }

  filter_expression     = "objectType=\"Network\""
  group_by_fields       = ["user_username"]
  is_prototype          = false
  name                  = "First User Login"
  name_expression       = "First User Login - {{ user_username }}"
  retention_window_size = "86400000"
  severity              = 1
  value_fields          = ["dstDevice_hostname"]
}
```

## Argument reference

The following arguments are supported:

- `baseline_type` - (Required) The baseline type. Current acceptable values are GLOBAL or PER_ENTITY
- `baseline_window_size` - (Optional) The baseline window size in milliseconds 
- `category` - (Optional) The category
- `description_expression` - (Required) The description of the generated Signals
- `enabled` - (Required) Whether the rule should generate Signals
- `entity_selectors` - (Required) The entities to generate Signals on
  + `entityType` - (Required) The type of the entity to generate the Signal on
  + `expression` - (Required) The expression or field name to generate the Signal on
- `filter_expression` - (Required) The expression for which records to match on
- `group_by_fields` - (Optional) A list of fields to group records by
- `is_prototype` - (Optional) Whether the generated Signals should be prototype Signals
- `name` - (Required) The name of the Rule 
- `name_expression` - (Required) The name of the generated Signals
- `retention_window_size` - (Optional) The retention window size in milliseconds
- `severity` - (Required) The severity of the generated Signals
- `summary_expression` - (Optional) The summary of the generated Signals
- `tags` - (Optional) The tags of the generated Signals
- `value_fields` - (Required) The value fields

The following attributes are exported:

- `id` - The internal ID of the first seen rule.

## Import

First Seen Rules can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_first_seen_rule.first_seen_rule id
```