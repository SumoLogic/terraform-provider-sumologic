---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_insights_status"
description: |-
  Provides a CSE Insights Status
---

# insights_status
Provides a CSE Insights Status.

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


