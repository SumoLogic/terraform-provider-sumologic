---
layout: "sumologic"
page_title: "SumoLogic: sumologic_ingest_budget_v2"
description: |-
  Provides a Sumologic Ingest Budget v2
---

# sumologic_ingest_budget_v2
Provides a [Sumologic Ingest Budget v2][1].

## Example Usage
```hcl
resource "sumologic_ingest_budget_v2" "budget" {
  name            = "testBudget"
  scope           = "_sourceCategory=*prod*nginx*"
  budget_type     = "dailyVolume"
  capacity_bytes  = 30000000000
  description     = "For testing purposes"
  timezone        = "Etc/UTC"
  action          = "keepCollecting"
  reset_time      = "00:00"
  audit_threshold = 85
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) Display name of the ingest budget. This must be unique across all of the ingest budgets
  * `scope` - (Required) A scope is a constraint that will be used to identify the messages on which budget needs to be applied. A scope is consists of key and value separated by =. The field must be enabled in the fields table.
  * `capacity_bytes` - (Required) Capacity of the ingest budget, in bytes. It takes a few minutes for Collectors to stop collecting when capacity is reached. We recommend setting a soft limit that is lower than your needed hard limit. The capacity bytes unit varies based on the budgetType field. For `dailyVolume` budgetType the capacity specified is in bytes/day whereas for `minuteVolume` budgetType its bytes/min.
  * `description` - (Optional) The description of the collector.
  * `budget_type` - (Optional) The type of budget. Supported values are:  * `dailyVolume` * `minuteVolume`. Default value is `dailyVolume`.
  * `timezone` - (Optional) The time zone to use for this collector. The value follows the [tzdata][2] naming convention. Defaults to `Etc/UTC`
  * `action` - (Optional) Action to take when ingest budget's capacity is reached. All actions are audited. Supported values are `stopCollecting` and `keepCollecting`.
  * `reset_time` - (Optional) Reset time of the ingest budget in HH:MM format. Defaults to `00:00`
  * `audit_threshold` - (Optional) The threshold as a percentage of when an ingest budget's capacity usage is logged in the Audit Index.
  
The following attributes are exported:

  * `id` - The internal ID of the ingest budget. 

## Import
Ingest budgets can be imported using the budget ID or name, e.g.:

```hcl
terraform import sumologic_ingest_budget_v2.budget 00000000000123AB
terraform import sumologic_ingest_budget_v2.dev_budget "DevelopersBudget"
```

[1]: https://help.sumologic.com/Beta/Metadata_Ingest_Budgets
[2]: https://en.wikipedia.org/wiki/Tz_database
