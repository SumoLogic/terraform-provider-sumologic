---
layout: 'sumologic' page_title: 'SumoLogic: sumologic_slo' description: |- Provides the ability to create, read, delete,
and update SLOs.
---

# sumologic_slo

Provides the ability to create, read, delete, and update SLOs.

## Example SLO

```hcl
resource "sumologic_slo" "slo_tf_test" {
  name        = "slo-tf-test1"
  description = "example SLO created with terraform"
  parent_id   = "0000000000000001"
  signal_type = "Latency"
  service     = "auth"
  application = "login"
  compliance {
    compliance_type = "Rolling"
    size            = "7d"
    target          = 99
    timezone        = "Asia/Kolkata"
  }
  indicator {
    evaluation_type = "Window"
    op              = "LessThan"
    query_type      = "Metrics"
    size            = "1m"
    threshold       = 200
    queries {
      query_group_type = "Threshold"
      query_group {
        row_id = "A"
        query  = "metric=request_time_p90 service=auth api=login"
      }
    }

  }
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) The name of the SLO. The name must be alphanumeric.
- `description` - (Optional) The description of the SLO.
- `parent_id` - (Optional) The ID of the SLO Folder that contains this SLO. Defaults to the root folder.
- `signal_type` - (Required) The type of SLO. Valid values are `Latency`, `Error`,`Throughput`,`Availability`
  , `Other`. Defaults to `Latency`.
- `service` - (Optional) The notifications the SLO will send when the respective trigger condition is met.
- `application` - (Optional) Whether to group notifications for individual items that meet the trigger condition.
  Defaults to true.
- `compliance` - (Required) The compliance settings for the SLO.
    - `compliance_type` - (Required) The type of compliance to use. Valid values are `Rolling` or `Calendar`.
    - `target` - (Required) The target value to use, must be a number between 0 and 100.
    - `timezone` - (Required) Time zone for the SLO compliance. Follow the format in the [IANA Time Zone Database][3].
    - `size` - (Required) The size of the compliance period to use, minimum of `1d` and maximum of `14d`.
- `indicator` - (Required) The service level indicator on which SLO is to be defined.
    - `evaluation_type` - (Required) Evaluate SLI using successful/total windows, or occurrence of successful events
      over entire compliance period. Valid values are `Window` or `Request`.
    - `op` - (Required) The operator used to define a successful window or event. Valid values are `LessThan`, `LessThanOrEqual`, `GreaterThan`
      , `GreaterThanOrEqual`.
    - `size` - (Required) The size of the window to use, minimum of `1m` and maximum of `1h`. Only applicable for Window based evaluation.
    - `query_type` - (Required) The type of query to use. Valid values are `Metrics` or `Logs`.
    - `threshold` - (Required) Threshold for classifying window as successful or unsuccessful, i.e. the minimum value for `good windows / total events`.
    - `aggregation` - (Optional) Aggregation function applied over each window to arrive at SLI. Valid values are `Avg`
      , `Sum`, `Count`, `Max`, `Min` and `p[1-99]`.
    - `queries` - (Required) The queries to use.
        - `query_group_type` - (Required) The type of query. Valid values are `Successful`, `Unsuccessful`, `Total`
          , `Threshold`.
        - `query_group` - (Required) List of queries to use.
            - `row_id` - (Required) The row ID to use.
            - `query` - (Required) The query string to use.
            - `use_row_count` - (Required) Whether to use the row count. Defaults to false.
            - `field` - (Optional) Field of log query output to compare against. To be used only for logs based data
              type when `use_row_count` is false.

[1]: https://help.sumologic.com/?cid=10020

[2]: slo_folder.html.markdown

[3]: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List
