---
layout: 'sumologic'
page_title: 'SumoLogic: sumologic_monitor'
description: |-
  Provides the ability to create, read, delete, and update monitors.
---

# sumologic_monitor

Provides the ability to create, read, delete, and update [Monitors][1].
If Fine Grain Permission (FGP) feature is enabled with Monitors Content at one's Sumo Logic account, one can also set those permission details under this monitor resource. For further details about FGP, please see this [Monitor Permission document][3].

## Example Logs Monitor with FGP

NOTE:
- `obj_permission` are added at one of the monitor's to showcase how Fine Grain Permissions (FGP) are associated with two roles.

```hcl
resource "sumologic_role" "tf_test_role_01" {
  name        = "tf_test_role_01"
  description = "Testing resource sumologic_role"
  capabilities = [
    "viewAlerts",
    "viewMonitorsV2",
    "manageMonitorsV2"
  ]
}
resource "sumologic_role" "tf_test_role_02" {
  name        = "tf_test_role_02"
  description = "Testing resource sumologic_role"
  capabilities = [
    "viewAlerts",
    "viewMonitorsV2",
    "manageMonitorsV2"
  ]
}
resource "sumologic_monitor" "tf_logs_monitor_1" {
  name         = "Terraform Logs Monitor"
  description  = "tf logs monitor"
  type         = "MonitorsLibraryMonitor"
  is_disabled  = false
  content_type = "Monitor"
  monitor_type = "Logs"
  evaluation_delay = "5m"
  tags = {
    "team" = "monitoring"
    "application" = "sumologic"
  }

  queries {
    row_id = "A"
    query  = "_sourceCategory=event-action info"
  }

  trigger_conditions {
    logs_static_condition {
      critical {
        time_range = "15m"
        frequency = "5m"
        alert {
          threshold      = 40.0
          threshold_type = "GreaterThan"
        }
        resolution {
          threshold      = 40.0
          threshold_type = "LessThanOrEqual"
        }
      }
    }
  }

  notifications {
    notification {
      connection_type = "Email"
      recipients = [
        "abc@example.com",
      ]
      subject      = "Monitor Alert: {{TriggerType}} on {{Name}}"
      time_zone    = "PST"
      message_body = "Triggered {{TriggerType}} Alert on {{Name}}: {{QueryURL}}"
    }
    run_for_trigger_types = ["Critical", "ResolvedCritical"]
  }
  notifications {
    notification {
      connection_type = "Webhook"
      connection_id   = "0000000000ABC123"
    }
    run_for_trigger_types = ["Critical", "ResolvedCritical"]
  }
  playbook = "{{Name}} should be fixed in 24 hours when {{TriggerType}} is triggered."
  alert_name = "Alert {{ResultJson.my_field}} from {{Name}}"
  notification_group_fields = ["_sourceHost"]
  obj_permission {
    subject_type = "role"
    subject_id = sumologic_role.tf_test_role_01.id
    permissions = ["Read","Update"]
  }
  obj_permission {
    subject_type = "role"
    subject_id = sumologic_role.tf_test_role_02.id
    permissions = ["Read"]
  }
}
```

## Example Metrics Monitor

```hcl
resource "sumologic_monitor" "tf_metrics_monitor_1" {
  name         = "Terraform Metrics Monitor"
  description  = "tf metrics monitor"
  type         = "MonitorsLibraryMonitor"
  is_disabled  = false
  content_type = "Monitor"
  monitor_type = "Metrics"
  evaluation_delay = "1m"
  tags = {
    "team" = "monitoring"
    "application" = "sumologic"
  }

  queries {
    row_id = "A"
    query  = "metric=CPU* _sourceCategory=event-action"
  }

  trigger_conditions {
    metrics_static_condition {
      critical {
        time_range = "15m"
        occurrence_type = "Always"
        alert {
          threshold      = 40.0
          threshold_type = "GreaterThan"
          min_data_points = 5
        }
        resolution {
          threshold      = 40.0
          threshold_type = "LessThanOrEqual"
        }
      }
    }
  }
  notifications {
    notification {
      connection_type = "Email"
      recipients      = ["abc@example.com"]
      subject         = "Triggered {{TriggerType}} Alert on Monitor {{Name}}"
      time_zone       = "PST"
      message_body    = "Triggered {{TriggerType}} Alert on {{Name}}: {{QueryURL}}"
    }
    run_for_trigger_types = ["Critical", "ResolvedCritical"]
  }
  playbook = "test playbook"
  notification_group_fields = ["metric"]
}
```


