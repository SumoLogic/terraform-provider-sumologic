---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_insights_resolution"
description: |-
  Provides a Sumologic CSE Insights Resolution. When an insight gets closed, a resolution indicates why it got closed.
---

# sumologic_cse_insights_resolution
Provides a Sumologic CSE Insights Resolution. When an insight gets closed, a resolution indicates why it got closed.

See [Manage Custom Insight Resolutions](https://www.sumologic.com/help/docs/cse/administration/manage-custom-insight-resolutions/) for how custom resolutions are used when closing insights.

The underlying API is documented in the [Cloud SIEM API reference](https://api.sumologic.com/docs/sec/) (Insight Resolutions).

## Example Usage

### Basic
```hcl
resource "sumologic_cse_insights_resolution" "insights_resolution" {
  name        = "Closed by automation"
  description = "Insight closed without further action"
  parent      = "NoAction"
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) The name of the insights resolution.
- `description` - (Optional) The description of the insights resolution.
- `parent` - (Optional) The built-in resolution this custom resolution falls under. Allowed values: `Resolved`, `FalsePositive`, `NoAction`, `Duplicate` (case sensitive, no spaces); these are the four built-in resolutions "Resolved", "False Positive", "No Action" and "Duplicate" in the UI. The new resolution appears as a sub-resolution under its parent. Any other value is rejected at plan time.

The following attributes are exported:

- `id` - The internal ID of the insights resolution.

## Import

An insights resolution can be imported using its ID, e.g.:
```hcl
terraform import sumologic_cse_insights_resolution.insights_resolution <resolution_id>
```
