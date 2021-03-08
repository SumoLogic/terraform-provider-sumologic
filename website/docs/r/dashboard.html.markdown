---
layout: "sumologic"
page_title: "SumoLogic: sumologic_dashboard"
description: |-
  Provides a Sumologic Dashboard (New)
---

# sumologic_dashboard
Provides a [Sumologic Dashboard (New)][1].

## Example Usage
```hcl
data "sumologic_personal_folder" "personalFolder" {}

resource "sumologic_dashboard" "api-dashboard" {
	title = "Api Health Dashboard"
	description = "Demo dashboard description"
	folder_id = data.sumologic_personal_folder.personalFolder.id
	refresh_interval = 60
	theme = "Dark"

	time_range {
		begin_bounded_time_range {
            from {
                literal_time_range {
                    range_name = "today"
                }
            }
        }
	}

	## text panel
	panel {
		text_panel {
			key = "text-panel-01"
			title = "Api Health"
			visual_settings = jsonencode(
				{
					"text": {
						"verticalAlignment": "top",
						"horizontalAlignment": "left",
						"fontSize": 12
					}
				}
			)
			keep_visual_settings_consistent_with_parent = true
			text = <<-EOF
				## Api Health Monitoring

				Use this dashboard to monitor API service health. It contains following panels:

				1. API errors: Errors in last 12 hours
				3. API 5xx: Count of 5xx response
				3. CPU utilization: CPU utilization in last 60 mins
			EOF
		}
	}

	## search panel - log query
	panel {
		sumo_search_panel {
			key = "search-panel-01"
			title = "Api Errors by Host"
			description = "Errors in api service since last 12 hours"
			# stacked time series
			visual_settings = jsonencode(
				{
					"general": {
							"mode": "timeSeries",
							"type": "area",
							"displayType": "stacked",
							"markerSize": 5,
							"lineDashType": "solid",
							"markerType": "square",
							"lineThickness": 1
					},
					"title": {
							"fontSize": 14
					},
					"legend": {
							"enabled": true,
							"verticalAlign": "bottom",
							"fontSize": 12,
							"maxHeight": 50,
							"showAsTable": false,
							"wrap": true
					}
				}
			)
			keep_visual_settings_consistent_with_parent = true
			query {
				query_string = "_sourceCategory=api error | timeslice 1h | count by _timeslice, _sourceHost | transpose row _timeslice column _sourceHost"
				query_type = "Logs"
				query_key = "A"
			}
			time_range {
				begin_bounded_time_range {
					from {
						relative_time_range {
							relative_time = "-12h"
						}
					}
				}
			}
		}
	}

	## search panel - metrics query
	panel {
		sumo_search_panel {
			key = "metrics-panel-01"
			title = "Api 5xx Response Count"
			description = "Count of 5xx response from api service"
			# pie chart
			visual_settings = jsonencode(
				{
					"general": {
						"mode": "distribution",
						"type": "pie",
						"displayType": "default",
						"fillOpacity": 1,
						"startAngle": 270,
						"innerRadius": "40%",
						"maxNumOfSlices": 10,
						"aggregationType": "sum"
					},
					"title": {
						"fontSize": 14
					},
				}
			)
			keep_visual_settings_consistent_with_parent = true
			query {
				query_string = "_sourceCategory=api metric=Api-5xx"
				query_type = "Metrics"
				query_key = "A"
				metrics_query_mode = "Advanced"
			}
			time_range {
				begin_bounded_time_range {
					from {
						literal_time_range {
							range_name = "today"
						}
					}
				}
			}
		}
	}

	## search panel - multiple metrics queries
	panel {
		sumo_search_panel {
			key = "metrics-panel-02"
			title = "CPU Utilization"
			description = "CPU utilization in api service"
			# time series with line of dash dot type
			visual_settings = jsonencode(
				{
					"general": {
						"mode": "timeSeries",
						"type": "line",
						"displayType": "smooth",
						"markerSize": 5,
						"lineDashType": "dashDot",
						"markerType": "none",
						"lineThickness": 1
					},
					"title": {
						"fontSize": 14
					},
				}
			)
			keep_visual_settings_consistent_with_parent = true
			query {
				query_string = "metric=Proc_CPU nite-api-1"
				query_type = "Metrics"
				query_key = "A"
				metrics_query_mode = "Basic"
				metrics_query_data {
					metric = "Proc_CPU"
					filter {
						key = "_sourcehost"
						negation = false
						value = "nite-api-1"
					}
					aggregation_type = "None"
				}
			}
			query {
				query_string = "metric=Proc_CPU nite-api-2"
				query_type = "Metrics"
				query_key = "B"
				metrics_query_mode = "Basic"
				metrics_query_data {
					metric = "Proc_CPU"
					filter {
						key = "_sourcehost"
						negation = false
						value = "nite-api-2"
					}
					aggregation_type = "None"
				}
			}
			time_range {
				begin_bounded_time_range {
					from {
						relative_time_range {
							relative_time = "-1h"
						}
					}
				}
			}
		}
	}

	## layout
	layout {
		grid {
			layout_structure {
				key = "text-panel-01"
				structure = "{\"height\":5,\"width\":24,\"x\":0,\"y\":0}"
			}
			layout_structure {
				key = "search-panel-01"
				structure = "{\"height\":10,\"width\":12,\"x\":0,\"y\":5}"
			}
			layout_structure {
				key = "metrics-panel-01"
				structure = "{\"height\":10,\"width\":12,\"x\":12,\"y\":5}"
			}
			layout_structure {
				key = "metrics-panel-02"
				structure = "{\"height\":10,\"width\":24,\"x\":0,\"y\":25}"
			}
		}
	}

	## variables
	variable {
		name = "_sourceHost"
		display_name = "Source Host"
		default_value = "nite-api-1"
		source_definition {
			csv_variable_source_definition {
				values = "nite-api-1,nite-api-2"
			}
		}
		allow_multi_select = true
		include_all_option = true
		hide_from_ui = false
	}
}
```

