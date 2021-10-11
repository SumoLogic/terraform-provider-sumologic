---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_aggregation_rule"
description: |-
  Provides a CSE Aggregation Rule
---

# sumologic_cse_aggregation_rule
Provides a Sumo Logic CSE [Aggregation Rule](https://help.sumologic.com/Cloud_SIEM_Enterprise/CSE_Rules/09_Write_an_Aggregation_Rule).

## Example Usage
```hcl
resource "sumologic_cse_aggregation_rule" "aggregation_rule" {
  aggregation_functions {
    name = "distinct_eventid_count"
    function = "count_distinct"
    arguments = ["metadata_deviceEventId"]
  }
  description_expression = "Signal description"
  enabled = true
  entity_selectors {
    entity_type = "_ip"
    expression = "srcDevice_ip"
  }
  group_by_entity = true
  group_by_fields = ["dstDevice_hostname"]
  match_expression = "objectType = \"Network\""
  is_prototype = false
  name = "Aggregation Rule Example"
  name_expression = "Signal name"
  severity_mapping {
    type = "constant"
    default = 5
  }
  summary_expression = "Signal summary"
  tags = ["_mitreAttackTactic:TA0009"]
  trigger_expression = "distinct_eventid_count > 5"
  window_size = "T30M"
}
```

## Argument reference

The following arguments are supported:

- `aggregation_functions` - (Required) One or more named aggregation functions
  + `name` - (Required) The name to use to reference the result in the trigger_expression
  + `function` - (Required) The function to aggregate with
  + `arguments` - (Required) One or more expressions to pass as arguments to the function
- `description_expression` - (Required) The description of the generated Signals
- `enabled` - (Required) Whether the rule should generate Signals
- `entity_selectors` - (Required) The entities to generate Signals on
  + `entityType` - (Required) The type of the entity to generate the Signal on.
  + `expression` - (Required) The expression or field name to generate the Signal on.
- `group_by_entity` - (Optional; defaults to true) Whether to group records by the specified entity fields
- `group_by_fields` - (Optional) A list of fields to group records by
- `is_prototype` - (Optional) Whether the generated Signals should be prototype Signals
- `match_expression` - (Required) The expression for which records to match on
- `name` - (Required) The name of the Rule
- `name_expression` - (Required) The name of the generated Signals
- `severity_mapping` - (Required) The configuration of how the severity of the Signals should be mapped from the Records
  + `type` - (Required) Whether to set a constant severity ("constant"), set the severity based on the direct value of a record field ("fieldValue"), or map a record field value to a severity ("fieldValueMapping").
  + `default` - (Optional) The severity to use in the "constant" case or to fall back to if the field used by "fieldValue"/"fieldValueMapping" is not populated.
  + `field` - (Optional) The field to use in the "fieldValue"/"fieldValueMapping" cases.
  + `mapping` - (Optional) The map of record values to severities to use in the "fieldValueMapping" case
    - `type` - (Required) Must be set to "eq" currently
    - `from` - (Required) The record value to map from
    - `to` - (Required) The severity value to map to
- `summary_expression` - (Optional) The summary of the generated Signals
- `tags` - (Required) The tags of the generated Signals
- `trigger_expression` - (Required) The expression to determine whether a Signal should be created based on the aggregation results
- `window_size` - (Required) How long of a window to aggregate records for. Current acceptable values are T05M, T10M, T30M, T60M, T24H, T12H, or T05D.

The following attributes are exported:

- `id` - The internal ID of the aggregation rule.

## Import

Aggregation Rules can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_aggregation_rule.aggregation_rule id
```