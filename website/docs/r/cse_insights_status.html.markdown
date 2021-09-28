---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_insights_status"
description: |-
  Provides a Sumologic CSE Insights Status. Preconfigured Insight statuses cannot be edited or deleted; you can however create custom statuses.
---

# sumologic_cse_insights_status
Provides a Sumologic CSE Insights Status. Preconfigured Insight statuses cannot be edited or deleted; you can however create custom statuses.

## Example Usage
```hcl
resource "sumologic_cse_insights_status" "insights_status" {
  name = "New Name"
  description = "New description"
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) The name of the insights status.
- `description` - (Required) The description of the insights status.


The following attributes are exported:

- `id` - The internal ID of the insights status.

## Import

Insights Status can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_insights_status.insights_status id
```