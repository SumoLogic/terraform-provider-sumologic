---
layout: "sumologic"
page_title: "SumoLogic: sumologic_metrics_search"
description: |-
  Provides a Sumologic Metrics Search
---

# sumologic_metrics_search
Provides a [Sumologic Metrics Search][1].

## Example Usage
```hcl
data "sumologic_personal_folder" "personalFolder" {}

resource "sumologic_metrics_search" "example_metrics_search" {
    title = "Demo Metrics Search"
    description = "Demo search description"
    parent_id = data.sumologic_personal_folder.personalFolder.id
    metrics_queries {
	    row_id = "A"
		query = "metric=cpu_idle | avg"
	}
    desired_quantization_in_secs = 0

    time_range {
        begin_bounded_time_range {
            from {
                relative_time_range {
                    relative_time = "-30m"
                }
            }
        }
    }
}
```

## Argument reference

The following arguments are supported:

- `title` - (Required) Title of the search.
- `description` - (Required) Description of the search.
- `parent_id` - (Required) The identifier of the folder to create the log search in.
- `log_query` - Log query used to add an overlay to the chart.
- `metrics_queries` - (Required) Array of objects [MetricsSearchQuery](#schema-for-metrics_search_query). Metrics queries, up to the maximum of six.
- `time_range` - (Block List, Max: 1, Required) Time range of the log search. See [time range schema](#schema-for-time_range)
- `desired_quantization_in_secs` - (Optional) Desired quantization in seconds. Default value is `0`.

### Schema for `metrics_search_query`
- `row_id` - Row id for the query row, A to Z letter.
- `query` - A metric query consists of a metric, one or more filters and optionally, one or more [Metrics Operators](https://help.sumologic.com/?cid=10144).
Strictly speaking, both filters and operators are optional.
Most of the [Metrics Operators](https://help.sumologic.com/?cid=10144) are allowed in the query string except `fillmissing`, `outlier`, `quantize` and `timeshift`.
In practice, your metric queries will almost always contain filters that narrow the scope of your query.
For more information about the query language see [Metrics Queries](https://help.sumologic.com/?cid=1079).

### Schema for `time_range`
- `complete_literal_time_range` - (Block List, Max: 1, Optional) Literal time range. See
[complete_literal_time_range schema](#schema-for-complete_literal_time_range) for details.
- `begin_bounded_time_range` - (Block List, Max: 1, Optional) Bounded time range. See
[begin_bounded_time_range schema](#schema-for-begin_bounded_time_range) schema for details.

### Schema for `complete_literal_time_range`
- `range_name` - (Required) Name of complete literal time range. One of `today`, `yesterday`, `previous_week`, and
    `previous_month`.

### Schema for `begin_bounded_time_range`
- `from` - (Block List, Max: 1, Required) Start boundary of bounded time range. See
[time_range_boundary schema](#schema-for-time_range_boundary) for details.
- `to` - (Block List, Max: 1, Optional) End boundary of bounded time range. See
[time_range_boundary schema](#schema-for-time_range_boundary) for details.

### Schema for `time_range_boundary`
- `epoch_time_range` - (Block List, Optional) Time since the epoch.
    - `epoch_millis` - (Required) Time as a number of milliseconds since the epoch.

- `iso8601_time_range` - (Block List, Optional) Time in ISO 8601 format.
    - `iso8601_time` - (Required) Time as a string in ISO 8601 format.

- `relative_time_range` - (Block List, Optional) Time in relative format.
    - `relative_time` - (Required) Relative time as a string consisting of following elements:
      1. `-` (optional): minus sign indicates time in the past,
      2. `<number>`: number of time units,
      3. `<time_unit>`: time unit; possible values are: `w` (week), `d` (day), `h` (hour), `m` (minute), `s` (second).

      Multiple pairs of `<number><time_unit>` may be provided, and they may be in any order. For example,
      `-2w5d3h` points to the moment in time 2 weeks, 5 days and 3 hours ago.

- `literal_time_range` - (Block List, Optional) Time in literal format.
    - `range_name` - (Required) One of `now`, `second`, `minute`, `hour`, `day`, `today`, `week`, `month`, `year`.

## Attributes reference
In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the log search.


## Import
A metrics search can be imported using it's identifier, e.g.:
```hcl
terraform import sumologic_metrics_search.example_search 0000000007FFD79D
```

[1]: https://help.sumologic.com/docs/metrics/metrics-queries/metrics-explorer/