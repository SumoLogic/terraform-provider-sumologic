
# sumologic_scheduled_view
Provides a [Sumologic Scheduled View][1].

## Example Usage
```hcl
resource "sumologic_scheduled_view" "failed_connections" {
  index_name = "failed_connections"
  query = <<QUERY
_view=connections connectionStats
| parse "connectionStats.CS *" as body
| json field=body "exitCode", "isHttp2"
| lookup org_name from shared/partners on partner_id=partnerid
| timeslice 10m
QUERY
  start_time = "2019-09-01T00:00:00Z"
  retention_period = 365
}
```

## Argument reference
The following arguments are supported:
- `index_name` - (Required) Name of the index (scheduled view).
- `query` - (Required) Log query defining the scheduled view.
- `start_time` - (Required) Starting date/time for log indexing.
- `retention_period` - (Optional) Number of days to keep the scheduled view data for.
- `data_forwarding_id` - (Optional) ID of a data forwarding configuration to be used by the scheduled view.

## Attributes reference
The following attributes are exported:
- `id` - The internal ID of the scheduled view.

[Back to Index][0]

[0]: ../README.md
[1]: https://help.sumologic.com/Manage/Scheduled-Views
