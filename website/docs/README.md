# Table of Contents

#### [Basic usage and authentication][0]

#### Data Sources
  + [sumologic_caller_identity][10]
  + [sumologic_collector][11]
  + [sumologic_personal_folder][12]
  + [sumologic_role][13]

#### Resources
##### Sources
  + [sumologic_collector][20]
  + [sumologic_http_source][21]
  + [sumologic_polling_source][22]
  + [sumologic_cloudsyslog_source][23]
  
##### Ingest Budgets
  + [sumologic_collector_ingest_budget_assignment][24]
  + [sumologic_ingest_budget][25]

##### User / Roles
  + [sumologic_role][30]
  + [sumologic_user][31]
  
##### Content
  + [sumologic_scheduled_view][40]
  + [sumologic_partition][41]
  + [sumologic_field_extraction_rule][42]
  + [sumologic_folder][43]
  + [sumologic_content][44]
  + [sumologic_connection][45]

##### Monitors
  + [sumologic_monitor][46]
  + [sumologic_monitor_folder][47]
  

##### Cloud SIEM Enterprise (CSE)
+ [sumologic_cse_network_block][50]
+ [sumologic_cse_insights_resolution][51]
+ [sumologic_cse_insights_status][52]
+ [sumologic_cse_insights_configuration][53]

#### Common Source Properties

The following properties are common to ALL sources and can be used to configure each source.

- `collector_id` - (Required) The ID of the collector to attach this source to.
- `name` - (Required) The name of the source. This is required, and has to be unique in the scope of the collector. Changing this will force recreation the source.
- `description` - (Optional) Description of the source.
- `category` - (Optional) The source category this source logs to.
- `host_name` - (Optional) The source host this source logs to.
- `timezone` - (Optional) The timezone assigned to the source. The value follows the [tzdata][2] naming convention.
- `automatic_date_parsing` - (Optional) Determines if timestamp information is parsed or not. Type true to enable automatic parsing of dates (the default setting); type false to disable. If disabled, no timestamp information is parsed at all.
- `multiline_processing_enabled` - (Optional) Type true to enable; type false to disable. The default setting is true. Consider setting to false to avoid unnecessary processing if you are collecting single message per line files (for example, Linux system.log). If you're working with multiline messages (for example, log4J or exception stack traces), keep this setting enabled.
- `use_autoline_matching` - (Optional) Type true to enable if you'd like message boundaries to be inferred automatically; type false to prevent message boundaries from being automatically inferred (equivalent to the Infer Boundaries option in the UI). The default setting is true.
- `manual_prefix_regexp` - (Optional) When using useAutolineMatching=false, type a regular expression that matches the first line of the message to manually create the boundary. Note that any special characters in the regex, such as backslashes or double quotes, must be escaped.
- `force_timezone` - (Optional) Type true to force the source to use a specific time zone, otherwise type false to use the time zone found in the logs. The default setting is false.
- `default_date_formats` - (Optional) Define formats for the dates present in your log messages. You can specify a locator regex to identify where timestamps appear in log lines. 
- `filters` - (Optional) If you'd like to add a filter to the source, type the name of the filter (Exclude, Include, Mask, Hash, or Forward. 
- `cutoff_timestamp` - (Optional) Only collect data more recent than this timestamp, specified as milliseconds since epoch (13 digit). 
- `cutoff_relative_time` - (Optional) Can be specified instead of cutoffTimestamp to provide a relative offset with respect to the current time. Example: use -1h, -1d, or -1w to collect data that's less than one hour, one day, or one week old, respectively.
- `fields` - (Optional) Map containing [key/value pairs][3], e.g.
```
resource "sumologic_http_source" "instrumentation-logs" {
   ...
   fields = {"origin": "instrumentation"}
}

```

[0]: index.html.markdown
[2]: https://en.wikipedia.org/wiki/Tz_database
[10]: d/caller_identity.html.markdown
[11]: d/collector.html.markdown
[12]: d/personal_folder.html.markdown
[13]: d/role.html.markdown
[20]: r/collector.html.markdown
[21]: r/http_source.html.markdown
[22]: r/polling_source.html.markdown
[23]: r/cloudsyslog_source.html.markdown
[24]: r/collector_ingest_budget_assignment.html.markdown
[25]: r/ingest_budget.html.markdown
[30]: r/role.html.markdown
[31]: r/user.html.markdown
[40]: r/scheduled_view.html.markdown
[41]: r/partition.html.markdown
[42]: r/field_extraction_rule.html.markdown
[43]: r/folder.html.markdown
[44]: r/content.html.markdown
[45]: r/connection.html.markdown
[46]: r/monitor.html.markdown
[47]: r/monitor_folder.html.markdown
[50]: r/cse_network_block.html.markdown
[51]: r/cse_insights_resolution.html.markdown
[52]: r/cse_insights_status.html.markdown
[53]: r/cse_insights_configuration.html.markdown