## Example SLO Monitors

```hcl
resource "sumologic_monitor" "tf_slo_monitor_1" {
  name         = "SLO SLI monitor"
  type         = "MonitorsLibraryMonitor"
  is_disabled  = false
  content_type = "Monitor"
  monitor_type = "Slo"
  slo_id = "0000000000000009"
  evaluation_delay = "5m"
  tags = {
    "team" = "monitoring"
    "application" = "sumologic"
  }

  trigger_conditions {
    slo_sli_condition {
      critical {
        sli_threshold =  99.5
      }
      warning {
        sli_threshold =  99.9
      }
    }
  }

  notifications {
    notification {
      connection_type = "Email"
      recipients      = ["abc@example.com"]
      subject      = "Monitor Alert: {{TriggerType}} on {{Name}}"
      time_zone    = "PST"
      message_body = "Triggered {{TriggerType}} Alert on {{Name}}: {{QueryURL}}"
    }
    run_for_trigger_types = ["Critical", "ResolvedCritical"]
  }
  playbook = "test playbook"
}

resource "sumologic_monitor" "tf_slo_monitor_2" {
  name         = "SLO Burn rate monitor"
  type         = "MonitorsLibraryMonitor"
  is_disabled  = false
  content_type = "Monitor"
  monitor_type = "Slo"
  slo_id = "0000000000000009"
  evaluation_delay = "5m"
  tags = {
    "team" = "monitoring"
    "application" = "sumologic"
  }

  trigger_conditions {
    slo_burn_rate_condition {
      critical {
        burn_rate {
            burn_rate_threshold =  50
            time_range = "1d"
        }
      }
      warning {
        burn_rate {
            burn_rate_threshold =  30
            time_range = "3d"
        }
        burn_rate {
            burn_rate_threshold =  20
            time_range = "4d"
        }
      }
    }
  }

  #...
}

```

## Example Logs Monitor with Webhook Connection and Folder

```hcl
resource "sumologic_monitor_folder" "tf_monitor_folder_1" {
  name        = "Terraform Managed Folder 1"
  description = "A folder for Monitors"
}

resource "sumologic_connection" "example_pagerduty_connection" {
  name            = "example_pagerduty_connection"
  description     = "PagerDuty connection for notifications from Monitors"
  type            = "WebhookConnection"
  webhook_type    = "PagerDuty"
  url             = "https://events.pagerduty.com/"
  default_payload = <<JSON
{
  "service_key": "pagerduty_api_integration_key",
  "event_type": "trigger",
  "description": "PagerDuty connection for notifications",
  "client": "Sumo Logic",
  "client_url": ""
}
JSON
}

resource "sumologic_monitor" "tf_logs_monitor_2" {
  name         = "Terraform Logs Monitor with Webhook Connection"
  description  = "tf logs monitor with webhook"
  type         = "MonitorsLibraryMonitor"
  parent_id    = sumologic_monitor_folder.tf_monitor_folder_1.id
  is_disabled  = false
  content_type = "Monitor"
  monitor_type = "Logs"
  tags = {
    "team" = "monitoring"
    "application" = "sumologic"
  }
  queries {
    row_id = "A"
    query  = "_sourceCategory=event-action info"
  }
  trigger_conditions {
    logs_static_condition {
      critical {
        time_range = "15m"
        frequency = "5m"
        alert {
          threshold      = 40.0
          threshold_type = "GreaterThan"
        }
        resolution {
          threshold      = 40.0
          threshold_type = "LessThanOrEqual"
          resolution_window = "5m"
        }
      }
    }
  }
  notifications {
    notification {
      connection_type = "Email"
      recipients = [
        "abc@example.com",
      ]
      subject      = "Monitor Alert: {{TriggerType}} on {{Name}}"
      time_zone    = "PST"
      message_body = "Triggered {{TriggerType}} Alert on {{Name}}: {{QueryURL}}"
    }
    run_for_trigger_types = ["Critical", "ResolvedCritical"]
  }
  notifications {
    notification {
      connection_type  = "PagerDuty"
      connection_id    = sumologic_connection.example_pagerduty_connection.id
      payload_override = <<JSON
{
  "service_key": "your_pagerduty_api_integration_key",
  "event_type": "trigger",
  "description": "Alert: Triggered {{TriggerType}} for Monitor {{Name}}",
  "client": "Sumo Logic",
  "client_url": "{{QueryUrl}}"
}
JSON
     resolution_payload_override = <<JSON
{
  "service_key": "your_pagerduty_api_integration_key",
  "event_type": "trigger",
  "description": "Alert: Resolved {{TriggerType}} for Monitor {{Name}}",
  "client": "Sumo Logic",
  "client_url": "{{QueryUrl}}"
}
JSON
    }
    run_for_trigger_types = ["Critical", "ResolvedCritical"]
  }
}
```


