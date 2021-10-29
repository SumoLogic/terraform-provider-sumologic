---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_match_rule"
description: |-
  Provides a CSE Match Rule
---

# sumologic_cse_match_rule
Provides a Sumo Logic CSE [Match Rule](https://help.sumologic.com/Cloud_SIEM_Enterprise/CSE_Rules/03_Write_a_Match_Rule).

## Example Usage
```hcl
resource "sumologic_cse_match_rule" "match_rule" {
  description_expression = "Signal description"
  enabled = true
  entity_selectors {
    entity_type = "_ip"
    expression = "srcDevice_ip"
  }
  expression = "objectType = \"Network\""
  is_prototype = false
  name = "Match Rule Example"
  name_expression = "Signal name"
  severity_mapping {
    type = "constant"
    default = 5
  }
  summary_expression = "Signal summary"
  tags = ["_mitreAttackTactic:TA0009"]
}
```

## Argument reference

The following arguments are supported:

- `description_expression` - (Required) The description of the generated Signals
- `enabled` - (Required) Whether the rule should generate Signals
- `entity_selectors` - (Required) The entities to generate Signals on
  + `entityType` - (Required) The type of the entity to generate the Signal on.
  + `expression` - (Required) The expression or field name to generate the Signal on.
- `expression` - (Required) The expression for which records to match on
- `is_prototype` - (Optional) Whether the generated Signals should be prototype Signals
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

The following attributes are exported:

- `id` - The internal ID of the match rule.

## Import

Match Rules can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_match_rule.match_rule id
```