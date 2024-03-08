---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_threshold_rule"
description: |-
  Provides a CSE Threshold Rule
---

# sumologic_cse_threshold_rule
Provides a Sumo Logic CSE [Threshold Rule](https://help.sumologic.com/Cloud_SIEM_Enterprise/CSE_Rules/05_Write_a_Threshold_Rule).

## Example Usage
```hcl
resource "sumologic_cse_threshold_rule" "threshold_rule" {
  count_distinct = true
  count_field = "dstDevice_hostname"
  description = "Signal description"
  enabled = true
  entity_selectors {
    entity_type = "_ip"
    expression = "srcDevice_ip"
  }
  expression = "objectType = \"Network\""
  group_by_fields = ["dstDevice_hostname"]
  is_prototype = false
  limit = 1000
  name = "Threshold Rule Example"
  severity = 5
  summary_expression = "Signal summary"
  tags = ["_mitreAttackTactic:TA0009"]
  window_size = "T30M"
}
```

## Argument reference

The following arguments are supported:

- `count_distinct` - (Optional; defaults to false) Whether to count distinct values of a field, as opposed to just counting the number of records
- `count_field` - (Optional) The field to count if `count_distinct` is set to true
- `description` - (Required) The description of the generated Signals
- `enabled` - (Required) Whether the rule should generate Signals
- `entity_selectors` - (Required) The entities to generate Signals on
  + `entityType` - (Required) The type of the entity to generate the Signal on.
  + `expression` - (Required) The expression or field name to generate the Signal on.
- `expression` - (Required) The expression for which records to match on
- `group_by_fields` - (Optional) A list of fields to group records by
- `is_prototype` - (Optional) Whether the generated Signals should be prototype Signals
- `limit` - (Required) A Signal will be fired when this many records/distinct field values are matched
- `name` - (Required) The name of the Rule and the generated Signals
- `severity` - (Required) The severity of the generated Signals
- `summary_expression` - (Optional) The summary of the generated Signals
- `tags` - (Required) The tags of the generated Signals
- `window_size` - (Required) How long of a window to aggregate records for. Current acceptable values are T05M, T10M, T30M, T60M, T24H, T12H, T05D or CUSTOM
  + `window_size_millis` - (Optional) Used only when `window_size` is set to CUSTOM. Window size in milliseconds ranging from 1 minute to 5 days ("60000" to "432000000").

The following attributes are exported:

- `id` - The internal ID of the threshold rule.

## Import

Threshold Rules can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_threshold_rule.threshold_rule id
```