## Example Logs Anomaly Monitor
```hcl
resource "sumologic_monitor" "tf_example_anomaly_monitor" {
  name = "Example Anomaly Monitor"
  description = "example anomaly monitor"
  type = "MonitorsLibraryMonitor"
  monitor_type = "Logs"
  is_disabled = false

  queries {
      row_id = "A"
      query = "_sourceCategory=api error | timeslice 5m | count by _sourceHost"
  }

  trigger_conditions {
    logs_anomaly_condition {
      field = "_count"
      anomaly_detector_type = "Cluster"
      critical {
        sensitivity = 0.4
        min_anomaly_count = 9
        time_range = "-3h"
      }
    }
  }

  notifications {
    notification {
      connection_type = "Email"
      recipients = [
        "anomaly@example.com",
      ]
      subject      = "Monitor Alert: {{TriggerType}} on {{Name}}"
      time_zone    = "PST"
      message_body = "Triggered {{TriggerType}} Alert on {{Name}}: {{QueryURL}}"
    }
    run_for_trigger_types = ["Critical", "ResolvedCritical"]
  }
}
```

## Example Metrics Anomaly Monitor
```hcl
resource "sumologic_monitor" "tf_example_metrics_anomaly_monitor" {
  name = "Example Metrics Anomaly Monitor"
  description = "example metrics anomaly monitor"
  type = "MonitorsLibraryMonitor"
  monitor_type = "Metrics"
  is_disabled = false

  queries {
      row_id = "A"
      query = "service=auth api=login metric=HTTP_5XX_Count | avg"
  }

  trigger_conditions {
    metrics_anomaly_condition {
      anomaly_detector_type = "Cluster"
      critical {
        sensitivity = 0.4
        min_anomaly_count = 9
        time_range = "-3h"
      }
    }
  }

  notifications {
    notification {
      connection_type = "Email"
      recipients = [
        "anomaly@example.com",
      ]
      subject      = "Monitor Alert: {{TriggerType}} on {{Name}}"
      time_zone    = "PST"
      message_body = "Triggered {{TriggerType}} Alert on {{Name}}: {{QueryURL}}"
    }
    run_for_trigger_types = ["Critical", "ResolvedCritical"]
  }
}
```

## Monitor Folders

NOTE: Monitor folders are considered a different resource from Library content folders. See [sumologic_monitor_folder][2] for more details.

## Argument reference

The following arguments are supported:

- `type` - (Optional) The type of object model. Valid value:
  - `MonitorsLibraryMonitor`
- `name` - (Required) The name of the monitor. The name must be alphanumeric.
- `description` - (Optional) The description of the monitor.
- `is_disabled` - (Optional) Whether or not the monitor is disabled. Disabled monitors will not run and will not generate or send notifications.
- `parent_id` - (Optional) The ID of the Monitor Folder that contains this monitor. Defaults to the root folder.
- `content_type` - (Optional) The type of the content object. Valid value:
  - `Monitor`
- `monitor_type` - (Required) The type of monitor. Valid values:
  - `Logs`: A logs query monitor.
  - `Metrics`: A metrics query monitor.
  - `Slo`: A SLO based monitor.
- `tags` - (Optional) A map defining tag keys and tag values for the Monitor.
- `evaluation_delay` - (Optional) Evaluation delay as a string consists of the following elements:
      1. `<number>`: number of time units,
      2. `<time_unit>`: time unit; possible values are: `h` (hour), `m` (minute), `s` (second).

      Multiple pairs of `<number><time_unit>` may be provided. For example,
      `2m50s` means 2 minutes and 50 seconds.
