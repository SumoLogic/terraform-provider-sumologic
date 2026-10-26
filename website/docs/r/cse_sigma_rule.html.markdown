---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_sigma_rule"
description: |-
  Provides a CSE Match Rule imported from a Sigma rule
---

# sumologic_cse_sigma_rule
Provides a Sumo Logic CSE [Match Rule](https://help.sumologic.com/Cloud_SIEM_Enterprise/CSE_Rules/03_Write_a_Match_Rule) created from a Sigma rule.

The Sigma rule is looked up by its ID and converted to a Sumo Logic match rule. The converted fields are copied into the new rule when it is created, and the rule is marked as originating from Sigma. Later changes to the Sigma rule are not applied to the existing rule.

## Example Usage
```hcl
resource "sumologic_cse_sigma_rule" "sigma_rule" {
  sigma_id = "637f689e-b4a5-4a86-be0e-0100a0a33ba2"
  enabled  = true
}
```

Any converted field can be overridden:
```hcl
resource "sumologic_cse_sigma_rule" "sigma_rule" {
  sigma_id     = "637f689e-b4a5-4a86-be0e-0100a0a33ba2"
  enabled      = true
  name         = "My Sigma Rule"
  is_prototype = false
  tags         = ["_mitreAttackTactic:TA0004"]
}
```

## Argument reference

The following arguments are supported:

- `sigma_id` - (Required) The ID (UUID) of the Sigma rule to import. Changing this forces a new rule to be created.
- `enabled` - (Required) Whether the rule should generate Signals
- `name` - (Optional) The name of the Rule. Defaults to the name of the Sigma rule.
- `name_expression` - (Optional) The name of the generated Signals. Defaults to the converted Sigma value.
- `description_expression` - (Optional) The description of the generated Signals. Defaults to the converted Sigma value.
- `summary_expression` - (Optional) The summary of the generated Signals. Defaults to the converted Sigma value.
- `expression` - (Optional) The expression for which records to match on. Defaults to the converted Sigma value.
- `is_prototype` - (Optional) Whether the generated Signals should be prototype Signals. Defaults to the converted Sigma value.
- `entity_selectors` - (Optional) The entities to generate Signals on. Defaults to the converted Sigma value.
  + `entity_type` - (Required) The type of the entity to generate the Signal on.
  + `expression` - (Required) The expression or field name to generate the Signal on.
- `severity_mapping` - (Optional) The configuration of how the severity of the Signals should be mapped from the Records. Defaults to the converted Sigma value.
  + `type` - (Required) Whether to set a constant severity ("constant"), set the severity based on the direct value of a record field ("fieldValue"), or map a record field value to a severity ("fieldValueMapping").
  + `default` - (Optional) The severity to use in the "constant" case or to fall back to if the field used by "fieldValue"/"fieldValueMapping" is not populated.
  + `field` - (Optional) The field to use in the "fieldValue"/"fieldValueMapping" cases.
  + `mapping` - (Optional) The map of record values to severities to use in the "fieldValueMapping" case
      - `type` - (Required) Must be set to "eq" currently
      - `from` - (Required) The record value to map from
      - `to` - (Required) The severity value to map to
- `tags` - (Optional) The tags of the generated Signals. Defaults to the converted Sigma value.
- `suppression_window_size` - (Optional) For how long to suppress Signal generation, in milliseconds. Must be greater than 0 and less than the global limit of 7 days.

The following attributes are exported:

- `id` - The internal ID of the rule.
- `category` - The category of the rule.
- `original_source` - The original source of the rule. Always `sigma`.

## Import

Importing is not supported, because the Sigma ID cannot be read back from an existing rule.
