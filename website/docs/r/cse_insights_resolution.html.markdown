---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_insights_resolution"
description: |-
  Provides a CSE Insights Resolution
---

# insights_resolution
Provides a CSE Insights Resolution.

## Example Usage
```hcl
resource "sumologic_cse_insights_resolution" "insights_resolution" {
  name = "New Name"
  description = "New description"
  parent = "No Action"
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) The name of the insights resolution.
- `description` - (Required) The description of the insights resolution.
- `parent` - (Required) The name of the built-in parent insights resolution. Supported values: "Resolved", "False Positive", "No Action", "Duplicate"


The following attributes are exported:

- `id` - The internal ID of the insights resolution.


