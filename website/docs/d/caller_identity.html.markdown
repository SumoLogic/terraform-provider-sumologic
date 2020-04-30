---
layout: "sumologic"
page_title: "SumoLogic: sumologic_caller_identity"
description: |-
  Provides an easy way to retrieve Sumo Logic auth details.
---

# sumologic_caller_identity
Provides an easy way to retrieve Sumo Logic auth details.


## Example Usage
```hcl
data "sumologic_caller_identity" "current" {}
```


## Attributes reference

The following attributes are exported:

- `access_id` - Sumo Logic access ID.
- `environment` - API endpoint environment.