- `slo_id` - (Optional) Identifier of the SLO definition for the monitor. This is only applicable & required for Slo `monitor_type`.
- `queries` - (Required if `monitor_type` is not `Slo`) All queries from the monitor.
- `trigger_conditions` - (Required if not using `triggers`) Defines the conditions of when to send notifications. NOTE: `trigger_conditions` supplants the `triggers` argument.
  - `resolution_window` - The resolution window that the recovery condition must be met in each evaluation that happens within this entire duration before the alert is recovered (resolved). If not specified, the time range of your trigger will be used.
- `triggers` - (Deprecated) Defines the conditions of when to send notifications.
- `notifications` - (Optional) The notifications the monitor will send when the respective trigger condition is met.
- `group_notifications` - (Optional) Whether or not to group notifications for individual items that meet the trigger condition. Defaults to true.
- `playbook` - (Optional - Beta) Notes such as links and instruction to help you resolve alerts triggered by this monitor. {{Markdown}} supported. It will be enabled only if available for your organization. Please contact your Sumo Logic account team to learn more.
- `alert_name` - (Optional) The display name when creating alerts. Monitor name will be used if `alert_name` is not provided. All template variables can be used in `alert_name` except `{{AlertName}}`, `{{AlertResponseURL}}`, `{{ResultsJson}}`, and `{{Playbook}}`.
- `notification_group_fields` - (Optional) The set of fields to be used to group alerts and notifications for a monitor. The value of this field will be considered only when 'groupNotifications' is true. The fields with very high cardinality such as `_blockid`, `_raw`, `_messagetime`, `_receipttime`, and `_messageid` are not allowed for Alert Grouping.
- `obj_permission` - (Optional) `obj_permission` construct represents a Permission Statement associated with this Monitor. A set of `obj_permission` constructs can be specified under a Monitor. An `obj_permission` construct can be used to control permissions Explicitly associated with a Monitor. But, it cannot be used to control permissions Inherited from a Parent / Ancestor. Default FGP would be still set to the Monitor upon creation (e.g. the creating user would have full permission), even if no `obj_permission` construct is specified at a Monitor and the FGP feature is enabled at the account.
    - `subject_type` - (Required) Valid values:
        - `role`
        - `org`
    - `subject_id` - (Required) A Role ID or the Org ID of the account
    - `permissions` - (Required) A Set of Permissions. Valid Permission Values:
        - `Read`
        - `Update`
        - `Delete`
        - `Manage`

Additional data provided in state:

- `id` - (Computed) The ID for this monitor.
- `status` - (Computed) The current status for this monitor. Values are:
  - `Critical`
  - `Warning`
  - `MissingData`
  - `Normal`
  - `Disabled`

## The `trigger_conditions` block
A `trigger_conditions` block configures conditions for sending notifications.
### Example
```hcl
trigger_conditions {
  logs_static_condition {
    field = "_count"
    critical {
      time_range = "15m"
      frequency = "5m"
      alert {
        threshold = 100
        threshold_type = "GreaterThan"
      }
      resolution {
        threshold = 90
        threshold_type = "LessThanOrEqual"
        resolution_window = "5m"
      }
    }
    warning {
      time_range = "30m"
      alert {
        threshold = 80
        threshold_type = "GreaterThan"
      }
      resolution {
        threshold = 75
        threshold_type = "LessThanOrEqual"
        resolution_window = "5m"
      }
    }
  }
  logs_missing_data_condition {
    time_range = "30m"
    frequency = "5m"
  }
}
```
### Arguments
A `trigger_conditions` block contains one or more subblocks of the following types:
- `logs_static_condition`
- `metrics_static_condition`
- `logs_outlier_condition`
- `metrics_outlier_condition`
- `logs_missing_data_condition`
- `metrics_missing_data_condition`
- `slo_sli_condition`
- `slo_burn_rate_condition`
- `logs_anomaly_condition`
- `metrics_anomaly_condition`

Subblocks should be limited to at most 1 missing data condition and at most 1 static / outlier condition.

