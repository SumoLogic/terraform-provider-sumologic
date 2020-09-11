---
layout: "sumologic"
page_title: "SumoLogic: sumologic_ingest_budget_v2"
description: |-
  Provides a Sumologic Ingest Budget v2
---

# sumologic_ingest_budget_v2
Provides a [Sumologic Ingest Budget v2][1].This API is in private beta and is not available until given access. To participate in the beta program contact your Sumo Logic account representative.

## Example Usage
```hcl
resource "sumologic_ingest_budget_v2" "budget" {
  name           = "testBudget"
  scope          = "_sourceCategory=*prod*nginx*"
  capacity_bytes = 30000000000
  description    = "For testing purposes"
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) Display name of the ingest budget. This must be unique across all of the ingest budgets
  * `scope` - (Required) A scope is a constraint that will be used to identify the messages on which budget needs to be applied. A scope is consists of key and value separated by =. The field must be enabled in the fields table.
  * `capacity_bytes` - (Required) Capacity of the ingest budget, in bytes.
  * `description` - (Optional) The description of the collector.
  * `timezone` - (Optional) The time zone to use for this collector. The value follows the [tzdata][2] naming convention. Defaults to `Etc/UTC`
  * `reset_time` - (Optional) Reset time of the ingest budget in HH:MM format. Defaults to `00:00`
  * `reset_time` - (Optional) Reset time of the ingest budget in HH:MM format. Defaults to `00:00`
  * `description` - (Optional) Description of the ingest budget.
  * `action` - (Optional) Action to take when ingest budget's capacity is reached. All actions are audited. Supported values are `stopCollecting` and `keepCollecting`.

The following attributes are exported:

  * `id` - The internal ID of the ingest budget. This can be used to assign collectors to the ingest budget.

## Import
Ingest budgets can be imported using the name, e.g.:

```hcl
terraform import sumologic_ingest_budget_v2.budget budgetName
```

[1]: https://help.sumologic.com/Beta/Metadata_Ingest_Budgets
[2]: https://en.wikipedia.org/wiki/Tz_database
