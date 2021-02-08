---
layout: "sumologic"
page_title: "SumoLogic: sumologic_scheduled_view"
description: |-
  Provides a Sumologic Scheduled View
---

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
- `query` - (Required) Log query defining the scheduled view. This value cannot be updated.
- `start_time` - (Required) Starting date/time for log indexing.
- `retention_period` - (Optional) Number of days to keep the scheduled view data for.
- `data_forwarding_id` - (Optional) ID of a data forwarding configuration to be used by the scheduled view.

The following attributes are exported:

- `id` - The internal ID of the scheduled view.

## Import
Scheduled Views can can be imported using the id. The list of scheduled views and their ids can be obtained using the Sumologic [scheduled views api][2].

```hcl
terraform import sumologic_scheduled_view.failed_connections 1234567890
```

[1]: https://help.sumologic.com/Manage/Scheduled-Views
[2]: https://api.sumologic.com/docs/#operation/listScheduledViews
