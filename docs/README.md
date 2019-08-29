# Table of Contents

#### [Basic usage and authentication][0]

#### Data Sources
  + [sumologic_caller_identity][10]
  + [sumologic_collector][11]

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
- `lookup_by_name` - (Optional) Configures an existent collector using the same 'name' or creates a new one if non existent. Defaults to false.
- `destroy` - (Optional) Whether or not to delete the collector in Sumo when it is removed from Terraform.  Defaults to true.

[0]: sumologic-provider.md
[2]: https://en.wikipedia.org/wiki/Tz_database
[10]: d/sumologic_caller_identity.md
[11]: d/sumologic_collector.md
[20]: r/sumologic_collector.md
[21]: r/sumologic_http_source.md
[22]: r/sumologic_polling_source.md
[23]: r/sumologic_cloudsyslog_source.md
[24]: r/sumologic_collector_ingest_budget_assignment.md
[25]: r/sumologic_ingest_budget.md
[30]: d/sumologic_role.md
[31]: d/sumologic_user.md
