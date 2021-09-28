---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_insights_configuration"
description: |-
    Provides the Sumologic CSE Insights Configuration for the whole organization. There can be only one configuration per organization.
---

# sumologic_cse_insights_configuration
Provides the Sumologic CSE Insights Configuration for the whole organization. There can be only one configuration per organization.

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
- `threshold` - (Optional) Detection threshold activity score.

The following attributes are exported:

- `ID` - The internal ID of the insights configuration.

## Import

Insights Configuration can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_insights_configuration.insights_configuration ID
```