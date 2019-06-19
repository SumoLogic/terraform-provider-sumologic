# sumologic_ingest_budget
Provides a [Sumologic Ingest Budget][1].

## Example Usage
```hcl
resource "sumologic_ingest_budget" "budget" {
  name           = "test"
  field_value    = "test"
  capacity_bytes = 30000000000
  description    = "For testing purposes"
}
```

## Argument reference
The following arguments are supported:
- `name` - (Required) Display name of the ingest budget. This must be unique across all of the ingest budgets
- `field_value` - (Required) Custom field value that is used to assign Collectors to the ingest budget.
- `capacity_bytes` - (Required) Capacity of the ingest budget, in bytes.
- `description` - (Optional) The description of the collector.
- `timezone` - (Optional) The time zone to use for this collector. The value follows the [tzdata][2] naming convention. Defaults to `Etc/UTC`
- `reset_time` - (Optional) Reset time of the ingest budget in HH:MM format. Defaults to `00:00`
- `reset_time` - (Optional) Reset time of the ingest budget in HH:MM format. Defaults to `00:00`
- `description` - (Optional) Description of the ingest budget.
- `action` - (Optional) Action to take when ingest budget's capacity is reached. All actions are audited. Supported values are `stopCollecting` and `keepCollecting`.
                        

## Attributes reference
The following attributes are exported:
- `id` - The internal ID of the ingest budget. This can be used to assign collectors to the ingest budget.

## Import
Ingest budgets can be imported using the name, e.g.:
```bash
terraform import sumologic_collector.test test
```

[Back to Index][0]

[0]: ../README.md
[1]: https://help.sumologic.com/Manage/Ingestion-and-Volume/Ingest_Budgets
[2]: https://en.wikipedia.org/wiki/Tz_database