## Argument reference

The following arguments are supported:

- `title` - (Required) Title of the dashboard.
- `description` - (Optional) Description of the dashboard.
- `folder_id` - (Optional) The identifier of the folder to save the dashboard in. By default it is saved in your
personal folder.
- `refresh_interval` - (Optional) Interval of time (in seconds) to automatically refresh the dashboard.
- `theme` - (Optional) Theme of the dashboard.
- `time_range` - (Block List, Max: 1, Required) Time range of the dashboard. See [time range schema](#schema-for-time_range)
for details.
- `panel` - (Block List, Optional) A list of panels in the dashboard. See [panel schema](#schema-for-panel) for details.
- `layout` - (Block List, Max: 1, Optional) Layout of the dashboard. See [layout schema](#schema-for-layout) for details.
- `variable` - (Block List, Optional) A list of variables for the dashboard. See [variable schema](#schema-for-variable)
for details.

## Attributes reference
In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the dashboard.

### Schema for `time_range`
- `complete_literal_time_range` - (Block List, Max: 1, Optional) Literal time range. See
[complete_literal_time_range schema](#schema-for-complete_literal_time_range) for details.
- `begin_bounded_time_range` - (Block List, Max: 1, Optional) Bounded time range. See
[begin_bounded_time_range schema](#schema-for-begin_bounded_time_range) schema for details.
schema for details.

### Schema for `complete_literal_time_range`
- `range_name` - (Required) Name of complete literal time range. One of `today`, `yesterday`, `previous_week`, `previous_month`.

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

### Schema for `panel`
- `text_panel` - (Block List, Max: 1, Optional) A text panel. See [text_panel schema](#schema-for-text_panel) for details.
- `sumo_search_panel` - (Block List, Max: 1, Optional) A search panel. See [sumo_search_panel schema](#schema-for-sumo_search_panel)
for details.

### Schema for `text_panel`
- `key` - (Required) Key for the panel. Used to create searches for the queries in the panel and configure the layout
of the panel in the dashboard.
- `title` - (Optional) Title of the panel.
- `visual_settings` - (Optional) Visual settings of the panel.
- `keep_visual_settings_consistent_with_parent` - (Optional) Keeps the visual settings, like series colors, consistent
with the settings of the parent panel.
- `text` - (Required) Text to display in the panel.

### Schema for `sumo_search_panel`
- `key` - (Required) Key for the panel. Used to create searches for the queries in the panel and configure the layout
of the panel in the dashboard.
- `title` - (Optional) Title of the panel.
- `visual_settings` - (Optional) Visual settings of the panel.
- `keep_visual_settings_consistent_with_parent` - (Optional) Keeps the visual settings, like series colors, consistent
with the settings of the parent panel.
- `query` - (Block List, Required) A list of queries for the panel. Can be log or metric query. See
[query schema](#schema-for-query) for details.
- `description` - (Optional) Description of the panel.
- `time_range` - (Block List, Max: 1, Optional) Time range of the panel. See [time_range schema](#schema-for-time_range)
for details.
- `coloring_rule` - (Block List, Optional) Coloring rules for the panel. See [coloring_rule schema](#schema-for-coloring_rule)
for details.
- `linked_dashboard` - (Block List, Optional) A list of linked dashboards. See
[linked_dashboard schema](#schema-for-linked_dashboard) for details.

### Schema for `query`
- `query_string` - (Required) The metrics or logs query.
- `query_type` - (Required) The type of the query. One of `Metrics` or `Logs`.
- `query_key` - (Required) The key for metric or log query. Used as an identifier for queries.
- `metric_query_mode` - (Optional) _Should only be specified for metric query_. The mode of the metric query.
One of `Basic` or `Advanced`.
- `metric_query_data` - (Optional) _Should only be specified for metric query_. Data format for the metric query. See
[metric_query_data schema](#schema-for-metric_query_data) for details.

### Schema for `metric_query_data`
- `metric` - (Required) The metric of the query.
- `aggregation_type` - (Optional) The type of aggregation. One of `Count`, `Minimum`, `Maximum`, `Sum`, `Average`, `None`.
- `group_by` - The field to group the results by.
- `filter` - (Block List, Required) A list of filters for the metrics query.
    - `key` - (Required) The key of the metrics filter.
    - `value` - (Required) The value of the metrics filter.
    - `negation` - (Optional) Whether or not the metrics filter is negated.
- `operator` - (Block List, Optional) A list of operator data for the metrics query.

### Schema for `operator`
- `operator_name` - (Required) The name of the metrics operator.
- `parameter` - (Block List, Required) A list of operator parameters for the operator data.
    - `key` - (Required) The key of the operator parameter.
    - `value` - (Required) The value of the operator parameter.

### Schema for `coloring_rule`
- `scope` - (Required) Regex string to match queries to apply coloring to.
- `single_series_aggregate_function` - (Required) Function to aggregate one series into one single value.
- `multiple_series_aggregate_function` - (Required) Function to aggregate the aggregate values of multiple time series
into one single value.
- `color_threshold` - (Block List, Optional) A list of color threshold object.
    - `color` - (Required) Color for the threshold.
    - `min` - (Optional) Absolute inclusive threshold to color by.
    - `max` - (Optional) Absolute exclusive threshold to color by.

### Schema for `linked_dashboard`
- `id` - (Required) Identifier of the linked dashboard.
- `relative_path` - (Optional) Relative path of the linked dashboard to the dashboard of the linking panel.
- `include_time_range` - (Optional) Include time range from the current dashboard to the linked dashboard. _Defaults to true_.
- `include_variables` - (Optional) Include variables from the current dashboard to the linked dashboard. _Defaults to true_.

### Schema for `layout`
- `grid` - (Block List, Max: 1, Optional) Panel layout for the dashboard.

### Schema for `grid`
- `layout_structure` - (Block List, Required) Layout structure for the panels in the dashboard.
    - `key` - (Required) The identifier of the panel that this structure applies to. It's same as `panel.key`.
    - `structure` - (Required) The structure of the panel.

### Schema for `variable`
- `name` - (Required) Name of the variable. The variable name is case-insensitive.
- `display_name` - (Optional) Display name of the variable shown in the UI. If this field is empty, the name field will be used.
- `default_value` - (Optional) Default value of the variable.
- `source_definition` - (Required) Source definition for variable values. See
[source_definition schema](#schema-for-source_definition) for details.
- `allow_multi_select` - (Optional) Allow multiple selections in the values dropdown.
- `include_all_option` - (Optional) Include an "All" option at the top of the variable's values dropdown. _Defaults to true._
- `hide_from_ui` - (Optional) Hide the variable in the dashboard UI.

### Schema for `source_definition`
- `log_query_variable_source_definition` - (Optional) Variable values from a log query.
    - `query` - (Required) A log query.
    - `field` - (Required) A field in log query to populate the variable values
- `metadata_variable_source_definition` - (Optional) Variable values from a metric query.
    - `filter` - (Required) Filter to search the catalog.
    - `key` - (Required) Return the values for this given key.
- `csv_variable_source_definition` - (Optional) Variable values in csv format.
    - `values` - (Required) A comma separated values for the variable.


## Import
Dashboard can be imported using the dashboard id, e.g.:
```hcl
terraform import sumologic_dashboard.example-dashboard q0IKwAK5t2qRI4sgiANwnS87k5S4twN2sCpTuZFSsz6ZmbENPsG7PnpqZygc
```

[1]: https://help.sumologic.com/Visualizations-and-Alerts/Dashboard_(New)
