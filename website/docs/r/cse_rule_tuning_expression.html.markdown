---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_rule_tuning_expression"
description: |-
  Provides a CSE Rule Tuning Expression
---

# rule_tuning_expression
Provides a CSE Rule Tuning Expression.

## Example Usage
```hcl
resource "sumologic_cse_rule_tuning_expression" "rule_tuning_expression" {
  name = "New Rule Tuning Name"
  description = "New Rule Tuning Description"
  expression = "accountId = 1234"
  enabled = "true"
  exclude = "true"
  is_global = "false"
  rule_ids = ["LEGACY-S00084"]
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) The name of the insights status.
- `description` - (Required) The description of the insights status.
- `expression` - (Required) Expression to match.
- `enabled` - (Required) Enabled flag.
- `exclude` - (Required) Set to true to exclude records that also match expression.
- `is_global` - (Required) Set to true if tuning expression intended to be global.
- `rule_ids` - (Required) List of rule ids, for the tuning expression to be applied. ( Empty if is_global set to true)


The following attributes are exported:

- `id` - The internal ID of the rule tuning expression.

## Import

Rule tuning expression can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_rule_tuning_expression.rule_tuning_expression id
```
