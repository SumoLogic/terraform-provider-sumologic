---
layout: "sumologic"
page_title: "SumoLogic: sumologic_collector_ingest_budget_assignment"
description: |-
  Assigns a Sumologic Collector to a Sumologic Ingest Budget.
---

# sumologic_collector_ingest_budget_assignment
Assigns a [Sumologic Collector][1] to a [Sumologic Ingest Budget][2].

## Example Usage
```hcl
resource "sumologic_collector" "collector" {
  name = "my-collector"
}

resource "sumologic_ingest_budget" "budget" {
  name           = "my-budget"
  field_value    = "my-budget"
  capacity_bytes = 30000000000
}

resource "sumologic_collector_ingest_budget_assignment" "assignment" {
  collector_id     = "${sumologic_collector.collector.id}"
  ingest_budget_id = "${sumologic_ingest_budget.budget.id}"
}
```

## Argument reference
The following arguments are supported:
- `collector_id` - (Required) ID of the collector to assign to the ingest budget.
- `ingest_budget_id` - (Required) ID of the ingest budget to assign the collector to.

## Attributes reference
The following attributes are exported:
- `id` - The internal ID of the assignment.

[Back to Index][0]

[0]: ../README.md
[1]: https://help.sumologic.com/Send_Data/Hosted_Collectors
[2]: https://help.sumologic.com/Manage/Ingestion-and-Volume/Ingest_Budgets
