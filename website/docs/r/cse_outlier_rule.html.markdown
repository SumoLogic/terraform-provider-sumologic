___
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_outlier_rule"
description: |-
  Provides a CSE Outlier Rule
---

# sumologic_cse_outlier_rule
Provides a Sumo Logic CSE [Outlier Rule](https://help.sumologic.com/docs/cse/rules/write-outlier-rule/).

## Example Usage
```hcl
resource "sumologic_cse_outlier_rule" "outlier_rule" {
  name = "Spike in Login Failures from a User"
  enabled = true
  severity = 4
  is_prototype = false
  summary_expression = "Excessive count of failure login events identified for user: {{user_username}} based on daily historic activity"

  aggregation_functions {
    name = "current" 
    function = "count"
    arguments = ["true"]
  }
  group_by_fields = ["user_username"]
  
  window_size = "T24H"
  baseline_window_size = 604800000  
  retention_window_size = 7776000000

  floor_value = 10

  name_expression = "Spike in Login Failures from a User"
  description_expression = "Detects excessive failed login attempts for the same username based on a daily outlier standard deviation for said user. This is designed to catch both slow and quick brute force type attacks using a user specific historic baseline. The minimum floor of failures expected by default is set to 10."
  match_expression = "objectType = 'Authentication'\nAND normalizedAction = 'logon'\nAND success = false"
  deviation_threshold = 2

  entity_selectors {
    entity_type = "_username"
    expression = "user_username"
  }

  floor_value            = 0
  deviation_threshold    = 3
  group_by_fields        = ["user_username"]
  is_prototype           = false
  match_expression       = "objectType=\"Authentication\" AND success=false"
  name                   = "Spike in Login Failures"
  name_expression        = "Spike in Login Failures - {{ user_username }}"
  retention_window_size  = "7776000000" // 90 days
  severity               = 1
  summary_expression     = "Spike in Login Failures - {{ user_username }}"
  window_size            = "T24H"
  suppression_window_size = 90000000
  tags = ["_mitreAttackTactic:TA0006", "_mitreAttackTechnique:T1110"]

}
```
## Argument Reference

The following arguments are supported:

- `aggregation_function` - (Required) One named aggregation functions
  + `name` - (Required) The name to use to reference the result
  + `function` - (Required) The function to aggregate with
  + `arguments` - (Required) One or more expressions to pass as arguments to the function
- `baseline_window_size` - (Required) The baseline window size in milliseconds
- `description_expression` - (Required) The description of the generated Signals
- `deviation_threshold` - (Required) The deviation threshold used to calculate the threshold to trigger signals
- `enabled` - (Required) Whether the rule should generate Signals
- `entity_selectors` - (Required) The entities to generate Signals on
  + `entityType` - (Required) The type of the entity to generate the Signal on
  + `expression` - (Required) The expression or field name to generate the Signal on
- `floor_value` - (Required) The minimum threshold to trigger signals
- `group_by_fields` - (Optional) A list of fields to group records by
- `is_prototype` - (Optional) Whether the generated Signals should be prototype Signals
- `match_expression` - (Required) The expression for which records to match on
- `name` - (Required) The name of the Rule
- `name_expression` - (Required) The name of the generated Signals
- `retention_window_size` - (Required) The retention window size in milliseconds
- `severity` - (Required) The severity of the generated Signals
- `summary_expression` - (Optional) The summary of the generated Signals
- `tags` - (Optional) The tags of the generated Signals
- `window_size` - (Required) The window size. Current acceptable values are T60M (1 hr) or  T24H (1 day)
- `suppression_window_size` - (Optional) For how long to suppress Signal generation, in milliseconds. Must be greater than `window_size` and less than the global limit of 7 days.

The following attributes are exported:

- `id` - The ID of the Outlier rule.

## Import

Outlier rules can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_outlier_rule.outlier_rule id
```