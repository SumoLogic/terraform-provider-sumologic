---
layout: 'sumologic' page_title: 'SumoLogic: sumologic_slo' description: |- Provides the ability to create, read, delete,
and update SLOs ([Beta feature](https://help.sumologic.com/Beta/SLO_Reliability_Management#beta-notices)).
---

# sumologic_slo

Provides the ability to create, read, delete, and update [SLOs][1].

## Example SLO

```hcl
resource "sumologic_slo" "slo_tf_window_metric_ratio" {
  name        = "login error rate"
  description = "per minute login error rate over rolling 7 days"
  parent_id   = "0000000000000001"
  signal_type = "Error"
  service     = "auth"
  application = "login"
  compliance {
      compliance_type = "Rolling"
      size            = "7d"
      target          = 95
      timezone        = "Asia/Kolkata"
  }
  indicator {
    window_based_evaluation {
      op         = "LessThan"
      query_type = "Metrics"
      size       = "1m"
      threshold  = 99.0
      queries {
        query_group_type = "Unsuccessful"
        query_group {
          row_id        = "A"
          query         = "service=auth api=login metric=HTTP_5XX_Count"
          use_row_count = false
        }
      }
      queries {
        query_group_type = "Total"
        query_group {
          row_id = "A"
          query  = "service=auth api=login metric=TotalRequests"
          use_row_count = false
        }
      }
    }
  }
}

resource "sumologic_slo" "slo_tf_window_based" {
  name        = "slo-tf-window-based"
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
    window_based_evaluation {
      op         = "LessThan"
      query_type = "Metrics"
      size       = "1m"
      threshold  = 200
      queries {
        query_group_type = "Threshold"
        query_group {
          row_id        = "A"
          query         = "metric=request_time_p90  service=auth api=login"
          use_row_count = false
        }
      }
    }
  }
}

resource "sumologic_slo" "slo_tf_request_based" {
  name        = "slo-tf-request-based"
  description = "example SLO created with terraform for request based SLI"
  parent_id   = sumologic_slo_folder.tf_slo_folder.id
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
    request_based_evaluation {
      op         = "LessThanOrEqual"
      query_type = "Logs"
      threshold  = 1
      queries {
        query_group_type = "Threshold"
        query_group {
          row_id        = "A"
          query         = <<QUERY
          cluster=sedemostaging namespace=warp004*
              | parse "Coffee preparation request time: * ms" as latency nodrop
              |  if(isBlank(latency), "false", "true") as hasLatency
              | where hasLatency = "true"
              |  if(isBlank(latency), 0.0, latency) as latency
              | latency/ 1000 as latency_sec
QUERY
          use_row_count = true
          field = "latency_sec"
        }
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
- `signal_type` - (Required) The type of SLO. Valid values are `Latency`, `Error`, `Throughput`, `Availability`
  , `Other`. Defaults to `Latency`.
- `service` - (Optional) Name of the service.
- `application` - (Optional) Name of the application.
  Defaults to true.
- `compliance` - (Required) The compliance settings for the SLO.
    - `compliance_type` - (Required) The type of compliance to use. Valid values are `Rolling` or `Calendar`.
    - `target` - (Required) The target value to use, must be a number between 0 and 100.
    - `timezone` - (Required) Time zone for the SLO compliance. Follow the format in the [IANA Time Zone Database][3].
    - `size` - (Required) The size of the compliance period to use.
      - For `Rolling` compliance type it must be a multiple of days e.g. `1d`, `2d`.
      - For `Calendar` compliance type the allowed values are `Week`, `Month`, `Quarter`.
    - `start_from` - Start of the calendar window. For `Week` its required and it would be the day of the week (for e.g. Sunday,
      Monday etc).  For `Quarter` its required, it would be the first month of the start of quarter (for e.g. January, February etc.). 
      For `Month` it's not required and is set to first day of the month.
- `indicator` - (Required) The service level indicator on which SLO is to be defined. more details on the difference
  b/w them can be found on
  the [slo help page](https://help.sumologic.com/Beta/SLO_Reliability_Management/Access_and_Create_SLOs)
    - [window_based_evaluation](#window_based_evaluation) - Evaluate SLI using successful/total windows.
    - [request_based_evaluation](#request_based_evaluation) - Evaluate SLI based on occurrence of successful
      events / total events over entire compliance period.

#### window_based_evaluation

- `size` - (Required) The size of the window to use, minimum of `1m` and maximum of `1h`. Only applicable for Window
  based evaluation.
- `query_type` - (Required) The type of query to use. Valid values are `Metrics` or `Logs`.
- `threshold` - (Required) Threshold for classifying window as successful or unsuccessful, i.e. the minimum value
  for `(good windows / total windows) * 100`.
- `op` - (Required) The operator used to define a successful window. Valid values are `LessThan`
  , `LessThanOrEqual`, `GreaterThan`
  , `GreaterThanOrEqual`.
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

#### request_based_evaluation

- `query_type` - (Required) The type of query to use. Valid values are `Metrics` or `Logs`.
- `threshold` - (Required) Compared against threshold query's raw data points to determine success criteria.
- `op` - (Required) Comparison function with threshold. Valid values are `LessThan`, `LessThanOrEqual`, `GreaterThan`
  , `GreaterThanOrEqual`.
- `queries` - (Required) The queries to use.
    - `query_group_type` - (Required) The type of query. Valid values are `Successful`, `Unsuccessful`, `Total`
      , `Threshold`.
    - `query_group` - (Required) List of queries to use.
        - `row_id` - (Required) The row ID to use.
        - `query` - (Required) The query string to use.
        - `use_row_count` - (Required) Whether to use the row count. Defaults to false.
        - `field` - (Optional) Field of log query output to compare against. To be used only for logs based data
          type when `use_row_count` is false.

[1]: https://help.sumologic.com/Beta/SLO_Reliability_Management

[2]: slo_folder.html.markdown

[3]: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List