Here is a summary of arguments for each condition type (fields which are not marked as `Required` are optional):
#### logs_static_condition
  - `field`
  - `critical`
    - `time_range` (Required) : Accepted format: Optional `-` sign followed by `<number>` followed by a `<time_unit>` character: `s` for seconds, `m` for minutes, `h` for hours, `d` for days. Examples: `30m`, `-12h`.
    - `frequency` Accepted format: Optional `-` sign followed by `<number>` followed by a `<time_unit>` character: `s` for seconds, `m` for minutes, `h` for hours, `d` for days. Examples: `1m`, `2m`, `10m`'.
    - `alert` (Required)
      - `threshold`
      - `threshold_type`
    - `resolution` (Required)
      - `threshold`
      - `threshold_type`
      - `resolution_window` Accepted format: `<number>` followed by a `<time_unit>` character: `s` for seconds, `m` for minutes, `h` for hours, `d` for days. Examples: `0s, 30m`.
  - `warning`
    - `time_range` (Required) :  Accepted format: Optional `-` sign followed by `<number>` followed by a `<time_unit>` character: `s` for seconds, `m` for minutes, `h` for hours, `d` for days. Examples: `30m`, `-12h`.
    - `alert` (Required)
      - `threshold`
      - `threshold_type`
    - `resolution` (Required)
      - `threshold`
      - `threshold_type`
      - `resolution_window` Accepted format: `<number>` followed by a `<time_unit>` character: `s` for seconds, `m` for minutes, `h` for hours, `d` for days. Examples: `0s, 30m`.
#### metrics_static_condition
  - `critical`
    - `time_range` (Required) :  Accepted format: Optional `-` sign followed by `<number>` followed by a `<time_unit>` character: `s` for seconds, `m` for minutes, `h` for hours, `d` for days. Examples: `30m`, `-12h`.
    - `occurrence_type` (Required)
    - `alert` (Required)
      - `threshold`
      - `threshold_type`
      - `min_data_points` (Optional)
    - `resolution` (Required)
      - `threshold`
      - `threshold_type`
      - `min_data_points` (Optional)
    - `warning`
    - `time_range` (Required) :  Accepted format: Optional `-` sign followed by `<number>` followed by a `<time_unit>` character: `s` for seconds, `m` for minutes, `h` for hours, `d` for days. Examples: `30m`, `-12h`.
    - `occurrence_type` (Required)
    - `alert` (Required)
      - `threshold`
      - `threshold_type`
      - `min_data_points` (Optional)
    - `resolution` (Required)
      - `threshold`
      - `threshold_type`
      - `min_data_points` (Optional)
#### logs_outlier_condition
  - `field`
  - `direction`
  - `critical`
     - `window`
     - `consecutive`
     - `threshold`
  - `warning`
     - `window`
     - `consecutive`
     - `threshold`
#### metrics_outlier_condition
  - `direction`
  - `critical`
     - `baseline_window`
     - `threshold`
  - `warning`
    - `baseline_window`
    - `threshold`
#### logs_missing_data_condition
  - `time_range` (Required) :  Accepted format: Optional `-` sign followed by `<number>` followed by a `<time_unit>` character: `s` for seconds, `m` for minutes, `h` for hours, `d` for days. Examples: `30m`, `-12h`.
  - `frequency` Accepted format: Optional `-` sign followed by `<number>` followed by a `<time_unit>` character: `s` for seconds, `m` for minutes, `h` for hours, `d` for days. Examples: `1m`, `2m`, `10m`'.
#### metrics_missing_data_condition
  - `time_range` (Required) :  Accepted format: Optional `-` sign followed by `<number>` followed by a `<time_unit>` character: `s` for seconds, `m` for minutes, `h` for hours, `d` for days. Examples: `30m`, `-12h`.
#### slo_sli_condition
  - `critical`
    - `sli_threshold` (Required) : The remaining SLI error budget threshold percentage [0,100).
  - `warning`
    - `sli_threshold` (Required)

#### slo_burn_rate_condition
  - `critical`
    - `time_range` (Deprecated) : The relative time range for the burn rate percentage evaluation.  Accepted format: Optional `-` sign followed by `<number>` followed by a `<time_unit>` character: `s` for seconds, `m` for minutes, `h` for hours, `d` for days. Examples: `30m`, `-12h`.
    - `burn_rate_threshold` (Deprecated) : The burn rate percentage threshold.
    - `burn_rate` (Required if above two fields are not present): Block to specify burn rate threshold and time range for the condition.
      - `burn_rate_threshold` (Required): The burn rate percentage threshold.
      - `time_range` (Required): The relative time range for the burn rate percentage evaluation.  Accepted format: Optional `-` sign followed by `<number>` followed by a `<time_unit>` character: `s` for seconds, `m` for minutes, `h` for hours, `d` for days. Examples: `30m`, `-12h`.
  - `warning`
    - `time_range` (Deprecated) :  Accepted format: Optional `-` sign followed by `<number>` followed by a `<time_unit>` character: `s` for seconds, `m` for minutes, `h` for hours, `d` for days. Examples: `30m`, `-12h`.
    - `burn_rate_threshold` (Deprecated)
    - `burn_rate` (Required if above two fields are not present): Block to specify burn rate threshold and time range for the condition.
      - `burn_rate_threshold` (Required): The burn rate percentage threshold.
      - `time_range` (Required): The relative time range for the burn rate percentage evaluation.  Accepted format: Optional `-` sign followed by `<number>` followed by a `<time_unit>` character: `s` for seconds, `m` for minutes, `h` for hours, `d` for days. Examples: `30m`, `-12h`.

