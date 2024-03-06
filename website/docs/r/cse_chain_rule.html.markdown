---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_chain_rule"
description: |-
  Provides a CSE Chain Rule
---

# sumologic_cse_chain_rule
Provides a Sumo Logic CSE [Chain Rule](https://help.sumologic.com/Cloud_SIEM_Enterprise/CSE_Rules/07_Write_a_Chain_Rule).

## Example Usage
```hcl
resource "sumologic_cse_chain_rule" "chain_rule" {
  description = "Signal description"
  enabled = true
  entity_selectors {
    entity_type = "_username"
    expression = "user_username"
  }
  expressions_and_limits {
    expression = "success = false"
    limit = 5
  }
  expressions_and_limits {
    expression = "success = true"
    limit = 1
  }
  group_by_fields = []
  is_prototype = false
  ordered = true
  name = "Chain Rule Example"
  severity = 5
  summary_expression = "Signal summary"
  tags = ["_mitreAttackTactic:TA0009"]
  window_size = "T30M"
}
```

## Argument reference

The following arguments are supported:

- `description` - (Required) The description of the generated Signals
- `enabled` - (Required) Whether the rule should generate Signals
- `entity_selectors` - (Required) The entities to generate Signals on
  + `entityType` - (Required) The type of the entity to generate the Signal on.
  + `expression` - (Required) The expression or field name to generate the Signal on.
- `expressions_and_limits` - (Required) The list of expressions and associated limits to make up the conditions of the chain rule
  + `expression` - (Required) The expression for which records to match on
  + `limit` - (Required) How many times this expression must match for the Signal to fire
- `group_by_fields` - (Optional) A list of fields to group records by
- `is_prototype` - (Optional) Whether the generated Signals should be prototype Signals
- `ordered` - (Optional; defaults to false) Whether the records matching the expressions must be in the same chronological order as the expressions are listed in the rule
- `name` - (Required) The name of the Rule and the generated SignalS
- `severity` - (Required) The severity of the generated Signals
- `summary_expression` - (Optional) The summary of the generated Signals
- `tags` - (Required) The tags of the generated Signals
- `window_size` - (Required) How long of a window to aggregate records for. Current acceptable values are T05M, T10M, T30M, T60M, T24H, T12H, T05D or CUSTOM
  + `window_size_millis` - (Optional) Used only when `window_size` is set to CUSTOM. Window size in milliseconds ranging from 1 minute to 5 days ("60000" to "432000000").

The following attributes are exported:

- `id` - The internal ID of the chain rule.

## Import

Chain Rules can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_chain_rule.chain_rule id
```