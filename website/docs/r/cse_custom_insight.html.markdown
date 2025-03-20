---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_custom_insight"
description: |-
  Provides a CSE Custom Insight
---

# sumologic_cse_custom_insight
Provides a Sumo Logic CSE Custom Insight.

## Example Usage
```hcl
resource "sumologic_cse_custom_insight" "custom_insight" {
  description = "Insight description"
  enabled = true
  ordered = true
  name = "Custom Insight Example"
  rule_ids = ["MATCH-S00001", "THRESHOLD-U00005"]
  severity = "HIGH"
  signal_match_strategy = "ENTITY"
  dynamic_severity {
    minimum_signal_severity = 8
    insight_severity = "CRITICAL"
  }
  signal_names = ["Some Signal Name", "Wildcard Signal Name *"]
  tags = ["_mitreAttackTactic:TA0009"]
}
```

## Argument reference

The following arguments are supported:

- `description` - (Required) The description of the generated Insights
- `enabled` - (Required) Whether the Custom Insight should generate Insights
- `ordered` - (Required) Whether the signals matching the rule IDs/signal names must be in the same chronological order as they are listed in the Custom Insight
- `name` - (Required) The name of the Custom Insight and the generated Insights
- `rule_ids` - (Optional) The Rule IDs to match to generate an Insight (exactly one of rule_ids or signal_names must be specified)
- `severity` - (Required) The severity of the generated Insights (CRITICAL, HIGH, MEDIUM, or LOW)
- `signal_match_strategy` - (Optional) The signal match strategy to use when generating insights (ENTITY, STRICT)
- `dynamic_severity` - (Optional) The severity of the generated Insight that is based on the severity of the Signals that trigger the Insight.
  + `minimum_signal_severity` - (Required) minimum Signal severity as the threshold for an Insight severity level
  + `insight_severity` - (Required) The severity of the generated Insight (CRITICAL, HIGH, MEDIUM, or LOW)
- `signal_names` - (Optional) The Signal names to match to generate an Insight (exactly one of rule_ids or signal_names must be specified)
- `tags` - (Required) The tags of the generated Insights

The following attributes are exported:

- `id` - The internal ID of the chain rule.

## Import

Custom Insights can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_custom_insight.custom_insight id
```