---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_insights_configuration"
description: |-
  Provides a CSE Insights Configuration
---

# sumologic_cse_insights_configuration
Provides a CSE Insights Configuration.

## Example Usage
```hcl
resource "sumologic_cse_insights_configuration" "insights_configuration" {
  lookback_days = 13
  threshold = 12
}
```

## Argument reference

The following arguments are supported:

- `lookback_days` - (Optional) Detection window expressed in days.
- `threshold` - (Optional) Detection threshold.

The following attributes are exported:

- `id` - The internal ID of the insights configuration.

## Import

Insights Configuration can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_insights_configuration.insights_configuration id
```