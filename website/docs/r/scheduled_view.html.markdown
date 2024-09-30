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
  lifecycle {
    prevent_destroy = true
    ignore_changes = [index_id]
  }
}
```

## Argument reference

The following arguments are supported:

~> For attributes that force a new resource, if the value is updated, it will destroy the resource and recreate it which may incur significant costs. We advise customers to set the `lifecycle` attribute `prevent_destroy` to `true` to avoid accidentally destroying and recreating expensive resources.

- `index_name` - (Required, Forces new resource) Name of the index (scheduled view).
- `query` - (Required, Forces new resource) Log query defining the scheduled view.
- `start_time` - (Required, Forces new resource) Start timestamp in UTC in RFC3339 format.
- `retention_period` - (Optional) The number of days to retain data in the scheduled view, or -1 to use the default value for your account. Only relevant if your account has multi-retention. enabled.
- `data_forwarding_id` - (Optional) An optional ID of a data forwarding configuration to be used by the scheduled view.
- `parsing_mode` - (Optional, Forces new resource) Default to `Manual`. Define the parsing mode to scan the JSON format log messages. Possible values are: `AutoParse` - In AutoParse mode, the system automatically figures out fields to parse based on the search query. `Manual` - While in the Manual mode, no fields are parsed out automatically. For more information see Dynamic Parsing.
- `reduce_retention_period_immediately` - (Optional) This is required on update if the newly specified retention period is less than the existing retention period. In such a situation, a value of true says that data between the existing retention period and the new retention period should be deleted immediately; if false, such data will be deleted after seven days. This property is optional and ignored if the specified retentionPeriod is greater than or equal to the current retention period.

The following attributes are exported:

- `id` - The internal ID of the scheduled view.
- `index_id` - The Index ID of the scheduled view. It remains unchanged during resource updates, and any manual modifications will be disregarded. While it’s not mandatory, we recommend to ignore this via `ignore_changes = [index_id]`.

## Import
Scheduled Views can can be imported using the id. The list of scheduled views and their ids can be obtained using the Sumologic [scheduled views api][2].

```hcl
terraform import sumologic_scheduled_view.failed_connections 1234567890
```

[1]: https://help.sumologic.com/Manage/Scheduled-Views
[2]: https://api.sumologic.com/docs/#operation/listScheduledViews