#### logs_anomaly_condition
  - `field`: The name of the field that the trigger condition will alert on. The trigger could compare the value of specified field with the threshold. If field is not specified, monitor would default to result count instead.
  - `anomaly_detector_type`: The type of anomaly model that will be used for evaluating this monitor. Possible values are: `Cluster`.
  - `critical`
    - `sensitivity`: The triggering sensitivity of the anomaly model used for this monitor.
    - `min_anomaly_count` (Required) : The minimum number of anomalies required to exist in the current time range for the condition to trigger.
    - `time_range` (Required) : The relative time range for anomaly evaluation.  Accepted format: Optional `-` sign followed by `<number>` followed by a `<time_unit>` character: `s` for seconds, `m` for minutes, `h` for hours, `d` for days. Examples: `30m`, `-12h`.

#### metrics_anomaly_condition
- `anomaly_detector_type`: The type of anomaly model that will be used for evaluating this monitor. Possible values are: `Cluster`.
- `critical`
    - `sensitivity`: The triggering sensitivity of the anomaly model used for this monitor.
    - `min_anomaly_count` (Required) : The minimum number of anomalies required to exist in the current time range for the condition to trigger.
    - `time_range` (Required) : The relative time range for anomaly evaluation.  Accepted format: Optional `-` sign followed by `<number>` followed by a `<time_unit>` character: `s` for seconds, `m` for minutes, `h` for hours, `d` for days. Examples: `30m`, `-12h`.

## The `triggers` block
The `triggers` block is deprecated. Please use `trigger_conditions` to specify notification conditions.

Here's an example logs monitor that uses `triggers` to specify trigger conditions:
```hcl
resource "sumologic_monitor" "tf_logs_monitor_1" {
  name         = "Terraform Logs Monitor"
  description  = "tf logs monitor"
  type         = "MonitorsLibraryMonitor"
  is_disabled  = false
  content_type = "Monitor"
  monitor_type = "Logs"
  queries {
    row_id = "A"
    query  = "_sourceCategory=event-action info"
  }
  triggers {
    threshold_type   = "GreaterThan"
    threshold        = 40.0
    time_range       = "15m"
    occurrence_type  = "ResultCount"
    trigger_source   = "AllResults"
    trigger_type     = "Critical"
    detection_method = "StaticCondition"
  }
  triggers {
    threshold_type   = "LessThanOrEqual"
    threshold        = 40.0
    time_range       = "15m"
    occurrence_type  = "ResultCount"
    trigger_source   = "AllResults"
    trigger_type     = "ResolvedCritical"
    detection_method = "StaticCondition"
    resolution_window = "5m"
  }
  notifications {
    notification {
      connection_type = "Email"
      recipients = [
        "abc@example.com",
      ]
      subject      = "Monitor Alert: {{TriggerType}} on {{Name}}"
      time_zone    = "PST"
      message_body = "Triggered {{TriggerType}} Alert on {{Name}}: {{QueryURL}}"
    }
    run_for_trigger_types = ["Critical", "ResolvedCritical"]
  }
  notifications {
    notification {
      connection_type = "Webhook"
      connection_id   = "0000000000ABC123"
    }
    run_for_trigger_types = ["Critical", "ResolvedCritical"]
  }
}
```

## Import

Monitors can be imported using the monitor ID, such as:

```hcl
terraform import sumologic_monitor.test 1234567890
```

[1]: https://help.sumologic.com/?cid=10020
[2]: monitor_folder.html.markdown
[3]: https://help.sumologic.com/Visualizations-and-Alerts/Alerts/Monitors#configure-permissions-for-a-monitor
