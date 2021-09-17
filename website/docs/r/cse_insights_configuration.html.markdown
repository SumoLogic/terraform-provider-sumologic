---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_insight_status"
description: |-
  Provides a CSE Insight Status
---

# insight_status
Provides a [CSE Insight Status].

## Example Usage
```hcl
resource "sumologic_cse_insight_status" "insight_status" {
  name = "New Name"
  description = "New description"
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) The name of the insight status.
- `description` - (Required) The description of the insight status.


The following attributes are exported:

- `id` - The internal ID of the insight status.


