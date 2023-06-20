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
  lookback_days = 13.0
  threshold = 12.0
  global_signal_suppression_window = 48.0
}
```

## Argument reference

The following arguments are supported:

- `lookback_days` - (Optional) Detection window expressed in days.
- `threshold` - (Optional) Detection threshold activity score.
- `global_signal_suppression_window` - (Optional) Detection global signal suppression window expressed in hours.

The following attributes are exported:

- `ID` - The internal ID of the insights configuration.

## Import

Insights Configuration can be imported using the id `cse-insights-configuration`:

~> **NOTE:** Only `cse-insights-configuration` id should be used when importing hte insights configuration. Using any other id may have unintended consequences.

```hcl
terraform import sumologic_cse_insights_configuration.insights_configuration cse-insights-configuration
```