---
layout: "sumologic"
page_title: "SumoLogic: sumologic_scan_budget"
description: |-
  Provides a Sumologic Scan Budget
---

# sumologic_scan_budget
Provides a [Sumologic Scan Budget][1].

## Example Usage
```hcl
resource "sumologic_scan_budget" "budget" {
  name          = "TestBudget"
  capacity      = 10
  unit          = "GB"
  budget_type   = "ScanBudget"
  window        = "Query"
  applicable_on = "PerEntity"
  group_by      = "User"
  action        = "StopScan"
  status        = "active"
  scope {
      included_users = ["000000000000011C"]
      excluded_users = []
      included_roles = []
      excluded_roles = ["0000000000000196"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Display name of the scan budget. This must be unique across all the scan budgets.
* `capacity` - (Required) Capacity of the scan budget. Only whole numbers are supported.
* `unit` - (Required) Unit of the capacity. Supported values are: `MB`, `GB` and `TB`.
* `budget_type` - (Required) Type of the budget. Supported values are: `ScanBudget`.
* `window` - (Required) Window of the budget. Supported values are: `Query`, `Daily`, `Weekly` and `Monthly`.
* `applicable_on` - (Required) Grouping of the budget. Supported values are: `PerEntity` and `Sum`.
* `group_by` - (Required) Grouping Entity of the budget. Supported values are: `User`.
* `action` - (Required) Action to be taken if the budget is breached. Supported values are: `StopForeGroundScan` and `Warn`.
* `scope` - (Required) Scope of the budget.
* `status` - (Required) Signifies the state of the budget. Supported values are: `active` and `inactive`.

The following attributes are exported:

* `id` - The internal ID of the budget.

## Import
Scan budgets can be imported using the budget ID, e.g.:

```hcl
terraform import sumologic_scan_budget.budget 00000000000123AB
```

[1]: https://help.sumologic.com/docs/manage/manage-subscription/usage-management/
